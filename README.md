# 🏥 Online Doctor Appointment System

A simple, efficient web-based platform for booking medical appointments online. Built with Go, PostgreSQL, and integrated with Kaspi payment system.

## 📋 Table of Contents

- [Features](#features)
- [Technology Stack](#technology-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Database Setup](#database-setup)
- [Running the Application](#running-the-application)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [Demo Accounts](#demo-accounts)
- [Screenshots](#screenshots)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

### For Patients
- 👤 User registration and login
- 🔍 Search and filter doctors by specialty
- 📅 Book appointments with real-time availability
- 💳 Flexible payment options (Kaspi Pay or Pay Later)
- 📊 View appointment history and status
- 📝 Add appointment notes

### For Doctors
- 📈 Dashboard with appointment statistics
- 📋 View and manage appointments
- ✅ Confirm, cancel, or complete appointments
- 👥 Access patient contact information
- 💼 Manage professional profile

### For Administrators
- 📊 System overview and statistics
- 👨‍⚕️ View all registered doctors
- 👥 View all registered patients
- 🔧 User management capabilities

## 🛠️ Technology Stack

**Backend:**
- Go (Golang) 1.21+
- Gorilla Mux (HTTP routing)
- PostgreSQL 13+
- bcrypt (password hashing)

**Frontend:**
- HTML5
- CSS3
- Vanilla JavaScript

**Payment Integration:**
- Kaspi Payment Gateway

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.21 or higher)
- [PostgreSQL](https://www.postgresql.org/download/) (version 13 or higher)
- Git

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/online-doctor-appointment.git
cd online-doctor-appointment
```

### 2. Install Go Dependencies

```bash
go mod download
```

### 3. Configure Environment Variables

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=doctor_appointment

# Server Configuration
PORT=8080
SECRET_KEY=your_secret_key_change_in_production

# Email Configuration (Optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
```

## 🗄️ Database Setup

### 1. Create Database

```bash
# Connect to PostgreSQL
sudo -u postgres psql

# Create database
CREATE DATABASE doctor_appointment;

# Exit PostgreSQL
\q
```

### 2. Run Database Migrations

```bash
# Run schema
psql -U postgres -d doctor_appointment -f sql/schema.sql

# Run seed data (optional - includes demo accounts)
psql -U postgres -d doctor_appointment -f sql/seed.sql
```

**Or use the automated setup script:**

```bash
chmod +x setup_db.sh
./setup_db.sh
```

## ▶️ Running the Application

### Development Mode

```bash
go run cmd/server/main.go
```

The application will be available at `http://localhost:8080`

### Production Build

```bash
# Build the application
go build -o bin/app cmd/server/main.go

# Run the built binary
./bin/app
```

## 📱 Usage

### Accessing the Application

1. Open your browser and navigate to `http://localhost:8080`
2. Register a new account or use demo accounts
3. Login and start using the platform

### Demo Accounts

The system includes pre-seeded demo accounts:

**Patient Account:**
- Email: `patient1@email.com`
- Password: `password123`

**Doctor Account:**
- Email: `dr.smith@hospital.com`
- Password: `password123`

**Admin Account:**
- Email: `admin@hospital.com`
- Password: `password123`

## 📁 Project Structure

```
online-doctor-appointment/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── database/
│   │   └── db.go                # Database connection
│   ├── handlers/
│   │   ├── auth.go              # Authentication handlers
│   │   ├── patient.go           # Patient handlers
│   │   ├── doctor.go            # Doctor handlers
│   │   └── admin.go             # Admin handlers
│   └── models/
│       ├── user.go              # User model
│       ├── doctor.go            # Doctor model
│       └── appointment.go       # Appointment model
├── static/
│   ├── css/
│   │   └── style.css            # Styles
│   └── js/
│       └── main.js              # JavaScript
├── sql/
│   ├── schema.sql               # Database schema
│   └── seed.sql                 # Sample data
├── .env.example                 # Environment template
├── .gitignore                   # Git ignore rules
├── go.mod                       # Go module file
├── go.sum                       # Go dependencies
└── README.md                    # This file
```

## 🔌 API Endpoints

### Public Routes
- `GET /` - Home page
- `GET /login` - Login page
- `POST /login` - Login handler
- `GET /register` - Registration page
- `POST /register` - Registration handler
- `POST /logout` - Logout handler

### Patient Routes (Protected)
- `GET /dashboard/patient` - Patient dashboard
- `GET /dashboard/patient/book` - Book appointment page
- `POST /dashboard/patient/book` - Submit booking
- `GET /dashboard/patient/appointments` - View appointments

### Doctor Routes (Protected)
- `GET /dashboard/doctor` - Doctor dashboard
- `GET /dashboard/doctor/appointments` - View appointments
- `POST /dashboard/doctor/appointment/:id/update` - Update status

### Admin Routes (Protected)
- `GET /dashboard/admin` - Admin dashboard
- `GET /dashboard/admin/doctors` - View all doctors
- `GET /dashboard/admin/patients` - View all patients

### API Endpoints
- `GET /api/doctors` - Get all doctors (JSON)
- `GET /api/doctors/:specialty` - Get doctors by specialty
- `GET /api/available-slots/:doctorId/:date` - Get available slots

## 🧪 Testing

### Run Tests
```bash
go test ./...
```

### Test Coverage
```bash
go test -cover ./...
```

### Manual Testing
```bash
# Test if server is running
curl http://localhost:8080

# Test API endpoints
curl http://localhost:8080/api/doctors
```

## 🌐 Access from Other Devices

### Local Network Access

1. Find your local IP:
```bash
# Linux/macOS
ifconfig | grep "inet "

# Windows
ipconfig
```

2. Access from other devices on same WiFi:
```
http://YOUR_IP:8080
```

### Public Access (Using ngrok)

```bash
# Download ngrok
curl -O https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-darwin-amd64.zip
unzip ngrok-v3-stable-darwin-amd64.zip

# Start tunnel
./ngrok http 8080
```

## 🎯 Target Audience

1. **Tech-Savvy Professionals (25-40)** - Quick online booking
2. **Busy Parents (30-45)** - Family healthcare management
3. **Senior Citizens (55+)** - Simple, accessible interface

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 🐛 Known Issues

- Language support limited to English (Kazakh/Russian planned)
- No mobile app (web-only)
- Basic payment integration (demo mode)

## 🔮 Future Enhancements

- [ ] Multi-language support (Kazakh, Russian)
- [ ] Video consultation feature
- [ ] SMS/Email notifications
- [ ] Reviews and ratings system
- [ ] Medical records management
- [ ] Mobile application (iOS/Android)
- [ ] Advanced search filters
- [ ] Doctor availability calendar view

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👥 Authors

- **Your Name** - [GitHub Profile](https://github.com/yourusername)

## 🙏 Acknowledgments

- Inspired by modern telemedicine platforms
- Built for Kazakhstan healthcare market
- Thanks to the Go and PostgreSQL communities

## 📞 Support

For support, email support@yourdomain.com or open an issue in the repository.

## 📊 Project Status

**Status:** ✅ MVP Complete

**Version:** 1.0.0

**Last Updated:** November 2024

---

Made with ❤️ for better healthcare access