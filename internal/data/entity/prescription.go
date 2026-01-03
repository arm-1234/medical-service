package entity

import (
	"encoding/json"
	"time"
)

type Prescription struct {
	ID                     string    `gorm:"primaryKey;type:varchar(36)"`
	AppointmentID          string    `gorm:"type:varchar(36);index"`
	PatientID              string    `gorm:"type:varchar(36);not null;index"`
	PatientName            string    `gorm:"type:varchar(200)"`
	DoctorID               string    `gorm:"type:varchar(36);not null;index"`
	DoctorName             string    `gorm:"type:varchar(200)"`
	Medications            string    `gorm:"type:text"`
	Diagnosis              string    `gorm:"type:text"`
	AdditionalInstructions string    `gorm:"type:text"`
	PrescriptionDate       time.Time `gorm:"type:datetime;not null"`
	ValidUntil             time.Time `gorm:"type:datetime;not null"`
	IsActive               bool      `gorm:"type:boolean;default:true"`
	CreatedAt              time.Time `gorm:"autoCreateTime"`
}

func (Prescription) TableName() string {
	return "prescriptions"
}

type Medication struct {
	MedicationName string `json:"medication_name"`
	Dosage         string `json:"dosage"`
	Frequency      string `json:"frequency"`
	Duration       string `json:"duration"`
	Route          string `json:"route"`
	Instructions   string `json:"instructions"`
	Quantity       int32  `json:"quantity"`
}

func MarshalMedications(meds []*Medication) string {
	if meds == nil {
		return "[]"
	}
	b, _ := json.Marshal(meds)
	return string(b)
}

func UnmarshalMedications(s string) []*Medication {
	if s == "" {
		return nil
	}
	var meds []*Medication
	json.Unmarshal([]byte(s), &meds)
	return meds
}
