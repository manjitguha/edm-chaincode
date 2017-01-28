package main

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//=============================================================================================================================
// Appointment - Defines the structure of Appointment entity
//=============================================================================================================================
type Appointment struct{
    AppointmentId string `json:"appointmentId"`
    ProviderId string `json:"providerId"`
    ReferralProviderId string `json:"referralProviderId"`
    PatientId string `json:"patientId"`
    PharmacyId string `json:"pharmacyId"`
    SecretoryId string `json:"secretoryId"`
    LaboratoryId string `json:"laboratoryId"`
    AppointmentTime string `json:"appointmentTime"`
    DiagnosisNotes string `json:"diagnosisNotes"`
    PrescriptionNotes string `json:"prescriptionNotes"`
    Status string `json:"status"`
}


//=============================================================================================================================
// Patient - Defines the structure of Patient entity
//=============================================================================================================================
type Patient struct{
    PatientId string `json:"patientId"`
    uuidArray []string `json:"uuidArray"`
}

//=============================================================================================================================
// Provider - Defines the structure of Provider entity
//=============================================================================================================================
type Provider struct{
    ProviderId string `json:"providerId"`
    uuidArray []string `json:"uuidArray"`
}


//=============================================================================================================================
// Secretory - Defines the structure of Secretory entity
//=============================================================================================================================
type Secretory struct{
    SecretoryId string `json:"secretoryId"`
    uuidArray []string `json:"uuidArray"`
}

//=============================================================================================================================
// Pharmacy - Defines the structure of Pharmacy entity
//=============================================================================================================================
type Pharmacy struct{
    PharmacyId string `json:"pharmacyId"`
    uuidArray []string `json:"uuidArray"`
}


//=============================================================================================================================
// Laboratory - Defines the structure of Pharmacy entity
//=============================================================================================================================
type Laboratory struct{
    LaboratoryId string `json:"laboratoryId"`
    uuidArray []string `json:"uuidArray"`
}



type ActiveUUIDs struct{
    uuidArray []string `json:"uuidArray"`
}

const (
    PAYER = "PAYER"
    PROVIDER = "PROVIDER"
    PHARMACY = "PHARMACY"
    PATIENT = "PATIENT"
    SECRETARY = "SECRETARY"
    UNAUTHORIZED = "UNAUTHORIZED"
)