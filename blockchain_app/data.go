package blockchain_app

import (
	"encoding/json"
	"strings"

	dbm "github.com/tendermint/tmlibs/db"
)

const (
	// CarAssetPrefix represents the db prefix key for car assets
	CarAssetPrefix = "Car:"
)

// State holds the application state (database, and last block height and app hash for handshake)
// https://tendermint.readthedocs.io/en/master/app-development.html#handshake
type State struct {
	db               dbm.DB
	LastBlockHeight  int64  `json:"last_block_height"`   //the last block for which the app ran Commit successfully
	LastBlockAppHash []byte `json:"last_block_app_hash"` //the response from the last successfull Commit
	Size             int64  `json:"size"`                //we use it to calculate the app hash
}

// Populates last block from data storage
func loadState(db dbm.DB) State {
	stateData := db.Get([]byte("stateData"))
	var state State
	if len(stateData) != 0 {
		err := json.Unmarshal(stateData, &state)

		if err != nil {
			panic(err)
		}
	}
	state.db = db

	return state
}

// Saves last block to data storage
func saveState(state State) {
	stateData, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set([]byte("stateData"), stateData)
}

// GetAllCars retrieves all cars records saved in state db
func (state *State) GetAllCars() (map[string]AssetCar, error) {

	iter := state.db.Iterator(nil, nil)

	var data = make(map[string]AssetCar)

	//iterate over all keys that start with CarAssetPrefix
	for ; iter.Valid(); iter.Next() {
		if strings.HasPrefix(string(iter.Key()), CarAssetPrefix) {

			// we decode the data so we can create a proper json in response, if we use iter.Value() directly we get an unusable json
			var asset AssetCar
			err := json.Unmarshal(iter.Value(), &asset)

			if err != nil {
				return data, err
			}

			data[string(iter.Key())] = asset
		}
	}
	return data, nil
}

// GetCar returns a car record from db, return empty if no record found
func (state *State) GetCar(ID string) []byte {
	//search for car record
	key := CarAssetPrefix + ID // something like Car:car3
	return state.db.Get([]byte(key))
}

// SaveCar saves a car record on db
func (state *State) SaveCar(ID string, data []byte) {
	key := []byte(CarAssetPrefix + ID) // something like Car:car3
	value := data

	state.db.Set(key, value)
}
