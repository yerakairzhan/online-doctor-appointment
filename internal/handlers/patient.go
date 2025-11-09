package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"online-doctor-appointment/internal/database"
	"online-doctor-appointment/internal/models"

	"github.com/gorilla/mux"
)

// PatientDashboardHandler serves the patient dashboard
func PatientDashboardHandler(w http.ResponseWriter, r *http.Request) {
	userID, userType, email := GetCurrentUser(r)
	if userType != "patient" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get user info
	user, err := models.GetUserByID(database.DB, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Get recent appointments
	appointments, err := models.GetAppointmentsByPatientID(database.DB, userID)
	if err != nil {
		appointments = []models.Appointment{} // Empty if error
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Patient Dashboard - Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <div class="dashboard">
            <div class="dashboard-header">
                <h2>Patient Dashboard</h2>
                <div class="user-info">
                    <span>Welcome, ` + user.GetFullName() + `</span>
                    <span>` + email + `</span>
                    <form method="POST" action="/logout" style="margin-top: 10px;">
                        <button type="submit" class="btn btn-secondary">Logout</button>
                    </form>
                </div>
            </div>

            <div class="dashboard-content">
                <div class="card">
                    <h3>Quick Actions</h3>
                    <div class="action-buttons">
                        <a href="/dashboard/patient/book" class="btn btn-primary">Book New Appointment</a>
                        <a href="/dashboard/patient/appointments" class="btn btn-info">View All Appointments</a>
                    </div>
                </div>

                <div class="card">
                    <h3>Recent Appointments</h3>`

	if len(appointments) == 0 {
		tmpl += `<p>No appointments found. <a href="/dashboard/patient/book">Book your first appointment</a></p>`
	} else {
		tmpl += `<table class="table">
                        <thead>
                            <tr>
                                <th>Doctor</th>
                                <th>Specialty</th>
                                <th>Date</th>
                                <th>Time</th>
                                <th>Status</th>
                            </tr>
                        </thead>
                        <tbody>`

		for i, appointment := range appointments {
			if i >= 5 { // Show only first 5
				break
			}
			tmpl += fmt.Sprintf(`
                            <tr>
                                <td>Dr. %s</td>
                                <td>%s</td>
                                <td>%s</td>
                                <td>%s</td>
                                <td><span class="status %s">%s</span></td>
                            </tr>`,
				appointment.Doctor.User.GetFullName(),
				appointment.Doctor.Specialty,
				appointment.AppointmentDate,
				appointment.AppointmentTime,
				appointment.Status,
				appointment.Status)
		}

		tmpl += `</tbody></table>`
	}

	tmpl += `
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

// BookAppointmentPageHandler serves the appointment booking page
func BookAppointmentPageHandler(w http.ResponseWriter, r *http.Request) {
	_, userType, email := GetCurrentUser(r)
	if userType != "patient" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get all doctors
	doctors, err := models.GetAllDoctors(database.DB)
	if err != nil {
		http.Error(w, "Error loading doctors", http.StatusInternalServerError)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Book Appointment - Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .payment-section {
            background: #f8f9fa;
            border: 2px solid #e9ecef;
            border-radius: 10px;
            padding: 20px;
            margin-top: 20px;
        }
        .payment-info {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 15px;
        }
        .fee-display {
            font-size: 1.2rem;
            font-weight: bold;
            color: #28a745;
        }
        .kaspi-button {
            background: linear-gradient(135deg, #00a651 0%, #00d4aa 100%);
            color: white;
            border: none;
            padding: 12px 30px;
            border-radius: 8px;
            font-weight: 600;
            cursor: pointer;
            font-size: 1rem;
            transition: all 0.3s ease;
            text-decoration: none;
            display: inline-block;
            text-align: center;
        }
        .kaspi-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(0, 166, 81, 0.4);
        }
        .kaspi-logo {
            width: 20px;
            height: 20px;
            margin-right: 8px;
        }
        .payment-options {
            display: flex;
            gap: 15px;
            align-items: center;
            flex-wrap: wrap;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="dashboard">
            <div class="dashboard-header">
                <h2>Book New Appointment</h2>
                <div class="user-info">
                    <span>` + email + `</span>
                    <a href="/dashboard/patient" class="btn btn-secondary">← Back to Dashboard</a>
                </div>
            </div>

            <form method="POST" action="/dashboard/patient/book" class="booking-form" id="bookingForm">
                <div class="form-group">
                    <label for="doctor_id">Select Doctor:</label>
                    <select id="doctor_id" name="doctor_id" required onchange="updateFee()">
                        <option value="">Choose a doctor...</option>`

	for _, doctor := range doctors {
		tmpl += fmt.Sprintf(`
                        <option value="%d" data-fee="%.2f" data-name="Dr. %s" data-specialty="%s">Dr. %s - %s ($%.2f)</option>`,
			doctor.ID,
			doctor.ConsultationFee,
			doctor.User.GetFullName(),
			doctor.Specialty,
			doctor.User.GetFullName(),
			doctor.Specialty,
			doctor.ConsultationFee)
	}

	// Get tomorrow's date as minimum
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	tmpl += `
                    </select>
                </div>

                <div class="form-group">
                    <label for="appointment_date">Appointment Date:</label>
                    <input type="date" id="appointment_date" name="appointment_date" 
                           min="` + tomorrow + `" required>
                </div>

                <div class="form-group">
                    <label for="appointment_time">Appointment Time:</label>
                    <select id="appointment_time" name="appointment_time" required>
                        <option value="">Select time...</option>
                        <option value="09:00">09:00 AM</option>
                        <option value="10:00">10:00 AM</option>
                        <option value="11:00">11:00 AM</option>
                        <option value="12:00">12:00 PM</option>
                        <option value="13:00">01:00 PM</option>
                        <option value="14:00">02:00 PM</option>
                        <option value="15:00">03:00 PM</option>
                        <option value="16:00">04:00 PM</option>
                        <option value="17:00">05:00 PM</option>
                    </select>
                </div>

                <div class="form-group">
                    <label for="notes">Notes (optional):</label>
                    <textarea id="notes" name="notes" rows="4" 
                              placeholder="Describe your symptoms or reason for visit..."></textarea>
                </div>

                <!-- Payment Section -->
                <div class="payment-section" id="paymentSection" style="display: none;">
                    <h3>💳 Payment Information</h3>
                    <div class="payment-info">
                        <div>
                            <strong>Consultation Fee:</strong>
                            <span class="fee-display" id="consultationFee">$0.00</span>
                        </div>
                        <div>
                            <strong>Doctor:</strong>
                            <span id="selectedDoctor">-</span>
                        </div>
                    </div>
                    <p><small>💡 You can pay now via Kaspi or pay at the clinic during your visit.</small></p>
                    <div class="payment-options">
                        <button type="submit" class="btn btn-primary">Book Appointment (Pay Later)</button>
                        <span style="margin: 0 10px; color: #666;">OR</span>
                        <a href="#" id="kaspiPayButton" class="kaspi-button" onclick="payWithKaspi(event)">
                            🏦 Pay with Kaspi
                        </a>
                    </div>
                </div>

                <!-- Default button when no doctor selected -->
                <div id="defaultButton">
                    <button type="submit" class="btn btn-primary">Book Appointment</button>
                </div>
            </form>
        </div>
    </div>

    <script>
        function updateFee() {
            const doctorSelect = document.getElementById('doctor_id');
            const paymentSection = document.getElementById('paymentSection');
            const defaultButton = document.getElementById('defaultButton');
            const feeDisplay = document.getElementById('consultationFee');
            const doctorDisplay = document.getElementById('selectedDoctor');
            
            if (doctorSelect.value) {
                const selectedOption = doctorSelect.options[doctorSelect.selectedIndex];
                const fee = selectedOption.dataset.fee;
                const doctorName = selectedOption.dataset.name;
                const specialty = selectedOption.dataset.specialty;
                
                feeDisplay.textContent = '$' + parseFloat(fee).toFixed(2);
                doctorDisplay.textContent = doctorName + ' (' + specialty + ')';
                
                paymentSection.style.display = 'block';
                defaultButton.style.display = 'none';
            } else {
                paymentSection.style.display = 'none';
                defaultButton.style.display = 'block';
            }
        }

        function payWithKaspi(event) {
            event.preventDefault();
            
            const doctorSelect = document.getElementById('doctor_id');
            const dateInput = document.getElementById('appointment_date');
            const timeInput = document.getElementById('appointment_time');
            
            if (!doctorSelect.value || !dateInput.value || !timeInput.value) {
                alert('Please fill in all required fields before proceeding to payment.');
                return;
            }
            
            const selectedOption = doctorSelect.options[doctorSelect.selectedIndex];
            const fee = selectedOption.dataset.fee;
            const doctorName = selectedOption.dataset.name;
            const appointmentDate = dateInput.value;
            const appointmentTime = timeInput.value;
            
            // Create payment description
            const description = 'Medical consultation with ' + doctorName + ' on ' + appointmentDate + ' at ' + appointmentTime;
            
            // Redirect to Kaspi payment (this is a demo URL - replace with actual Kaspi integration)
            const kaspiUrl = generateKaspiPaymentUrl(fee, description);
            
            // In a real application, you would:
            // 1. First save the appointment with "pending_payment" status
            // 2. Then redirect to Kaspi
            // 3. Handle the callback to confirm payment
            
            // For now, we'll show the Kaspi payment link
            if (confirm('Proceed to Kaspi payment for $' + fee + '?')) {
                window.open(kaspiUrl, '_blank');
                // Optionally submit the form after payment
                // document.getElementById('bookingForm').submit();
            }
        }

        function generateKaspiPaymentUrl(amount, description) {
            // This is a simplified Kaspi payment URL structure
            // In production, you would use official Kaspi Payment API
            const baseUrl = 'https://kaspi.kz/pay';
            const merchantId = 'DEMO_MERCHANT'; // Replace with your actual merchant ID
            const orderId = 'ORDER_' + Date.now();
            
            const params = new URLSearchParams({
                'amount': amount,
                'currency': 'KZT', // Assuming Kazakhstani Tenge
                'description': description,
                'merchant_id': merchantId,
                'order_id': orderId,
                'return_url': window.location.origin + '/dashboard/patient/appointments',
                'cancel_url': window.location.href
            });
            
            // Note: This is a demo URL structure
            // For real Kaspi integration, you need to:
            // 1. Register as a Kaspi merchant
            // 2. Use their official API endpoints
            // 3. Implement proper authentication and callbacks
            
            return baseUrl + '?' + params.toString();
        }

        // Auto-update fee when page loads if doctor is pre-selected
        document.addEventListener('DOMContentLoaded', function() {
            updateFee();
        });
    </script>
</body>
</html>
`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// BookAppointmentHandler handles appointment booking form submission
func BookAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	userID, userType, _ := GetCurrentUser(r)
	if userType != "patient" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	doctorID, _ := strconv.Atoi(r.FormValue("doctor_id"))
	appointmentDate := r.FormValue("appointment_date")
	appointmentTime := r.FormValue("appointment_time")
	notes := r.FormValue("notes")

	// Create appointment
	appointment := &models.Appointment{
		PatientID:       userID,
		DoctorID:        doctorID,
		AppointmentDate: appointmentDate,
		AppointmentTime: appointmentTime,
		Notes:           notes,
	}

	err := models.CreateAppointment(database.DB, appointment)
	if err != nil {
		http.Error(w, "Failed to book appointment. Time slot may be unavailable.", http.StatusBadRequest)
		return
	}

	// Redirect to appointments page with success
	http.Redirect(w, r, "/dashboard/patient/appointments", http.StatusSeeOther)
}

// PatientAppointmentsHandler shows all patient appointments
func PatientAppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	userID, userType, email := GetCurrentUser(r)
	if userType != "patient" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get all appointments
	appointments, err := models.GetAppointmentsByPatientID(database.DB, userID)
	if err != nil {
		http.Error(w, "Error loading appointments", http.StatusInternalServerError)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My Appointments - Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <div class="dashboard">
            <div class="dashboard-header">
                <h2>My Appointments</h2>
                <div class="user-info">
                    <span>` + email + `</span>
                    <a href="/dashboard/patient" class="btn btn-secondary">← Back to Dashboard</a>
                </div>
            </div>

            <div class="card">
                <div class="action-buttons">
                    <a href="/dashboard/patient/book" class="btn btn-primary">Book New Appointment</a>
                </div>
            </div>

            <div class="card">
                <h3>All Appointments</h3>`

	if len(appointments) == 0 {
		tmpl += `<p>No appointments found. <a href="/dashboard/patient/book">Book your first appointment</a></p>`
	} else {
		tmpl += `<table class="table">
                    <thead>
                        <tr>
                            <th>Doctor</th>
                            <th>Specialty</th>
                            <th>Date</th>
                            <th>Time</th>
                            <th>Status</th>
                            <th>Fee</th>
                            <th>Notes</th>
                        </tr>
                    </thead>
                    <tbody>`

		for _, appointment := range appointments {
			tmpl += fmt.Sprintf(`
                        <tr>
                            <td>Dr. %s</td>
                            <td>%s</td>
                            <td>%s</td>
                            <td>%s</td>
                            <td><span class="status %s">%s</span></td>
                            <td>$%.2f</td>
                            <td>%s</td>
                        </tr>`,
				appointment.Doctor.User.GetFullName(),
				appointment.Doctor.Specialty,
				appointment.AppointmentDate,
				appointment.AppointmentTime,
				appointment.Status,
				appointment.Status,
				appointment.Doctor.ConsultationFee,
				appointment.Notes)
		}

		tmpl += `</tbody></table>`
	}

	tmpl += `
            </div>
        </div>
    </div>
</body>
</html>
`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// GetDoctorsHandler returns all doctors as JSON
func GetDoctorsHandler(w http.ResponseWriter, r *http.Request) {
	doctors, err := models.GetAllDoctors(database.DB)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error loading doctors"})
		return
	}

	respondWithJSON(w, http.StatusOK, doctors)
}

// GetDoctorsBySpecialtyHandler returns doctors filtered by specialty
func GetDoctorsBySpecialtyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	specialty := vars["specialty"]

	doctors, err := models.GetDoctorsBySpecialty(database.DB, specialty)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error loading doctors"})
		return
	}

	respondWithJSON(w, http.StatusOK, doctors)
}

// GetAvailableSlotsHandler returns available time slots for a doctor on a specific date
func GetAvailableSlotsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorIDStr := vars["doctorId"]
	date := vars["date"]

	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid doctor ID"})
		return
	}

	slots, err := models.GetAvailableTimeSlots(database.DB, doctorID, date)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error loading time slots"})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"doctor_id": doctorID,
		"date":      date,
		"slots":     slots,
	})
}

// PaymentSuccessHandler handles successful payments from Kaspi
func PaymentSuccessHandler(w http.ResponseWriter, r *http.Request) {
	_, userType, email := GetCurrentUser(r)
	if userType != "patient" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get payment parameters from URL (in real implementation, verify with Kaspi)
	orderId := r.URL.Query().Get("order_id")
	amount := r.URL.Query().Get("amount")

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Payment Successful - Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .success-card {
            background: #d4edda;
            border: 2px solid #c3e6cb;
            border-radius: 10px;
            padding: 30px;
            text-align: center;
            color: #155724;
        }
        .success-icon {
            font-size: 4rem;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="dashboard">
            <div class="dashboard-header">
                <h2>Payment Successful</h2>
                <div class="user-info">
                    <span>` + email + `</span>
                </div>
            </div>

            <div class="success-card">
                <div class="success-icon">✅</div>
                <h3>Payment Completed Successfully!</h3>
                <p><strong>Order ID:</strong> ` + orderId + `</p>
                <p><strong>Amount:</strong> $` + amount + `</p>
                <p>Your appointment has been confirmed and payment processed via Kaspi.</p>
                <p>You will receive a confirmation email shortly.</p>
                
                <div class="action-buttons" style="margin-top: 30px;">
                    <a href="/dashboard/patient/appointments" class="btn btn-primary">View My Appointments</a>
                    <a href="/dashboard/patient" class="btn btn-secondary">Back to Dashboard</a>
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

// PaymentFailureHandler handles failed payments from Kaspi
func PaymentFailureHandler(w http.ResponseWriter, r *http.Request) {
	_, userType, email := GetCurrentUser(r)
	if userType != "patient" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get error details from URL
	errorMessage := r.URL.Query().Get("error")
	if errorMessage == "" {
		errorMessage = "Payment was cancelled or failed"
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Payment Failed - Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .error-card {
            background: #f8d7da;
            border: 2px solid #f5c6cb;
            border-radius: 10px;
            padding: 30px;
            text-align: center;
            color: #721c24;
        }
        .error-icon {
            font-size: 4rem;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="dashboard">
            <div class="dashboard-header">
                <h2>Payment Failed</h2>
                <div class="user-info">
                    <span>` + email + `</span>
                </div>
            </div>

            <div class="error-card">
                <div class="error-icon">❌</div>
                <h3>Payment Could Not Be Processed</h3>
                <p>` + errorMessage + `</p>
                <p>Don't worry! You can still book the appointment and pay at the clinic, or try the payment again.</p>
                
                <div class="action-buttons" style="margin-top: 30px;">
                    <a href="/dashboard/patient/book" class="btn btn-primary">Try Again</a>
                    <a href="/dashboard/patient" class="btn btn-secondary">Back to Dashboard</a>
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
