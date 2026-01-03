package server

import (
	v1 "github.com/arm-1234/common-protos/medical/v1/service"
	"github.com/arm-1234/medical-service/internal/conf"
	"github.com/arm-1234/medical-service/internal/service"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewGRPCServer(
	c *conf.Server,
	patient *service.PatientService,
	doctor *service.DoctorService,
	appointment *service.AppointmentService,
	prescription *service.PrescriptionService,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterPatientServiceServer(srv, patient)
	v1.RegisterDoctorServiceServer(srv, doctor)
	v1.RegisterAppointmentServiceServer(srv, appointment)
	v1.RegisterPrescriptionServiceServer(srv, prescription)
	return srv
}
