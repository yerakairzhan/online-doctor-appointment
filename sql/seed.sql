-- Insert admin user
INSERT INTO users (email, password_hash, first_name, last_name, phone, user_type) VALUES
    ('admin@hospital.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Admin', 'User', '1234567890', 'admin');

-- Insert sample doctors
INSERT INTO users (email, password_hash, first_name, last_name, phone, user_type) VALUES
                                                                                      ('dr.smith@hospital.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'John', 'Smith', '1234567891', 'doctor'),
                                                                                      ('dr.johnson@hospital.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Sarah', 'Johnson', '1234567892', 'doctor'),
                                                                                      ('dr.brown@hospital.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Michael', 'Brown', '1234567893', 'doctor');

-- Insert sample patients
INSERT INTO users (email, password_hash, first_name, last_name, phone, user_type) VALUES
                                                                                      ('patient1@email.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Alice', 'Wilson', '1234567894', 'patient'),
                                                                                      ('patient2@email.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Bob', 'Davis', '1234567895', 'patient');

-- Insert doctor details
INSERT INTO doctors (user_id, specialty, experience_years, education, about, consultation_fee) VALUES
                                                                                                   (2, 'Cardiology', 10, 'MD from Harvard Medical School', 'Specialized in heart diseases and cardiovascular surgery', 150.00),
                                                                                                   (3, 'Dermatology', 8, 'MD from Johns Hopkins', 'Expert in skin conditions and cosmetic procedures', 100.00),
                                                                                                   (4, 'General Practice', 5, 'MD from State University', 'Family medicine and general health consultations', 80.00);

-- Insert doctor availability (Monday to Friday, 9 AM to 5 PM)
INSERT INTO doctor_availability (doctor_id, day_of_week, start_time, end_time) VALUES
-- Dr. Smith (Cardiology)
(1, 1, '09:00', '17:00'), -- Monday
(1, 2, '09:00', '17:00'), -- Tuesday
(1, 3, '09:00', '17:00'), -- Wednesday
(1, 4, '09:00', '17:00'), -- Thursday
(1, 5, '09:00', '17:00'), -- Friday

-- Dr. Johnson (Dermatology)
(2, 1, '10:00', '16:00'), -- Monday
(2, 2, '10:00', '16:00'), -- Tuesday
(2, 3, '10:00', '16:00'), -- Wednesday
(2, 4, '10:00', '16:00'), -- Thursday
(2, 5, '10:00', '16:00'), -- Friday

-- Dr. Brown (General Practice)
(3, 1, '08:00', '18:00'), -- Monday
(3, 2, '08:00', '18:00'), -- Tuesday
(3, 3, '08:00', '18:00'), -- Wednesday
(3, 4, '08:00', '18:00'), -- Thursday
(3, 5, '08:00', '18:00'), -- Friday
(3, 6, '09:00', '13:00'); -- Saturday