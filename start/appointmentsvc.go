package main

import (
    "errors"
    "fmt"
    "encoding/json"
    "os/exec"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

func (t *SimpleChaincode) createAppointment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var p Patient
    var patientId, patientFirstName, patientLastName string
    var err error
    fmt.Println("running createAppointment()")


    if len(args) != 7 {
        return nil, errors.New("Incorrect number of arguments. Expecting 7. name of the variable and value to set")
    }

    uuid, err := t.getUUID()

    if err != nil { 
        fmt.Printf("createAppointment: Error getting UUID: %s", err); 
        return nil, errors.New("createAppointment: Error getting UUID") 
    }

    patientId = string(uuid)
    patientFirstName = args[0]
    patientLastName=args[1]
    patientId_json :=  "\"patientId\":\""+patientId+"\", "      
    patientFirstName_json := "\"patientFirstName\":\""+patientFirstName+"\","
    patientLastName_json := "\"patientLastName\":\""+patientLastName+"\","    
   
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
        return nil, errors.New("Error converting Appointment record") 
    }

    err = stub.PutState(p.PatientId, bytes)

    if err != nil { 
        fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
        return nil, errors.New("Error storing Appointment record") 
    }

    return bytes, nil
}

func (t *SimpleChaincode) getUUID()([]byte, error){
    uuid, err := exec.Command("uuidgen").Output()
    if err != nil { 
        fmt.Printf("getUUID: Error getting UUID: %s", err); 
        return nil, errors.New("getUUID: Error getting UUID") 
    }
    return uuid, nil
}