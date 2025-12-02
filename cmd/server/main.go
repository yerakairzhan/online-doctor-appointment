package main

import (
	"log"
	"net/http"
	"os"

	"online-doctor-appointment/internal/database"
	"online-doctor-appointment/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Create router
	router := mux.NewRouter()

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Public routes
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	router.HandleFunc("/login", handlers.LoginPageHandler).Methods("GET")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/register", handlers.RegisterPageHandler).Methods("GET")
	router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")

	// API routes for getting data
	router.HandleFunc("/api/doctors", handlers.GetDoctorsHandler).Methods("GET")
	router.HandleFunc("/api/doctors/{specialty}", handlers.GetDoctorsBySpecialtyHandler).Methods("GET")

	router.HandleFunc("/payment/success", handlers.PaymentSuccessHandler).Methods("GET")
	router.HandleFunc("/payment/failure", handlers.PaymentFailureHandler).Methods("GET")

	// Protected routes (require authentication)
	protected := router.PathPrefix("/dashboard").Subrouter()
	protected.Use(handlers.AuthMiddleware)
	protected.HandleFunc("/payment/success", handlers.PaymentSuccessHandler).Methods("GET")
	protected.HandleFunc("/payment/failure", handlers.PaymentFailureHandler).Methods("GET")

	// Patient routes
	protected.HandleFunc("/patient", handlers.PatientDashboardHandler).Methods("GET")
	protected.HandleFunc("/patient/book", handlers.BookAppointmentPageHandler).Methods("GET")
	protected.HandleFunc("/patient/book", handlers.BookAppointmentHandler).Methods("POST")
	protected.HandleFunc("/patient/appointments", handlers.PatientAppointmentsHandler).Methods("GET")

	// Doctor routes
	protected.HandleFunc("/doctor", handlers.DoctorDashboardHandler).Methods("GET")
	protected.HandleFunc("/doctor/appointments", handlers.DoctorAppointmentsHandler).Methods("GET")
	protected.HandleFunc("/doctor/appointment/{id}/update", handlers.UpdateAppointmentStatusHandler).Methods("POST")

	// Admin routes
	protected.HandleFunc("/admin", handlers.AdminDashboardHandler).Methods("GET")
	protected.HandleFunc("/admin/doctors", handlers.AdminDoctorsHandler).Methods("GET")
	protected.HandleFunc("/admin/patients", handlers.AdminPatientsHandler).Methods("GET")

	// Chatbot routes (can be accessed by all authenticated users)
	protected.HandleFunc("/chatbot", handlers.ChatbotPageHandler).Methods("GET")
	// Chatbot API endpoint
	router.HandleFunc("/api/chatbot", handlers.ChatbotAPIHandler).Methods("POST")

	// API routes for available time slots
	router.HandleFunc("/api/available-slots/{doctorId}/{date}", handlers.GetAvailableSlotsHandler).Methods("GET")

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Visit http://localhost:%s to access the application", port)

	// Start server
	log.Fatal(http.ListenAndServe(":"+port, router))
}
