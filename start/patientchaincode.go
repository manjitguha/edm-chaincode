package main

import (
    "errors"
    "fmt"
    "log"
    "encoding/json"
    "os/exec"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//=============================================================================================================================
// Patient - Defines the structure of Patient entity
//=============================================================================================================================
type Patient struct{
    PatientId string `json:"patientId"`
    PatientFirstName string `json:"patientFirstName"`
    PatientLastName string `json:"patientLastName"`
    Address Address `json:"address"`
}


//=============================================================================================================================
// Provider - Defines the structure of Provider entity
//=============================================================================================================================
type Provider struct{
    ProviderId string `json:"providerId"`
    ProviderFirstName string `json:"providerFirstName"`
    ProviderLastName string `json:"providerLastName"`
    Address Address `json:"address"`
}

//=============================================================================================================================
// Appointment - Defines the structure of Appointment entity
//=============================================================================================================================
type Appointment struct{
    AppointmentId string `json:"appointmentId"`

}

//=============================================================================================================================
// Address - Defines the structure of Address entity
//=============================================================================================================================
type Address struct{
    AddressLine1 string `json:"addressLine1"`
    AddressLine2 string `json:"addressLine2"`
    City string `json:"city"`
    State string `json:"state"`
    Zip string `json:"zip"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    err := stub.PutState("hello_world", []byte(args[0]))
    if err != nil {
        return nil, err
    }

    return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "write" {
        return t.write(stub, args)
    } else if function == "createPatient" {
        return t.createPatient(stub, args)
    }  


    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation")
}


func (t *SimpleChaincode) createPatient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var p Patient
    var patientId, patientFirstName, patientLastName,addressLine1,addressLine2,city,state,zip string
    var err error
    fmt.Println("running createPatient()")


    if len(args) != 7 {
        return nil, errors.New("Incorrect number of arguments. Expecting 7. name of the variable and value to set")
    }

    uuid, err := exec.Command("uuidgen").Output()
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    patientId = string(uuid)
    patientFirstName = args[0]
    patientLastName=args[1]
    addressLine1=args[2]
    addressLine2=args[3]
    city=args[4]
    state=args[5]
    zip=args[6]
    
    patientId_json :=  "\"patientId\":\""+patientId+"\", "      
    patientFirstName_json := "\"patientFirstName\":\""+patientFirstName+"\","
    patientLastName_json := "\"patientLastName\":\""+patientLastName+"\","    
    addressLine1_json := "\"addressLine1\":\""+addressLine1+"\","    
    addressLine2_json := "\"addressLine2\":\""+addressLine2+"\","    
    city_json := "\"city\":\""+city+"\","    
    state_json := "\"state\":\""+state+"\","    
    zip_json := "\"zip\":\""+zip+"\""    

    address_json := "\"address\":{"+addressLine1_json+addressLine2_json+city_json+state_json+zip_json+"}"

    patient_json := "{"+patientId_json+patientFirstName_json+patientLastName_json+address_json+"}"

    err = json.Unmarshal([]byte(patient_json), &p)  

    bytes, err  := t.save_changes(stub, p)

    //err = stub.PutState(patientId, []byte(patient_json))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }

    fmt.Println("Patient created successfully")

    return bytes, nil
}

func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, p Patient) ([]byte, error) {
    bytes, err := json.Marshal(p)
    if err != nil { 
        fmt.Printf("SAVE_CHANGES: Error converting patient record: %s", err); 
        return nil, errors.New("Error converting patient record") 
    }

    err = stub.PutState(p.PatientId, bytes)

    if err != nil { 
        fmt.Printf("SAVE_CHANGES: Error storing patient record: %s", err); 
        return nil, errors.New("Error storing patient record") 
    }

    return bytes, nil
}



func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var name, value string
    var err error
    fmt.Println("running write()")

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
    }

    name = args[0]                            //rename for fun
    value = args[1]
    err = stub.PutState(name, []byte(value))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}





// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
    if function == "read" {                            //read a variable
        return t.read(stub, args)
    } else if function == "getPatient" {
        return t.getPatient(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}

func (t *SimpleChaincode) getPatient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var patientId, jsonResp string
    var err error
    var p Patient

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting ID of the patient to query")
    }

    patientId = args[0]
    valAsbytes, err := stub.GetState(patientId)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + patientId + "\"}"
        return nil, errors.New(jsonResp)
    }

    err = json.Unmarshal(valAsbytes, &p);

    if err != nil { 
        fmt.Printf("getPatient: Corrupt Patient record "+string(valAsbytes)+": %s", err); 
        return nil, errors.New("getPatient: Invalid patient object") 
    }

    bytes, err := json.Marshal(p)

    if err != nil { 
        return nil, errors.New("getPatient: Invalid patient object") 
    }

    return bytes, nil
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var name, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
    }

    name = args[0]
    valAsbytes, err := stub.GetState(name)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}


