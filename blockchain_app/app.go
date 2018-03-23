package blockchain_app

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
)

// Return codes for tendermint core
const (
	CodeTypeOK            uint32 = 0
	CodeTypeEncodingError uint32 = 1
	CodeTypeBadData       uint32 = 2
	CodeTypeUnauthorized  uint32 = 3
)

// BlockChainApplication represents the block chain app
type BlockChainApplication struct {
	types.BaseApplication

	state State
}

// NewBlockChainApplication starts a blockchain application with state loaded from db
func NewBlockChainApplication() *BlockChainApplication {
	//TODO: find a better storage dir
	dir, err := ioutil.TempDir("/tmp", "abci-car")
	if err != nil {
		panic(err)
	}

	name := "car"
	db, err := dbm.NewGoLevelDB(name, dir)

	if err != nil {
		panic(err)
	}

	state := loadState(db)
	return &BlockChainApplication{state: state}
}

// Info returns the app state
// It is used by the Query connection
// http://tendermint.readthedocs.io/projects/tools/en/master/app-development.html?highlight=Info#query-connection
// https://tendermint.readthedocs.io/en/master/app-development.html#handshake
func (app *BlockChainApplication) Info(req types.RequestInfo) types.ResponseInfo {
	return types.ResponseInfo{LastBlockHeight: app.state.LastBlockHeight, LastBlockAppHash: app.state.LastBlockAppHash}
}

// SetOption is not implemented
func (app *BlockChainApplication) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	// Not yet implemented, does nothing
	return types.ResponseSetOption{}
}

// DeliverTx submits the data to put into the blockchain
// http://tendermint.readthedocs.io/projects/tools/en/master/app-development.html#delivertx
// eg transaction: original json: {"operation": "createCar", "data": {"ID": "car3", "Make": "Toyota", "Model": "Prius", "Colour": "blue", "Owner": "Tomoko"}}
// base64encoded: eyJvcGVyYXRpb24iOiAiY3JlYXRlQ2FyIiwgImRhdGEiOiB7IklEIjogImNhcjMiLCAiTWFrZSI6ICJUb3lvdGEiLCAiTW9kZWwiOiAiUHJpdXMiLCAiQ29sb3VyIjogImJsdWUiLCAiT3duZXIiOiAiVG9tb2tvIn19
func (app *BlockChainApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {

	transaction, err := decodePayload(tx)

	if err != nil {
		return types.ResponseDeliverTx{
			Code: CodeTypeEncodingError,
			Log:  err.(*ErrTransactionDecoding).message}
	}

	log.Println("We are in deliver txt with transaction", transaction)

	if transaction.Operation == OpCreateCar {
		log.Println("operation is create car")

		return createCar(app, transaction)
	} else if transaction.Operation == OpChangeCarOwner {
		log.Println("operation is change car owner")

		return changeCarOwner(app, transaction)
	}

	return types.ResponseDeliverTx{
		Code: CodeTypeBadData,
		Log:  fmt.Sprintf("Invalid operation %v", transaction.Operation)}
}

// CheckTx validates the transaction
func (app *BlockChainApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	log.Println("We are in check txt")

	transaction, err := decodePayload(tx)

	if err != nil {
		return types.ResponseCheckTx{
			Code: CodeTypeEncodingError,
			Log:  err.(*ErrTransactionDecoding).message}
	}

	if transaction.Operation == OpCreateCar {
		log.Println("operation to check is create car")
		// no validation here
		return types.ResponseCheckTx{Code: CodeTypeOK}
	} else if transaction.Operation == OpChangeCarOwner {
		log.Println("operation to check is change car owner")

		return validateChangeCarOwner(app, transaction)
	}

	return types.ResponseCheckTx{Code: CodeTypeOK}
}

// Commit saves the application states and returns the LastBlockAppHash
// http://tendermint.readthedocs.io/projects/tools/en/master/app-development.html#commit
func (app *BlockChainApplication) Commit() (resp types.ResponseCommit) {
	log.Println("We are in commit")

	//update and save app state
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.LastBlockAppHash = appHash
	app.state.LastBlockHeight++
	saveState(app.state)

	return types.ResponseCommit{Data: app.state.LastBlockAppHash}
}

// Query returns responses to state database query searches
// http://tendermint.readthedocs.io/projects/tools/en/master/app-development.html#query-connection
func (app *BlockChainApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	log.Printf("We are in query function with data: %s, path: %s", reqQuery.Data, reqQuery.Path)

	switch reqQuery.Path {
	// query looks like localhost:46657/abci_query?path="allCars"'
	case "allCars":
		return getAllCars(app, reqQuery)
	default:
		return types.ResponseQuery{Code: CodeTypeBadData, Value: []byte(cmn.Fmt("Invalid query path: %v", reqQuery.Path))}
	}
}

func decodePayload(tx []byte) (Transaction, error) {
	var transaction Transaction

	payload, err := base64.StdEncoding.DecodeString(string(tx))

	log.Println("Payload received: ", string(payload), string(tx))

	if err != nil {
		return transaction, &ErrTransactionDecoding{fmt.Sprintf("Cannot decode payload %s, we got %v", string(tx), err)}
	}

	//try to unserialize the data
	err = json.Unmarshal(payload, &transaction)

	if err != nil {
		return transaction, &ErrTransactionDecoding{message: fmt.Sprintf("Cannot json unserialize %s, we got %v", string(payload), err)}
	}

	return transaction, nil
}

func validateChangeCarOwner(app *BlockChainApplication, transaction Transaction) types.ResponseCheckTx {
	var changeOwnerPayload ChangeOwnerPayload

	err := json.Unmarshal(transaction.Data, &changeOwnerPayload)

	if err != nil {
		return types.ResponseCheckTx{
			Code: CodeTypeEncodingError,
			Log:  fmt.Sprintf("Cannot json unserialize %s into  change car payload, we got %v", string(transaction.Data), err)}
	}

	record := app.state.GetCar(changeOwnerPayload.AssetID)

	if len(record) == 0 {
		return types.ResponseCheckTx{
			Code: CodeTypeBadData,
			Log:  fmt.Sprintf("Asset with id %s was not found", changeOwnerPayload.AssetID)}
	}

	return types.ResponseCheckTx{Code: CodeTypeOK}
}

func createCar(app *BlockChainApplication, transaction Transaction) types.ResponseDeliverTx {
	var transactionDbData DbData

	var asset AssetCar
	err := json.Unmarshal(transaction.Data, &asset)

	if err != nil {
		return types.ResponseDeliverTx{
			Code: CodeTypeEncodingError,
			Log:  fmt.Sprintf("Cannot json unserialize %s into asset car, we got %v", string(transaction.Data), err)}
	}

	app.state.SaveCar(asset.ID, transaction.Data)
	app.state.Size++

	log.Println("we have saved the data in db")

	return types.ResponseDeliverTx{Code: CodeTypeOK, Data: transactionDbData.Value}
}

func changeCarOwner(app *BlockChainApplication, transaction Transaction) types.ResponseDeliverTx {
	var changeOwnerPayload ChangeOwnerPayload

	err := json.Unmarshal(transaction.Data, &changeOwnerPayload)

	if err != nil {
		return types.ResponseDeliverTx{
			Code: CodeTypeEncodingError,
			Log:  fmt.Sprintf("Cannot json unserialize %s into  change car payload, we got %v", string(transaction.Data), err)}
	}

	//search for car record
	record := app.state.GetCar(changeOwnerPayload.AssetID)

	if len(record) == 0 {
		return types.ResponseDeliverTx{
			Code: CodeTypeBadData,
			Log:  fmt.Sprintf("Asset with id %s was not found", changeOwnerPayload.AssetID)}
	}

	var asset AssetCar

	err = json.Unmarshal(record, &asset)

	if err != nil {
		return types.ResponseDeliverTx{
			Code: CodeTypeBadData,
			Log:  fmt.Sprintf("Cannot decode asset from db record %s", string(record))}
	}

	//update owner
	asset.Owner = changeOwnerPayload.NewOwner

	record, err = json.Marshal(asset)

	if err != nil {
		return types.ResponseDeliverTx{
			Code: CodeTypeBadData,
			Log:  fmt.Sprintf("Cannot encode asset %v", asset)}
	}

	app.state.SaveCar(asset.ID, record)
	app.state.Size++

	log.Println("we have saved the data in db")

	return types.ResponseDeliverTx{Code: CodeTypeOK, Data: record}

}

func getAllCars(app *BlockChainApplication, reqQuery types.RequestQuery) types.ResponseQuery {
	allCars, err := app.state.GetAllCars()

	if err != nil {
		return types.ResponseQuery{Code: CodeTypeBadData, Value: []byte(cmn.Fmt("Cannot retrieve data from database: %v", err))}
	}

	responseData, err := json.Marshal(allCars)

	if err != nil {
		return types.ResponseQuery{Code: CodeTypeBadData, Value: []byte(cmn.Fmt("Cannot serialize data: %v", err))}
	}

	log.Println("We have found", string(responseData))

	return types.ResponseQuery{Code: CodeTypeOK, Value: responseData}
}
