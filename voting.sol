// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;
pragma experimental ABIEncoderV2;


contract Stakeholders{
 
    struct Stakeholder{
        //uint id; // this especific process, containing id and quantity
        string name; // the product
        uint timestamp; // when it was applied, just in case it is not the same date than token creation
        uint [] involvedproducts; // products used by stakeholder
        string description; // other info
        address maker; // who applied this proces
        address own;
        bool active;
        string hashIPFS; // hash of the elements of the struct, for auditing AND IPFS 
    }
    struct TempStakeholder {
        string name;
        address entityAddress;
        uint timestamp; // when it was applied, just in case it is not the same date than token creation
        
        string description;
        address maker;
    }
    // mapping(address => bool) private voted;
    mapping(address => TempStakeholder) private temps;
    mapping(address => Stakeholder) private stakeholderChanges;
    mapping(address => bool) private authorizedEntities;
    mapping(address => bool) private votingRights;

    uint private stakeholderCount;
    uint private votesRequired;

    event updateEvent(
        string NewMessage
    );
    event changeStatusEvent();
   
    event NewStakeholderEvent(
    address from,
    string Entity,
    address newEntity,
    string body
    );

    modifier check {
        require(authorizedEntities[msg.sender], "User not authorized");
        _;
    }

    constructor() {
        address deployer = msg.sender;
    authorizedEntities[deployer] = true;
    votingRights[deployer] = true;

    stakeholderChanges[deployer] = Stakeholder({
        name: "Manufacturer",
        timestamp: block.timestamp,
        description: "Manufactures several components CPU, RAM, and chipsets.",
        active: true,
        maker: deployer,
        own: deployer,
        involvedproducts: new uint[](0),
        hashIPFS: ""
    });

    stakeholderCount = 1;
    votesRequired = 0; // Initialize the votes required (50% of the total stakeholders)
}

    function addStakeholder(string memory _name, address entityAddress, string memory _description) public check {
        temps[entityAddress] = TempStakeholder({
            name: _name,
            
            timestamp:block.timestamp,
            entityAddress: entityAddress,
            description: _description,
            maker: msg.sender
        });
        votesRequired++;
        emit NewStakeholderEvent(msg.sender,_name, entityAddress, _description);
    }

    function voteToAddStakeholder(address _addy) public check {
        require(votingRights[msg.sender], "User does not have voting rights");
        // require(!voted[msg.sender], "User has already voted");
        // Add the voter's approval
        votesRequired++;
        // voted[msg.sender] = true;
        // Check if the votes reach the required threshold
        if (votesRequired * 2 > stakeholderCount) {
            addStakeholderInternal(_addy); // If yes, add the new stakeholder
            votesRequired=0;

        }
        
    }

    function addStakeholderInternal(address _addy) private {
        authorizedEntities[_addy] = true;
        votingRights[_addy] = true; // New stakeholders get voting rights by default
        stakeholderCount++;
        stakeholderChanges[_addy] = Stakeholder({
            name: temps[_addy].name,
            description: temps[_addy].description,
            timestamp:temps[_addy].timestamp,
            active: true,
            maker: temps[_addy].maker,
            own: _addy,
            involvedproducts: new uint[](0),
            hashIPFS: ""
        });
        emit updateEvent("New entity added");
    }

    function addStakeholderProduct(uint _id) public check {

        
        stakeholderChanges[msg.sender].involvedproducts.push(_id);
        emit updateEvent("New product added"); // trigger event 
    }
    
    // get the products managed by the stakeholder
    function getStakeholdersProduct (address _addy) public view check returns (uint [] memory)  {
        require(authorizedEntities[_addy], "User does not exist");
        
        return stakeholderChanges[_addy].involvedproducts;
    }

    function changeStatus (bool _active) public check {
        stakeholderChanges[msg.sender].active = _active;
        emit changeStatusEvent(); // trigger event 
    }

    function getStakeholder (address _id) public view check returns (Stakeholder memory)  {
        
        return stakeholderChanges[_id];
    }
    
    // returns global number of status, needed to iterate the mapping and to know info.
    function getNumberOfStakeholders () public view check returns (uint){    
        //tx.origin
        return stakeholderCount;
    }

}
