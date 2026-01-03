package data

import (
	"context"

	"github.com/arm-1234/medical-service/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrescriptionRepo interface {
	Create(ctx context.Context, prescription *entity.Prescription) error
	Get(ctx context.Context, id string) (*entity.Prescription, error)
	GetByPatientID(ctx context.Context, patientID string, filters map[string]interface{}) ([]*entity.Prescription, error)
	GetByDoctorID(ctx context.Context, doctorID string, filters map[string]interface{}) ([]*entity.Prescription, error)
	GetByAppointmentID(ctx context.Context, appointmentID string) (*entity.Prescription, error)
}

type prescriptionRepo struct {
	data *Data
	log  *log.Helper
}

func NewPrescriptionRepo(data *Data, logger log.Logger) PrescriptionRepo {
	return &prescriptionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *prescriptionRepo) Create(ctx context.Context, prescription *entity.Prescription) error {
	if prescription.ID == "" {
		prescription.ID = uuid.New().String()
	}

	if err := r.data.db.WithContext(ctx).Create(prescription).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to create prescription: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("created prescription with ID: %s", prescription.ID)
	return nil
}

func (r *prescriptionRepo) Get(ctx context.Context, id string) (*entity.Prescription, error) {
	var prescription entity.Prescription

	if err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&prescription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to get prescription: %v", err)
		return nil, err
	}

	return &prescription, nil
}

func (r *prescriptionRepo) GetByPatientID(ctx context.Context, patientID string, filters map[string]interface{}) ([]*entity.Prescription, error) {
	var prescriptions []*entity.Prescription
	query := r.data.db.WithContext(ctx).Where("patient_id = ?", patientID)

	if fromDate, ok := filters["from_date"].(string); ok && fromDate != "" {
		query = query.Where("prescription_date >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"].(string); ok && toDate != "" {
		query = query.Where("prescription_date <= ?", toDate)
	}

	query = query.Order("prescription_date DESC")

	if err := query.Find(&prescriptions).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to get patient prescriptions: %v", err)
		return nil, err
	}

	return prescriptions, nil
}

func (r *prescriptionRepo) GetByDoctorID(ctx context.Context, doctorID string, filters map[string]interface{}) ([]*entity.Prescription, error) {
	var prescriptions []*entity.Prescription
	query := r.data.db.WithContext(ctx).Where("doctor_id = ?", doctorID)

	if fromDate, ok := filters["from_date"].(string); ok && fromDate != "" {
		query = query.Where("prescription_date >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"].(string); ok && toDate != "" {
		query = query.Where("prescription_date <= ?", toDate)
	}

	query = query.Order("prescription_date DESC")

	if err := query.Find(&prescriptions).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to get doctor prescriptions: %v", err)
		return nil, err
	}

	return prescriptions, nil
}

func (r *prescriptionRepo) GetByAppointmentID(ctx context.Context, appointmentID string) (*entity.Prescription, error) {
	var prescription entity.Prescription

	if err := r.data.db.WithContext(ctx).Where("appointment_id = ?", appointmentID).First(&prescription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to get prescription by appointment: %v", err)
		return nil, err
	}

	return &prescription, nil
}
