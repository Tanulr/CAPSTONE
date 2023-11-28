package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Check struct {
	ID 		uint64      `json:"id"` // string
	Pack    string   `json:"pack"`
}


// InitLedger run by manurfacturer. Sets assets list. Whoever runs this is the 
// default (single) manufacturer
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// asset {ChipID(primary key), "asset1", ChipName: "Intel Core i9 10900K", Owner: "Tomoko", Quantity: 5, Value: 300},

	//assets := []Asset{}
	// gotta research on what getCreator(), getID(), getMSPID() returns
	manufacturers := []Check{
		{ID: 1, Pack: ctx.GetClientIdentity().GetCreator()},
		{ID: 2, Pack: ctx.GetClientIdentity().GetID()},
		{ID: 3, Pack: ctx.GetClientIdentity().GetMSPID()}
	}

	

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	for _, manufacturer := range manufacturers {
		manufacturerJSON, err := json.Marshal(manufacturer)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(manufacturer.ID, manufacturerJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}


// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Check, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange(1, 4)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var checks []*Check
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var check Check
		err = json.Unmarshal(queryResponse.Value, &check)
		if err != nil {
			return nil, err
		}
		assets = append(checks, &check)
	}

	return checks, nil
}