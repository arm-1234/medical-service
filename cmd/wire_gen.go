
//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject

package main

import (
	"github.com/arm-1234/medical-service/internal/biz"
	"github.com/arm-1234/medical-service/internal/conf"
	"github.com/arm-1234/medical-service/internal/data"
	"github.com/arm-1234/medical-service/internal/server"
	"github.com/arm-1234/medical-service/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)


func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	patientRepo := data.NewPatientRepo(dataData, logger)
	medicalRecordRepo := data.NewMedicalRecordRepo(dataData, logger)
	patientHandler := biz.NewPatientHandler(patientRepo, medicalRecordRepo, logger)
	patientService := service.NewPatientService(patientHandler, logger)
	doctorRepo := data.NewDoctorRepo(dataData, logger)
	doctorHandler := biz.NewDoctorHandler(doctorRepo, logger)
	doctorService := service.NewDoctorService(doctorHandler, logger)
	appointmentRepo := data.NewAppointmentRepo(dataData, logger)
	appointmentHandler := biz.NewAppointmentHandler(appointmentRepo, patientRepo, doctorRepo, logger)
	appointmentService := service.NewAppointmentService(appointmentHandler, logger)
	prescriptionRepo := data.NewPrescriptionRepo(dataData, logger)
	prescriptionHandler := biz.NewPrescriptionHandler(prescriptionRepo, patientRepo, doctorRepo, logger)
	prescriptionService := service.NewPrescriptionService(prescriptionHandler, logger)
	grpcServer := server.NewGRPCServer(confServer, patientService, doctorService, appointmentService, prescriptionService)
	httpServer := server.NewHTTPServer(confServer, patientService, doctorService, appointmentService, prescriptionService)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
