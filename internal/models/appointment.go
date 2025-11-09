package models

import (
	"database/sql"
	"time"
)

type Appointment struct {
	ID              int       `json:"id"`
	PatientID       int       `json:"patient_id"`
	DoctorID        int       `json:"doctor_id"`
	AppointmentDate string    `json:"appointment_date"` // YYYY-MM-DD format
	AppointmentTime string    `json:"appointment_time"` // HH:MM format
	Status          string    `json:"status"`           // pending, confirmed, cancelled, completed
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Embedded information
	Patient *User   `json:"patient,omitempty"`
	Doctor  *Doctor `json:"doctor,omitempty"`
}

// CreateAppointment inserts a new appointment
func CreateAppointment(db *sql.DB, appointment *Appointment) error {
	query := `
		INSERT INTO appointments (patient_id, doctor_id, appointment_date, appointment_time, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, status, created_at, updated_at
	`

	err := db.QueryRow(query, appointment.PatientID, appointment.DoctorID,
		appointment.AppointmentDate, appointment.AppointmentTime, appointment.Notes).Scan(
		&appointment.ID, &appointment.Status, &appointment.CreatedAt, &appointment.UpdatedAt)

	return err
}

// GetAppointmentByID retrieves an appointment by ID
func GetAppointmentByID(db *sql.DB, appointmentID int) (*Appointment, error) {
	appointment := &Appointment{}
	query := `
		SELECT a.id, a.patient_id, a.doctor_id, a.appointment_date, 
		       a.appointment_time, a.status, a.notes, a.created_at, a.updated_at,
		       u.first_name, u.last_name, u.email, u.phone,
		       d.specialty, d.consultation_fee,
		       du.first_name, du.last_name
		FROM appointments a
		JOIN users u ON a.patient_id = u.id
		JOIN doctors d ON a.doctor_id = d.id
		JOIN users du ON d.user_id = du.id
		WHERE a.id = $1
	`

	var patient User
	var doctor Doctor
	var doctorUser User

	err := db.QueryRow(query, appointmentID).Scan(
		&appointment.ID, &appointment.PatientID, &appointment.DoctorID,
		&appointment.AppointmentDate, &appointment.AppointmentTime,
		&appointment.Status, &appointment.Notes, &appointment.CreatedAt, &appointment.UpdatedAt,
		&patient.FirstName, &patient.LastName, &patient.Email, &patient.Phone,
		&doctor.Specialty, &doctor.ConsultationFee,
		&doctorUser.FirstName, &doctorUser.LastName,
	)

	if err != nil {
		return nil, err
	}

	patient.ID = appointment.PatientID
	doctorUser.ID = doctor.UserID
	doctor.User = &doctorUser
	appointment.Patient = &patient
	appointment.Doctor = &doctor

	return appointment, nil
}

// GetAppointmentsByPatientID retrieves all appointments for a patient
func GetAppointmentsByPatientID(db *sql.DB, patientID int) ([]Appointment, error) {
	query := `
		SELECT a.id, a.patient_id, a.doctor_id, a.appointment_date, 
		       a.appointment_time, a.status, a.notes, a.created_at, a.updated_at,
		       d.specialty, d.consultation_fee,
		       du.first_name, du.last_name
		FROM appointments a
		JOIN doctors d ON a.doctor_id = d.id
		JOIN users du ON d.user_id = du.id
		WHERE a.patient_id = $1
		ORDER BY a.appointment_date DESC, a.appointment_time DESC
	`

	rows, err := db.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var appointment Appointment
		var doctor Doctor
		var doctorUser User

		err := rows.Scan(
			&appointment.ID, &appointment.PatientID, &appointment.DoctorID,
			&appointment.AppointmentDate, &appointment.AppointmentTime,
			&appointment.Status, &appointment.Notes, &appointment.CreatedAt, &appointment.UpdatedAt,
			&doctor.Specialty, &doctor.ConsultationFee,
			&doctorUser.FirstName, &doctorUser.LastName,
		)
		if err != nil {
			return nil, err
		}

		doctor.User = &doctorUser
		appointment.Doctor = &doctor
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

// GetAppointmentsByDoctorID retrieves all appointments for a doctor
func GetAppointmentsByDoctorID(db *sql.DB, doctorID int) ([]Appointment, error) {
	query := `
		SELECT a.id, a.patient_id, a.doctor_id, a.appointment_date, 
		       a.appointment_time, a.status, a.notes, a.created_at, a.updated_at,
		       u.first_name, u.last_name, u.email, u.phone
		FROM appointments a
		JOIN users u ON a.patient_id = u.id
		WHERE a.doctor_id = $1
		ORDER BY a.appointment_date DESC, a.appointment_time DESC
	`

	rows, err := db.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var appointment Appointment
		var patient User

		err := rows.Scan(
			&appointment.ID, &appointment.PatientID, &appointment.DoctorID,
			&appointment.AppointmentDate, &appointment.AppointmentTime,
			&appointment.Status, &appointment.Notes, &appointment.CreatedAt, &appointment.UpdatedAt,
			&patient.FirstName, &patient.LastName, &patient.Email, &patient.Phone,
		)
		if err != nil {
			return nil, err
		}

		patient.ID = appointment.PatientID
		appointment.Patient = &patient
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

// UpdateAppointmentStatus updates the status of an appointment
func UpdateAppointmentStatus(db *sql.DB, appointmentID int, status string) error {
	query := `UPDATE appointments SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.Exec(query, status, appointmentID)
	return err
}

// GetAvailableTimeSlots gets available time slots for a doctor on a specific date
func GetAvailableTimeSlots(db *sql.DB, doctorID int, date string) ([]string, error) {
	// First, get the day of week for the date
	dayQuery := `SELECT EXTRACT(DOW FROM DATE $1)`
	var dayOfWeek int
	err := db.QueryRow(dayQuery, date).Scan(&dayOfWeek)
	if err != nil {
		return nil, err
	}

	// Get doctor's availability for that day
	availQuery := `
		SELECT start_time, end_time 
		FROM doctor_availability 
		WHERE doctor_id = $1 AND day_of_week = $2 AND is_active = true
	`

	var startTime, endTime string
	err = db.QueryRow(availQuery, doctorID, dayOfWeek).Scan(&startTime, &endTime)
	if err != nil {
		return []string{}, nil // No availability that day
	}

	// Get booked appointments for that day
	bookedQuery := `
		SELECT appointment_time 
		FROM appointments 
		WHERE doctor_id = $1 AND appointment_date = $2 AND status != 'cancelled'
	`

	rows, err := db.Query(bookedQuery, doctorID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookedTimes := make(map[string]bool)
	for rows.Next() {
		var bookedTime string
		err := rows.Scan(&bookedTime)
		if err != nil {
			return nil, err
		}
		bookedTimes[bookedTime] = true
	}

	// Generate available time slots (every hour from start to end)
	var availableSlots []string
	// This is a simplified version - you might want to make it more sophisticated
	start, _ := time.Parse("15:04", startTime)
	end, _ := time.Parse("15:04", endTime)

	for current := start; current.Before(end); current = current.Add(time.Hour) {
		timeStr := current.Format("15:04")
		if !bookedTimes[timeStr] {
			availableSlots = append(availableSlots, timeStr)
		}
	}

	return availableSlots, nil
}
