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


    if len(args) != 11 {
        return nil, errors.New("Incorrect number of arguments. Expecting 11. name of the variable and value to set")
    }


    appointment.AppointmentId = args[0]
    appointment.Patient.PatientId = args[1]
    appointment.Patient.PatientFirstName = args[2]
    appointment.Patient.PatientLastName = args[3]
    appointment.Provider.ProviderId = args[4]
    appointment.Provider.ProviderFirstName = args[5]
    appointment.Provider.ProviderLastName = args[6]
    appointment.AppointmentTime = args[7]
    appointment.DiagnosisNotes = args[8]
    appointment.PrescriptionNotes = args[9]
    appointment.Status = args[10]

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
    log.Println("Saving Appointment")
   
    log.Println("AppointmentId = %s", appointment.AppointmentId)
   
    err = stub.PutState(appointment.AppointmentId, bytes)

    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
        return nil, err
    }
    log.Println(bytes)
    return bytes, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) getAppointment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    log.Println("Printing Length")
    var appointment Appointment
    var key, role, jsonResp string
    var err error

    log.Println("Length = %d", len(args))
    log.Println("After Prining Length")
    
    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    role = args[0]
    key = args[1]

    log.Println("Key = %s", key)
    log.Println("Role = %s", role)
    
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        log.Fatalln("Failed to open log file", ":", err)
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    log.Println(valAsbytes)


    err = json.Unmarshal(valAsbytes, &appointment)  

    if err != nil {
        return nil, err
    }

    if role != PAYER && role != PROVIDER && role != PHARMACY && role != PATIENT && role != SECRETARY {
        jsonResp = "{\"Error\":\"Role doesn't exist " + role + "\"}"
        return nil, errors.New(jsonResp)
    } else if role == SECRETARY {
        appointment.DiagnosisNotes = UNAUTHORIZED
        appointment.PrescriptionNotes = UNAUTHORIZED
    } else if role  == PHARMACY {
        appointment.DiagnosisNotes = UNAUTHORIZED
    }


    bytes, err := json.Marshal(appointment)
   
    if err != nil {
        return nil, err
    }

    return bytes, nil
}


