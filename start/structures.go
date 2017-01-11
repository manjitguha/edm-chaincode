package main

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
}


//=============================================================================================================================
// Provider - Defines the structure of Provider entity
//=============================================================================================================================
type Provider struct{
    ProviderId string `json:"providerId"`
    ProviderFirstName string `json:"providerFirstName"`
    ProviderLastName string `json:"providerLastName"`
}

//=============================================================================================================================
// Appointment - Defines the structure of Appointment entity
//=============================================================================================================================
type Appointment struct{
    AppointmentId string `json:"appointmentId"`
    Provider Provider `json:"provider"`
    Patient Patient `json:"patient"`
    AppointmentTime string `json:"appointmentTime"`
}