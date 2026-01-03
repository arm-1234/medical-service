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

type PatientHandler struct {
	repo       data.PatientRepo
	recordRepo data.MedicalRecordRepo
	log        *log.Helper
}

func NewPatientHandler(repo data.PatientRepo, recordRepo data.MedicalRecordRepo, logger log.Logger) *PatientHandler {
	return &PatientHandler{
		repo:       repo,
		recordRepo: recordRepo,
		log:        log.NewHelper(logger),
	}
}

func (h *PatientHandler) RegisterPatient(ctx context.Context, req *requestpb.RegisterPatientRequest) (*responsepb.PatientResponse, error) {
	ctx, span := otel.Trace(ctx, "PatientHandler.RegisterPatient")
	defer span.End()

	if req.FirstName == "" || req.LastName == "" {
		h.log.WithContext(ctx).Errorf("First name and last name are required")
		return nil, fmt.Errorf("first_name and last_name are required")
	}
	if req.Email == "" || req.PhoneNumber == "" {
		h.log.WithContext(ctx).Errorf("Email and phone number are required")
		return nil, fmt.Errorf("email and phone_number are required")
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

	existing, err = h.repo.GetByPhone(ctx, req.PhoneNumber)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to check existing phone: %v", err)
		return nil, fmt.Errorf("failed to check existing phone: %w", err)
	}
	if existing != nil {
		h.log.WithContext(ctx).Errorf("Phone number already registered: %s", req.PhoneNumber)
		return nil, fmt.Errorf("phone number already registered")
	}

	patient := &entity.Patient{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		DateOfBirth: req.DateOfBirth,
		Gender:      int32(req.Gender),
		BloodGroup:  int32(req.BloodGroup),
	}

	if req.Address != nil {
		addr := &entity.Address{
			Street:  req.Address.Street,
			City:    req.Address.City,
			State:   req.Address.State,
			ZipCode: req.Address.ZipCode,
			Country: req.Address.Country,
		}
		patient.Address = entity.MarshalAddress(addr)
	}

	if err := h.repo.Create(ctx, patient); err != nil {
		h.log.WithContext(ctx).Errorf("Failed to create patient: %v", err)
		return nil, fmt.Errorf("failed to create patient: %w", err)
	}

	return h.entityToProto(patient), nil
}

func (h *PatientHandler) GetPatient(ctx context.Context, id string) (*responsepb.PatientResponse, error) {
	ctx, span := otel.Trace(ctx, "PatientHandler.GetPatient")
	defer span.End()

	if id == "" {
		h.log.WithContext(ctx).Errorf("Patient ID is required")
		return nil, fmt.Errorf("patient_id is required")
	}

	patient, err := h.repo.Get(ctx, id)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get patient: %v", err)
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}

	if patient == nil {
		return nil, fmt.Errorf("patient not found")
	}

	return h.entityToProto(patient), nil
}

func (h *PatientHandler) UpdatePatient(ctx context.Context, id string, req *requestpb.UpdatePatientRequest) (*responsepb.PatientResponse, error) {
	ctx, span := otel.Trace(ctx, "PatientHandler.UpdatePatient")
	defer span.End()

	if id == "" {
		h.log.WithContext(ctx).Errorf("Patient ID is required")
		return nil, fmt.Errorf("patient_id is required")
	}

	patient, err := h.repo.Get(ctx, id)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get patient: %v", err)
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		h.log.WithContext(ctx).Errorf("Patient not found: %s", id)
		return nil, fmt.Errorf("patient not found")
	}

	if req.FirstName != nil {
		patient.FirstName = req.GetFirstName()
	}
	if req.LastName != nil {
		patient.LastName = req.GetLastName()
	}
	if req.DateOfBirth != nil {
		patient.DateOfBirth = req.GetDateOfBirth()
	}
	if req.Gender != nil {
		patient.Gender = int32(req.GetGender())
	}
	if req.BloodGroup != nil {
		patient.BloodGroup = int32(req.GetBloodGroup())
	}
	if req.PhoneNumber != nil {
		patient.PhoneNumber = req.GetPhoneNumber()
	}
	if req.Email != nil {
		patient.Email = req.GetEmail()
	}

	if req.Address != nil {
		addr := &entity.Address{
			Street:  req.Address.Street,
			City:    req.Address.City,
			State:   req.Address.State,
			ZipCode: req.Address.ZipCode,
			Country: req.Address.Country,
		}
		patient.Address = entity.MarshalAddress(addr)
	}

	if err := h.repo.Update(ctx, patient); err != nil {
		h.log.WithContext(ctx).Errorf("Failed to update patient: %v", err)
		return nil, fmt.Errorf("failed to update patient: %w", err)
	}

	return h.entityToProto(patient), nil
}

func (h *PatientHandler) SearchPatients(ctx context.Context, req *requestpb.SearchPatientsRequest) ([]*responsepb.PatientResponse, error) {
	ctx, span := otel.Trace(ctx, "PatientHandler.SearchPatients")
	defer span.End()

	filters := make(map[string]interface{})

	if req.GetName() != "" {
		filters["name"] = req.GetName()
	}
	if req.GetEmail() != "" {
		filters["email"] = req.GetEmail()
	}
	if req.GetPhoneNumber() != "" {
		filters["phone_number"] = req.GetPhoneNumber()
	}
	if req.GetPatientId() != "" {
		filters["patient_id"] = req.GetPatientId()
	}

	patients, err := h.repo.Search(ctx, filters)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to search patients: %v", err)
		return nil, fmt.Errorf("failed to search patients: %w", err)
	}

	var responses []*responsepb.PatientResponse
	for _, p := range patients {
		responses = append(responses, h.entityToProto(p))
	}

	return responses, nil
}

func (h *PatientHandler) GetMedicalHistory(ctx context.Context, id string, fromDate, toDate string) (*responsepb.MedicalHistoryResponse, error) {
	ctx, span := otel.Trace(ctx, "PatientHandler.GetMedicalHistory")
	defer span.End()

	patient, err := h.repo.Get(ctx, id)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get patient: %v", err)
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		h.log.WithContext(ctx).Errorf("Patient not found: %s", id)
		return nil, fmt.Errorf("patient not found")
	}

	filters := make(map[string]interface{})
	if fromDate != "" {
		filters["from_date"] = fromDate
	}
	if toDate != "" {
		filters["to_date"] = toDate
	}

	records, err := h.recordRepo.GetByPatientID(ctx, id, filters)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to fetch medical records: %v", err)
		return nil, fmt.Errorf("failed to fetch medical records: %w", err)
	}

	var protoRecords []*responsepb.MedicalRecord
	for _, record := range records {
		protoRecords = append(protoRecords, h.medicalRecordToProto(record))
	}

	return &responsepb.MedicalHistoryResponse{
		PatientId: patient.ID,
		Records:   protoRecords,
	}, nil
}

func (h *PatientHandler) entityToProto(patient *entity.Patient) *responsepb.PatientResponse {
	response := &responsepb.PatientResponse{
		PatientId:   patient.ID,
		FirstName:   patient.FirstName,
		LastName:    patient.LastName,
		Email:       patient.Email,
		PhoneNumber: patient.PhoneNumber,
		DateOfBirth: patient.DateOfBirth,
		Gender:      commonpb.Gender(patient.Gender),
		BloodGroup:  commonpb.BloodGroup(patient.BloodGroup),
	}

	if patient.Address != "" {
		addr := entity.UnmarshalAddress(patient.Address)
		if addr != nil {
			response.Address = &commonpb.Address{
				Street:  addr.Street,
				City:    addr.City,
				State:   addr.State,
				ZipCode: addr.ZipCode,
				Country: addr.Country,
			}
		}
	}

	return response
}

func (h *PatientHandler) medicalRecordToProto(record *entity.MedicalRecord) *responsepb.MedicalRecord {
	return &responsepb.MedicalRecord{
		RecordId:  record.ID,
		DoctorId:  record.DoctorID,
		VisitDate: record.VisitDate.Format("2006-01-02"),
		Diagnosis: record.Diagnosis,
		Notes:     record.Notes,
	}
}
