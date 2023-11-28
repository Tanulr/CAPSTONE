package main

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
	Quantity       uint64    `json:"quantity"`
	VerifyValue 		   string    `json:"value"`
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
	// # ProductVerifications map[string]ProductVerification `json:"productVerifications"`
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
	manufacturer := Check{ID: 1, Pack: ctx.GetClientIdentity().GetMSPID()}
	
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
// # type ProductVerification struct {
// 	AValue        uint64              `json:"aValue"`
// 	BValue        uint64              `json:"bValue"`
// 	OriginalValue uint64              `json:"originalValue"`
// 	Status        VerificationStatus `json:"status"`
// }
/*
type Verification struct {
	Entity  	  string   `json:"Entity"` //getCreator() / getID()
	Verification  string   `json:"Verify"`
}
type VerificationStatus int

const (
	Dispute
	NotVerified VerificationStatus = iota
	Verified
)*/

// Client drafts order. Checks if order is valid.
func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface, txn, assetID, deliveryEntity, receiver, verify string, escrowAmount, deliveryStake uint64) error {
	// Check if asset exists
	asset, err := s.ReadAsset(ctx, assetID)
	if err != nil {
		return "", err
	}
	// Check is client owns asset
	if asset.Owner != ctx.GetClientIdentity().GetMSPID() {
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
		// # ProductVerifications map[string]ProductVerification `json:"productVerifications"`
		InitiateDelivery: false,
		Receiver: receiver,
		Sender: ctx.GetClientIdentity().GetMSPID(),
		StartDelivery: false,
		TransactionCompleted: false,
		TxnID: txn, //Txn1	
		Verify: verify,
	}
	// RECHECK
	flag := Flag {

	}
	escrowJSON, err := json.Marshal(escrow)
	if err != nil {
		return err
	}
	err := ctx.GetStub().PutState(txn, escrowJSON)
	if err != nil {
		return err
	}

	return err 
}

// delivery PROCESS started by receiver
func (s *SmartContract) StartDelivery(ctx contractapi.TransactionContextInterface, txn string, decision bool) error {
	// how to access stuct data using object. or else we'll just add txnID, access getState()
	escrowJSON, err := ctx.GetStub().GetState(txn)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if ctx.GetClientIdentity().GetMSPID() != escrowJson[receiver] {
		return fmt.Errorf("Only receiver (B) can start the delivery")
	}
	if escrowJSON.TransactionCompleted {
		return fmt.Errorf("Transaction is already completed")
	}
	escrowJSON[startDelivery] = decision
	err := ctx.GetStub().PutState(txn, escrowJSON)
	if err != nil {
		return err
	}

	return nil 
}

// ran by sender, to say they gave product to delivery
func (s *SmartContract) InitiateDelivery(ctx contractapi.TransactionContextInterface, txn string, decision bool) error {
	escrowJSON, err := ctx.GetStub().GetState(txn)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if ctx.GetClientIdentity().GetMSPID() != escrowJson[sender] {
		return fmt.Errorf("Only sender (A) can start the delivery")
	}
	if escrowJSON.TransactionCompleted {
		return fmt.Errorf("Transaction is already completed")
	}
	escrowJSON[initiateDelivery] = decision
	err := ctx.GetStub().PutState(txn, escrowJSON)
	if err != nil {
		return err
	}

	return nil 
}

// ran by delivery to say they finished their job of delivery
func (s *SmartContract) ConfirmDelivery(ctx contractapi.TransactionContextInterface, txn string, decision bool) error {
	escrowJSON, err := ctx.GetStub().GetState(txn)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if ctx.GetClientIdentity().GetMSPID() != escrowJson[delivery] {
		return fmt.Errorf("Only delivery entity (D) can start the delivery")
	}
	if escrowJSON[transactionCompleted] {
		return fmt.Errorf("Transaction is already completed")
	}
	escrowJSON[confirmDelivery] = decision
	err := ctx.GetStub().PutState(txn, escrowJSON)
	if err != nil {
		return err
	}

	return nil
}

// receiver runs this to verify with manufacture verify (in asset) and sender (in escrowContract)
func (s *SmartContract) VerifyProduct(ctx contractapi.TransactionContextInterface, txn string, bValue string) error {
	escrowJSON, err := ctx.GetStub().GetState(txn)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	assetJSON, err := ctx.GetStub().GetState(escrowJSON[assetID])
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if ctx.GetClientIdentity().GetMSPID() != escrowJson[receiver] {
		return fmt.Errorf("Only receiver (B) can start the delivery")
	}
	if escrowJSON[transactionCompleted] {
		return fmt.Errorf("Transaction is already completed")
	}	
	if escrowJSON[disputeFlag] {
		return fmt.Errorf("Dispute is ongoing")
	}

	aValue := escrowJSON[verify]
	originalValue := assetJSON[value]
	// Compare values with the original manufacturer values
	if originalValue == aValue && aValue == bValue {
		// s.ProductVerifications[s.Receiver] = ProductVerification{originalValue, aValue, bValue, Verified}
		// # If all values match, release funds to sender A and delivery entity D
		
		escrowJSON[transactionCompleted] = true
	} else if aValue == originalValue && bValue != aValue {
		// D is malicious, delivery stake to A, return escrow to B, and flag D
		escrowJSON[disputeFlag] = true
		// # s.ProductVerifications[s.Receiver] = ProductVerification{originalValue, aValue, bValue, Dispute}
		// ctx.GetStub().Transfer(s.Sender, s.DeliveryStake)
		// ctx.GetStub().Transfer(s.Receiver, s.EscrowAmount)
		// Flag D (potentially ban D from the network)
		escrowJSON[transactionCompleted] = false
	} else {
		// A is malicious, refund delivery stake to D, return escrow to B, and flag A
		escrowJSON[disputeFlag] = true
		// # s.ProductVerifications[s.Receiver] = ProductVerification{originalValue, aValue, bValue, Dispute}
		// ctx.GetStub().Transfer(s.DeliveryEntity, s.DeliveryStake)
		// ctx.GetStub().Transfer(s.Receiver, s.EscrowAmount)
		// Flag A (potentially ban A from the network)
		escrowJSON[transactionCompleted] = false
	}
	assetJSON[owner] := ctx.GetClientIdentity().GetMSPID() // new owner is receiver
	err := ctx.GetStub().PutState(txn, escrowJSON)
	if err != nil {
		return err
	}
	err := ctx.GetStub().PutState(txn, assetJSON)
	if err != nil {
		return err
	}
	return nil
}

// func (s *SmartContract) GetContractBalance(ctx contractapi.TransactionContextInterface) (uint64, error) {
// 	return ctx.GetStub().GetState(s.ContractID)
// }

// func main() {
// 	contract := new(EscrowContract)
// 	contract.TransactionContextHandler = &contractapi.ExampleCC{}
// 	contractapi.NewChaincode(contract)
// }