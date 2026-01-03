package biz

import (
	"context"
	"fmt"
	"time"

	requestpb "github.com/arm-1234/common-protos/medical/v1/request"
	responsepb "github.com/arm-1234/common-protos/medical/v1/response"
	"github.com/arm-1234/medical-service/internal/data"
	"github.com/arm-1234/medical-service/internal/data/entity"
	"github.com/arm-1234/medical-service/internal/pkg/otel"
	"github.com/go-kratos/kratos/v2/log"
)

type PrescriptionHandler struct {
	repo        data.PrescriptionRepo
	patientRepo data.PatientRepo
	doctorRepo  data.DoctorRepo
	log         *log.Helper
}

func NewPrescriptionHandler(
	repo data.PrescriptionRepo,
	patientRepo data.PatientRepo,
	doctorRepo data.DoctorRepo,
	logger log.Logger,
) *PrescriptionHandler {
	return &PrescriptionHandler{
		repo:        repo,
		patientRepo: patientRepo,
		doctorRepo:  doctorRepo,
		log:         log.NewHelper(logger),
	}
}

func (h *PrescriptionHandler) CreatePrescription(ctx context.Context, req *requestpb.CreatePrescriptionRequest) (*responsepb.PrescriptionResponse, error) {
	ctx, span := otel.Trace(ctx, "PrescriptionHandler.CreatePrescription")
	defer span.End()

	if req.PatientId == "" || req.DoctorId == "" {
		h.log.WithContext(ctx).Errorf("Patient ID and doctor ID are required")
		return nil, fmt.Errorf("patient_id and doctor_id are required")
	}
	if len(req.Medications) == 0 {
		h.log.WithContext(ctx).Errorf("At least one medication is required")
		return nil, fmt.Errorf("at least one medication is required")
	}

	patient, err := h.patientRepo.Get(ctx, req.PatientId)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get patient: %v", err)
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		h.log.WithContext(ctx).Errorf("Patient not found: %s", req.PatientId)
		return nil, fmt.Errorf("patient not found")
	}

	doctor, err := h.doctorRepo.Get(ctx, req.DoctorId)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get doctor: %v", err)
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}
	if doctor == nil {
		h.log.WithContext(ctx).Errorf("Doctor not found: %s", req.DoctorId)
		return nil, fmt.Errorf("doctor not found")
	}

	var medications []*entity.Medication
	for _, med := range req.Medications {
		medications = append(medications, &entity.Medication{
			MedicationName: med.MedicationName,
			Dosage:         med.Dosage,
			Frequency:      med.Frequency,
			Duration:       med.Duration,
			Route:          med.Route,
			Instructions:   med.Instructions,
			Quantity:       med.Quantity,
		})
	}

	validityDays := req.ValidityDays
	if validityDays <= 0 {
		validityDays = 30
	}
	now := time.Now()
	validUntil := now.AddDate(0, 0, int(validityDays))

	prescription := &entity.Prescription{
		AppointmentID:          req.AppointmentId,
		PatientID:              req.PatientId,
		PatientName:            patient.FirstName + " " + patient.LastName,
		DoctorID:               req.DoctorId,
		DoctorName:             doctor.FirstName + " " + doctor.LastName,
		Medications:            entity.MarshalMedications(medications),
		Diagnosis:              req.Diagnosis,
		AdditionalInstructions: req.AdditionalInstructions,
		PrescriptionDate:       now,
		ValidUntil:             validUntil,
		IsActive:               true,
	}

	if err := h.repo.Create(ctx, prescription); err != nil {
		h.log.WithContext(ctx).Errorf("Failed to create prescription: %v", err)
		return nil, fmt.Errorf("failed to create prescription: %w", err)
	}

	return h.entityToProto(prescription), nil
}

func (h *PrescriptionHandler) GetPrescription(ctx context.Context, id string) (*responsepb.PrescriptionResponse, error) {
	ctx, span := otel.Trace(ctx, "PrescriptionHandler.GetPrescription")
	defer span.End()

	if id == "" {
		h.log.WithContext(ctx).Errorf("Prescription ID is required")
		return nil, fmt.Errorf("prescription_id is required")
	}

	prescription, err := h.repo.Get(ctx, id)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get prescription: %v", err)
		return nil, fmt.Errorf("failed to get prescription: %w", err)
	}

	if prescription == nil {
		return nil, fmt.Errorf("prescription not found")
	}

	if prescription.ValidUntil.Before(time.Now()) {
		prescription.IsActive = false
	}

	return h.entityToProto(prescription), nil
}

func (h *PrescriptionHandler) GetPatientPrescriptions(ctx context.Context, req *requestpb.GetPatientPrescriptionsRequest) (*responsepb.PatientPrescriptionsResponse, error) {
	ctx, span := otel.Trace(ctx, "PrescriptionHandler.GetPatientPrescriptions")
	defer span.End()

	if req.PatientId == "" {
		h.log.WithContext(ctx).Errorf("Patient ID is required")
		return nil, fmt.Errorf("patient_id is required")
	}

	filters := make(map[string]interface{})
	if req.FromDate != nil {
		filters["from_date"] = req.GetFromDate()
	}
	if req.ToDate != nil {
		filters["to_date"] = req.GetToDate()
	}

	prescriptions, err := h.repo.GetByPatientID(ctx, req.PatientId, filters)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get patient prescriptions: %v", err)
		return nil, fmt.Errorf("failed to get patient prescriptions: %w", err)
	}

	var protoPrescriptions []*responsepb.PrescriptionResponse
	for _, p := range prescriptions {
		if p.ValidUntil.Before(time.Now()) {
			p.IsActive = false
		}
		protoPrescriptions = append(protoPrescriptions, h.entityToProto(p))
	}

	return &responsepb.PatientPrescriptionsResponse{
		Prescriptions: protoPrescriptions,
	}, nil
}

func (h *PrescriptionHandler) GetDoctorPrescriptions(ctx context.Context, req *requestpb.GetDoctorPrescriptionsRequest) (*responsepb.DoctorPrescriptionsResponse, error) {
	ctx, span := otel.Trace(ctx, "PrescriptionHandler.GetDoctorPrescriptions")
	defer span.End()

	if req.DoctorId == "" {
		h.log.WithContext(ctx).Errorf("Doctor ID is required")
		return nil, fmt.Errorf("doctor_id is required")
	}

	filters := make(map[string]interface{})
	if req.FromDate != nil {
		filters["from_date"] = req.GetFromDate()
	}
	if req.ToDate != nil {
		filters["to_date"] = req.GetToDate()
	}

	prescriptions, err := h.repo.GetByDoctorID(ctx, req.DoctorId, filters)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get doctor prescriptions: %v", err)
		return nil, fmt.Errorf("failed to get doctor prescriptions: %w", err)
	}

	var protoPrescriptions []*responsepb.PrescriptionResponse
	for _, p := range prescriptions {
		if p.ValidUntil.Before(time.Now()) {
			p.IsActive = false
		}
		protoPrescriptions = append(protoPrescriptions, h.entityToProto(p))
	}

	return &responsepb.DoctorPrescriptionsResponse{
		Prescriptions: protoPrescriptions,
	}, nil
}

func (h *PrescriptionHandler) entityToProto(prescription *entity.Prescription) *responsepb.PrescriptionResponse {
	entityMeds := entity.UnmarshalMedications(prescription.Medications)
	var protoMeds []*requestpb.Medication
	for _, med := range entityMeds {
		protoMeds = append(protoMeds, &requestpb.Medication{
			MedicationName: med.MedicationName,
			Dosage:         med.Dosage,
			Frequency:      med.Frequency,
			Duration:       med.Duration,
			Route:          med.Route,
			Instructions:   med.Instructions,
			Quantity:       med.Quantity,
		})
	}

	return &responsepb.PrescriptionResponse{
		PrescriptionId:         prescription.ID,
		AppointmentId:          prescription.AppointmentID,
		PatientId:              prescription.PatientID,
		PatientName:            prescription.PatientName,
		DoctorId:               prescription.DoctorID,
		DoctorName:             prescription.DoctorName,
		Medications:            protoMeds,
		Diagnosis:              prescription.Diagnosis,
		AdditionalInstructions: prescription.AdditionalInstructions,
		PrescriptionDate:       prescription.PrescriptionDate.Format("2006-01-02"),
		ValidUntil:             prescription.ValidUntil.Format("2006-01-02"),
		IsActive:               prescription.IsActive,
		CreatedAt:              prescription.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
