package service

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewPatientService, NewDoctorService, NewAppointmentService, NewPrescriptionService)
