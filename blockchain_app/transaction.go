package blockchain_app

import "encoding/json"

const (
	OpCreateCar      string = "createCar"
	OpChangeCarOwner string = "changeCarOwner"
)

type DbData struct {
	Key   []byte
	Value []byte
}

// AssetCar represents the car data that is saved on the ledger
type AssetCar struct {
	ID     string `json:"id"`
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

// ChangeOwnerPayload represents the data of the changeOwner Transaction
type ChangeOwnerPayload struct {
	AssetID  string `json:"asset_id"`
	NewOwner string `json:"new_owner"`
}

// Transaction represents the type of transactions that can be applied on an asset
type Transaction struct {
	Operation string
	Data      json.RawMessage
}

// ErrTransactionDecoding represents an error occured when trying to decode a transaction payload
type ErrTransactionDecoding struct {
	message string
}

func (err *ErrTransactionDecoding) Error() string {
	return err.message
}
