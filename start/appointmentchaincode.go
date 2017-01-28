package main

import (
    "errors"
    "fmt"
    "encoding/json"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "log"
)


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

    uuidArray := []string{}
   
    UUIDsBytes, err := json.Marshal(uuidArray)
    log.Println("Saving")
    err = stub.PutState("activeUUIDs", UUIDsBytes)
    log.Println("Saved")
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
    } else if function == "upsertAppointment" {                           
        return t.upsertAppointment(stub, args)
    }

    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation")
}



// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    if function == "getUUID" {
        return t.getUUID()
    } else if function == "getAppointment" {
        return t.getAppointment(stub,args)
    } else if function == "getActiveUUIDs" {
        return t.getActiveUUIDs(stub,args)
    } else if function == "getActiveUUIDsForID" {
        return t.getActiveUUIDsForID(stub,args)
    }


    
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}
