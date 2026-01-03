package data

import (
	"context"
	"fmt"

	"github.com/arm-1234/medical-service/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MedicalRecordRepo interface {
	Create(ctx context.Context, record *entity.MedicalRecord) error
	Get(ctx context.Context, id string) (*entity.MedicalRecord, error)
	GetByPatientID(ctx context.Context, patientID string, filters map[string]interface{}) ([]*entity.MedicalRecord, error)
	Update(ctx context.Context, record *entity.MedicalRecord) error
	Delete(ctx context.Context, id string) error
}

type medicalRecordRepo struct {
	data *Data
	log  *log.Helper
}

func NewMedicalRecordRepo(data *Data, logger log.Logger) MedicalRecordRepo {
	return &medicalRecordRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *medicalRecordRepo) Create(ctx context.Context, record *entity.MedicalRecord) error {
	if record.ID == "" {
		record.ID = uuid.New().String()
	}

	if err := r.data.db.WithContext(ctx).Create(record).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to create medical record: %v", err)
		return err
	}
	return nil
}

func (r *medicalRecordRepo) Get(ctx context.Context, id string) (*entity.MedicalRecord, error) {
	var record entity.MedicalRecord
	if err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to get medical record: %v", err)
		return nil, err
	}
	return &record, nil
}

func (r *medicalRecordRepo) GetByPatientID(ctx context.Context, patientID string, filters map[string]interface{}) ([]*entity.MedicalRecord, error) {
	var records []*entity.MedicalRecord
	query := r.data.db.WithContext(ctx).Where("patient_id = ?", patientID)

	if fromDate, ok := filters["from_date"]; ok {
		query = query.Where("visit_date >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"]; ok {
		query = query.Where("visit_date <= ?", toDate)
	}
	if recordType, ok := filters["record_type"]; ok {
		query = query.Where("record_type = ?", recordType)
	}

	if err := query.Order("visit_date DESC").Find(&records).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to get medical records: %v", err)
		return nil, fmt.Errorf("failed to get medical records: %w", err)
	}

	return records, nil
}

func (r *medicalRecordRepo) Update(ctx context.Context, record *entity.MedicalRecord) error {
	if err := r.data.db.WithContext(ctx).Save(record).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to update medical record: %v", err)
		return err
	}
	return nil
}

func (r *medicalRecordRepo) Delete(ctx context.Context, id string) error {
	if err := r.data.db.WithContext(ctx).Delete(&entity.MedicalRecord{}, "id = ?", id).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to delete medical record: %v", err)
		return err
	}
	return nil
}
