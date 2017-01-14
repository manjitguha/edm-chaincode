package main

import (
    "errors"
    "fmt"
    "encoding/json"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "log"
)

func (t *SimpleChaincode) upsertAppointment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var appointment Appointment
    
    var err error
    fmt.Println("running createAppointment()")


    if len(args) != 10 {
        return nil, errors.New("Incorrect number of arguments. Expecting 10. name of the variable and value to set")
    }


    appointment.AppointmentId = args[0]
    appointment.Patient.PatientId = args[1]
    appointment.Patient.PatientFirstName = args[2]
    appointment.Patient.PatientLastName = args[3]
    appointment.Provider.ProviderId = args[4]
    appointment.Provider.ProviderFirstName = args[5]
    appointment.Provider.ProviderLastName = args[6]
    appointment.AppointmentTime = args[7]

    if args[8] != "" {
        appointment.DiagnosisNotes = args[8]
    }

    if args[9] != "" {
        appointment.PrescriptionNotes = args[9]
    }
    
    bytes, err  := t.save_changes(stub, appointment)

    if err != nil {
        return nil, err
    }

    fmt.Println("Appointment created successfully")

    return bytes, nil
}

func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, appointment Appointment) ([]byte, error) {
    bytes, err := json.Marshal(appointment)
    if err != nil { 
        fmt.Printf("save_changes: Error converting Appointment record: %s", err); 
        return nil, err 
    }

    err = stub.PutState(appointment.AppointmentId, bytes)

    if err != nil { 
        fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
        return nil, err
    }

    return bytes, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) getAppointment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    log.Println("Printing Length")
    // var appointment Appointment
    var key, jsonResp string
    var err error

    log.Println("Length = %d", len(args))
    log.Println("After Prining Length")
    
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]

    log.Println("Key = %s", key)

    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

/*    err = json.Unmarshal(valAsbytes, &appointment)  

    if err != nil {
        return nil, err
    }

    bytes, err := json.Marshal(appointment)
   
    if err != nil {
        return nil, err
    }
*/

    return valAsbytes, nil
}


