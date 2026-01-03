package data

import (
	"context"
	"time"

	"github.com/arm-1234/medical-service/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentRepo interface {
	Create(ctx context.Context, appointment *entity.Appointment) error
	Get(ctx context.Context, id string) (*entity.Appointment, error)
	Update(ctx context.Context, appointment *entity.Appointment) error
	GetByPatientID(ctx context.Context, patientID string, filters map[string]interface{}) ([]*entity.Appointment, error)
	GetByDoctorID(ctx context.Context, doctorID string, filters map[string]interface{}) ([]*entity.Appointment, error)
	GetByDoctorAndDate(ctx context.Context, doctorID string, date string) ([]*entity.Appointment, error)
	CheckConflict(ctx context.Context, doctorID, date, time string, excludeID string) (*entity.Appointment, error)
}

type appointmentRepo struct {
	data *Data
	log  *log.Helper
}

func NewAppointmentRepo(data *Data, logger log.Logger) AppointmentRepo {
	return &appointmentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *appointmentRepo) Create(ctx context.Context, appointment *entity.Appointment) error {
	if appointment.ID == "" {
		appointment.ID = uuid.New().String()
	}

	if err := r.data.db.WithContext(ctx).Create(appointment).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to create appointment: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("created appointment with ID: %s", appointment.ID)
	return nil
}

func (r *appointmentRepo) Get(ctx context.Context, id string) (*entity.Appointment, error) {
	var appointment entity.Appointment

	if err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to get appointment: %v", err)
		return nil, err
	}

	return &appointment, nil
}

func (r *appointmentRepo) Update(ctx context.Context, appointment *entity.Appointment) error {
	if err := r.data.db.WithContext(ctx).Save(appointment).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to update appointment: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("updated appointment with ID: %s", appointment.ID)
	return nil
}

func (r *appointmentRepo) GetByPatientID(ctx context.Context, patientID string, filters map[string]interface{}) ([]*entity.Appointment, error) {
	var appointments []*entity.Appointment
	query := r.data.db.WithContext(ctx).Where("patient_id = ?", patientID)

	if status, ok := filters["status"].(int32); ok && status > 0 {
		query = query.Where("status = ?", status)
	}
	if fromDate, ok := filters["from_date"].(string); ok && fromDate != "" {
		query = query.Where("appointment_date >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"].(string); ok && toDate != "" {
		query = query.Where("appointment_date <= ?", toDate)
	}

	query = query.Order("appointment_date DESC, appointment_time DESC")

	if err := query.Find(&appointments).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to get patient appointments: %v", err)
		return nil, err
	}

	return appointments, nil
}

func (r *appointmentRepo) GetByDoctorID(ctx context.Context, doctorID string, filters map[string]interface{}) ([]*entity.Appointment, error) {
	var appointments []*entity.Appointment
	query := r.data.db.WithContext(ctx).Where("doctor_id = ?", doctorID)

	if status, ok := filters["status"].(int32); ok && status > 0 {
		query = query.Where("status = ?", status)
	}
	if date, ok := filters["date"].(string); ok && date != "" {
		query = query.Where("appointment_date = ?", date)
	}

	query = query.Order("appointment_date DESC, appointment_time DESC")

	if err := query.Find(&appointments).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to get doctor appointments: %v", err)
		return nil, err
	}

	return appointments, nil
}

func (r *appointmentRepo) GetByDoctorAndDate(ctx context.Context, doctorID string, date string) ([]*entity.Appointment, error) {
	var appointments []*entity.Appointment

	query := r.data.db.WithContext(ctx).
		Where("doctor_id = ?", doctorID).
		Where("appointment_date = ?", date).
		Where("status NOT IN (?)", []int32{entity.AppointmentStatusCancelled})

	if err := query.Order("appointment_time ASC").Find(&appointments).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to get appointments by doctor and date: %v", err)
		return nil, err
	}

	return appointments, nil
}

func (r *appointmentRepo) CheckConflict(ctx context.Context, doctorID, date, appointmentTime string, excludeID string) (*entity.Appointment, error) {
	var appointment entity.Appointment

	query := r.data.db.WithContext(ctx).
		Where("doctor_id = ?", doctorID).
		Where("appointment_date = ?", date).
		Where("appointment_time = ?", appointmentTime).
		Where("status NOT IN (?)", []int32{entity.AppointmentStatusCancelled})

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to check appointment conflict: %v", err)
		return nil, err
	}

	return &appointment, nil
}

func FormatTimePointer(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05Z")
}
