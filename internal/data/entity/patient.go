package entity

import (
	"encoding/json"
	"time"
)

type Patient struct {
	ID               string    `gorm:"primaryKey;type:varchar(36)"`
	FirstName        string    `gorm:"type:varchar(100);not null"`
	LastName         string    `gorm:"type:varchar(100);not null"`
	Email            string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PhoneNumber      string    `gorm:"type:varchar(20);uniqueIndex;not null"`
	DateOfBirth      string    `gorm:"type:varchar(10)"`
	Gender           int32     `gorm:"type:int"`
	BloodGroup       int32     `gorm:"type:int"`
	Address          string    `gorm:"type:text"`
	MedicalHistory   string    `gorm:"type:text"`
	EmergencyContact string    `gorm:"type:text"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

func (Patient) TableName() string {
	return "patients"
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

func MarshalAddress(addr *Address) string {
	if addr == nil {
		return ""
	}
	b, _ := json.Marshal(addr)
	return string(b)
}

func UnmarshalAddress(s string) *Address {
	if s == "" {
		return nil
	}
	var addr Address
	json.Unmarshal([]byte(s), &addr)
	return &addr
}
