# Ride Hauling Application

A modern, scalable ride-hailing platform built with Go, gRPC, and microservices architecture. This application provides a robust foundation for ride-hailing services with features like real-time trip tracking, driver management, and efficient routing.

## 🚀 Features

- **Microservices Architecture**
  - API Gateway for unified access
  - Driver Service for driver management
  - Trip Service for ride handling
  - Shared components and protocols

- **Real-time Trip Management**
  - Trip creation and tracking
  - Route optimization using OSRM
  - Real-time location updates
  - Trip status management

- **Driver Management**
  - Driver registration and verification
  - Real-time driver status
  - Location tracking
  - Ride assignment

- **Technical Features**
  - gRPC for efficient service communication
  - Graceful shutdown handling
  - Structured logging
  - Configuration management
  - Health checks and monitoring
  - Keepalive connection management

## 🏗️ Architecture

```
├── services/
│   ├── api-gateway/     # API Gateway service
│   ├── driver-service/  # Driver management service
│   └── trip-service/    # Trip management service
├── shared/
│   ├── pkg/            # Shared packages
│   ├── proto/          # Protocol buffer definitions
│   └── types/          # Common type definitions
├── web/                # Web interface
└── proto/             # Protocol definitions
```

## 🛠️ Prerequisites

- Go 1.21 or higher
- Protocol Buffers compiler
- Make (for build automation)
- Docker (optional, for containerization)

## 🚀 Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/ride-hauling-app.git
   cd ride-hauling-app
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure the application**
   - Copy `config.yaml.example` to `config.yaml`
   - Update the configuration values as needed

4. **Build the services**
   ```bash
   make build
   ```

5. **Run the services**
   ```bash
   # Run all services
   make run

   # Or run individual services
   make run-api-gateway
   make run-driver-service
   make run-trip-service
   ```

## 📝 Configuration

The application uses a YAML configuration file (`config.yaml`) with the following structure:

```yaml
api-gateway:
  port: 8080

driver-service:
  port: 8081

trip-service:
  port: 8082
```

## 🔧 Development

### Project Structure

- `services/`: Contains all microservices
  - `api-gateway/`: API Gateway service
  - `driver-service/`: Driver management service
  - `trip-service/`: Trip management service
- `shared/`: Shared components
  - `pkg/`: Common packages
  - `proto/`: Protocol buffer definitions
  - `types/`: Common type definitions
- `web/`: Web interface
- `proto/`: Protocol definitions

### Building

```bash
# Build all services
make build

# Build specific service
make build-api-gateway
make build-driver-service
make build-trip-service
```

### Running Tests

```bash
# Run all tests
make test

# Run specific service tests
make test-api-gateway
make test-driver-service
make test-trip-service
```

## 🔍 API Documentation

The application uses gRPC for service communication. API documentation can be found in the `shared/proto` directory:

- `trip/`: Trip service API definitions
- `driver/`: Driver service API definitions
- `rider/`: Rider service API definitions

## 🚀 Deployment

### Docker Deployment

```bash
# Build Docker images
make docker-build

# Run with Docker Compose
docker-compose up
```

### Kubernetes Deployment

```bash
# Apply Kubernetes configurations
kubectl apply -f k8s/
```

## 📊 Monitoring

The application includes built-in monitoring capabilities:

- Health check endpoints
- Structured logging
- Metrics collection
- Performance monitoring

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- Your Name - Initial work

## 🙏 Acknowledgments

- OSRM for routing capabilities
- gRPC for efficient service communication
- Go community for excellent tools and libraries 