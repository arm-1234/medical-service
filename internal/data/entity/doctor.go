package entity

import (
	"encoding/json"
	"time"
)

type Doctor struct {
	ID                 string    `gorm:"primaryKey;type:varchar(36)"`
	FirstName          string    `gorm:"type:varchar(100);not null"`
	LastName           string    `gorm:"type:varchar(100);not null"`
	Email              string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PhoneNumber        string    `gorm:"type:varchar(20);uniqueIndex;not null"`
	Specialization     int32     `gorm:"type:int;not null"`
	LicenseNumber      string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	YearsOfExperience  int32     `gorm:"type:int;default:0"`
	Qualifications     string    `gorm:"type:text"`
	Languages          string    `gorm:"type:text"`
	ConsultationFee    int32     `gorm:"type:int;default:0"`
	IsAvailable        bool      `gorm:"type:boolean;default:true"`
	AverageRating      float32   `gorm:"type:float;default:0"`
	TotalConsultations int32     `gorm:"type:int;default:0"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

func (Doctor) TableName() string {
	return "doctors"
}

func MarshalStringArray(arr []string) string {
	if arr == nil {
		return "[]"
	}
	b, _ := json.Marshal(arr)
	return string(b)
}

func UnmarshalStringArray(s string) []string {
	if s == "" {
		return nil
	}
	var arr []string
	json.Unmarshal([]byte(s), &arr)
	return arr
}

type DoctorAvailability struct {
	ID                  string    `gorm:"primaryKey;type:varchar(36)"`
	DoctorID            string    `gorm:"type:varchar(36);not null;index"`
	DayOfWeek           string    `gorm:"type:varchar(20);not null"`
	StartTime           string    `gorm:"type:varchar(10);not null"`
	EndTime             string    `gorm:"type:varchar(10);not null"`
	SlotDurationMinutes int32     `gorm:"type:int;default:30"`
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
}

func (DoctorAvailability) TableName() string {
	return "doctor_availability"
}
