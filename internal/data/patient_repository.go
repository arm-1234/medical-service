package data

import (
	"context"

	"github.com/arm-1234/medical-service/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PatientRepo interface {
	Create(ctx context.Context, patient *entity.Patient) error
	Get(ctx context.Context, id string) (*entity.Patient, error)
	Update(ctx context.Context, patient *entity.Patient) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, filters map[string]interface{}) ([]*entity.Patient, error)
	GetByEmail(ctx context.Context, email string) (*entity.Patient, error)
	GetByPhone(ctx context.Context, phone string) (*entity.Patient, error)
}

type patientRepo struct {
	data *Data
	log  *log.Helper
}

func NewPatientRepo(data *Data, logger log.Logger) PatientRepo {
	return &patientRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *patientRepo) Create(ctx context.Context, patient *entity.Patient) error {
	if patient.ID == "" {
		patient.ID = uuid.New().String()
	}

	if err := r.data.db.WithContext(ctx).Create(patient).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to create patient: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("created patient with ID: %s", patient.ID)
	return nil
}

func (r *patientRepo) Get(ctx context.Context, id string) (*entity.Patient, error) {
	var patient entity.Patient

	if err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to get patient: %v", err)
		return nil, err
	}

	return &patient, nil
}

func (r *patientRepo) Update(ctx context.Context, patient *entity.Patient) error {
	if err := r.data.db.WithContext(ctx).Save(patient).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to update patient: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("updated patient with ID: %s", patient.ID)
	return nil
}

func (r *patientRepo) Delete(ctx context.Context, id string) error {
	if err := r.data.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Patient{}).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to delete patient: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("deleted patient with ID: %s", id)
	return nil
}

func (r *patientRepo) Search(ctx context.Context, filters map[string]interface{}) ([]*entity.Patient, error) {
	var patients []*entity.Patient
	query := r.data.db.WithContext(ctx)

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("first_name LIKE ? OR last_name LIKE ?", "%"+name+"%", "%"+name+"%")
	}
	if email, ok := filters["email"].(string); ok && email != "" {
		query = query.Where("email = ?", email)
	}
	if phone, ok := filters["phone_number"].(string); ok && phone != "" {
		query = query.Where("phone_number = ?", phone)
	}
	if patientID, ok := filters["patient_id"].(string); ok && patientID != "" {
		query = query.Where("id = ?", patientID)
	}

	if err := query.Find(&patients).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to search patients: %v", err)
		return nil, err
	}

	return patients, nil
}

func (r *patientRepo) GetByEmail(ctx context.Context, email string) (*entity.Patient, error) {
	var patient entity.Patient

	if err := r.data.db.WithContext(ctx).Where("email = ?", email).First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to get patient by email: %v", err)
		return nil, err
	}

	return &patient, nil
}

func (r *patientRepo) GetByPhone(ctx context.Context, phone string) (*entity.Patient, error) {
	var patient entity.Patient

	if err := r.data.db.WithContext(ctx).Where("phone_number = ?", phone).First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to get patient by phone: %v", err)
		return nil, err
	}

	return &patient, nil
}
