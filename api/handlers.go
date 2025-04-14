package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"store-service/config"
	"store-service/models"
	"store-service/proto"

	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Barber struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func checkAvailability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	barberID := r.URL.Query().Get("barberId")
	date := r.URL.Query().Get("date")
	timeStr := r.URL.Query().Get("time")
	cfg := config.LoadConfig()

	// Use the gRPC server address from the configuration
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := proto.NewBarberServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &proto.BarberAvailabilityRequest{
		BarberId: barberID,
		Date:     date,
		Time:     timeStr,
	}

	resp, err := client.CheckAvailability(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("gRPC call failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Error != "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"available": resp.Available,
		"error":     resp.Error,
	})
}

func getAvailableSlots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	barberID := r.URL.Query().Get("barberId")
	date := r.URL.Query().Get("date")
	cfg := config.LoadConfig()

	// Use the gRPC server address from the configuration
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := proto.NewBarberServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &proto.AvailableSlotsRequest{
		BarberId: barberID,
		Date:     date,
	}

	resp, err := client.GetAvailableSlots(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("gRPC call failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Error != "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"slots": resp.Slots,
		"error": resp.Error,
	})
}

func bookAppointment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		UserID   string `json:"userId"`
		BarberID string `json:"barberId"`
		Date     string `json:"date"`
		Time     string `json:"time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	cfg := config.LoadConfig()

	// Use the gRPC server address from the configuration
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := proto.NewBarberServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &proto.BookAppointmentRequest{
		UserId:   body.UserID,
		BarberId: body.BarberID,
		Date:     body.Date,
		Time:     body.Time,
	}

	resp, err := client.BookAppointment(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("gRPC call failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Error != "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": resp.Success,
		"error":   resp.Error,
	})
}

func listBarbers(w http.ResponseWriter, r *http.Request, db models.DBEngine) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userStores []models.UserStore
	if err := db.GormDB.Where("role = ?", models.Staff).Find(&userStores).Error; err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch barbers: %v", err), http.StatusInternalServerError)
		return
	}

	var barbers []Barber
	for _, userStore := range userStores {
		var userProfile models.UserProfile
		if err := db.GormDB.First(&userProfile, "account_id = ?", userStore.UserID).Error; err != nil {
			log.Printf("Error fetching user profile: %v", err)
			continue
		}
		barbers = append(barbers, Barber{
			ID:   userProfile.ID.String(),
			Name: userProfile.FullName,
		})
	}

	/*
		barbers := []Barber{
			{ID: uuid.New().String(), Name: "Barber 1"},
			{ID: uuid.New().String(), Name: "Barber 2"},
			{ID: uuid.New().String(), Name: "Barber 3"},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(barbers)*/
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(barbers)
}

func listServices(w http.ResponseWriter, r *http.Request, db models.DBEngine) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var storeServices []models.StoreService
	if err := db.GormDB.Find(&storeServices).Error; err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch services: %v", err), http.StatusInternalServerError)
		return
	}

	var services []Service
	for _, storeService := range storeServices {
		services = append(services, Service{
			ID:    storeService.ID.String(),
			Name:  storeService.ServiceName,
			Price: storeService.ServicePrice,
		})
	}

	services := []Service{
		{ID: uuid.New().String(), Name: "Haircut", Price: 25.0},
		{ID: uuid.New().String(), Name: "Beard Trim", Price: 15.0},
		{ID: uuid.New().String(), Name: "Shave", Price: 20.0},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}

func SetupRoutes(mux *http.ServeMux, db models.DBEngine) {
	mux.HandleFunc("/availability", checkAvailability)
	mux.HandleFunc("/slots", getAvailableSlots)
	mux.HandleFunc("/appointment", bookAppointment)
	mux.HandleFunc("/barbers", func(w http.ResponseWriter, r *http.Request) { listBarbers(w, r, db) })
	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) { listServices(w, r, db) })
	log.Println("Rest routes created")
}
