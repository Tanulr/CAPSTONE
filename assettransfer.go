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

type Asset struct {
	ChipID         string `json:"id"` // UNIQUE asset#
	ChipName       string `json:"chipName"`
	Owner          string `json:"owner"`
	Quantity       uint64    `json:"quantity"`
	VerifyValue 		   string    `json:"value"`
}

type Check struct {
	ID 		string      `json:"id"` // string
	Pack    string   `json:"pack"`
}



// InitLedger run by manurfacturer. Sets assets list. Whoever runs this is the 
// default (single) manufacturer. [invoke]
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// asset {ChipID(primary key), "asset1", ChipName: "Intel Core i9 10900K", Owner: "Tomoko", Quantity: 5, Value: 300},

	assets := []Asset{}
	// gotta research on what getCreator(), getID(), getMSPID() returns
	manufacturer := Check{ID: 1, Pack: ctx.GetClientIdentity().GetMSPID()}

	

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
	
	manufacturerJSON, err := json.Marshal(manufacturer)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(manufacturer.ID, manufacturerJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

// return txn callers identity [query]
func (s *SmartContract) getIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	id, err = ctx.GetClientIdentity().GetMSPID()
	return id, err
}

// CreateAsset issues a new asset to the world state with given details. [invoke]
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, ID string, Name string, Qty int, Val string) error {
	manufacturer := ctx.GetClientIdentity().GetMSPID()
	if manufacturer != ctx.GetStub().GetState("1"){
		return fmt.Errorf("this entity is not a manufacturer and cannot create assets")
	}

	exists, err := s.AssetExists(ctx, ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", ID)
	}
	
	asset := Asset{
		ChipID: 	ID,
		ChipName:   Name,
		Owner: 		manufacturer,
		Quantity:   Qty,
		Value: 		Val
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(ChipID, assetJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

// ReadAsset returns the asset stored in the world state with given id. [query]
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// // UpdateAsset updates an existing asset in the world state with provided parameters.
// func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
// 	exists, err := s.AssetExists(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("the asset %s does not exist", id)
// 	}

// 	// overwriting original asset with new asset
// 	asset := Asset{
	// correct naming
// 		ID:             id,
// 		Color:          color,
// 		Size:           size,
// 		Owner:          owner,
// 		AppraisedValue: appraisedValue,
// 	}
// 	assetJSON, err := json.Marshal(asset)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(id, assetJSON)
// }

// // DeleteAsset deletes an given asset from the world state.
// func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
// 	exists, err := s.AssetExists(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("the asset %s does not exist", id)
// 	}

// 	return ctx.GetStub().DelState(id)
// }

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}
/*
func (s *SmartContract) ManufacturerExists(ctx contractapi.TransactionContextInterface, id int) (bool, error) {
	manufacturerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return manufacturerJSON != nil, nil
}*/

// TransferAsset updates the owner field of   with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldOwner := asset.Owner
	asset.Owner = newOwner

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}