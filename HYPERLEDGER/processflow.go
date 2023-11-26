package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type EscrowContract struct { // initiated by A (sender)
	DeliveryEntity    string `json:"deliveryEntity"`
	DeliveryStake     uint64 `json:"deliveryStake"`
	DisputeFlag       bool   `json:"disputeFlag"`
	EscrowAmount      uint64 `json:"escrowAmount"`
	ProductVerifications map[string]ProductVerification `json:"productVerifications"`
	Receiver          string `json:"receiver"`
	Sender            string `json:"sender"` 
	TransactionCompleted bool   `json:"transactionCompleted"`
	TxnID			  string `json: "txnID"` //Txn1	
	AssetID           string `json: "assetID"`
}

// type ProductVerification struct {
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
func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface, txn, assetID, manufacturer, deliveryEntity, receiver string, escrowAmount, deliveryStake uint64) error {
	// Check if asset exists
	asset, err := s.ReadAsset(ctx, assetID)
	if err != nil {
		return "", err
	}
	// Check is client owns asset
	if asset.Owner != ctx.GetClientIdentity().GetCreator() {
		return fmt.Errorf("Client doesnt own asset %v", err)
	}
	// Check is delivery, receiver exist in same channel


	// create a transaction for this asset
	escrow := EscrowContract {
		DeliveryEntity: deliveryEntity,
		DeliveryStake: deliveryStake,
		DisputeFlag: false,
		EscrowAmount: escrowAmount,
		Receiver: receiver,
		Sender: ctx.GetClientIdentity().GetID(),
		TransactionCompleted: false,
		TxnID: txn, 
	}
	escrowJSON, err := json.Marshal(escrow)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(txn, escrowJSON), nil
}

func (s *SmartContract) StartDelivery(ctx contractapi.TransactionContextInterface) error {
	// how to access stuct data using object. or else we'll just add txnID, access getState()
	if ctx.GetClientIdentity().GetID() != s.Receiver {
		return fmt.Errorf("Only receiver (B) can start the delivery")
	}
	if s.TransactionCompleted {
		return fmt.Errorf("Transaction is already completed")
	}
	// google if line is correct
	s.EscrowAmount = ctx.GetStub().GetArgs()[0]
	return nil
}

func (s *SmartContract) InitiateDelivery(ctx contractapi.TransactionContextInterface) error {
	if ctx.GetClientIdentity().GetID() != s.Sender {
		return fmt.Errorf("Only sender (A) can initiate the delivery")
	}
	if s.EscrowAmount == 0 {
		return fmt.Errorf("Escrow amount not paid")
	}
	// Stake amount is deposited by delivery entity D
	return nil
}

func (s *SmartContract) ConfirmDelivery(ctx contractapi.TransactionContextInterface) error {
	if ctx.GetClientIdentity().GetID() != s.DeliveryEntity {
		return fmt.Errorf("Only delivery entity (D) can confirm delivery")
	}
	if s.TransactionCompleted {
		return fmt.Errorf("Transaction is already completed")
	}
	s.DeliveryStake = ctx.GetStub().GetArgs()[0]
	// Implement the logic for confirming successful delivery here.
	s.TransactionCompleted = true
	return nil
}

func (s *SmartContract) VerifyProduct(ctx contractapi.TransactionContextInterface, originalValue, aValue, bValue uint64) error {
	if ctx.GetClientIdentity().GetID() != s.Receiver {
		return fmt.Errorf("Only receiver (B) can verify the product")
	}
	if !s.TransactionCompleted {
		return fmt.Errorf("Delivery hasn't reached yet")
	}
	if s.DisputeFlag {
		return fmt.Errorf("Dispute is ongoing")
	}

	// Compare values with the original manufacturer values
	if originalValue == aValue && aValue == bValue {
		s.ProductVerifications[s.Receiver] = ProductVerification{originalValue, aValue, bValue, Verified}
		// If all values match, release funds to sender A and delivery entity D
		ctx.GetStub().Transfer(s.Sender, s.EscrowAmount)
		ctx.GetStub().Transfer(s.DeliveryEntity, s.DeliveryStake)
		s.TransactionCompleted = false
	} else if aValue == originalValue && bValue != aValue {
		// D is malicious, delivery stake to A, return escrow to B, and flag D
		s.DisputeFlag = true
		s.ProductVerifications[s.Receiver] = ProductVerification{originalValue, aValue, bValue, Dispute}
		ctx.GetStub().Transfer(s.Sender, s.DeliveryStake)
		ctx.GetStub().Transfer(s.Receiver, s.EscrowAmount)
		// Flag D (potentially ban D from the network)
		s.TransactionCompleted = false
	} else {
		// A is malicious, refund delivery stake to D, return escrow to B, and flag A
		s.DisputeFlag = true
		s.ProductVerifications[s.Receiver] = ProductVerification{originalValue, aValue, bValue, Dispute}
		ctx.GetStub().Transfer(s.DeliveryEntity, s.DeliveryStake)
		ctx.GetStub().Transfer(s.Receiver, s.EscrowAmount)
		// Flag A (potentially ban A from the network)
		s.TransactionCompleted = false
	}

	return nil
}

func (s *SmartContract) GetContractBalance(ctx contractapi.TransactionContextInterface) (uint64, error) {
	return ctx.GetStub().GetState(s.ContractID)
}

func main() {
	contract := new(EscrowContract)
	contract.TransactionContextHandler = &contractapi.ExampleCC{}
	contractapi.NewChaincode(contract)
}
