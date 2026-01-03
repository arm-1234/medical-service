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

type PrescriptionService struct {
	pb.UnimplementedPrescriptionServiceServer

	handler *biz.PrescriptionHandler
	log     *log.Helper
}

func NewPrescriptionService(handler *biz.PrescriptionHandler, logger log.Logger) *PrescriptionService {
	return &PrescriptionService{
		handler: handler,
		log:     log.NewHelper(logger),
	}
}

func (s *PrescriptionService) CreatePrescription(ctx context.Context, req *requestpb.CreatePrescriptionRequest) (*responsepb.PrescriptionResponse, error) {
	ctx, span := otel.Trace(ctx, "PrescriptionService.CreatePrescription")
	defer span.End()

	s.log.Infof("CreatePrescription request: patient=%s, doctor=%s", req.PatientId, req.DoctorId)
	return s.handler.CreatePrescription(ctx, req)
}

func (s *PrescriptionService) GetPrescription(ctx context.Context, req *requestpb.GetPrescriptionRequest) (*responsepb.PrescriptionResponse, error) {
	ctx, span := otel.Trace(ctx, "PrescriptionService.GetPrescription")
	defer span.End()

	s.log.Infof("GetPrescription request: %s", req.PrescriptionId)
	return s.handler.GetPrescription(ctx, req.PrescriptionId)
}

func (s *PrescriptionService) GetPatientPrescriptions(ctx context.Context, req *requestpb.GetPatientPrescriptionsRequest) (*responsepb.PatientPrescriptionsResponse, error) {
	ctx, span := otel.Trace(ctx, "PrescriptionService.GetPatientPrescriptions")
	defer span.End()

	s.log.Infof("GetPatientPrescriptions request: %s", req.PatientId)
	return s.handler.GetPatientPrescriptions(ctx, req)
}

func (s *PrescriptionService) GetDoctorPrescriptions(ctx context.Context, req *requestpb.GetDoctorPrescriptionsRequest) (*responsepb.DoctorPrescriptionsResponse, error) {
	ctx, span := otel.Trace(ctx, "PrescriptionService.GetDoctorPrescriptions")
	defer span.End()

	s.log.Infof("GetDoctorPrescriptions request: %s", req.DoctorId)
	return s.handler.GetDoctorPrescriptions(ctx, req)
}
