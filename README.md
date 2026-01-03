# Medical Service

Healthcare management microservice built with Go Kratos, featuring patient management, appointments, prescriptions, and medical records

## 🚀 Quick Start

```bash
git clone https://github.com/arm-1234/medical-service.git
cd medical-service

mysql -u root -p -e "CREATE DATABASE medical_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

make init
make wire
make run
```

## ✨ Features

### Core Services
- **Patients** - Register, update, search, medical history
- **Doctors** - Profile management, specializations, availability scheduling
- **Appointments** - Book, reschedule, cancel, conflict detection
- **Prescriptions** - Create, track, validity management
- **Medical Records** - Diagnosis tracking, visit history

### Technical Features
- ✅ Clean Architecture (Service → Handler → Repository)
- ✅ OpenTelemetry distributed tracing
- ✅ Protocol Buffers (gRPC + HTTP)
- ✅ Wire dependency injection
- ✅ GORM with MySQL

## 📦 Tech Stack

- **Framework**: [Go Kratos v2](https://go-kratos.dev/)
- **Database**: MySQL 8.0+ with GORM
- **API**: Protocol Buffers (dual gRPC/HTTP)
- **Observability**: OpenTelemetry (traces + logs)
- **DI**: Google Wire

## 📁 Structure

```
medical-service/
├── cmd/                          # Entry point + Wire
├── configs/config.yaml           # Configuration
├── internal/
│   ├── biz/                     # Business logic handlers
│   ├── data/                    # Repositories + entities
│   ├── service/                 # gRPC/HTTP service layer
│   ├── server/                  # Server setup
│   └── pkg/otel/                # OpenTelemetry utilities
└── third_party/                 # Proto dependencies
```

## ⚙️ Configuration

Edit `configs/config.yaml`:

```yaml
server:
  http:
    addr: 0.0.0.0:8000
  grpc:
    addr: 0.0.0.0:9000

data:
  database:
    driver: mysql
    source: root:password@tcp(127.0.0.1:3306)/medical_db?charset=utf8mb4&parseTime=True&loc=Local
```

## 🔌 API Examples

```bash
curl -X POST http://localhost:8000/v1/medical/patients \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone_number": "+1234567890",
    "date_of_birth": "1990-01-15",
    "gender": "MALE",
    "blood_group": "O_POSITIVE"
  }'

curl http://localhost:8000/v1/medical/patients/{patient_id}

curl -X POST http://localhost:8000/v1/medical/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": "...",
    "doctor_id": "...",
    "appointment_date": "2026-01-10",
    "appointment_time": "10:00",
    "consultation_type": "IN_PERSON",
    "reason_for_visit": "Regular checkup"
  }'
```
## 💻 Development

```bash
make init     # Install tools (first time)
make wire     # Generate DI files
make run      # Start service
make test     # Run tests
make fmt      # Format code
make tidy     # Tidy modules
make clean    # Clean artifacts
make dev      # Wire + Run (⭐ recommended)
make check    # Format + Test (before commit)
make help     # Show all commands
```

## 🐳 Docker

```bash
make docker-build
make docker-run
```

## 📊 Observability

### Distributed Tracing
Every request flows through:
1. **Service Layer** → Creates trace span, logs request
2. **Handler Layer** → Creates child span, context-aware logging
3. **Repository Layer** → Uses context for log correlation

### Log Example
```go
ctx, span := otel.Trace(ctx, "DoctorService.GetDoctor")
defer span.End()

doctor, err := h.repo.Get(ctx, id)
if err != nil {
    h.log.WithContext(ctx).Errorf("Failed to get doctor: %v", err)
    return nil, fmt.Errorf("failed to get doctor: %w", err)
}
```

All logs automatically include trace ID and span ID for correlation.

## 🏗️ Architecture

**3-Layer Clean Architecture:**

```
Service → Handler → Repository
  ↓         ↓          ↓
Proto    Business    Database
Layer     Logic       Access
```

- **Service** - gRPC/HTTP endpoints, request validation, tracing entry
- **Handler** - Business logic, validation, orchestration
- **Repository** - Data access, GORM operations, entity mapping
