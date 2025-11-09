package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"online-doctor-appointment/internal/database"
	"online-doctor-appointment/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// Session management (simple in-memory store for MVP)
var sessions = make(map[string]*SessionData)

type SessionData struct {
	UserID   int
	UserType string
	Email    string
	ExpireAt time.Time
}

// Generate simple session token (in production, use proper JWT or secure sessions)
func generateSessionToken() string {
	return generateRandomString(32)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// HomeHandler serves the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <header class="header">
            <h1>🏥 Online Doctor Appointment</h1>
            <p>Book appointments with qualified doctors online</p>
        </header>
        
        <div class="home-content">
            <div class="hero-section">
                <h2>Welcome to Our Medical Platform</h2>
                <p>Easy, fast, and secure way to book appointments with healthcare professionals</p>
                
                <div class="action-buttons">
                    <a href="/login" class="btn btn-primary">Login</a>
                    <a href="/register" class="btn btn-secondary">Register</a>
                </div>
            </div>
            
            <div class="features">
                <div class="feature-card">
                    <h3>👨‍⚕️ Qualified Doctors</h3>
                    <p>Connect with experienced and certified medical professionals</p>
                </div>
                <div class="feature-card">
                    <h3>📅 Easy Booking</h3>
                    <p>Simple and intuitive appointment booking system</p>
                </div>
                <div class="feature-card">
                    <h3>🔒 Secure Platform</h3>
                    <p>Your health data is protected with industry-standard security</p>
                </div>
            </div>
        </div>
    </div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// LoginPageHandler serves the login page
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <div class="auth-form">
            <h2>Login</h2>
            <form method="POST" action="/login">
                <div class="form-group">
                    <label for="email">Email:</label>
                    <input type="email" id="email" name="email" required>
                </div>
                <div class="form-group">
                    <label for="password">Password:</label>
                    <input type="password" id="password" name="password" required>
                </div>
                <button type="submit" class="btn btn-primary">Login</button>
            </form>
            <p><a href="/register">Don't have an account? Register here</a></p>
            <p><a href="/">← Back to Home</a></p>
            
            <div class="demo-accounts">
                <h4>Demo Accounts:</h4>
                <p><strong>Patient:</strong> patient1@email.com / password123</p>
                <p><strong>Doctor:</strong> dr.smith@hospital.com / password123</p>
                <p><strong>Admin:</strong> admin@hospital.com / password123</p>
            </div>
        </div>
    </div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// LoginHandler handles login form submission
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get user from database
	user, err := models.GetUserByEmail(database.DB, email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Create session
	sessionToken := generateSessionToken()
	sessions[sessionToken] = &SessionData{
		UserID:   user.ID,
		UserType: user.UserType,
		Email:    user.Email,
		ExpireAt: time.Now().Add(24 * time.Hour), // Session expires in 24 hours
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// Redirect based on user type
	switch user.UserType {
	case "patient":
		http.Redirect(w, r, "/dashboard/patient", http.StatusSeeOther)
	case "doctor":
		http.Redirect(w, r, "/dashboard/doctor", http.StatusSeeOther)
	case "admin":
		http.Redirect(w, r, "/dashboard/admin", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// RegisterPageHandler serves the registration page
func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Register - Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <div class="auth-form">
            <h2>Register</h2>
            <form method="POST" action="/register">
                <div class="form-group">
                    <label for="first_name">First Name:</label>
                    <input type="text" id="first_name" name="first_name" required>
                </div>
                <div class="form-group">
                    <label for="last_name">Last Name:</label>
                    <input type="text" id="last_name" name="last_name" required>
                </div>
                <div class="form-group">
                    <label for="email">Email:</label>
                    <input type="email" id="email" name="email" required>
                </div>
                <div class="form-group">
                    <label for="phone">Phone:</label>
                    <input type="tel" id="phone" name="phone">
                </div>
                <div class="form-group">
                    <label for="password">Password:</label>
                    <input type="password" id="password" name="password" required minlength="6">
                </div>
                <div class="form-group">
                    <label for="user_type">I am a:</label>
                    <select id="user_type" name="user_type" required>
                        <option value="">Select...</option>
                        <option value="patient">Patient</option>
                        <option value="doctor">Doctor</option>
                    </select>
                </div>
                <button type="submit" class="btn btn-primary">Register</button>
            </form>
            <p><a href="/login">Already have an account? Login here</a></p>
            <p><a href="/">← Back to Home</a></p>
        </div>
    </div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// RegisterHandler handles registration form submission
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	userType := r.FormValue("user_type")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Create user
	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		Phone:        phone,
		UserType:     userType,
	}

	err = models.CreateUser(database.DB, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Email already exists or registration failed", http.StatusBadRequest)
		return
	}

	// If user is a doctor, create doctor profile
	if userType == "doctor" {
		doctor := &models.Doctor{
			UserID:          user.ID,
			Specialty:       "General Practice", // Default specialty
			ExperienceYears: 0,
			Education:       "",
			About:           "",
			ConsultationFee: 50.00, // Default fee
		}
		err = models.CreateDoctor(database.DB, doctor)
		if err != nil {
			log.Printf("Error creating doctor profile: %v", err)
		}
	}

	// Redirect to login page with success message
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Delete session from memory
		delete(sessions, cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AuthMiddleware checks if user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, exists := sessions[cookie.Value]
		if !exists || session.ExpireAt.Before(time.Now()) {
			// Session expired or doesn't exist
			delete(sessions, cookie.Value)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user info to request context (we'll use a simple approach)
		r.Header.Set("X-User-ID", string(rune(session.UserID)))
		r.Header.Set("X-User-Type", session.UserType)
		r.Header.Set("X-User-Email", session.Email)

		next.ServeHTTP(w, r)
	})
}

// GetCurrentUser extracts current user from request
func GetCurrentUser(r *http.Request) (int, string, string) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, "", ""
	}

	session, exists := sessions[cookie.Value]
	if !exists {
		return 0, "", ""
	}

	return session.UserID, session.UserType, session.Email
}

// Helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
