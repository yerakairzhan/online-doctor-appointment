-- Users table (both patients and doctors)
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       first_name VARCHAR(100) NOT NULL,
                       last_name VARCHAR(100) NOT NULL,
                       phone VARCHAR(20),
                       user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('patient', 'doctor', 'admin')),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Doctors table (additional info for doctors)
CREATE TABLE doctors (
                         id SERIAL PRIMARY KEY,
                         user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                         specialty VARCHAR(100) NOT NULL,
                         experience_years INTEGER DEFAULT 0,
                         education TEXT,
                         about TEXT,
                         consultation_fee DECIMAL(10,2) DEFAULT 0.00,
                         is_active BOOLEAN DEFAULT true,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Doctor availability
CREATE TABLE doctor_availability (
                                     id SERIAL PRIMARY KEY,
                                     doctor_id INTEGER REFERENCES doctors(id) ON DELETE CASCADE,
                                     day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 0 AND 6), -- 0=Sunday, 6=Saturday
                                     start_time TIME NOT NULL,
                                     end_time TIME NOT NULL,
                                     is_active BOOLEAN DEFAULT true,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Appointments table
CREATE TABLE appointments (
                              id SERIAL PRIMARY KEY,
                              patient_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                              doctor_id INTEGER REFERENCES doctors(id) ON DELETE CASCADE,
                              appointment_date DATE NOT NULL,
                              appointment_time TIME NOT NULL,
                              status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed')),
                              notes TEXT,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              UNIQUE(doctor_id, appointment_date, appointment_time)
);


-- Create table for chat logs (optional)
CREATE TABLE chat_logs (
                           id SERIAL PRIMARY KEY,
                           user_id INTEGER REFERENCES users(id),
                           message TEXT,
                           response TEXT,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_type ON users(user_type);
CREATE INDEX idx_doctors_specialty ON doctors(specialty);
CREATE INDEX idx_appointments_date ON appointments(appointment_date);
CREATE INDEX idx_appointments_patient ON appointments(patient_id);
CREATE INDEX idx_appointments_doctor ON appointments(doctor_id);

-- Update timestamp function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_doctors_updated_at BEFORE UPDATE ON doctors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_appointments_updated_at BEFORE UPDATE ON appointments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();