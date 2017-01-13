package main

import (
    "errors"
    "fmt"
    "encoding/json"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

func (t *SimpleChaincode) createAppointment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var appointment Appointment
    var appointmentId, patientId, patientFirstName, patientLastName, providerId, providerFirstName, providerLastName, appointmentTime string
    
    var err error
    fmt.Println("running createAppointment()")


    if len(args) != 8 {
        return nil, errors.New("Incorrect number of arguments. Expecting 8. name of the variable and value to set")
    }

   
    appointmentId = args[0]
    patientId = args[1]
    patientFirstName = args[2]
    patientLastName=args[3]
    providerId = args[4]
    providerFirstName = args[5]
    providerLastName = args[6]
    appointmentTime = args[7]

    appointment.AppointmentId = appointmentId
    appointment.Patient.PatientId = patientId
    appointment.Patient.PatientFirstName = patientFirstName
    appointment.Patient.PatientLastName = patientLastName
    appointment.Provider.ProviderId = providerId
    appointment.Provider.ProviderFirstName = providerFirstName
    appointment.Provider.ProviderLastName = providerLastName
    appointment.AppointmentTime = appointmentTime

  /*  appointmentId_json :=  "\"appointmentId\":\""+appointmentId+"\", "   

    patientId_json :=  "\"patientId\":\""+patientId+"\", "      
    patientFirstName_json := "\"patientFirstName\":\""+patientFirstName+"\","
    patientLastName_json := "\"patientLastName\":\""+patientLastName+"\""    
 
    providerId_json :=  "\"providerId\":\""+providerId+"\", "      
    providerFirstName_json := "\"providerFirstName\":\""+providerFirstName+"\","
    providerLastName_json := "\"providerLastName\":\""+providerLastName+"\""    
      
    appointmentTime_json :=  "\"appointmentTime\":\""+appointmentTime+"\""    


    patient_json := "\"patient\":{"+patientId_json+patientFirstName_json+patientLastName_json+"},"
    provider_json := "\"provider\":{"+providerId_json+providerFirstName_json+providerLastName_json+"},"
 
    appointment_json := "{"+appointmentId_json+patient_json+provider_json+appointmentTime_json+"}"
   

    err = json.Unmarshal([]byte(appointment_json), &appointment)  
*/


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
    var appointment Appointment
    var key, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    err = json.Unmarshal(valAsbytes, &appointment)  

    if err != nil {
        return nil, err
    }

    bytes, err := json.Marshal(appointment)
   
    if err != nil {
        return nil, err
    }


    return bytes, nil
}


