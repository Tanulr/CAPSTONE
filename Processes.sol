pragma solidity ^0.5.0;
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


contract Processes{

    struct Process {
        uint id; // 
	string name;       
        uint timestamp; // when it was applied, just in case it is not the same date than token creation
        string description; // other info
	bool active;
        address maker; // who applied this proces
        string hashIPFS; // hash of the elements of the struct, for auditing AND IPFS 
	uint [] involvedproducts; // id of products that use this process
    }

    mapping(uint => Process) private processChanges; //

    
    uint private productsCount;
    uint private pedigreeCount;
    uint private processCount;
    uint private stakeholdersCount;


    // events, since SC is for global accounts it does not have too much sense but is left here 
    event updateEvent ( // triggers update complete
    );
    
    event changeStatusEvent ( // triggers status change
    );

    address constant public stakeholder = 0xE0f5206BBD039e7b0592d8918820024e2a7437b9; // who registers the token into system. 
    address constant public stakeholder2 = 0xE0F5206bbd039e7b0592d8918820024E2A743222;

    constructor () public { // constructor, inserts new process in system. we map starting from id=1, hardcoded values of all
        addProcess("Repair_camera",1573564413,"substitute camera of a device"); // simmple example. needs products ID to be added.
        
    }
    
    // add process item to a product, updates hash
    function addProcess (string memory _name, uint _timestamp, string memory _description) public {

        processCount++;

        processChanges[processCount].id = processCount;
        processChanges[processCount].name = _name; 
        processChanges[processCount].timestamp = _timestamp; 
        processChanges[processCount].description = _description; 
        processChanges[processCount].active = true; 
        processChanges[processCount].maker = msg.sender;
        emit updateEvent(); // trigger event 
    
    }

    // we add a product to a process to keep track of it
    function addProcessProduct(uint _id) public {
        processChanges[processCount].involvedproducts.push(_id);
        emit updateEvent(); // trigger event 
    }

    function changeStatus (uint _id, bool _active) public { 
        require(_id > 0 && _id <= processCount); 

        processChanges[processCount].active = _active;
        emit changeStatusEvent(); // trigger event 
    }
    
    // get the products managed by the process
    function getStakeholdersProduct (uint _id) public view returns (uint [] memory)  {
        require(_id > 0 && _id <= processCount);  // security check avoid memory leaks
        require(msg.sender == processChanges[_id].maker);
        
        return processChanges[_id].involvedproducts;
    }

    function getProcess (uint _processId) public view returns (Process memory)  {
        require(_processId > 0 && _processId <= processCount); 
        require(msg.sender==processChanges[_processId].maker); // only if he is the author of the content
        
        return processChanges[_processId];
    }
    
    // returns global number of stories, needed to iterate the mapping and to know info.
    function getNumberOfProcesses () public view returns (uint){
    //tx.origin 
        return processCount;
    }
    

    
}
