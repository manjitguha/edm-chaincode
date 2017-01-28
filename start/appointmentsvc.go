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
    fmt.Println("Changing Blockchain createAppointment()")
    fmt.Println("Changing Blockchain During Meeting()")


    if len(args) != 8 {
        return nil, errors.New("Incorrect number of arguments. Expecting 7. name of the variable and value to set")
    }

    appointment.AppointmentId = args[0]
    appointment.PatientId = args[1]
    appointment.ProviderId = args[2]
    appointment.AppointmentDate = args[3]
    appointment.AppointmentTime = args[4]
    appointment.DiagnosisNotes = args[5]
    appointment.PrescriptionNotes = args[6]
    appointment.Status = args[7]

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
   
    log.Println("AppointmentId = ", appointment.AppointmentId)
    
    err = stub.PutState(appointment.AppointmentId, bytes)

    log.Println("Saving Appointment - PutState ", appointment.AppointmentId)
   

    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
        return nil, err
    }
    
    if appointment.ProviderId != "" {
        log.Println("Inside appointment.ProviderId != \"\" ", appointment.ProviderId)
   
        providerBytes, err  := t.saveUUIDsForProvider(stub, appointment)
        if err != nil { 
            log.Println(err)
            fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
            return nil, err
        }
        log.Println(providerBytes)
    }
    /*if appointment.PatientId != nil {
        activeUUIDsBytes, err  = t.saveUUIDsForPatient(stub, appointment)
    }
    if appointment.PharmacyId != nil {
        activeUUIDsBytes, err  = t.saveUUIDsForPharmacy(stub, appointment)
    }
    if appointment.SecretoryId != nil {
        activeUUIDsBytes, err  = t.saveUUIDsForSecretory(stub, appointment)
    }
    if appointment.LaboratoryId != nil {
        activeUUIDsBytes, err  = t.saveUUIDsForLaboratory(stub, appointment)
    }
    if appointment.ReferralProviderId != nil {
        activeUUIDsBytes, err  = t.saveUUIDsReferralProvider(stub, appointment)
    }
*/
    if err != nil {
        return nil, err
    }

    return bytes, nil
}


func (t *SimpleChaincode) saveUUIDsForProvider(stub shim.ChaincodeStubInterface, appointment Appointment) ([]byte, error) {
    var provider Provider
    log.Println("Inside saveUUIDsForProvider", appointment.ProviderId)

    providerBytes, err := stub.GetState(appointment.ProviderId);
    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error fetching activeUUIDs: %s", err); 
        return nil, err
    }

    log.Println("Before unmarshalling", len(providerBytes))
    if len(providerBytes) >0 {
        err = json.Unmarshal(providerBytes, &provider)
        if provider.AppointmentSlotMap[appointment.AppointmentDate].AppointmentDate == ""{
            var dateSlot DateSlot
            dateSlot.AppointmentDate = appointment.AppointmentDate
            dateSlot.TimeSlotMap = make(map[string]string);
            dateSlot.TimeSlotMap[appointment.AppointmentTime] = appointment.AppointmentId
            provider.AppointmentSlotMap[appointment.AppointmentDate] = dateSlot
        }else if provider.UUIDMap[appointment.AppointmentId] == ""{
            provider.AppointmentSlotMap[appointment.AppointmentDate].TimeSlotMap[appointment.AppointmentTime] = appointment.AppointmentId
        }else {
            provider.AppointmentSlotMap[appointment.AppointmentDate].TimeSlotMap[provider.UUIDMap[appointment.AppointmentId]] = ""
            provider.AppointmentSlotMap[appointment.AppointmentDate].TimeSlotMap[appointment.AppointmentTime] = appointment.AppointmentId
        }
    } else {
        provider.ProviderId = appointment.ProviderId
        provider.UUIDMap = make(map[string]string);
        provider.AppointmentSlotMap = make(map[string]DateSlot);
        var dateSlot DateSlot
        dateSlot.AppointmentDate = appointment.AppointmentDate
        dateSlot.TimeSlotMap = make(map[string]string);
        dateSlot.TimeSlotMap[appointment.AppointmentTime] = appointment.AppointmentId
        provider.AppointmentSlotMap[appointment.AppointmentDate] = dateSlot
    }
    log.Println("After unmarshalling")

    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error unmarshalling activeUUIDs: %s", err); 
        return nil, err
    }
    log.Println("Printing uuidMap")
    log.Println(provider.UUIDMap)

    provider.UUIDMap[appointment.AppointmentId] = appointment.AppointmentTime
    log.Println(provider)
    UUIDsBytes, err := json.Marshal(provider)
    log.Println("Saving")
    err = stub.PutState(appointment.ProviderId, UUIDsBytes)
    log.Println("Saved")
    
    if err != nil {
        return nil, err
    }

    return providerBytes, nil
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


    if role != PAYER && role != PROVIDER && role != PHARMACY && role != PATIENT && role != SECRETARY {
        jsonResp = "{\"Error\":\"Role doesn't exist " + role + "\"}"
        return nil, errors.New(jsonResp)
    } 
    
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
    if role == SECRETARY {
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


func (t *SimpleChaincode) getActiveUUIDs(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
    activeUUIDsBytes, err := stub.GetState("activeUUIDs");
    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error fetching activeUUIDs: %s", err); 
        return nil, err
    }
    return activeUUIDsBytes, nil 
}

func (t *SimpleChaincode) getActiveUUIDsForProviders(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
    providerId := args[0]
    providerBytes, err := stub.GetState(providerId);
    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error fetching activeUUIDs: %s", err); 
        return nil, err
    }
    return providerBytes, nil 
}
