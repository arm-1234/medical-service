package biz

import (
	"context"
	"fmt"

	commonpb "github.com/arm-1234/common-protos/medical/v1/common"
	requestpb "github.com/arm-1234/common-protos/medical/v1/request"
	responsepb "github.com/arm-1234/common-protos/medical/v1/response"
	"github.com/arm-1234/medical-service/internal/data"
	"github.com/arm-1234/medical-service/internal/data/entity"
	"github.com/arm-1234/medical-service/internal/pkg/otel"
	"github.com/go-kratos/kratos/v2/log"
)

type DoctorHandler struct {
	repo data.DoctorRepo
	log  *log.Helper
}

func NewDoctorHandler(repo data.DoctorRepo, logger log.Logger) *DoctorHandler {
	return &DoctorHandler{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (h *DoctorHandler) RegisterDoctor(ctx context.Context, req *requestpb.RegisterDoctorRequest) (*responsepb.DoctorResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorHandler.RegisterDoctor")
	defer span.End()

	if req.FirstName == "" || req.LastName == "" {
		h.log.WithContext(ctx).Errorf("First name and last name are required")
		return nil, fmt.Errorf("first_name and last_name are required")
	}
	if req.Email == "" || req.PhoneNumber == "" {
		h.log.WithContext(ctx).Errorf("Email and phone number are required")
		return nil, fmt.Errorf("email and phone_number are required")
	}
	if req.LicenseNumber == "" {
		h.log.WithContext(ctx).Errorf("License number is required")
		return nil, fmt.Errorf("license_number is required")
	}

	existing, err := h.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to check existing email: %v", err)
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existing != nil {
		h.log.WithContext(ctx).Errorf("Email already registered: %s", req.Email)
		return nil, fmt.Errorf("email already registered")
	}

	existing, err = h.repo.GetByLicense(ctx, req.LicenseNumber)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to check existing license: %v", err)
		return nil, fmt.Errorf("failed to check existing license: %w", err)
	}
	if existing != nil {
		h.log.WithContext(ctx).Errorf("License number already registered: %s", req.LicenseNumber)
		return nil, fmt.Errorf("license number already registered")
	}

	doctor := &entity.Doctor{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		Specialization:    int32(req.Specialization),
		LicenseNumber:     req.LicenseNumber,
		YearsOfExperience: req.YearsOfExperience,
		Qualifications:    entity.MarshalStringArray(req.Qualifications),
		Languages:         entity.MarshalStringArray(req.Languages),
		ConsultationFee:   req.ConsultationFee,
		IsAvailable:       true,
	}

	if err := h.repo.Create(ctx, doctor); err != nil {
		h.log.WithContext(ctx).Errorf("Failed to create doctor: %v", err)
		return nil, fmt.Errorf("failed to create doctor: %w", err)
	}

	return h.entityToProto(doctor), nil
}

func (h *DoctorHandler) GetDoctor(ctx context.Context, id string) (*responsepb.DoctorResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorHandler.GetDoctor")
	defer span.End()

	if id == "" {
		h.log.WithContext(ctx).Errorf("Doctor ID is required")
		return nil, fmt.Errorf("doctor_id is required")
	}

	doctor, err := h.repo.Get(ctx, id)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get doctor: %v", err)
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}

	if doctor == nil {
		return nil, fmt.Errorf("doctor not found")
	}

	return h.entityToProto(doctor), nil
}

func (h *DoctorHandler) UpdateDoctor(ctx context.Context, id string, req *requestpb.UpdateDoctorRequest) (*responsepb.DoctorResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorHandler.UpdateDoctor")
	defer span.End()

	if id == "" {
		h.log.WithContext(ctx).Errorf("Doctor ID is required")
		return nil, fmt.Errorf("doctor_id is required")
	}

	doctor, err := h.repo.Get(ctx, id)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get doctor: %v", err)
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}
	if doctor == nil {
		h.log.WithContext(ctx).Errorf("Doctor not found: %s", id)
		return nil, fmt.Errorf("doctor not found")
	}

	if req.PhoneNumber != nil {
		doctor.PhoneNumber = req.GetPhoneNumber()
	}
	if req.Email != nil {
		doctor.Email = req.GetEmail()
	}
	if req.ConsultationFee != nil {
		doctor.ConsultationFee = req.GetConsultationFee()
	}
	if req.IsAvailable != nil {
		doctor.IsAvailable = req.GetIsAvailable()
	}

	if err := h.repo.Update(ctx, doctor); err != nil {
		h.log.WithContext(ctx).Errorf("Failed to update doctor: %v", err)
		return nil, fmt.Errorf("failed to update doctor: %w", err)
	}

	return h.entityToProto(doctor), nil
}

func (h *DoctorHandler) SearchDoctors(ctx context.Context, req *requestpb.SearchDoctorsRequest) ([]*responsepb.DoctorResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorHandler.SearchDoctors")
	defer span.End()

	filters := make(map[string]interface{})

	if req.GetName() != "" {
		filters["name"] = req.GetName()
	}
	if req.Specialization != nil && req.GetSpecialization() > 0 {
		filters["specialization"] = int32(req.GetSpecialization())
	}
	if req.IsAvailable != nil {
		filters["is_available"] = req.GetIsAvailable()
	}

	doctors, err := h.repo.Search(ctx, filters)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to search doctors: %v", err)
		return nil, fmt.Errorf("failed to search doctors: %w", err)
	}

	var responses []*responsepb.DoctorResponse
	for _, d := range doctors {
		responses = append(responses, h.entityToProto(d))
	}

	return responses, nil
}

func (h *DoctorHandler) SetAvailability(ctx context.Context, doctorID string, slots []*requestpb.AvailabilitySlot) (*responsepb.DoctorAvailabilityResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorHandler.SetAvailability")
	defer span.End()

	if doctorID == "" {
		h.log.WithContext(ctx).Errorf("Doctor ID is required")
		return nil, fmt.Errorf("doctor_id is required")
	}

	doctor, err := h.repo.Get(ctx, doctorID)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get doctor: %v", err)
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}
	if doctor == nil {
		h.log.WithContext(ctx).Errorf("Doctor not found: %s", doctorID)
		return nil, fmt.Errorf("doctor not found")
	}

	var entitySlots []*entity.DoctorAvailability
	for _, slot := range slots {
		entitySlots = append(entitySlots, &entity.DoctorAvailability{
			DayOfWeek:           slot.DayOfWeek,
			StartTime:           slot.StartTime,
			EndTime:             slot.EndTime,
			SlotDurationMinutes: slot.SlotDurationMinutes,
		})
	}

	if err := h.repo.SetAvailability(ctx, doctorID, entitySlots); err != nil {
		h.log.WithContext(ctx).Errorf("Failed to set availability: %v", err)
		return nil, fmt.Errorf("failed to set availability: %w", err)
	}

	return h.GetDoctorAvailability(ctx, doctorID)
}

func (h *DoctorHandler) GetDoctorAvailability(ctx context.Context, doctorID string) (*responsepb.DoctorAvailabilityResponse, error) {
	ctx, span := otel.Trace(ctx, "DoctorHandler.GetDoctorAvailability")
	defer span.End()

	if doctorID == "" {
		h.log.WithContext(ctx).Errorf("Doctor ID is required")
		return nil, fmt.Errorf("doctor_id is required")
	}

	doctor, err := h.repo.Get(ctx, doctorID)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get doctor: %v", err)
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}
	if doctor == nil {
		h.log.WithContext(ctx).Errorf("Doctor not found: %s", doctorID)
		return nil, fmt.Errorf("doctor not found")
	}

	slots, err := h.repo.GetAvailability(ctx, doctorID)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get availability: %v", err)
		return nil, fmt.Errorf("failed to get availability: %w", err)
	}

	var protoSlots []*requestpb.AvailabilitySlot
	for _, slot := range slots {
		protoSlots = append(protoSlots, &requestpb.AvailabilitySlot{
			DayOfWeek:           slot.DayOfWeek,
			StartTime:           slot.StartTime,
			EndTime:             slot.EndTime,
			SlotDurationMinutes: slot.SlotDurationMinutes,
		})
	}

	return &responsepb.DoctorAvailabilityResponse{
		DoctorId:          doctorID,
		AvailabilitySlots: protoSlots,
	}, nil
}

func (h *DoctorHandler) entityToProto(doctor *entity.Doctor) *responsepb.DoctorResponse {
	return &responsepb.DoctorResponse{
		DoctorId:           doctor.ID,
		FirstName:          doctor.FirstName,
		LastName:           doctor.LastName,
		Email:              doctor.Email,
		PhoneNumber:        doctor.PhoneNumber,
		Specialization:     commonpb.Specialization(doctor.Specialization),
		LicenseNumber:      doctor.LicenseNumber,
		YearsOfExperience:  doctor.YearsOfExperience,
		Qualifications:     entity.UnmarshalStringArray(doctor.Qualifications),
		Languages:          entity.UnmarshalStringArray(doctor.Languages),
		ConsultationFee:    doctor.ConsultationFee,
		IsAvailable:        doctor.IsAvailable,
		AverageRating:      doctor.AverageRating,
		TotalConsultations: doctor.TotalConsultations,
		CreatedAt:          doctor.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:          doctor.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
