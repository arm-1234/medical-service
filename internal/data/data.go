package data

import (
	"github.com/arm-1234/medical-service/internal/conf"
	"github.com/arm-1234/medical-service/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewData, NewPatientRepo, NewMedicalRecordRepo, NewDoctorRepo, NewAppointmentRepo, NewPrescriptionRepo)

type Data struct {
	db *gorm.DB
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)

	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil, nil, err
	}

	log.Info("database connection established")

	if err := db.AutoMigrate(
		&entity.Patient{},
		&entity.MedicalRecord{},
		&entity.Doctor{},
		&entity.DoctorAvailability{},
		&entity.Appointment{},
		&entity.Prescription{},
	); err != nil {
		log.Errorf("failed to migrate tables: %v", err)
		return nil, nil, err
	}

	log.Info("database migrations completed")

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		log.Info("closing the data resources")
	}

	return &Data{db: db}, cleanup, nil
}
