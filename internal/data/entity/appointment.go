package entity

import (
	"time"
)

type Appointment struct {
	ID                 string     `gorm:"primaryKey;type:varchar(36)"`
	PatientID          string     `gorm:"type:varchar(36);not null;index"`
	PatientName        string     `gorm:"type:varchar(200)"`
	DoctorID           string     `gorm:"type:varchar(36);not null;index"`
	DoctorName         string     `gorm:"type:varchar(200)"`
	AppointmentDate    string     `gorm:"type:varchar(10);not null;index"`
	AppointmentTime    string     `gorm:"type:varchar(10);not null"`
	Status             int32      `gorm:"type:int;not null;default:1"`
	ConsultationType   int32      `gorm:"type:int;not null;default:1"`
	ReasonForVisit     string     `gorm:"type:text"`
	Notes              string     `gorm:"type:text"`
	Diagnosis          string     `gorm:"type:text"`
	CancelledAt        *time.Time `gorm:"type:datetime"`
	CancellationReason string     `gorm:"type:text"`
	CreatedAt          time.Time  `gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime"`
}

func (Appointment) TableName() string {
	return "appointments"
}

const (
	AppointmentStatusUnspecified = 0
	AppointmentStatusScheduled   = 1
	AppointmentStatusConfirmed   = 2
	AppointmentStatusInProgress  = 3
	AppointmentStatusCompleted   = 4
	AppointmentStatusCancelled   = 5
	AppointmentStatusNoShow      = 6
	AppointmentStatusRescheduled = 7
)

const (
	ConsultationTypeUnspecified = 0
	ConsultationTypeInPerson    = 1
	ConsultationTypeVideo       = 2
	ConsultationTypePhone       = 3
)
