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

type PatientService struct {
	pb.UnimplementedPatientServiceServer

	handler *biz.PatientHandler
	log     *log.Helper
}

func NewPatientService(handler *biz.PatientHandler, logger log.Logger) *PatientService {
	return &PatientService{
		handler: handler,
		log:     log.NewHelper(logger),
	}
}

func (s *PatientService) RegisterPatient(ctx context.Context, req *requestpb.RegisterPatientRequest) (*responsepb.PatientResponse, error) {
	ctx, span := otel.Trace(ctx, "PatientService.RegisterPatient")
	defer span.End()

	s.log.Infof("RegisterPatient request: %v", req)
	return s.handler.RegisterPatient(ctx, req)
}

func (s *PatientService) GetPatient(ctx context.Context, req *requestpb.GetPatientRequest) (*responsepb.PatientResponse, error) {
	ctx, span := otel.Trace(ctx, "PatientService.GetPatient")
	defer span.End()

	s.log.Infof("GetPatient request: %s", req.PatientId)
	return s.handler.GetPatient(ctx, req.PatientId)
}

func (s *PatientService) UpdatePatient(ctx context.Context, req *requestpb.UpdatePatientRequest) (*responsepb.PatientResponse, error) {
	ctx, span := otel.Trace(ctx, "PatientService.UpdatePatient")
	defer span.End()

	s.log.Infof("UpdatePatient request: %s", req.PatientId)
	return s.handler.UpdatePatient(ctx, req.PatientId, req)
}

func (s *PatientService) SearchPatients(ctx context.Context, req *requestpb.SearchPatientsRequest) (*responsepb.SearchPatientsResponse, error) {
	ctx, span := otel.Trace(ctx, "PatientService.SearchPatients")
	defer span.End()

	s.log.Infof("SearchPatients request: %v", req)

	patients, err := s.handler.SearchPatients(ctx, req)
	if err != nil {
		return nil, err
	}

	return &responsepb.SearchPatientsResponse{
		Patients: patients,
	}, nil
}

func (s *PatientService) GetMedicalHistory(ctx context.Context, req *requestpb.GetMedicalHistoryRequest) (*responsepb.MedicalHistoryResponse, error) {
	ctx, span := otel.Trace(ctx, "PatientService.GetMedicalHistory")
	defer span.End()

	s.log.Infof("GetMedicalHistory request: %s", req.PatientId)
	return s.handler.GetMedicalHistory(ctx, req.PatientId, req.GetFromDate(), req.GetToDate())
}
