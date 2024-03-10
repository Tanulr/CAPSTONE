/*
SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	ChipID         string `json:"id"` // UNIQUE asset#
	ChipName       string `json:"chipName"`
	Owner          string `json:"owner"`
	Quantity       uint64 `json:"quantity"`
	VerifyValue    string `json:"value"`
}

type Check struct {
	ID 		string      `json:"id"` // string
	Pack    string   `json:"pack"`
}

type EscrowContract struct { // initiated by A (sender)
	AssetID           string `json: "assetID"`
	ConfirmDelivery   bool	`json: "confirmDelivery"`
	Delivery    	  string `json:"delivery"`
	DeliveryStake     uint64 `json:"deliveryStake"`
	DisputeFlag       bool	`json:"disputeFlag"`
	EscrowAmount      uint64 `json:"escrowAmount"`
	InitiateDelivery  bool	`json: "initiateDelivery"`
	Receiver          string `json:"receiver"`
	Sender            string `json:"sender"` 
	StartDelivery 	  bool	`json: "startDelivery"`
	TransactionCompleted	bool    `json:"transactionCompleted"`
	TxnID			  string `json: "txnID"` //Txn1	
	Verify		   	  string `json:"value"`
}

// InitLedger run by manurfacturer. Sets assets list. Whoever runs this is the 
// default (single) manufacturer. [invoke]
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// asset {ChipID(primary key), "asset1", ChipName: "Intel Core i9 10900K", Owner: "Tomoko", Quantity: 5, Value: 300},

	//assets := []Asset{}
	x, _ := ctx.GetClientIdentity().GetMSPID()
	manufacturer := Check{ID: "1", Pack: x}
	
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
	id, err := ctx.GetClientIdentity().GetMSPID()
	return id, err
}

// CreateAsset issues a new asset to the world state with given details. [invoke]
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, ID string, Name string, Qty uint64, Val string) error {
	manufacturer, _ := ctx.GetClientIdentity().GetMSPID()
	xy, _ := ctx.GetStub().GetState("1")
	var x Check
	err := json.Unmarshal(xy, &x)
	if manufacturer != x.Pack{
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
		VerifyValue: 		Val,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(ID, assetJSON)
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



// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
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

// #################################################


// Client drafts order. Checks if order is valid. [invoke]
func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface, txn, assetID, deliveryEntity, receiver, verify string, escrowAmount, deliveryStake uint64) error {
	// Check if asset exists
	asset, err := s.ReadAsset(ctx, assetID)
	if err != nil {
		return err
	}
	// Check is client owns asset
	x, err := ctx.GetClientIdentity().GetMSPID()
	if asset.Owner !=  x {
		return fmt.Errorf("Client doesnt own asset %v", err)
	}
	// # Check is delivery, receiver exist in same channel


	// create a transaction for this asset
	escrow := EscrowContract {
		AssetID: assetID,
		ConfirmDelivery: false,
		Delivery: deliveryEntity,
		DeliveryStake: deliveryStake,
		DisputeFlag: false,
		EscrowAmount: escrowAmount,
		InitiateDelivery: false,
		Receiver: receiver,
		Sender: x,
		StartDelivery: false,
		TransactionCompleted: false,
		TxnID: txn, //Txn1	
		Verify: verify,
	}

	escrowJSON, err := json.Marshal(escrow)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(txn, escrowJSON)
	if err != nil {
		return err
	}

	return err 
}

// delivery PROCESS started by receiver. [invoke]
func (s *SmartContract) StartDelivery(ctx contractapi.TransactionContextInterface, txn string, decision bool) error {
	// how to access stuct data using object. or else we'll just add txnID, access getState()
	escrow, err := ctx.GetStub().GetState(txn)
	var escrowJSON EscrowContract
	err = json.Unmarshal(escrow, &escrowJSON)
	
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	x, err := ctx.GetClientIdentity().GetMSPID()
	if x != escrowJSON.Receiver {
		return fmt.Errorf("Only receiver (B) can start the delivery")
	}
	if escrowJSON.TransactionCompleted {
		return fmt.Errorf("Transaction is already completed")
	}
	escrowJSON.StartDelivery = decision
	escrow, err = json.Marshal(escrowJSON)
	err = ctx.GetStub().PutState(txn, escrow)
	if err != nil {
		return err
	}

	return nil 
}

// ran by sender, to say they gave product to delivery. [invoke]
func (s *SmartContract) InitiateDelivery(ctx contractapi.TransactionContextInterface, txn string, decision bool) error {

	escrow, err := ctx.GetStub().GetState(txn)
	var escrowJSON EscrowContract
	err = json.Unmarshal(escrow, &escrowJSON)
	
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	x, err := ctx.GetClientIdentity().GetMSPID()
	if x != escrowJSON.Sender {
		return fmt.Errorf("Only sender (A) can start the delivery")
	}
	if escrowJSON.TransactionCompleted {
		return fmt.Errorf("Transaction is already completed")
	}
	escrowJSON.InitiateDelivery = decision
	escrow, err = json.Marshal(escrowJSON)
	err = ctx.GetStub().PutState(txn, escrow)
	if err != nil {
		return err
	}

	return nil 
}

// ran by delivery to say they finished their job of delivery. [invoke]
func (s *SmartContract) ConfirmDelivery(ctx contractapi.TransactionContextInterface, txn string, decision bool) error {

	escrow, err := ctx.GetStub().GetState(txn)
	var escrowJSON EscrowContract
	err = json.Unmarshal(escrow, &escrowJSON)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	x, _ := ctx.GetClientIdentity().GetMSPID() 
	if x != escrowJSON.Delivery {
		return fmt.Errorf("Only delivery entity (D) can start the delivery")
	}
	if escrowJSON.TransactionCompleted {
		return fmt.Errorf("Transaction is already completed")
	}
	escrowJSON.ConfirmDelivery = decision
	escrow, err = json.Marshal(escrowJSON)
	err = ctx.GetStub().PutState(txn, escrow)
	if err != nil {
		return err
	}

	return nil
}

// receiver runs this to verify with manufacture verify (in asset) and sender (in escrowContract) [invoke]
func (s *SmartContract) VerifyProduct(ctx contractapi.TransactionContextInterface, txn string, bValue string) error {
	escrow, err := ctx.GetStub().GetState(txn)
	var escrowJSON EscrowContract
	err = json.Unmarshal(escrow, &escrowJSON)
	
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	asset, err := ctx.GetStub().GetState(escrowJSON.AssetID)
	var assetJSON Asset
	err = json.Unmarshal(asset, &assetJSON)
	
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	x, _ := ctx.GetClientIdentity().GetMSPID()
	if x != escrowJSON.Receiver {
		return fmt.Errorf("Only receiver (B) can start the delivery")
	}
	if escrowJSON.TransactionCompleted {
		return fmt.Errorf("Transaction is already completed")
	}	
	if escrowJSON.DisputeFlag {
		return fmt.Errorf("Dispute is ongoing")
	}

	aValue := escrowJSON.Verify
	originalValue := assetJSON.VerifyValue
	// Compare values with the original manufacturer values
	if originalValue == aValue && aValue == bValue {
		escrowJSON.TransactionCompleted = true
	} else if aValue == originalValue && bValue != aValue {
		// D is malicious, delivery stake to A, return escrow to B, and flag D
		escrowJSON.DisputeFlag = true
		escrowJSON.TransactionCompleted = false
	} else {
		// A is malicious, refund delivery stake to D, return escrow to B, and flag A
		escrowJSON.DisputeFlag = true
		escrowJSON.TransactionCompleted = false
	}
	assetJSON.Owner, _ = ctx.GetClientIdentity().GetMSPID() // new owner is receiver
	escrow, err = json.Marshal(escrowJSON)
	err = ctx.GetStub().PutState(txn, escrow)
	if err != nil {
		return err
	}
	asset, err = json.Marshal(assetJSON)
	err = ctx.GetStub().PutState(txn, asset)
	if err != nil {
		return err
	}
	return fmt.Errorf("Product verified and is authentic. Owner is now %s", x)
}
