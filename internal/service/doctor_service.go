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

type DoctorService struct {
	pb.UnimplementedDoctorServiceServer

	handler *biz.DoctorHandler
	log     *log.Helper
}

func NewDoctorService(handler *biz.DoctorHandler, logger log.Logger) *DoctorService {
	return &DoctorService{
		handler: handler,
		log:     log.NewHelper(logger),
	}
}

func (s *DoctorService) RegisterDoctor(ctx context.Context, req *requestpb.RegisterDoctorRequest) (*responsepb.DoctorResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorService.RegisterDoctor")
	defer span.End()

	s.log.Infof("RegisterDoctor request: %v", req)
	return s.handler.RegisterDoctor(ctx, req)
}

func (s *DoctorService) GetDoctor(ctx context.Context, req *requestpb.GetDoctorRequest) (*responsepb.DoctorResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorService.GetDoctor")
	defer span.End()

	s.log.Infof("GetDoctor request: %s", req.DoctorId)
	return s.handler.GetDoctor(ctx, req.DoctorId)
}

func (s *DoctorService) UpdateDoctor(ctx context.Context, req *requestpb.UpdateDoctorRequest) (*responsepb.DoctorResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorService.UpdateDoctor")
	defer span.End()

	s.log.Infof("UpdateDoctor request: %s", req.DoctorId)
	return s.handler.UpdateDoctor(ctx, req.DoctorId, req)
}

func (s *DoctorService) SearchDoctors(ctx context.Context, req *requestpb.SearchDoctorsRequest) (*responsepb.SearchDoctorsResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorService.SearchDoctors")
	defer span.End()

	s.log.Infof("SearchDoctors request: %v", req)

	doctors, err := s.handler.SearchDoctors(ctx, req)
	if err != nil {
		return nil, err
	}

	return &responsepb.SearchDoctorsResponse{
		Doctors: doctors,
	}, nil
}

func (s *DoctorService) SetAvailability(ctx context.Context, req *requestpb.SetAvailabilityRequest) (*responsepb.DoctorAvailabilityResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorService.SetAvailability")
	defer span.End()

	s.log.Infof("SetAvailability request for doctor: %s", req.DoctorId)
	return s.handler.SetAvailability(ctx, req.DoctorId, req.AvailabilitySlots)
}

func (s *DoctorService) GetDoctorAvailability(ctx context.Context, req *requestpb.GetDoctorAvailabilityRequest) (*responsepb.DoctorAvailabilityResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorService.GetDoctorAvailability")
	defer span.End()

	s.log.Infof("GetDoctorAvailability request for doctor: %s", req.DoctorId)
	return s.handler.GetDoctorAvailability(ctx, req.DoctorId)
}
