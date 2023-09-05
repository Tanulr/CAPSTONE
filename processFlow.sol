// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;


contract EscrowContract {
    address public manufacturer; // Address of the manufacturer (M)
    address public sender;      // Address of sender A
    address public deliveryEntity;  // Address of delivery entity D
    address public receiver;    // Address of receiver B

    uint256 public escrowAmount; // Escrow amount paid by B
    uint256 public deliveryStake; // Stake amount deposited by D

    bool public transactionCompleted; // Flag indicating whether the transaction is completed
    bool public disputeFlag; // Flag to track disputes

    enum VerificationStatus { NotVerified, Verified, Dispute }

    struct ProductVerification {
        uint256 originalValue;
        uint256 aValue;
        uint256 bValue;
        VerificationStatus status;
    }

    mapping(address => ProductVerification) public productVerifications;

    constructor(address _manufacturer, address _deliveryEntity, address _receiver, uint256 _escrowAmount) {
        manufacturer = _manufacturer;
        sender = msg.sender;          // The creator of the contract is A
        deliveryEntity = _deliveryEntity;
        escrowAmount = _escrowAmount;
        receiver=_receiver;
        transactionCompleted = false;
        disputeFlag = false;
    }

    // Receiver B pays the escrow amount to start the delivery
    // B runs
    function startDelivery() public payable {
        require(msg.sender == receiver, "Only receiver (B) can start the delivery");
        require(!transactionCompleted, "Transaction is already completed");

        // Set the escrow amount to the value sent with this transaction
        escrowAmount = msg.value;
    }

    // Sender A initiates the delivery by depositing the escrow amount
    // A runs
    function initiateDelivery() public payable {
        require(msg.sender == sender, "Only sender (A) can initiate the delivery");
        // require(!transactionCompleted, "Transaction is already completed");
        require(escrowAmount > 0, "Escrow amount not paid");

        // Stake amount is deposited by delivery entity D
        
    }

    // Delivery entity D confirms delivery to receiver B
    // D runs
    function confirmDelivery() public payable {
        require(msg.sender == deliveryEntity, "Only delivery entity (D) can confirm delivery");
        require(!transactionCompleted, "Transaction is already completed");
        deliveryStake = msg.value;
        // Implement the logic for confirming successful delivery here.

        // If successful, finalize the transaction
        transactionCompleted = true;
    }

    event Display(uint cash);

    // Receiver B verifies the product and compares values
    // B runs
    function verifyProduct(uint256 _originalValue, uint256 _aValue, uint256 _bValue) public payable{
        require(msg.sender == receiver, "Only receiver (B) can verify the product");
        require(transactionCompleted, "Delivery hasn't reached yet");
        require(!disputeFlag, "Dispute is ongoing");

        // Compare values with the original manufacturer values
        if (_originalValue == _aValue && _aValue == _bValue) {
            productVerifications[receiver] = ProductVerification(_originalValue, _aValue, _bValue, VerificationStatus.Verified);

            // If all values match, release funds to sender A and delivery entity D
            emit Display(address(sender).balance);
            payable(sender).transfer(escrowAmount);
            emit Display(address(sender).balance);

            emit Display(address(deliveryEntity).balance);
            payable(deliveryEntity).transfer(deliveryStake);
            emit Display(address(deliveryEntity).balance);
            transactionCompleted= false;
        } else if (_aValue == _originalValue && _bValue != _aValue) {
            //D is malicious, delivery stake to A, return escrow to B, and flag D
            disputeFlag = true;
            productVerifications[receiver] = ProductVerification(_originalValue, _aValue, _bValue, VerificationStatus.Dispute);


            emit Display(address(sender).balance);
            payable(sender).transfer(deliveryStake);
            emit Display(address(sender).balance);

            emit Display(address(receiver).balance);
            payable(receiver).transfer(escrowAmount);
            emit Display(address(receiver).balance);
            // Flag D (potentially ban D from the network)
            transactionCompleted= false;
        } else {
            // A is malicious, refund delivery stake to D, return escrow to B, and flag A
            disputeFlag = true;
            productVerifications[receiver] = ProductVerification(_originalValue, _aValue, _bValue, VerificationStatus.Dispute);

            payable(deliveryEntity).transfer(deliveryStake);
            payable(receiver).transfer(escrowAmount);
            // Flag A (potentially ban A from the network)
            transactionCompleted= false;
        }
    }

    // Function to check the contract's current balance
    function getContractBalance() public view returns (uint256) {
        return address(this).balance;
    }
}