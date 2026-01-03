package service

import (
	"context"

	requestpb "github.com/arm-1234/common-protos/medical/v1/request"
	responsepb "github.com/arm-1234/common-protos/medical/v1/response"
	pb "github.com/arm-1234/common-protos/medical/v1/service"
	"github.com/arm-1234/medical-service/internal/biz"
	"github.com/arm-1234/medical-service/internal/pkg/otel"
	"github.com/go-kratos/kratos/v2/log"
)

type AppointmentService struct {
	pb.UnimplementedAppointmentServiceServer

	handler *biz.AppointmentHandler
	log     *log.Helper
}

func NewAppointmentService(handler *biz.AppointmentHandler, logger log.Logger) *AppointmentService {
	return &AppointmentService{
		handler: handler,
		log:     log.NewHelper(logger),
	}
}

func (s *AppointmentService) BookAppointment(ctx context.Context, req *requestpb.BookAppointmentRequest) (*responsepb.AppointmentResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentService.BookAppointment")
	defer span.End()

	s.log.Infof("BookAppointment request: patient=%s, doctor=%s", req.PatientId, req.DoctorId)
	return s.handler.BookAppointment(ctx, req)
}

func (s *AppointmentService) GetAppointment(ctx context.Context, req *requestpb.GetAppointmentRequest) (*responsepb.AppointmentResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentService.GetAppointment")
	defer span.End()

	s.log.Infof("GetAppointment request: %s", req.AppointmentId)
	return s.handler.GetAppointment(ctx, req.AppointmentId)
}

func (s *AppointmentService) CancelAppointment(ctx context.Context, req *requestpb.CancelAppointmentRequest) (*responsepb.AppointmentResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentService.CancelAppointment")
	defer span.End()

	s.log.Infof("CancelAppointment request: %s", req.AppointmentId)
	return s.handler.CancelAppointment(ctx, req.AppointmentId, req.CancellationReason)
}

func (s *AppointmentService) RescheduleAppointment(ctx context.Context, req *requestpb.RescheduleAppointmentRequest) (*responsepb.AppointmentResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentService.RescheduleAppointment")
	defer span.End()

	s.log.Infof("RescheduleAppointment request: %s", req.AppointmentId)
	return s.handler.RescheduleAppointment(ctx, req)
}

func (s *AppointmentService) CompleteAppointment(ctx context.Context, req *requestpb.CompleteAppointmentRequest) (*responsepb.AppointmentResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentService.CompleteAppointment")
	defer span.End()

	s.log.Infof("CompleteAppointment request: %s", req.AppointmentId)
	return s.handler.CompleteAppointment(ctx, req)
}

func (s *AppointmentService) GetAvailableSlots(ctx context.Context, req *requestpb.GetAvailableSlotsRequest) (*responsepb.AvailableSlotsResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentService.GetAvailableSlots")
	defer span.End()

	s.log.Infof("GetAvailableSlots request: doctor=%s, date=%s", req.DoctorId, req.Date)
	return s.handler.GetAvailableSlots(ctx, req)
}

func (s *AppointmentService) GetPatientAppointments(ctx context.Context, req *requestpb.GetPatientAppointmentsRequest) (*responsepb.PatientAppointmentsResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentService.GetPatientAppointments")
	defer span.End()

	s.log.Infof("GetPatientAppointments request: %s", req.PatientId)
	return s.handler.GetPatientAppointments(ctx, req)
}

func (s *AppointmentService) GetDoctorAppointments(ctx context.Context, req *requestpb.GetDoctorAppointmentsRequest) (*responsepb.DoctorAppointmentsResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentService.GetDoctorAppointments")
	defer span.End()

	s.log.Infof("GetDoctorAppointments request: %s", req.DoctorId)
	return s.handler.GetDoctorAppointments(ctx, req)
}
