package main

import (
    "errors"
    "fmt"
    "github.com/hyperledger/fabric/core/chaincode/shim"
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

/*func main(){
    uuid, err := exec.Command("uuidgen").Output()
    fmt.Println(string(uuid))
    fmt.Println(err)
}*/

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
    }  else if function == "createAppointment" {
        return t.createAppointment(stub, args)
    }  

    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation")
}



// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}
