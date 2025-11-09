package handlers

import (
	"fmt"
	"net/http"

	"online-doctor-appointment/internal/database"
	"online-doctor-appointment/internal/models"
)

// AdminDashboardHandler serves the admin dashboard
func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	_, userType, email := GetCurrentUser(r)
	if userType != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get counts for dashboard
	// For simplicity, we'll do basic counts (in production, you'd want proper count queries)

	doctors, _ := models.GetAllDoctors(database.DB)
	patients, _ := models.GetAllPatients(database.DB)

	doctorCount := len(doctors)
	patientCount := len(patients)

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Dashboard - Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <div class="dashboard">
            <div class="dashboard-header">
                <h2>Admin Dashboard</h2>
                <div class="user-info">
                    <span>Administrator</span>
                    <span>` + email + `</span>
                    <form method="POST" action="/logout" style="margin-top: 10px;">
                        <button type="submit" class="btn btn-secondary">Logout</button>
                    </form>
                </div>
            </div>

            <div class="dashboard-content">
                <div class="card">
                    <h3>System Overview</h3>
                    <div class="stats-grid" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px;">
                        <div class="stat-card" style="background: #d4edda; padding: 20px; border-radius: 8px; text-align: center;">
                            <h4 style="margin: 0; color: #155724; font-size: 2rem;">` + fmt.Sprintf("%d", doctorCount) + `</h4>
                            <p style="margin: 5px 0 0 0; color: #155724;">Total Doctors</p>
                        </div>
                        <div class="stat-card" style="background: #d1ecf1; padding: 20px; border-radius: 8px; text-align: center;">
                            <h4 style="margin: 0; color: #0c5460; font-size: 2rem;">` + fmt.Sprintf("%d", patientCount) + `</h4>
                            <p style="margin: 5px 0 0 0; color: #0c5460;">Total Patients</p>
                        </div>
                    </div>
                </div>

                <div class="card">
                    <h3>Quick Actions</h3>
                    <div class="action-buttons">
                        <a href="/dashboard/admin/doctors" class="btn btn-primary">Manage Doctors</a>
                        <a href="/dashboard/admin/patients" class="btn btn-info">View Patients</a>
                    </div>
                </div>

                <div class="card">
                    <h3>Recent Doctors</h3>`

	if len(doctors) == 0 {
		tmpl += `<p>No doctors registered.</p>`
	} else {
		tmpl += `<table class="table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Specialty</th>
                                <th>Experience</th>
                                <th>Fee</th>
                                <th>Status</th>
                            </tr>
                        </thead>
                        <tbody>`

		for i, doctor := range doctors {
			if i >= 5 { // Show only first 5
				break
			}
			status := "Active"
			if !doctor.IsActive {
				status = "Inactive"
			}
			tmpl += fmt.Sprintf(`
                            <tr>
                                <td>Dr. %s</td>
                                <td>%s</td>
                                <td>%d years</td>
                                <td>$%.2f</td>
                                <td><span class="status %s">%s</span></td>
                            </tr>`,
				doctor.User.GetFullName(),
				doctor.Specialty,
				doctor.ExperienceYears,
				doctor.ConsultationFee,
				func() string {
					if doctor.IsActive {
						return "confirmed"
					}
					return "cancelled"
				}(),
				status)
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

// AdminDoctorsHandler shows all doctors for admin management
func AdminDoctorsHandler(w http.ResponseWriter, r *http.Request) {
	_, userType, email := GetCurrentUser(r)
	if userType != "admin" {
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
    <title>Manage Doctors - Admin Dashboard</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <div class="dashboard">
            <div class="dashboard-header">
                <h2>Manage Doctors</h2>
                <div class="user-info">
                    <span>Administrator</span>
                    <span>` + email + `</span>
                    <a href="/dashboard/admin" class="btn btn-secondary">← Back to Dashboard</a>
                </div>
            </div>

            <div class="card">
                <h3>All Doctors (` + fmt.Sprintf("%d", len(doctors)) + `)</h3>`

	if len(doctors) == 0 {
		tmpl += `<p>No doctors registered in the system.</p>`
	} else {
		tmpl += `<table class="table">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Name</th>
                            <th>Email</th>
                            <th>Phone</th>
                            <th>Specialty</th>
                            <th>Experience</th>
                            <th>Fee</th>
                            <th>Status</th>
                            <th>Joined</th>
                        </tr>
                    </thead>
                    <tbody>`

		for _, doctor := range doctors {
			status := "Active"
			statusClass := "confirmed"
			if !doctor.IsActive {
				status = "Inactive"
				statusClass = "cancelled"
			}
			tmpl += fmt.Sprintf(`
                        <tr>
                            <td>%d</td>
                            <td>Dr. %s</td>
                            <td>%s</td>
                            <td>%s</td>
                            <td>%s</td>
                            <td>%d years</td>
                            <td>$%.2f</td>
                            <td><span class="status %s">%s</span></td>
                            <td>%s</td>
                        </tr>`,
				doctor.ID,
				doctor.User.GetFullName(),
				doctor.User.Email,
				doctor.User.Phone,
				doctor.Specialty,
				doctor.ExperienceYears,
				doctor.ConsultationFee,
				statusClass,
				status,
				doctor.CreatedAt.Format("2006-01-02"))
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

// AdminPatientsHandler shows all patients for admin
func AdminPatientsHandler(w http.ResponseWriter, r *http.Request) {
	_, userType, email := GetCurrentUser(r)
	if userType != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get all patients
	patients, err := models.GetAllPatients(database.DB)
	if err != nil {
		http.Error(w, "Error loading patients", http.StatusInternalServerError)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>View Patients - Admin Dashboard</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <div class="dashboard">
            <div class="dashboard-header">
                <h2>View Patients</h2>
                <div class="user-info">
                    <span>Administrator</span>
                    <span>` + email + `</span>
                    <a href="/dashboard/admin" class="btn btn-secondary">← Back to Dashboard</a>
                </div>
            </div>

            <div class="card">
                <h3>All Patients (` + fmt.Sprintf("%d", len(patients)) + `)</h3>`

	if len(patients) == 0 {
		tmpl += `<p>No patients registered in the system.</p>`
	} else {
		tmpl += `<table class="table">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Name</th>
                            <th>Email</th>
                            <th>Phone</th>
                            <th>Joined</th>
                        </tr>
                    </thead>
                    <tbody>`

		for _, patient := range patients {
			tmpl += fmt.Sprintf(`
                        <tr>
                            <td>%d</td>
                            <td>%s</td>
                            <td>%s</td>
                            <td>%s</td>
                            <td>%s</td>
                        </tr>`,
				patient.ID,
				patient.GetFullName(),
				patient.Email,
				patient.Phone,
				patient.CreatedAt.Format("2006-01-02"))
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
