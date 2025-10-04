# 📞 VoIP API with FreeSWITCH + Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![FreeSWITCH](https://img.shields.io/badge/FreeSWITCH-Compatible-green.svg)](https://freeswitch.org/)

A production-ready RESTful API built with Go and the Gin framework that integrates seamlessly with [FreeSWITCH](https://freeswitch.org/) via ESL (Event Socket Layer). This API provides complete call management, real-time monitoring, and comprehensive CDR (Call Detail Records) tracking for your VoIP infrastructure.

---

## ✨ Features

- 📲 **Call Initiation** - Originate calls between SIP users programmatically
- 🔄 **Real-time Monitoring** - Track active calls and receive hangup events via ESL
- 📊 **CDR Management** - Query and paginate through call detail records
- 🎯 **Event-Driven** - Built-in ESL event listener for `CHANNEL_HANGUP` and other events
- 🏗️ **Clean Architecture** - Organized code structure with separation of concerns
- ⚡ **High Performance** - Efficient Go concurrency for handling multiple calls
- 🔒 **Database Integration** - PostgreSQL support for CDR storage
- 📝 **Comprehensive Logging** - Detailed logging for debugging and monitoring

---

## 📋 Table of Contents

- [Architecture](#-architecture)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Configuration](#️-configuration)
- [API Documentation](#-api-documentation)
- [Project Structure](#-project-structure)
- [Usage Examples](#-usage-examples)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🏛️ Architecture

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Client    │◄───────►│   VoIP API   │◄───────►│ FreeSWITCH  │
│ (REST API)  │         │   (Go/Gin)   │         │    (ESL)    │
└─────────────┘         └──────┬───────┘         └─────────────┘
                               │
                               ▼
                        ┌──────────────┐
                        │  PostgreSQL  │
                        │    (CDRs)    │
                        └──────────────┘
```

---

## 🔧 Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.21 or higher ([Download](https://golang.org/dl/))
- **FreeSWITCH** 1.10+ with ESL enabled ([Installation Guide](https://freeswitch.org/confluence/display/FREESWITCH/Installation))
- **PostgreSQL** 12+ ([Download](https://www.postgresql.org/download/))
- **Git** for version control

### FreeSWITCH Configuration

Ensure ESL is enabled in your FreeSWITCH configuration:

```xml
<!-- /etc/freeswitch/autoload_configs/event_socket.conf.xml -->
<configuration name="event_socket.conf" description="Socket Client">
  <settings>
    <param name="nat-map" value="false"/>
    <param name="listen-ip" value="127.0.0.1"/>
    <param name="listen-port" value="8021"/>
    <param name="password" value="ClueCon"/>
  </settings>
</configuration>
```

---

## 📦 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Vishaltalsaniya-7/voip-api.git
cd voip-api
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up Database

Create a PostgreSQL database and run migrations (if applicable):

```bash
psql -U postgres -c "CREATE DATABASE fusionpbx;"
```

If using FusionPBX, follow their [installation guide](https://www.fusionpbx.com/app/www/installation.php) for database setup.

### 4. Configure Environment Variables

Copy the example environment file and customize it:

```bash
cp .env.example .env
```

Edit `.env` with your configuration (see [Configuration](#️-configuration) section).

### 5. Run the Application

```bash
go run main.go
```

The API server will start on `http://localhost:8080` (or your configured port).

---

## ⚙️ Configuration

Create a `.env` file in the project root with the following variables:

```env
# Database Configuration
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=fusionpbx
DB_PASSWORD=your_secure_password
DB_NAME=fusionpbx
DB_SSLMODE=prefer

# FreeSWITCH ESL Configuration
FREESWITCH_HOST=127.0.0.1
FREESWITCH_PORT=8021
FREESWITCH_PASSWORD=ClueCon

# Server Configuration
SERVER_PORT=8080
GIN_MODE=release  # Use 'debug' for development

# Optional: Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Configuration Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `DB_HOST` | PostgreSQL server address | `127.0.0.1` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database username | `fusionpbx` |
| `DB_PASSWORD` | Database password | *required* |
| `DB_NAME` | Database name | `fusionpbx` |
| `FREESWITCH_HOST` | FreeSWITCH server address | `127.0.0.1` |
| `FREESWITCH_PORT` | FreeSWITCH ESL port | `8021` |
| `FREESWITCH_PASSWORD` | ESL password | `ClueCon` |
| `SERVER_PORT` | API server port | `8080` |

---

## 📡 API Documentation

### Base URL
```
http://localhost:8080
```

### Authentication
Currently, the API does not require authentication. For production use, implement JWT or API key authentication.

---

### 📞 Initiate Call

Originates a new call between two SIP users.

**Endpoint:** `POST /call`

**Request Body:**
```json
{
  "caller": "1001",
  "callee": "1002"
}
```

**Response (200 OK):**
```json
{
  "call_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "Call initiated"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "caller and callee are required"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/call \
  -H "Content-Type: application/json" \
  -d '{
    "caller": "1001",
    "callee": "1002"
  }'
```

---

### 📊 Get Call Status

Retrieve the current status of a call.

**Endpoint:** `GET /call/:uuid`

**URL Parameters:**
- `uuid` (string, required) - The unique call identifier

**Response (200 OK):**
```json
{
  "call_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "ACTIVE",
  "duration": 45,
  "caller": "1001",
  "callee": "1002"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/call/550e8400-e29b-41d4-a716-446655440000
```

---

### 📁 Get CDRs (Call Detail Records)

Retrieve paginated call detail records.

**Endpoint:** `GET /cdrs`

**Query Parameters:**
- `page` (integer, optional) - Page number (default: 1)
- `limit` (integer, optional) - Records per page (default: 10, max: 100)

**Response (200 OK):**
```json
{
  "cdrs": [
    {
      "xml_cdr_uuid": "550e8400-e29b-41d4-a716-446655440000",
      "caller_id_number": "1001",
      "destination_number": "1002",
      "start_stamp": "2024-01-15T10:30:00Z",
      "answer_stamp": "2024-01-15T10:30:05Z",
      "end_stamp": "2024-01-15T10:35:00Z",
      "duration": 300,
      "billsec": 295,
      "hangup_cause": "NORMAL_CLEARING",
      "direction": "outbound"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 150
  }
}
```

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/cdrs?page=1&limit=10"
```

---

### 🔍 Advanced CDR Filtering (Future Enhancement)

```bash
GET /cdrs?caller=1001&start_date=2024-01-01&end_date=2024-01-31&status=completed
```

---

## 📂 Project Structure

```
voip-api/
├── config/                 # Configuration management
│   └── config.go          # Load and parse .env file
├── controller/            # HTTP request handlers
│   ├── call.go           # Call initiation controller
│   └── cdr.go            # CDR retrieval controller
├── database/              # Database connection
│   └── database.go       # PostgreSQL initialization
├── manager/               # Business logic layer
│   └── esl.go            # FreeSWITCH ESL management
├── middleware/            # HTTP middleware (future)
│   ├── auth.go           # Authentication middleware
│   └── logger.go         # Request logging
├── models/                # Data models
│   └── cdr.go            # CDR database model
├── request/               # API request DTOs
│   └── call.go           # Call request struct
├── response/              # API response DTOs
│   └── cdr.go            # CDR response struct
├── utils/                 # Utility functions
│   └── helpers.go        # Common helper functions
├── tests/                 # Unit and integration tests
│   ├── controller_test.go
│   └── manager_test.go
├── .env.example           # Example environment file
├── .gitignore            # Git ignore rules
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── main.go               # Application entry point
├── Makefile              # Build and deployment tasks
├── Dockerfile            # Docker container config
├── docker-compose.yml    # Docker Compose setup
└── README.md             # This file
```

---

## 💡 Usage Examples

### Example 1: Making a Call

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func makeCall() {
    url := "http://localhost:8080/call"
    payload := map[string]string{
        "caller": "1001",
        "callee": "1002",
    }
    
    jsonData, _ := json.Marshal(payload)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
}
```

### Example 2: Fetching CDRs with Python

```python
import requests

response = requests.get('http://localhost:8080/cdrs', params={
    'page': 1,
    'limit': 50
})

cdrs = response.json()
for cdr in cdrs['cdrs']:
    print(f"Call from {cdr['caller_id_number']} to {cdr['destination_number']}")
```

### Example 3: Real-time Call Monitoring

```javascript
// Using JavaScript/Node.js
const axios = require('axios');

async function monitorCall(callId) {
    try {
        const response = await axios.get(`http://localhost:8080/call/${callId}`);
        console.log('Call Status:', response.data);
    } catch (error) {
        console.error('Error:', error.message);
    }
}

setInterval(() => monitorCall('call-uuid-here'), 5000);
```

---

## 🛠️ Development

### Running in Development Mode

```bash
# Enable debug mode
export GIN_MODE=debug
go run main.go
```

### Building the Application

```bash
# Build for current platform
go build -o voip-api main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o voip-api-linux main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o voip-api.exe main.go
```

### Using Makefile

```bash
# Build the application
make build

# Run tests
make test

# Run the application
make run

# Clean build artifacts
make clean
```

---

## 🧪 Testing

### Run All Tests

```bash
go test ./... -v
```

### Run Tests with Coverage

```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Individual Packages

```bash
go test ./controller -v
go test ./manager -v
```

---

## 🐳 Deployment

### Using Docker

**Build Docker Image:**
```bash
docker build -t voip-api:latest .
```

**Run Container:**
```bash
docker run -d \
  -p 8080:8080 \
  --env-file .env \
  --name voip-api \
  voip-api:latest
```

### Using Docker Compose

```bash
docker-compose up -d
```

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  voip-api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - FREESWITCH_HOST=freeswitch
    depends_on:
      - postgres
      - freeswitch
    restart: unless-stopped

  postgres:
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: fusionpbx
      POSTGRES_USER: fusionpbx
      POSTGRES_PASSWORD: secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  freeswitch:
    image: drachtio/drachtio-freeswitch-mrf:latest
    ports:
      - "8021:8021"
    restart: unless-stopped

volumes:
  postgres_data:
```

### Systemd Service (Linux)

Create `/etc/systemd/system/voip-api.service`:

```ini
[Unit]
Description=VoIP API Service
After=network.target postgresql.service

[Service]
Type=simple
User=voip
WorkingDirectory=/opt/voip-api
EnvironmentFile=/opt/voip-api/.env
ExecStart=/opt/voip-api/voip-api
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**Enable and Start:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable voip-api
sudo systemctl start voip-api
sudo systemctl status voip-api
```

---

## 🔧 Troubleshooting

### Common Issues

#### 1. Cannot Connect to FreeSWITCH ESL

**Error:** `Failed to connect to FreeSWITCH: dial tcp 127.0.0.1:8021: connect: connection refused`

**Solution:**
- Verify FreeSWITCH is running: `systemctl status freeswitch`
- Check ESL is enabled in `event_socket.conf.xml`
- Ensure port 8021 is not blocked by firewall

#### 2. Database Connection Failed

**Error:** `Failed to ping database: pq: password authentication failed`

**Solution:**
- Verify PostgreSQL credentials in `.env`
- Check PostgreSQL is accepting connections
- Ensure database exists: `psql -U postgres -c "\l"`

#### 3. Calls Not Initiating

**Error:** `No Call-ID returned from FreeSWITCH`

**Solution:**
- Check SIP users exist and are registered
- Verify originate command syntax
- Check FreeSWITCH logs: `tail -f /var/log/freeswitch/freeswitch.log`

### Debugging

Enable verbose logging:

```bash
export LOG_LEVEL=debug
export GIN_MODE=debug
go run main.go
```

Check FreeSWITCH console:
```bash
fs_cli -x "console loglevel debug"
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/AmazingFeature`)
3. **Commit your changes** (`git commit -m 'Add some AmazingFeature'`)
4. **Push to the branch** (`git push origin feature/AmazingFeature`)
5. **Open a Pull Request**

### Code Style

- Follow Go best practices and conventions
- Use `gofmt` for formatting
- Add comments for exported functions
- Write unit tests for new features

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [FreeSWITCH](https://freeswitch.org/) - Open source telephony platform
- [Gin Framework](https://gin-gonic.com/) - HTTP web framework
- [go-eventsocket](https://github.com/fiorix/go-eventsocket) - FreeSWITCH ESL client
- [FusionPBX](https://www.fusionpbx.com/) - Multi-tenant PBX system

---

## 👨‍💻 Author

**Vishal Talsaniya**

- GitHub: [@Vishaltalsaniya-7](https://github.com/Vishaltalsaniya-7)
- LinkedIn: [Vishal Talsaniya](https://www.linkedin.com/in/vishal-talsaniya/)
- Email: vishaltalsaniya991@gmail.com

---

## 📞 Support

If you encounter any issues or have questions:

1. Check the [Troubleshooting](#-troubleshooting) section
2. Search [existing issues](https://github.com/Vishaltalsaniya-7/voip-api/issues)
3. Create a [new issue](https://github.com/Vishaltalsaniya-7/voip-api/issues/new)

---

## 🗺️ Roadmap

- [ ] JWT Authentication
- [ ] WebSocket for real-time call updates
- [ ] Call recording management
- [ ] SMS integration
- [ ] Advanced analytics dashboard
- [ ] Multi-tenant support
- [ ] Rate limiting
- [ ] GraphQL API

---

## 📊 Project Stats

![GitHub stars](https://img.shields.io/github/stars/Vishaltalsaniya-7/voip-api?style=social)
![GitHub forks](https://img.shields.io/github/forks/Vishaltalsaniya-7/voip-api?style=social)
![GitHub issues](https://img.shields.io/github/issues/Vishaltalsaniya-7/voip-api)
![GitHub pull requests](https://img.shields.io/github/issues-pr/Vishaltalsaniya-7/voip-api)

---

<div align="center">
  <p>Made with ❤️ by Vishal Talsaniya</p>
  <p>⭐ Star this repo if you find it helpful!</p>
</div>