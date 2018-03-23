package blockchain_app

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	dbm "github.com/tendermint/tmlibs/db"
)

func createStateWithEmptyDb() State {
	//create empty db
	dir, err := ioutil.TempDir("/tmp", "abci-car-testing")
	if err != nil {
		panic(err)
	}

	name := "car-testing"
	db, err := dbm.NewGoLevelDB(name, dir)

	if err != nil {
		panic(err)
	}

	state := State{db: db, LastBlockHeight: 0, LastBlockAppHash: []byte("")}

	return state
}

func prepareTestData() (AssetCar, []byte) {
	assetCar := AssetCar{ID: "car1"}

	assetCarData, _ := json.Marshal(assetCar)

	return assetCar, assetCarData
}

func TestSaveCar(t *testing.T) {
	state := createStateWithEmptyDb()

	state.SaveCar("car1", []byte("car-data"))
	carData := state.GetCar("car1")

	if string(carData) != "car-data" {
		t.Errorf("Data was not properly saved in db, we got %s instead of %s", string(carData), "car-data")
	}

}

func TestGetAllCars(t *testing.T) {
	state := createStateWithEmptyDb()

	assetCar, assetCarData := prepareTestData()
	state.SaveCar(assetCar.ID, assetCarData)
	allCars, err := state.GetAllCars()

	if err != nil {
		panic(err)
	}

	if len(allCars) != 1 {
		t.Error("We should have got a map with 1 record, we got", len(allCars))
	}

	if allCars[CarAssetPrefix+assetCar.ID].ID != assetCar.ID {
		t.Errorf("The result is not correct, it should have been %s , we got %s", assetCar.ID, allCars[assetCar.ID].ID)
	}

}
