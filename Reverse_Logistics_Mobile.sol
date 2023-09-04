pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "./Processes.sol";
import "./Stakeholders.sol";

/* REQUIRE CLAUSES DEACTIVATED, TO BE DEFINED ACCORDING TO EACH SITUATION
SC FOR REVERSE LOGISTICS TRACKING AND AUDITING
SEVERAL STAKEHOLDERS AND PROCESSES ARE CONSIDERED, AND DEFINED IN PARALLELL SCS
FUNCTIONS CAN BE EXPANDED, THESE FULFILL BASIC REQUIREMENTS
*/

contract Reverse_Logistics_Mobile{

    
    struct Product {
        uint id;
        string name;
        uint quantity;
        string others;  // for QOs or conditions, location etc
        uint numberoftraces;
        uint numberofcomponents;
        uint [] tracesProduct; // the ID of the traces of the product
        //uint [] temperaturesProduct;
        uint [] componentsProduct;
        address maker; // who  updates
        string globalId; // global id in manufacturing 
        bytes32 hashIPFS; // refernce to manufacturing description, serial number, IMEI
    }
    // key is a uint, later corresponding to the product id
    // what we store (the value) is a Product
    // the information of this mapping is the set of products of the order.
    mapping(uint => Product) private products; // public, so that w can access with a free function 

    struct Trace {
        uint id;
        uint id_product;
        string location;
        string temp_owner; // refurbishing agency/retailer
        uint timestamp;
        address maker; // who  updates
    }

    mapping(uint => Trace) private traces; // public, so that w can access with a free function 
    //store products count
    // since mappings cant be looped and is difficult the have a count like array
    // we need a var to store the coutings  
    // useful also to iterate the mapping 


    struct Component {  // each component of the device
        uint id;
        uint id_product;
        string name_component; // the component
        bool haschanged; // if it has been changed
        string serial_number; // corresponding serial of the actual component
        uint timestamp;
        address maker; // who  updates
        bytes32 componentIPFS; //hash reference to manufacturing description of the actual component.
    }

    mapping(uint => Component) private components; // public, so that w can access with a free function 



    //uint private temperaturesCount;
    uint private productsCount;
    uint private tracesCount;
    uint private componentsCount;

    //declare address of the participants
    address constant public customer = 0xE0f5206BBD039e7b0592d8918820024e2a7437b9;
    address constant public wholesaler = 0xE0F5206bbd039E7B0592D8918820024E2a743445;
    address constant public distributor = 0xE0F5206bbd039e7b0592d8918820024E2A743222;
    address constant public manufacturer = 0x50e00dE2c5cC4e456Cf234FCb1A0eFA367ED016E;
    address constant public govermment = 0x1533234Bd32f59909E1D471CF0C9BC80C92c97d2;
    address constant public refurbisher = 0x395BE1C1Eb316f82781462C4C028893e51d8b2a5;

    bool private  triggered;
    bool private  delivery;
    bool private  received;


    //Processes public p;
    //Stakeholders public s;

    // event, voted event. this will trigger when we want
    //  when a vote is cast for example, in the vote function. 
    event triggeredEvent (  // triggers new accepted order 
    );

    event deliveryEvent (  // triggers delivery start
    );

    event receivedEvent ( // triggers order received by customer
    );

    event updateEvent ( // triggers product status change
    );


    constructor () public { // constructor, creates order. we map starting from id=1,  hardcoded values of all
        addProduct("Example",200, "Delivey in 3 days","5400AA","ADDeFFtt45045594xxE3948"); //
        addTrace(1,"some coordinates", "name or address of actual owner",1573564413);
        //addComponent();
        triggered=false;
        delivery=false;
        received=false;
    }


    //PRODUCT OPERATIONS******************************************
    // enables product creation
    // get product
    // get total for externally looping the mapping
    // update others.

    // add product to mapping. private because we dont want to be accesible or add products afterwards to our mapping. We only want
    // our contract to be able to do that, from constructor
    // otherwise the conditions of the accepted contract could change
    function addProduct (string memory _name, uint _quantity, string memory _others, string memory _globalID, bytes32 _hashIpfs) private {
        //require(msg.sender==vendor);

        productsCount ++; // inc count at the begining. represents ID also. 
        products[productsCount].id = productsCount; 
        products[productsCount].name = _name;
        products[productsCount].quantity = _quantity;
        products[productsCount].others = _others;
        products[productsCount].numberoftraces = 0;
        products[productsCount].numberofcomponents = 0; 
        products[productsCount].maker = msg.sender;
        products[productsCount].globalId = _globalID;
        products[productsCount].hashIPFS = _hashIpfs;
        // reference the mapping with the key (that is the count). We assign the value to 
        // the mapping, the count will be the ID.  
    }

    // returns the number of products, needed to iterate the mapping and to know info about the order.
    function getNumberOfProducts () public view returns (uint){
        //require(msg.sender==customer || msg.sender==wholesaler || msg.sender==distributor);
        
        return productsCount;
    }

     // only specific stakeholders, can be changed
    function UpdateProduct (uint _productId, string memory _others) public { 
        ///require(msg.sender==wholesaler || msg.sender==distributor); // example. 
        require(_productId > 0 && _productId <= productsCount); 

        products[_productId].others = _others;  // update conditions
        emit updateEvent(); // trigger event 
    }

    // function to check the contents of the contract, the customer will check it and later will trigger if correct
    // only customer can check it 
    // customer will loop outside for this, getting the number of products before with getNumberOfProducts
    function getProduct (uint _productId) public view returns (Product memory) {
        //require(msg.sender==wholesaler || msg.sender==customer);
        require(_productId > 0 && _productId <= productsCount); 

        return products[_productId];
    }

      function getProductGlobalID (uint _productId) public view returns (string memory) {
        //require(msg.sender==customer);
        require(_productId > 0 && _productId <= productsCount); 

        return products[_productId].globalId;
    }


      function getProductHistoric (uint _productId) public view returns (bytes32) {
        //require(msg.sender==customer);
        require(_productId > 0 && _productId <= productsCount); 

        return products[_productId].hashIPFS;
    }
    //TRACES and Component OPERATIONS********************************************
    // enables add trace to a product
    // enables total number of traces to loop
    // get a trace
    // gets the total number of traces of a product. for statistical purposes
    // get the list of traces of a product, that can be consulter afterwards using get a trace
    // the same for components

    function addTrace (uint _productId, string memory _location, string memory _temp_owner, uint _timestamp) public {  // acts as update location
        //require(msg.sender==wholesaler || msg.sender==distributor);
        require(_productId > 0 && _productId <= productsCount); // check if product exists
        
        tracesCount ++; // inc count at the begining. represents ID also. 
        traces[tracesCount] = Trace(tracesCount, _productId, _location,_temp_owner,_timestamp,msg.sender);
        products[_productId].tracesProduct.push(tracesCount); // we store the trace reference in the corresponding product
        products[_productId].numberoftraces++;
         //this will give us the set of ID traces about our productid
        emit updateEvent();
    }


    function addComponent (uint _productId, string memory _name_component, bool _haschanged, string memory _serial_number, uint _timestamp, string memory _temp_owner, bytes32 _componentIPFS) public {  // acts as update location
        //require(msg.sender==wholesaler || msg.sender==manufacturer);
        require(_productId > 0 && _productId <= productsCount); // check if product exists
        
        componentsCount ++; // inc count at the begining. represents ID also. 
        components[componentsCount] = Component(componentsCount, _productId, _name_component,_haschanged,_serial_number,_timestamp,msg.sender, _componentIPFS);
        products[_productId].componentsProduct.push(componentsCount); // we store the trace reference in the corresponding product
        products[_productId].numberofcomponents++;
         //this will give us the set of ID traces about our productid
        emit updateEvent();
    }
   
    // returns the number of traced locations
    //useful for generic statistical purposes
    function getNumberOfTraces () public view returns (uint) {
        //require(msg.sender==customer || msg.sender==wholesaler || msg.sender==distributor);
        
        return tracesCount;
    }

    function getNumberOfComponents () public view returns (uint) {
        //require(msg.sender==customer || msg.sender==wholesaler || msg.sender==distributor);
        
        return componentsCount;
    }


    // get a trace
    function getTrace (uint _traceId) public view returns (Trace memory)  {
        //require(msg.sender==customer );
        require(_traceId > 0 && _traceId <= tracesCount); 

        return traces[_traceId];
    }

    function getComponent (uint _componentId) public view returns (Component memory)  {
        //require(msg.sender==customer );
        require(_componentId > 0 && _componentId <= componentsCount); 

        return components[_componentId];
    }


    // returns the number of traced locations for specific product
    function getNumberOfTracesProduct (uint _productId) public view returns (uint) {
        //require(msg.sender==customer || msg.sender==wholesaler || msg.sender==distributor);
        require(_productId > 0 && _productId <= productsCount); // check if product exists
        
        return products[_productId].numberoftraces;
    }

    // returns the number of components for specific product
    function getNumberOfComponentsProduct (uint _productId) public view returns (uint) {
        //require(msg.sender==customer || msg.sender==wholesaler || msg.sender==distributor);
        require(_productId > 0 && _productId <= productsCount); // check if product exists
        
        return products[_productId].numberofcomponents;
    }


    // get the array of traces of a product, later we can loop them using getTrace to obtain the data
    function getTracesProduct (uint _productId) public view returns (uint [] memory)  {
        //require(msg.sender==customer );
        require(_productId > 0 && _productId <= productsCount); // check if product exists

        return products[_productId].tracesProduct;
    }

    // get array of components product
    function getComponentsProduct (uint _productId) public view returns (uint [] memory)  {
        //require(msg.sender==customer );
        require(_productId > 0 && _productId <= productsCount); // check if product exists

        return products[_productId].componentsProduct;
    }

    function updateComponentsHash (uint _componentId) public {

        components[_componentId].componentIPFS = keccak256(abi.encodePacked(block.number,msg.data, components[_componentId].id, components[_componentId].name_component, components[_componentId].serial_number, components[_componentId].timestamp, components[_componentId].maker));
    }


    //EVENT AND SC OPERATIONS********************************************************
    //  computes hash of transaction
    // several event triggers


    function retrieveHashProduct (uint _productId) public view returns (bytes32){ 
        //computehash according to unique characteristics
        // hash has to identify a unique transaction so timestamp and locations and products should be used.
        // this example hashes a transaction as a whole.
        return keccak256(abi.encodePacked(block.number,msg.data, products[_productId].id, products[_productId].name, products[_productId].quantity, products[_productId].others, products[_productId].numberoftraces, products[_productId].numberofcomponents, products[_productId].maker));

    }

    function retrieveHashComponent (uint _componentId) public view returns (bytes32){ 
        require(_componentId > 0 && _componentId <= componentsCount); 

        return components[_componentId].componentIPFS;
    }

     //this function triggers the contract, enables it since the customer accepts it 
    // only customer
    function triggerContract () public { 
        //require(msg.sender==customer);
        triggered=true;
        emit triggeredEvent(); // trigger event 

    }

    // only wholesaler
    function deliverOrder () public { 
        //require(msg.sender==wholesaler);
        delivery=true;
        emit deliveryEvent(); // trigger event 

    }

    //only customer
    function receivedOrder () public { 
        //require(msg.sender==customer);
        received=true;
        emit receivedEvent(); // trigger event 

    }

        // returns global number of stories, needed to iterate the mapping and to know info.
    // smart to smart comm
    function updateNumberOfProcesses (address addr) public view returns (uint){
        
        Processes p = Processes(addr);
        return p.getNumberOfProcesses();
       
    }
    
    // returns global number of status, needed to iterate the mapping and to know info.
    // smart to smart comm
    function updateNumberOfStakeholders (address addr) public view returns (uint){
        
        Stakeholders s = Stakeholders(addr);
        return s.getNumberOfStakeholders();
        
    }


}
