// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;
pragma experimental ABIEncoderV2;

/*
COPYRIGHT FRAN CASINO. 2019.
SECURITY CHECKS ARE COMMENTED FOR AN EASY USE TEST.
UNCOMMENT THE CODE FOR A FULLY FUNCTIONAL VERSION. 
YOU WILL NEED TO USE METAMASK OR OTHER EXTENSIONS TO USE THE REQUIRED ADDRESSES


ACTUALLY DATA ARE STORED IN THE SC. TO ENABLE IPFS, FUNCTIONS WILL NOT STORE the values and just the hash in the structs.
This can be changed in the code by calling the hash creation function. 
Nevertheless, the code is kept clear for the sake of understanding. 

*/

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

    mapping(address => Stakeholder) private stakeholderChanges; //
    mapping (address => bool) private authorizedEntities;

    
    //uint private productsCount;
    uint private stakeholderCount;

    // events, since SC is for global accounts it does not have too much sense but is left here 
    event updateEvent ( // triggers update complete
    );
    
    event changeStatusEvent ( // triggers status change
    );

    modifier check {
      require(authorizedEntities[msg.sender], "User not authorized");
      _;
    }

    // address constant public stakeholder = 0xE0f5206BBD039e7b0592d8918820024e2a7437b9; // who registers the product into system. 
    // address constant public stakeholder2 = 0xE0F5206bbd039e7b0592d8918820024E2A743222;

    constructor () { // constructor, inserts new token in system. we map starting from id=1, hardcoded values of all
        authorizedEntities[msg.sender] = true;
        addStakeholder("Manufacturer",msg.sender,1573564413,"Manufactures several components CPU, RAM  and chipsets. "); //
        
    }
    
    // add stakeholder to the list. checkers security disabled
    function addStakeholder (string memory _name, address entityAddress, uint _timestamp, string memory _description) public check{
        
        // require(authorized, "User not authorized to add entity");
        authorizedEntities[entityAddress] = true;
        stakeholderCount++;

        //stakeholderChanges[entityAddress].id = stakeholderCount;
        stakeholderChanges[entityAddress].name = _name; 
        stakeholderChanges[entityAddress].timestamp = _timestamp; 
        stakeholderChanges[entityAddress].description = _description; 
        stakeholderChanges[entityAddress].active = true; 
        stakeholderChanges[entityAddress].maker = msg.sender;
        stakeholderChanges[entityAddress].own = entityAddress;

        emit updateEvent(); // trigger event 
    }

    function addStakeholderProduct(uint _id) public check {

        
        stakeholderChanges[msg.sender].involvedproducts.push(_id);
        emit updateEvent(); // trigger event 
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
