package models

import (
	"database/sql"
	"time"
)

type Doctor struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Specialty       string    `json:"specialty"`
	ExperienceYears int       `json:"experience_years"`
	Education       string    `json:"education"`
	About           string    `json:"about"`
	ConsultationFee float64   `json:"consultation_fee"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Embedded user information
	User *User `json:"user,omitempty"`
}

type DoctorAvailability struct {
	ID        int       `json:"id"`
	DoctorID  int       `json:"doctor_id"`
	DayOfWeek int       `json:"day_of_week"` // 0=Sunday, 6=Saturday
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateDoctor inserts a new doctor into the database
func CreateDoctor(db *sql.DB, doctor *Doctor) error {
	query := `
		INSERT INTO doctors (user_id, specialty, experience_years, education, about, consultation_fee)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(query, doctor.UserID, doctor.Specialty, doctor.ExperienceYears,
		doctor.Education, doctor.About, doctor.ConsultationFee).Scan(
		&doctor.ID, &doctor.CreatedAt, &doctor.UpdatedAt)

	return err
}

// GetDoctorByUserID retrieves doctor information by user ID
func GetDoctorByUserID(db *sql.DB, userID int) (*Doctor, error) {
	doctor := &Doctor{}
	query := `
		SELECT d.id, d.user_id, d.specialty, d.experience_years, d.education, 
		       d.about, d.consultation_fee, d.is_active, d.created_at, d.updated_at,
		       u.email, u.first_name, u.last_name, u.phone
		FROM doctors d
		JOIN users u ON d.user_id = u.id
		WHERE d.user_id = $1
	`

	user := &User{}
	err := db.QueryRow(query, userID).Scan(
		&doctor.ID, &doctor.UserID, &doctor.Specialty, &doctor.ExperienceYears,
		&doctor.Education, &doctor.About, &doctor.ConsultationFee, &doctor.IsActive,
		&doctor.CreatedAt, &doctor.UpdatedAt,
		&user.Email, &user.FirstName, &user.LastName, &user.Phone,
	)

	if err != nil {
		return nil, err
	}

	user.ID = userID
	user.UserType = "doctor"
	doctor.User = user

	return doctor, nil
}

// GetDoctorByID retrieves doctor information by doctor ID
func GetDoctorByID(db *sql.DB, doctorID int) (*Doctor, error) {
	doctor := &Doctor{}
	query := `
		SELECT d.id, d.user_id, d.specialty, d.experience_years, d.education, 
		       d.about, d.consultation_fee, d.is_active, d.created_at, d.updated_at,
		       u.email, u.first_name, u.last_name, u.phone
		FROM doctors d
		JOIN users u ON d.user_id = u.id
		WHERE d.id = $1
	`

	user := &User{}
	err := db.QueryRow(query, doctorID).Scan(
		&doctor.ID, &doctor.UserID, &doctor.Specialty, &doctor.ExperienceYears,
		&doctor.Education, &doctor.About, &doctor.ConsultationFee, &doctor.IsActive,
		&doctor.CreatedAt, &doctor.UpdatedAt,
		&user.Email, &user.FirstName, &user.LastName, &user.Phone,
	)

	if err != nil {
		return nil, err
	}

	user.ID = doctor.UserID
	user.UserType = "doctor"
	doctor.User = user

	return doctor, nil
}

// GetAllDoctors retrieves all active doctors
func GetAllDoctors(db *sql.DB) ([]Doctor, error) {
	query := `
		SELECT d.id, d.user_id, d.specialty, d.experience_years, d.education, 
		       d.about, d.consultation_fee, d.is_active, d.created_at, d.updated_at,
		       u.email, u.first_name, u.last_name, u.phone
		FROM doctors d
		JOIN users u ON d.user_id = u.id
		WHERE d.is_active = true
		ORDER BY u.first_name, u.last_name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []Doctor
	for rows.Next() {
		var doctor Doctor
		var user User

		err := rows.Scan(
			&doctor.ID, &doctor.UserID, &doctor.Specialty, &doctor.ExperienceYears,
			&doctor.Education, &doctor.About, &doctor.ConsultationFee, &doctor.IsActive,
			&doctor.CreatedAt, &doctor.UpdatedAt,
			&user.Email, &user.FirstName, &user.LastName, &user.Phone,
		)
		if err != nil {
			return nil, err
		}

		user.ID = doctor.UserID
		user.UserType = "doctor"
		doctor.User = &user
		doctors = append(doctors, doctor)
	}

	return doctors, nil
}

// GetDoctorsBySpecialty retrieves doctors by specialty
func GetDoctorsBySpecialty(db *sql.DB, specialty string) ([]Doctor, error) {
	query := `
		SELECT d.id, d.user_id, d.specialty, d.experience_years, d.education, 
		       d.about, d.consultation_fee, d.is_active, d.created_at, d.updated_at,
		       u.email, u.first_name, u.last_name, u.phone
		FROM doctors d
		JOIN users u ON d.user_id = u.id
		WHERE d.specialty ILIKE $1 AND d.is_active = true
		ORDER BY u.first_name, u.last_name
	`

	rows, err := db.Query(query, "%"+specialty+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []Doctor
	for rows.Next() {
		var doctor Doctor
		var user User

		err := rows.Scan(
			&doctor.ID, &doctor.UserID, &doctor.Specialty, &doctor.ExperienceYears,
			&doctor.Education, &doctor.About, &doctor.ConsultationFee, &doctor.IsActive,
			&doctor.CreatedAt, &doctor.UpdatedAt,
			&user.Email, &user.FirstName, &user.LastName, &user.Phone,
		)
		if err != nil {
			return nil, err
		}

		user.ID = doctor.UserID
		user.UserType = "doctor"
		doctor.User = &user
		doctors = append(doctors, doctor)
	}

	return doctors, nil
}
