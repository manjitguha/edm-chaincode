/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
"errors"
"fmt"
"encoding/json"
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
    Address Address `json:address`
}

type Address struct{
    AddressId string `json:"addressId"`
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
    var patientId, patientFirstName, patientLastName,addressId,addressLine1,addressLine2,city,state,zip string
    var err error
    fmt.Println("running createPatient()")


    if len(args) != 7 {
        return nil, errors.New("Incorrect number of arguments. Expecting 9. name of the variable and value to set")
    }


    patientFirstName = args[0]
    patientLastName=args[1]
    patientId = patientFirstName+patientLastName

    addressLine1=args[2]
    addressLine2=args[3]
    city=args[4]
    state=args[5]
    zip=args[6]
    addressId = addressLine1+addressLine2+city+state+zip

    patientId_json :=  "\"patientId\":\""+patientId+"\", "      
    patientFirstName_json := "\"patientFirstName\":\""+patientFirstName+"\","
    patientLastName_json := "\"patientLastName\":\""+patientLastName+"\","    
    addressId_json := "\"addressId\":\""+addressId+"\","    
    addressLine1_json := "\"addressLine1\":\""+addressLine1+"\","    
    addressLine2_json := "\"addressLine2\":\""+addressLine2+"\","    
    city_json := "\"city\":\""+city+"\","    
    state_json := "\"state\":\""+state+"\","    
    zip_json := "\"zip\":\""+zip+"\""    

    address_json := "\"address\":{"+addressId_json+addressLine1_json+addressLine2_json+city_json+state_json+zip_json+"}"

    patient_json := "{"+patientId_json+patientFirstName_json+patientLastName_json+address_json+"}"

    err = json.Unmarshal([]byte(patient_json), &p)  

    _, err  = t.save_changes(stub, p)

    //err = stub.PutState(patientId, []byte(patient_json))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }

    fmt.Println("Patient created successfully")

    return nil, nil
}

func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, p Patient) (bool, error) {
    bytes, err := json.Marshal(p)
    if err != nil { 
        fmt.Printf("SAVE_CHANGES: Error converting patient record: %s", err); 
        return false, errors.New("Error converting patient record") 
    }

    err = stub.PutState(p.PatientId, bytes)

    if err != nil { 
        fmt.Printf("SAVE_CHANGES: Error storing patient record: %s", err); 
        return false, errors.New("Error storing patient record") 
    }

    return true, nil
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

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting ID of the patient to query")
    }

    patientId = args[0]
    valAsbytes, err := stub.GetState(patientId)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + patientId + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
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


