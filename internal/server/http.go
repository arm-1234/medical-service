package server

import (
	v1 "github.com/arm-1234/common-protos/medical/v1/service"
	"github.com/arm-1234/medical-service/internal/conf"
	"github.com/arm-1234/medical-service/internal/service"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(
	c *conf.Server,
	patient *service.PatientService,
	doctor *service.DoctorService,
	appointment *service.AppointmentService,
	prescription *service.PrescriptionService,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterPatientServiceHTTPServer(srv, patient)
	v1.RegisterDoctorServiceHTTPServer(srv, doctor)
	v1.RegisterAppointmentServiceHTTPServer(srv, appointment)
	v1.RegisterPrescriptionServiceHTTPServer(srv, prescription)
	return srv
}
