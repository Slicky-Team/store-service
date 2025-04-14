package grpc

import (
	"context"
	"errors"
	"store-service/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Server struct {
	proto.UnimplementedBarberServiceServer
	DB      *gorm.DB
	StoreId uuid.UUID
}

func NewServer(db *gorm.DB, storeId uuid.UUID) *Server {
	return &Server{
		DB:      db,
		StoreId: storeId,
	}
}

func (s *Server) CheckAvailability(ctx context.Context, req *proto.BarberAvailabilityRequest) (*proto.BarberAvailabilityResponse, error) {
	barberID, err := uuid.Parse(req.BarberId)
	if err != nil {
		return &proto.BarberAvailabilityResponse{Available: false, Error: "Invalid barber ID format"}, nil
	}

	_, err = s.findBarber(barberID, s.StoreId)
	if err != nil {
		return &proto.BarberAvailabilityResponse{Available: false, Error: "Barber not found"}, nil
	}

	dateTime, err := time.Parse("2006-01-02 15:04", req.Date+" "+req.Time)
	if err != nil {
		return &proto.BarberAvailabilityResponse{Available: false, Error: "Invalid date or time format"}, nil
	}

	// Check if there is any appointment in the database for the barber at the given time
	var appointment models.StoreAppointment
	result := s.DB.Where("user_store_id = ? AND date_time = ?", barberID, dateTime).First(&appointment)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &proto.BarberAvailabilityResponse{Available: false, Error: "Error checking availability"}, result.Error
	}

	if result.RowsAffected > 0 {
		return &proto.BarberAvailabilityResponse{Available: false, Error: "Barber is not available at this time"}, nil
	}

	return &proto.BarberAvailabilityResponse{Available: true, Error: ""}, nil
}

func (s *Server) GetAvailableSlots(ctx context.Context, req *proto.AvailableSlotsRequest) (*proto.AvailableSlotsResponse, error) {
	barberID, err := uuid.Parse(req.BarberId)
	if err != nil {
		return &proto.AvailableSlotsResponse{Slots: nil, Error: "Invalid barber ID format"}, nil
	}

	_, err = s.findBarber(barberID, s.StoreId)
	if err != nil {
		return &proto.AvailableSlotsResponse{Slots: nil, Error: "Barber not found"}, nil
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return &proto.AvailableSlotsResponse{Slots: nil, Error: "Invalid date format"}, nil
	}

	var availableSlots []string
	for hour := 9; hour <= 17; hour++ {
		for minute := 0; minute < 60; minute += 30 {
			dateTime := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, time.UTC)
			var appointment models.Appointment
			result := s.DB.Joins("UserStore").Where("user_store_id = ? AND date_time = ?", barberID, dateTime).First(&appointment)
			if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return &proto.AvailableSlotsResponse{Slots: nil, Error: "Error checking available slots"}, result.Error
			}
			if result.RowsAffected == 0 {
				availableSlots = append(availableSlots, dateTime.Format("15:04"))
			}
		}
	}

	return &proto.AvailableSlotsResponse{Slots: availableSlots, Error: ""}, nil
}

func (s *Server) BookAppointment(ctx context.Context, req *proto.BookAppointmentRequest) (*proto.BookAppointmentResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &proto.BookAppointmentResponse{Success: false, Error: "Invalid user ID format"}, nil
	}

	barberID, err := uuid.Parse(req.BarberId)
	if err != nil {
		return &proto.BookAppointmentResponse{Success: false, Error: "Invalid barber ID format"}, nil
	}

	_, err = s.findUser(userID)
	if err != nil {
		return &proto.BookAppointmentResponse{Success: false, Error: "User not found"}, nil
	}

	barber, err := s.findBarber(barberID, s.StoreId)
	if err != nil {
		return &proto.BookAppointmentResponse{Success: false, Error: "Barber not found"}, nil
	}

	dateTime, err := time.Parse("2006-01-02 15:04", req.Date+" "+req.Time)
	if err != nil {
		return &proto.BookAppointmentResponse{Success: false, Error: "Invalid date or time format"}, nil
	}

	var appointment models.StoreAppointment
	result := s.DB.Where("user_store_id = ? AND date_time = ?", barber.ID, dateTime).First(&appointment)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &proto.BookAppointmentResponse{Success: false, Error: "Error checking availability"}, result.Error
	}

	if result.RowsAffected > 0 {
		return &proto.BookAppointmentResponse{Success: false, Error: "Barber is not available at this time"}, nil
	}

	newAppointment := models.StoreAppointment{
		UserAccountID: userID,
		UserStoreID:   barber.ID,
		DateTime:      dateTime,
		StoreID:       s.StoreId,
		Status:        models.AppointmentStatusScheduled,
	}

	result = s.DB.Create(&newAppointment)
	if result.Error != nil {
		return &proto.BookAppointmentResponse{Success: false, Error: "Error booking appointment"}, result.Error
	}

	return &proto.BookAppointmentResponse{Success: true, Error: ""}, nil
}

func (s *Server) findBarber(barberID uuid.UUID, storeId uuid.UUID) (*models.UserStore, error) {
	var barber models.UserStore
	result := s.DB.Where("user_id = ? AND role = ? AND store_id = ?", barberID, models.Staff, storeId).First(&barber)
	if result.Error != nil {
		return nil, result.Error
	}
	return &barber, nil
}

func (s *Server) findUser(userId uuid.UUID) (*models.UserProfile, error) {
	var user models.UserProfile
	result := s.DB.Where("account_id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
