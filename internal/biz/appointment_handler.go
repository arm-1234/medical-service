package biz

import (
	"context"
	"fmt"
	"time"

	commonpb "github.com/arm-1234/common-protos/medical/v1/common"
	requestpb "github.com/arm-1234/common-protos/medical/v1/request"
	responsepb "github.com/arm-1234/common-protos/medical/v1/response"
	"github.com/arm-1234/medical-service/internal/data"
	"github.com/arm-1234/medical-service/internal/data/entity"
	"github.com/arm-1234/medical-service/internal/pkg/otel"
	"github.com/go-kratos/kratos/v2/log"
)

type AppointmentHandler struct {
	repo        data.AppointmentRepo
	patientRepo data.PatientRepo
	doctorRepo  data.DoctorRepo
	log         *log.Helper
}

func NewAppointmentHandler(
	repo data.AppointmentRepo,
	patientRepo data.PatientRepo,
	doctorRepo data.DoctorRepo,
	logger log.Logger,
) *AppointmentHandler {
	return &AppointmentHandler{
		repo:        repo,
		patientRepo: patientRepo,
		doctorRepo:  doctorRepo,
		log:         log.NewHelper(logger),
	}
}

func (h *AppointmentHandler) BookAppointment(ctx context.Context, req *requestpb.BookAppointmentRequest) (*responsepb.AppointmentResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentHandler.BookAppointment")
	defer span.End()

	if req.PatientId == "" || req.DoctorId == "" {
		h.log.WithContext(ctx).Errorf("Patient ID and doctor ID are required")
		return nil, fmt.Errorf("patient_id and doctor_id are required")
	}
	if req.AppointmentDate == "" || req.AppointmentTime == "" {
		h.log.WithContext(ctx).Errorf("Appointment date and time are required")
		return nil, fmt.Errorf("appointment_date and appointment_time are required")
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
	if !doctor.IsAvailable {
		h.log.WithContext(ctx).Errorf("Doctor is not available: %s", req.DoctorId)
		return nil, fmt.Errorf("doctor is not available")
	}

	conflict, err := h.repo.CheckConflict(ctx, req.DoctorId, req.AppointmentDate, req.AppointmentTime, "")
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to check appointment conflict: %v", err)
		return nil, fmt.Errorf("failed to check appointment conflict: %w", err)
	}
	if conflict != nil {
		h.log.WithContext(ctx).Errorf("Time slot already booked for doctor %s at %s %s", req.DoctorId, req.AppointmentDate, req.AppointmentTime)
		return nil, fmt.Errorf("time slot is already booked")
	}

	appointment := &entity.Appointment{
		PatientID:        req.PatientId,
		PatientName:      patient.FirstName + " " + patient.LastName,
		DoctorID:         req.DoctorId,
		DoctorName:       doctor.FirstName + " " + doctor.LastName,
		AppointmentDate:  req.AppointmentDate,
		AppointmentTime:  req.AppointmentTime,
		Status:           entity.AppointmentStatusScheduled,
		ConsultationType: int32(req.ConsultationType),
		ReasonForVisit:   req.ReasonForVisit,
		Notes:            req.Notes,
	}

	if err := h.repo.Create(ctx, appointment); err != nil {
		h.log.WithContext(ctx).Errorf("Failed to create appointment: %v", err)
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return h.entityToProto(appointment), nil
}

func (h *AppointmentHandler) GetAppointment(ctx context.Context, id string) (*responsepb.AppointmentResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentHandler.GetAppointment")
	defer span.End()

	if id == "" {
		h.log.WithContext(ctx).Errorf("Appointment ID is required")
		return nil, fmt.Errorf("appointment_id is required")
	}

	appointment, err := h.repo.Get(ctx, id)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get appointment: %v", err)
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	if appointment == nil {
		return nil, fmt.Errorf("appointment not found")
	}

	return h.entityToProto(appointment), nil
}

func (h *AppointmentHandler) CancelAppointment(ctx context.Context, id string, reason string) (*responsepb.AppointmentResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentHandler.CancelAppointment")
	defer span.End()

	if id == "" {
		h.log.WithContext(ctx).Errorf("Appointment ID is required")
		return nil, fmt.Errorf("appointment_id is required")
	}

	appointment, err := h.repo.Get(ctx, id)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get appointment: %v", err)
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}
	if appointment == nil {
		h.log.WithContext(ctx).Errorf("Appointment not found: %s", id)
		return nil, fmt.Errorf("appointment not found")
	}

	if appointment.Status == entity.AppointmentStatusCancelled {
		h.log.WithContext(ctx).Errorf("Appointment already cancelled: %s", id)
		return nil, fmt.Errorf("appointment is already cancelled")
	}
	if appointment.Status == entity.AppointmentStatusCompleted {
		h.log.WithContext(ctx).Errorf("Cannot cancel completed appointment: %s", id)
		return nil, fmt.Errorf("cannot cancel a completed appointment")
	}

	now := time.Now()
	appointment.Status = entity.AppointmentStatusCancelled
	appointment.CancelledAt = &now
	appointment.CancellationReason = reason

	if err := h.repo.Update(ctx, appointment); err != nil {
		h.log.WithContext(ctx).Errorf("Failed to cancel appointment: %v", err)
		return nil, fmt.Errorf("failed to cancel appointment: %w", err)
	}

	return h.entityToProto(appointment), nil
}

func (h *AppointmentHandler) RescheduleAppointment(ctx context.Context, req *requestpb.RescheduleAppointmentRequest) (*responsepb.AppointmentResponse, error) {
	ctx, span := otel.Trace(ctx, "AppointmentHandler.RescheduleAppointment")
	defer span.End()

	if req.AppointmentId == "" {
		h.log.WithContext(ctx).Errorf("Appointment ID is required")
		return nil, fmt.Errorf("appointment_id is required")
	}
	if req.NewAppointmentDate == "" || req.NewAppointmentTime == "" {
		h.log.WithContext(ctx).Errorf("New appointment date and time are required")
		return nil, fmt.Errorf("new_appointment_date and new_appointment_time are required")
	}

	appointment, err := h.repo.Get(ctx, req.AppointmentId)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to get appointment: %v", err)
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}
	if appointment == nil {
		h.log.WithContext(ctx).Errorf("Appointment not found: %s", req.AppointmentId)
		return nil, fmt.Errorf("appointment not found")
	}

	if appointment.Status == entity.AppointmentStatusCancelled {
		h.log.WithContext(ctx).Errorf("Cannot reschedule cancelled appointment: %s", req.AppointmentId)
		return nil, fmt.Errorf("cannot reschedule a cancelled appointment")
	}
	if appointment.Status == entity.AppointmentStatusCompleted {
		h.log.WithContext(ctx).Errorf("Cannot reschedule completed appointment: %s", req.AppointmentId)
		return nil, fmt.Errorf("cannot reschedule a completed appointment")
	}

	conflict, err := h.repo.CheckConflict(ctx, appointment.DoctorID, req.NewAppointmentDate, req.NewAppointmentTime, appointment.ID)
	if err != nil {
		h.log.WithContext(ctx).Errorf("Failed to check appointment conflict: %v", err)
		return nil, fmt.Errorf("failed to check appointment conflict: %w", err)
	}
	if conflict != nil {
		h.log.WithContext(ctx).Errorf("New time slot already booked: %s %s", req.NewAppointmentDate, req.NewAppointmentTime)
		return nil, fmt.Errorf("new time slot is already booked")
	}

	appointment.AppointmentDate = req.NewAppointmentDate
	appointment.AppointmentTime = req.NewAppointmentTime
	appointment.Status = entity.AppointmentStatusRescheduled
	if req.Reason != "" {
		appointment.Notes = appointment.Notes + "\nRescheduled: " + req.Reason
	}

	if err := h.repo.Update(ctx, appointment); err != nil {
		h.log.WithContext(ctx).Errorf("Failed to reschedule appointment: %v", err)
		return nil, fmt.Errorf("failed to reschedule appointment: %w", err)
	}

	return h.entityToProto(appointment), nil
}

func (h *AppointmentHandler) CompleteAppointment(ctx context.Context, req *requestpb.CompleteAppointmentRequest) (*responsepb.AppointmentResponse, error) {
	if req.AppointmentId == "" {
		return nil, fmt.Errorf("appointment_id is required")
	}

	appointment, err := h.repo.Get(ctx, req.AppointmentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}
	if appointment == nil {
		return nil, fmt.Errorf("appointment not found")
	}

	if appointment.Status == entity.AppointmentStatusCancelled {
		return nil, fmt.Errorf("cannot complete a cancelled appointment")
	}
	if appointment.Status == entity.AppointmentStatusCompleted {
		return nil, fmt.Errorf("appointment is already completed")
	}

	appointment.Status = entity.AppointmentStatusCompleted
	appointment.Diagnosis = req.Diagnosis
	if req.Notes != "" {
		appointment.Notes = req.Notes
	}

	if err := h.repo.Update(ctx, appointment); err != nil {
		h.log.Errorf("failed to complete appointment: %v", err)
		return nil, fmt.Errorf("failed to complete appointment: %w", err)
	}

	doctor, _ := h.doctorRepo.Get(ctx, appointment.DoctorID)
	if doctor != nil {
		doctor.TotalConsultations++
		_ = h.doctorRepo.Update(ctx, doctor)
	}

	return h.entityToProto(appointment), nil
}

func (h *AppointmentHandler) GetAvailableSlots(ctx context.Context, req *requestpb.GetAvailableSlotsRequest) (*responsepb.AvailableSlotsResponse, error) {
	if req.DoctorId == "" || req.Date == "" {
		return nil, fmt.Errorf("doctor_id and date are required")
	}

	doctor, err := h.doctorRepo.Get(ctx, req.DoctorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}
	if doctor == nil {
		return nil, fmt.Errorf("doctor not found")
	}

	existingAppointments, err := h.repo.GetByDoctorAndDate(ctx, req.DoctorId, req.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing appointments: %w", err)
	}

	bookedTimes := make(map[string]bool)
	for _, apt := range existingAppointments {
		bookedTimes[apt.AppointmentTime] = true
	}

	var slots []*responsepb.TimeSlot
	startHour := 9
	endHour := 17
	for hour := startHour; hour < endHour; hour++ {
		for minute := 0; minute < 60; minute += 30 {
			timeStr := fmt.Sprintf("%02d:%02d", hour, minute)
			endMinute := minute + 30
			endHourSlot := hour
			if endMinute >= 60 {
				endMinute = 0
				endHourSlot++
			}
			endTimeStr := fmt.Sprintf("%02d:%02d", endHourSlot, endMinute)

			slots = append(slots, &responsepb.TimeSlot{
				StartTime:   timeStr,
				EndTime:     endTimeStr,
				IsAvailable: !bookedTimes[timeStr],
			})
		}
	}

	return &responsepb.AvailableSlotsResponse{
		DoctorId:       req.DoctorId,
		DoctorName:     doctor.FirstName + " " + doctor.LastName,
		Date:           req.Date,
		AvailableSlots: slots,
	}, nil
}

func (h *AppointmentHandler) GetPatientAppointments(ctx context.Context, req *requestpb.GetPatientAppointmentsRequest) (*responsepb.PatientAppointmentsResponse, error) {
	if req.PatientId == "" {
		return nil, fmt.Errorf("patient_id is required")
	}

	filters := make(map[string]interface{})
	if req.Status != nil {
		filters["status"] = int32(req.GetStatus())
	}
	if req.FromDate != nil {
		filters["from_date"] = req.GetFromDate()
	}
	if req.ToDate != nil {
		filters["to_date"] = req.GetToDate()
	}

	appointments, err := h.repo.GetByPatientID(ctx, req.PatientId, filters)
	if err != nil {
		h.log.Errorf("failed to get patient appointments: %v", err)
		return nil, fmt.Errorf("failed to get patient appointments: %w", err)
	}

	var protoAppointments []*responsepb.AppointmentResponse
	for _, apt := range appointments {
		protoAppointments = append(protoAppointments, h.entityToProto(apt))
	}

	return &responsepb.PatientAppointmentsResponse{
		Appointments: protoAppointments,
	}, nil
}

func (h *AppointmentHandler) GetDoctorAppointments(ctx context.Context, req *requestpb.GetDoctorAppointmentsRequest) (*responsepb.DoctorAppointmentsResponse, error) {
	if req.DoctorId == "" {
		return nil, fmt.Errorf("doctor_id is required")
	}

	filters := make(map[string]interface{})
	if req.Status != nil {
		filters["status"] = int32(req.GetStatus())
	}
	if req.Date != nil {
		filters["date"] = req.GetDate()
	}

	appointments, err := h.repo.GetByDoctorID(ctx, req.DoctorId, filters)
	if err != nil {
		h.log.Errorf("failed to get doctor appointments: %v", err)
		return nil, fmt.Errorf("failed to get doctor appointments: %w", err)
	}

	var protoAppointments []*responsepb.AppointmentResponse
	for _, apt := range appointments {
		protoAppointments = append(protoAppointments, h.entityToProto(apt))
	}

	return &responsepb.DoctorAppointmentsResponse{
		Appointments: protoAppointments,
	}, nil
}

func (h *AppointmentHandler) entityToProto(appointment *entity.Appointment) *responsepb.AppointmentResponse {
	resp := &responsepb.AppointmentResponse{
		AppointmentId:    appointment.ID,
		PatientId:        appointment.PatientID,
		PatientName:      appointment.PatientName,
		DoctorId:         appointment.DoctorID,
		DoctorName:       appointment.DoctorName,
		AppointmentDate:  appointment.AppointmentDate,
		AppointmentTime:  appointment.AppointmentTime,
		Status:           commonpb.AppointmentStatus(appointment.Status),
		ConsultationType: commonpb.ConsultationType(appointment.ConsultationType),
		ReasonForVisit:   appointment.ReasonForVisit,
		Notes:            appointment.Notes,
		Diagnosis:        appointment.Diagnosis,
		CreatedAt:        appointment.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        appointment.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if appointment.CancelledAt != nil {
		cancelledAt := appointment.CancelledAt.Format("2006-01-02T15:04:05Z")
		resp.CancelledAt = &cancelledAt
	}
	if appointment.CancellationReason != "" {
		resp.CancellationReason = &appointment.CancellationReason
	}

	return resp
}
