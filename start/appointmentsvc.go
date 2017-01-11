package main

import (
    "errors"
    "fmt"
    "encoding/json"
    "crypto/rand"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

func (t *SimpleChaincode) createAppointment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var p Patient
    var patientId, patientFirstName, patientLastName string
    var err error
    fmt.Println("running createAppointment()")


    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
    }

    uuid, err := t.getUUID()

    if err != nil { 
        fmt.Printf("createAppointment: Error getting UUID: %s", err); 
        return nil, err
    }

    patientId = string(uuid)
    patientFirstName = args[0]
    patientLastName=args[1]
    patientId_json :=  "\"patientId\":\""+patientId+"\", "      
    patientFirstName_json := "\"patientFirstName\":\""+patientFirstName+"\","
    patientLastName_json := "\"patientLastName\":\""+patientLastName+"\""    
   
    patient_json := "{"+patientId_json+patientFirstName_json+patientLastName_json+"}"

    err = json.Unmarshal([]byte(patient_json), &p)  

    bytes, err  := t.save_changes(stub, p)

    if err != nil {
        return nil, err
    }

    fmt.Println("Appointment created successfully")

    return bytes, nil
}

func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, p Patient) ([]byte, error) {
    bytes, err := json.Marshal(p)
    if err != nil { 
        fmt.Printf("save_changes: Error converting Appointment record: %s", err); 
        return nil, err) 
    }

    err = stub.PutState(p.PatientId, bytes)

    if err != nil { 
        fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
        return nil, err) 
    }

    return bytes, nil
}

func (t *SimpleChaincode) getUUID()([]byte, error){
     b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        fmt.Println("Error: ", err)
        return nil, err 
    }
    uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
    byteArray := []byte(uuid)
    return byteArray, nil 
}