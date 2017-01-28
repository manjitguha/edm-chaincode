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


    if len(args) != 12 {
        return nil, errors.New("Incorrect number of arguments. Expecting 12. name of the variable and value to set")
    }

    appointment.AppointmentId = args[0]
    appointment.PatientId = args[1]
    appointment.ProviderId = args[2]
    appointment.ReferralProviderId = args[3]
    appointment.PharmacyId = args[4]
    appointment.SecretoryId = args[5]
    appointment.LaboratoryId = args[6]
    appointment.AppointmentDate = args[7]
    appointment.AppointmentTime = args[8]
    appointment.DiagnosisNotes = args[9]
    appointment.PrescriptionNotes = args[10]
    appointment.Status = args[11]

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
    if appointment.PatientId != "" {
        log.Println("Inside appointment.PatientId != \"\" ", appointment.PatientId)
   
        patientBytes, err  := t.saveUUIDsForPatient(stub, appointment)
        if err != nil { 
            log.Println(err)
            fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
            return nil, err
        }
        log.Println(patientBytes)
    }
    if appointment.PharmacyId != "" {
        log.Println("Inside appointment.PharmacyId != \"\" ", appointment.PharmacyId)
   
        pharmacyBytes, err  := t.saveUUIDsForPharmacy(stub, appointment)
        if err != nil { 
            log.Println(err)
            fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
            return nil, err
        }
        log.Println(pharmacyBytes)
    }
    if appointment.SecretoryId != "" {
        log.Println("Inside appointment.SecretoryId != \"\" ", appointment.SecretoryId)
   
        secretoryBytes, err  := t.saveUUIDsForSecretory(stub, appointment)
        if err != nil { 
            log.Println(err)
            fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
            return nil, err
        }
        log.Println(secretoryBytes)
    }
    if appointment.LaboratoryId != "" {
        log.Println("Inside appointment.LaboratoryId != \"\" ", appointment.LaboratoryId)
   
        laboratoryBytes, err  := t.saveUUIDsForLaboratory(stub, appointment)
        if err != nil { 
            log.Println(err)
            fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
            return nil, err
        }
        log.Println(laboratoryBytes)
    }
    if appointment.ReferralProviderId != "" {
        log.Println("Inside appointment.ReferralProviderId != \"\" ", appointment.ReferralProviderId)
        
        referralProviderBytes, err  := t.saveUUIDsForReferralProvider(stub, appointment)
        if err != nil { 
            log.Println(err)
            fmt.Printf("save_changes: Error storing Appointment record: %s", err); 
            return nil, err
        }
        log.Println(referralProviderBytes)

    }

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
        }else if provider.UUIDMap[appointment.AppointmentId] != ""{
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


func (t *SimpleChaincode) saveUUIDsForReferralProvider(stub shim.ChaincodeStubInterface, appointment Appointment) ([]byte, error) {
    var referralProvider ReferralProvider
    log.Println("Inside saveUUIDsForProvider", appointment.ReferralProviderId)

    referralProviderBytes, err := stub.GetState(appointment.ReferralProviderId);
    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error fetching activeUUIDs: %s", err); 
        return nil, err
    }

    log.Println("Before unmarshalling", len(referralProviderBytes))
    if len(referralProviderBytes) >0 {
        err = json.Unmarshal(referralProviderBytes, &referralProvider)
        if referralProvider.AppointmentSlotMap[appointment.AppointmentDate].AppointmentDate == ""{
            var dateSlot DateSlot
            dateSlot.AppointmentDate = appointment.AppointmentDate
            dateSlot.TimeSlotMap = make(map[string]string);
            dateSlot.TimeSlotMap[appointment.AppointmentTime] = appointment.AppointmentId
            referralProvider.AppointmentSlotMap[appointment.AppointmentDate] = dateSlot
        }else if referralProvider.UUIDMap[appointment.AppointmentId] != ""{
            referralProvider.AppointmentSlotMap[appointment.AppointmentDate].TimeSlotMap[appointment.AppointmentTime] = appointment.AppointmentId
        }else {
            referralProvider.AppointmentSlotMap[appointment.AppointmentDate].TimeSlotMap[referralProvider.UUIDMap[appointment.AppointmentId]] = ""
            referralProvider.AppointmentSlotMap[appointment.AppointmentDate].TimeSlotMap[appointment.AppointmentTime] = appointment.AppointmentId
        }
    } else {
        referralProvider.ReferralProviderId = appointment.ReferralProviderId
        referralProvider.UUIDMap = make(map[string]string);
        referralProvider.AppointmentSlotMap = make(map[string]DateSlot);
        var dateSlot DateSlot
        dateSlot.AppointmentDate = appointment.AppointmentDate
        dateSlot.TimeSlotMap = make(map[string]string);
        dateSlot.TimeSlotMap[appointment.AppointmentTime] = appointment.AppointmentId
        referralProvider.AppointmentSlotMap[appointment.AppointmentDate] = dateSlot
    }
    log.Println("After unmarshalling")

    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error unmarshalling activeUUIDs: %s", err); 
        return nil, err
    }
    log.Println("Printing uuidMap")
    log.Println(referralProvider.UUIDMap)

    referralProvider.UUIDMap[appointment.AppointmentId] = appointment.AppointmentTime
    log.Println(referralProvider)
    UUIDsBytes, err := json.Marshal(referralProvider)
    log.Println("Saving")
    err = stub.PutState(appointment.ReferralProviderId, UUIDsBytes)
    log.Println("Saved")
    
    if err != nil {
        return nil, err
    }

    return referralProviderBytes, nil
}

func (t *SimpleChaincode) saveUUIDsForPatient(stub shim.ChaincodeStubInterface, appointment Appointment) ([]byte, error) {
    var patient Patient
    log.Println("Inside saveUUIDsForPatient", appointment.PatientId)

    patientBytes, err := stub.GetState(appointment.PatientId);
    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error fetching activeUUIDs: %s", err); 
        return nil, err
    }

    log.Println("Before unmarshalling", len(patientBytes))
    if len(patientBytes) >0 {
        err = json.Unmarshal(patientBytes, &patient)
    } else{
        patient.PatientId = appointment.PatientId
        patient.UUIDMap = make(map[string]string);
    }
    log.Println("After unmarshalling")

    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error unmarshalling activeUUIDs: %s", err); 
        return nil, err
    }
    log.Println("Printing uuidMap")
    log.Println(patient.UUIDMap)

    patient.UUIDMap[appointment.AppointmentId] = appointment.AppointmentTime
    log.Println(patient)
    UUIDsBytes, err := json.Marshal(patient)
    log.Println("Saving")
    err = stub.PutState(appointment.PatientId, UUIDsBytes)
    log.Println("Saved")
    
    if err != nil {
        return nil, err
    }

    return patientBytes, nil
}

func (t *SimpleChaincode) saveUUIDsForPharmacy(stub shim.ChaincodeStubInterface, appointment Appointment) ([]byte, error) {
    var pharmacy Pharmacy
    log.Println("Inside saveUUIDsForPatient", appointment.PharmacyId)

    pharmacyBytes, err := stub.GetState(appointment.PharmacyId);
    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error fetching activeUUIDs: %s", err); 
        return nil, err
    }

    log.Println("Before unmarshalling", len(pharmacyBytes))
    if len(pharmacyBytes) >0 {
        err = json.Unmarshal(pharmacyBytes, &pharmacy)
    } else {
        pharmacy.PharmacyId = appointment.PharmacyId
        pharmacy.UUIDMap = make(map[string]string);
    } 
    log.Println("After unmarshalling")

    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error unmarshalling activeUUIDs: %s", err); 
        return nil, err
    }
    log.Println("Printing uuidMap")
    log.Println(pharmacy.UUIDMap)

    pharmacy.UUIDMap[appointment.AppointmentId] = appointment.AppointmentTime
    log.Println(pharmacy)
    UUIDsBytes, err := json.Marshal(pharmacy)
    log.Println("Saving")
    err = stub.PutState(appointment.PharmacyId, UUIDsBytes)
    log.Println("Saved")
    
    if err != nil {
        return nil, err
    }

    return pharmacyBytes, nil
}

func (t *SimpleChaincode) saveUUIDsForSecretory(stub shim.ChaincodeStubInterface, appointment Appointment) ([]byte, error) {
    var secretory Secretory
    log.Println("Inside saveUUIDsForSecretory", appointment.SecretoryId)

    secretoryBytes, err := stub.GetState(appointment.SecretoryId);
    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error fetching activeUUIDs: %s", err); 
        return nil, err
    }

    log.Println("Before unmarshalling", len(secretoryBytes))
    if len(secretoryBytes) >0 {
        err = json.Unmarshal(secretoryBytes, &secretory)
    } else {
        secretory.SecretoryId = appointment.SecretoryId
        secretory.UUIDMap = make(map[string]string);
    } 
    log.Println("After unmarshalling")

    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error unmarshalling activeUUIDs: %s", err); 
        return nil, err
    }
    log.Println("Printing uuidMap")
    log.Println(secretory.UUIDMap)

    secretory.UUIDMap[appointment.AppointmentId] = appointment.AppointmentTime
    log.Println(secretory)
    UUIDsBytes, err := json.Marshal(secretory)
    log.Println("Saving")
    err = stub.PutState(appointment.SecretoryId, UUIDsBytes)
    log.Println("Saved")
    
    if err != nil {
        return nil, err
    }

    return secretoryBytes, nil
}

func (t *SimpleChaincode) saveUUIDsForLaboratory(stub shim.ChaincodeStubInterface, appointment Appointment) ([]byte, error) {
    var laboratory Laboratory
    log.Println("Inside saveUUIDsForLaboratory", appointment.LaboratoryId)

    laboratoryBytes, err := stub.GetState(appointment.LaboratoryId);
    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error fetching activeUUIDs: %s", err); 
        return nil, err
    }

    log.Println("Before unmarshalling", len(laboratoryBytes))
    if len(laboratoryBytes) >0 {
        err = json.Unmarshal(laboratoryBytes, &laboratory)
    } else {
        laboratory.LaboratoryId = appointment.LaboratoryId
        laboratory.UUIDMap = make(map[string]string);
    } 
    log.Println("After unmarshalling")

    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error unmarshalling activeUUIDs: %s", err); 
        return nil, err
    }
    log.Println("Printing uuidMap")
    log.Println(laboratory.UUIDMap)

    laboratory.UUIDMap[appointment.AppointmentId] = appointment.AppointmentTime
    log.Println(laboratory)
    UUIDsBytes, err := json.Marshal(laboratory)
    log.Println("Saving")
    err = stub.PutState(appointment.LaboratoryId, UUIDsBytes)
    log.Println("Saved")
    
    if err != nil {
        return nil, err
    }

    return laboratoryBytes, nil
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

func (t *SimpleChaincode) getActiveUUIDsForID(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
    id := args[0]
    bytes, err := stub.GetState(id);
    if err != nil { 
        log.Println(err)
        fmt.Printf("save_changes: Error fetching activeUUIDs: %s", err); 
        return nil, err
    }
    return bytes, nil 
}
