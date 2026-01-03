package entity

import (
	"time"
)

type MedicalRecord struct {
	ID            string     `gorm:"primaryKey;type:varchar(36)"`
	PatientID     string     `gorm:"type:varchar(36);not null;index"`
	DoctorID      string     `gorm:"type:varchar(36);index"`
	VisitDate     time.Time  `gorm:"type:datetime;not null"`
	Diagnosis     string     `gorm:"type:text"`
	Symptoms      string     `gorm:"type:text"`
	Treatment     string     `gorm:"type:text"`
	Prescriptions string     `gorm:"type:text"`
	LabResults    string     `gorm:"type:text"`
	VitalSigns    string     `gorm:"type:text"`
	Notes         string     `gorm:"type:text"`
	FollowUpDate  *time.Time `gorm:"type:datetime"`
	RecordType    string     `gorm:"type:varchar(50)"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
}

func (MedicalRecord) TableName() string {
	return "medical_records"
}

type VitalSigns struct {
	Temperature      float64 `json:"temperature,omitempty"`
	BloodPressure    string  `json:"blood_pressure,omitempty"`
	HeartRate        int32   `json:"heart_rate,omitempty"`
	RespiratoryRate  int32   `json:"respiratory_rate,omitempty"`
	OxygenSaturation int32   `json:"oxygen_saturation,omitempty"`
	Weight           float64 `json:"weight,omitempty"`
	Height           float64 `json:"height,omitempty"`
}
