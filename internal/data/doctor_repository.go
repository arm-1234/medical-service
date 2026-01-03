package data

import (
	"context"

	"github.com/arm-1234/medical-service/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DoctorRepo interface {
	Create(ctx context.Context, doctor *entity.Doctor) error
	Get(ctx context.Context, id string) (*entity.Doctor, error)
	Update(ctx context.Context, doctor *entity.Doctor) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, filters map[string]interface{}) ([]*entity.Doctor, error)
	GetByEmail(ctx context.Context, email string) (*entity.Doctor, error)
	GetByLicense(ctx context.Context, license string) (*entity.Doctor, error)
	SetAvailability(ctx context.Context, doctorID string, slots []*entity.DoctorAvailability) error
	GetAvailability(ctx context.Context, doctorID string) ([]*entity.DoctorAvailability, error)
}

type doctorRepo struct {
	data *Data
	log  *log.Helper
}

func NewDoctorRepo(data *Data, logger log.Logger) DoctorRepo {
	return &doctorRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *doctorRepo) Create(ctx context.Context, doctor *entity.Doctor) error {
	if doctor.ID == "" {
		doctor.ID = uuid.New().String()
	}

	if err := r.data.db.WithContext(ctx).Create(doctor).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to create doctor: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("created doctor with ID: %s", doctor.ID)
	return nil
}

func (r *doctorRepo) Get(ctx context.Context, id string) (*entity.Doctor, error) {
	var doctor entity.Doctor

	if err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&doctor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to get doctor: %v", err)
		return nil, err
	}

	return &doctor, nil
}

func (r *doctorRepo) Update(ctx context.Context, doctor *entity.Doctor) error {
	if err := r.data.db.WithContext(ctx).Save(doctor).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to update doctor: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("updated doctor with ID: %s", doctor.ID)
	return nil
}

func (r *doctorRepo) Delete(ctx context.Context, id string) error {
	if err := r.data.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Doctor{}).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to delete doctor: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("deleted doctor with ID: %s", id)
	return nil
}

func (r *doctorRepo) Search(ctx context.Context, filters map[string]interface{}) ([]*entity.Doctor, error) {
	var doctors []*entity.Doctor
	query := r.data.db.WithContext(ctx)

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("first_name LIKE ? OR last_name LIKE ?", "%"+name+"%", "%"+name+"%")
	}
	if spec, ok := filters["specialization"].(int32); ok && spec > 0 {
		query = query.Where("specialization = ?", spec)
	}
	if isAvailable, ok := filters["is_available"].(bool); ok {
		query = query.Where("is_available = ?", isAvailable)
	}

	if err := query.Find(&doctors).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to search doctors: %v", err)
		return nil, err
	}

	return doctors, nil
}

func (r *doctorRepo) GetByEmail(ctx context.Context, email string) (*entity.Doctor, error) {
	var doctor entity.Doctor

	if err := r.data.db.WithContext(ctx).Where("email = ?", email).First(&doctor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to get doctor by email: %v", err)
		return nil, err
	}

	return &doctor, nil
}

func (r *doctorRepo) GetByLicense(ctx context.Context, license string) (*entity.Doctor, error) {
	var doctor entity.Doctor

	if err := r.data.db.WithContext(ctx).Where("license_number = ?", license).First(&doctor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("failed to get doctor by license: %v", err)
		return nil, err
	}

	return &doctor, nil
}

func (r *doctorRepo) SetAvailability(ctx context.Context, doctorID string, slots []*entity.DoctorAvailability) error {
	if err := r.data.db.WithContext(ctx).Where("doctor_id = ?", doctorID).Delete(&entity.DoctorAvailability{}).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to delete existing availability: %v", err)
		return err
	}

	for _, slot := range slots {
		if slot.ID == "" {
			slot.ID = uuid.New().String()
		}
		slot.DoctorID = doctorID
		if err := r.data.db.WithContext(ctx).Create(slot).Error; err != nil {
			r.log.WithContext(ctx).Errorf("failed to create availability slot: %v", err)
			return err
		}
	}

	r.log.WithContext(ctx).Infof("set availability for doctor: %s with %d slots", doctorID, len(slots))
	return nil
}

func (r *doctorRepo) GetAvailability(ctx context.Context, doctorID string) ([]*entity.DoctorAvailability, error) {
	var slots []*entity.DoctorAvailability

	if err := r.data.db.WithContext(ctx).Where("doctor_id = ?", doctorID).Find(&slots).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to get availability: %v", err)
		return nil, err
	}

	return slots, nil
}
