package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"online-doctor-appointment/internal/database"
	"online-doctor-appointment/internal/models"

	"github.com/gorilla/mux"
)

// DoctorDashboardHandler serves the doctor dashboard
func DoctorDashboardHandler(w http.ResponseWriter, r *http.Request) {
	userID, userType, email := GetCurrentUser(r)
	if userType != "doctor" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get doctor info
	doctor, err := models.GetDoctorByUserID(database.DB, userID)
	if err != nil {
		http.Error(w, "Doctor profile not found", http.StatusInternalServerError)
		return
	}

	// Get recent appointments
	appointments, err := models.GetAppointmentsByDoctorID(database.DB, doctor.ID)
	if err != nil {
		appointments = []models.Appointment{} // Empty if error
	}

	// Count appointments by status
	pendingCount := 0
	confirmedCount := 0
	completedCount := 0

	for _, apt := range appointments {
		switch apt.Status {
		case "pending":
			pendingCount++
		case "confirmed":
			confirmedCount++
		case "completed":
			completedCount++
		}
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Doctor Dashboard - Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <div class="dashboard">
            <div class="dashboard-header">
                <h2>Doctor Dashboard</h2>
                <div class="user-info">
                    <span>Dr. ` + doctor.User.GetFullName() + `</span>
                    <span>` + doctor.Specialty + `</span>
                    <span>` + email + `</span>
                    <form method="POST" action="/logout" style="margin-top: 10px;">
                        <button type="submit" class="btn btn-secondary">Logout</button>
                    </form>
                </div>
            </div>

            <div class="dashboard-content">
                <div class="card">
                    <h3>Profile Information</h3>
                    <p><strong>Specialty:</strong> ` + doctor.Specialty + `</p>
                    <p><strong>Experience:</strong> ` + fmt.Sprintf("%d", doctor.ExperienceYears) + ` years</p>
                    <p><strong>Consultation Fee:</strong> $` + fmt.Sprintf("%.2f", doctor.ConsultationFee) + `</p>
                    <p><strong>About:</strong> ` + doctor.About + `</p>
                </div>

                <div class="card">
                    <h3>Appointment Statistics</h3>
                    <div class="stats-grid" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 20px;">
                        <div class="stat-card" style="background: #fff3cd; padding: 20px; border-radius: 8px; text-align: center;">
                            <h4 style="margin: 0; color: #856404;">` + fmt.Sprintf("%d", pendingCount) + `</h4>
                            <p style="margin: 5px 0 0 0; color: #856404;">Pending</p>
                        </div>
                        <div class="stat-card" style="background: #d4edda; padding: 20px; border-radius: 8px; text-align: center;">
                            <h4 style="margin: 0; color: #155724;">` + fmt.Sprintf("%d", confirmedCount) + `</h4>
                            <p style="margin: 5px 0 0 0; color: #155724;">Confirmed</p>
                        </div>
                        <div class="stat-card" style="background: #d1ecf1; padding: 20px; border-radius: 8px; text-align: center;">
                            <h4 style="margin: 0; color: #0c5460;">` + fmt.Sprintf("%d", completedCount) + `</h4>
                            <p style="margin: 5px 0 0 0; color: #0c5460;">Completed</p>
                        </div>
                    </div>
                </div>

                <div class="card">
                    <h3>Quick Actions</h3>
                    <div class="action-buttons">
                        <a href="/dashboard/doctor/appointments" class="btn btn-primary">View All Appointments</a>
                    </div>
                </div>

                <div class="card">
                    <h3>Recent Appointments</h3>`

	if len(appointments) == 0 {
		tmpl += `<p>No appointments scheduled.</p>`
	} else {
		tmpl += `<table class="table">
                        <thead>
                            <tr>
                                <th>Patient</th>
                                <th>Date</th>
                                <th>Time</th>
                                <th>Status</th>
                                <th>Notes</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>`

		for i, appointment := range appointments {
			if i >= 5 { // Show only first 5
				break
			}
			tmpl += fmt.Sprintf(`
                            <tr>
                                <td>%s</td>
                                <td>%s</td>
                                <td>%s</td>
                                <td><span class="status %s">%s</span></td>
                                <td>%s</td>
                                <td>`,
				appointment.Patient.GetFullName(),
				appointment.AppointmentDate,
				appointment.AppointmentTime,
				appointment.Status,
				appointment.Status,
				appointment.Notes)

			// Add action buttons based on status
			if appointment.Status == "pending" {
				tmpl += fmt.Sprintf(`
                                    <form method="POST" action="/dashboard/doctor/appointment/%d/update" style="display: inline;">
                                        <input type="hidden" name="status" value="confirmed">
                                        <button type="submit" class="btn btn-success" style="padding: 5px 10px; font-size: 0.8rem;">Confirm</button>
                                    </form>
                                    <form method="POST" action="/dashboard/doctor/appointment/%d/update" style="display: inline; margin-left: 5px;">
                                        <input type="hidden" name="status" value="cancelled">
                                        <button type="submit" class="btn btn-danger" style="padding: 5px 10px; font-size: 0.8rem;">Cancel</button>
                                    </form>`,
					appointment.ID, appointment.ID)
			} else if appointment.Status == "confirmed" {
				tmpl += fmt.Sprintf(`
                                    <form method="POST" action="/dashboard/doctor/appointment/%d/update" style="display: inline;">
                                        <input type="hidden" name="status" value="completed">
                                        <button type="submit" class="btn btn-info" style="padding: 5px 10px; font-size: 0.8rem;">Complete</button>
                                    </form>`,
					appointment.ID)
			}

			tmpl += `</td></tr>`
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

// DoctorAppointmentsHandler shows all doctor appointments
func DoctorAppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	userID, userType, email := GetCurrentUser(r)
	if userType != "doctor" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get doctor info
	doctor, err := models.GetDoctorByUserID(database.DB, userID)
	if err != nil {
		http.Error(w, "Doctor profile not found", http.StatusInternalServerError)
		return
	}

	// Get all appointments
	appointments, err := models.GetAppointmentsByDoctorID(database.DB, doctor.ID)
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
    <title>My Appointments - Doctor Dashboard</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <div class="dashboard">
            <div class="dashboard-header">
                <h2>My Appointments</h2>
                <div class="user-info">
                    <span>Dr. ` + doctor.User.GetFullName() + `</span>
                    <span>` + email + `</span>
                    <a href="/dashboard/doctor" class="btn btn-secondary">← Back to Dashboard</a>
                </div>
            </div>

            <div class="card">
                <h3>All Appointments</h3>`

	if len(appointments) == 0 {
		tmpl += `<p>No appointments scheduled.</p>`
	} else {
		tmpl += `<table class="table">
                    <thead>
                        <tr>
                            <th>Patient</th>
                            <th>Contact</th>
                            <th>Date</th>
                            <th>Time</th>
                            <th>Status</th>
                            <th>Notes</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>`

		for _, appointment := range appointments {
			tmpl += fmt.Sprintf(`
                        <tr>
                            <td>%s</td>
                            <td>%s<br>%s</td>
                            <td>%s</td>
                            <td>%s</td>
                            <td><span class="status %s">%s</span></td>
                            <td>%s</td>
                            <td>`,
				appointment.Patient.GetFullName(),
				appointment.Patient.Email,
				appointment.Patient.Phone,
				appointment.AppointmentDate,
				appointment.AppointmentTime,
				appointment.Status,
				appointment.Status,
				appointment.Notes)

			// Add action buttons based on status
			if appointment.Status == "pending" {
				tmpl += fmt.Sprintf(`
                            <form method="POST" action="/dashboard/doctor/appointment/%d/update" style="display: inline;">
                                <input type="hidden" name="status" value="confirmed">
                                <button type="submit" class="btn btn-success" style="padding: 5px 10px; font-size: 0.8rem;">Confirm</button>
                            </form>
                            <form method="POST" action="/dashboard/doctor/appointment/%d/update" style="display: inline; margin-left: 5px;">
                                <input type="hidden" name="status" value="cancelled">
                                <button type="submit" class="btn btn-danger" style="padding: 5px 10px; font-size: 0.8rem;">Cancel</button>
                            </form>`,
					appointment.ID, appointment.ID)
			} else if appointment.Status == "confirmed" {
				tmpl += fmt.Sprintf(`
                            <form method="POST" action="/dashboard/doctor/appointment/%d/update" style="display: inline;">
                                <input type="hidden" name="status" value="completed">
                                <button type="submit" class="btn btn-info" style="padding: 5px 10px; font-size: 0.8rem;">Complete</button>
                            </form>`,
					appointment.ID)
			}

			tmpl += `</td></tr>`
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

// UpdateAppointmentStatusHandler handles appointment status updates
func UpdateAppointmentStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, userType, _ := GetCurrentUser(r)
	if userType != "doctor" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	appointmentIDStr := vars["id"]
	appointmentID, err := strconv.Atoi(appointmentIDStr)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	newStatus := r.FormValue("status")
	if newStatus == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	// Verify that this appointment belongs to the current doctor
	appointment, err := models.GetAppointmentByID(database.DB, appointmentID)
	if err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	doctor, err := models.GetDoctorByUserID(database.DB, userID)
	if err != nil || doctor.ID != appointment.DoctorID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Update appointment status
	err = models.UpdateAppointmentStatus(database.DB, appointmentID, newStatus)
	if err != nil {
		http.Error(w, "Failed to update appointment", http.StatusInternalServerError)
		return
	}

	// Redirect back to appointments page
	http.Redirect(w, r, "/dashboard/doctor/appointments", http.StatusSeeOther)
}
