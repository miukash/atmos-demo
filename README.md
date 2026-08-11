# Atmos Demo

Python Exporter for replaying telemetry data and visualizing it with InfluxDB and Grafana.

---

# Overview

Atmos Demo replays telemetry data from PKL files, converts the data into time-series metrics, stores them in InfluxDB, and visualizes them with Grafana.

---

# Architecture

```text
Telemetry Data (PKL)
        │
        ▼
Python Exporter
        │
        ▼
InfluxDB
        │
        ▼
Grafana Dashboard
```

---

# Prerequisites

- Git
- Python 3.12 or later
- Docker

---

# Setup

## 1. Install Docker

### Windows

Install **Docker Desktop for Windows**.

https://www.docker.com/products/docker-desktop/

After installation, ensure that Docker Desktop is running and **WSL 2 Integration** is enabled.

Verify the installation.

```bash
docker --version
docker compose version
```

### macOS

Install **Docker Desktop for Mac**.

https://www.docker.com/products/docker-desktop/

Verify the installation.

```bash
docker --version
docker compose version
```

### Linux (Ubuntu)

Install Docker Engine and Docker Compose.

```bash
sudo apt update
sudo apt install docker.io docker-compose-plugin
```

Verify the installation.

```bash
docker --version
docker compose version
```

---

## 2. Clone Repository

```bash
mkdir atmosdemo
cd atmosdemo
git clone <repository-url> src
cd src
```

---

## 3. Prepare Telemetry Data

Copy the telemetry data from the USB storage into the `data` directory.

```bash
scp -r /path/to/telemetry/ ../data
```

---

## 4. Create a Python Virtual Environment

### Linux / WSL

```bash
python -m venv .venv
source .venv/bin/activate
```

---

## 5. Install Python Dependencies

```bash
pip install -r requirements.txt
```

(Optional) Run unit tests.

```bash
PYTHONPATH=. pytest
```

---

## 6. Start InfluxDB and Grafana

### Using Docker Compose (Recommended)

```bash
docker compose up -d
```

Verify that the containers are running.

```bash
docker ps
```

Stop the containers.

```bash
docker compose down
```

### Without Docker Compose

#### Pull Images

```bash
docker pull influxdb:2
docker pull grafana/grafana:latest
```

#### Start InfluxDB

```bash
docker run -d \
  --name influxdb \
  -p 8086:8086 \
  influxdb:2
```

#### Start Grafana

```bash
docker run -d \
  --name grafana \
  -p 3000:3000 \
  grafana/grafana:latest
```

---

## 7. Run the Exporter

```bash
python cmd/atmosdemo.py
```

---

## 8. Access Grafana and InfluxDB

### Grafana

URL

```
http://localhost:3000
```

Default credentials

```
Username: admin
Password: admin
```

### InfluxDB

URL

```
http://localhost:8086
```

If this is your first launch, complete the initial InfluxDB setup.

---

## 9. Configure Grafana Dashboard

1. Add InfluxDB as a data source.
2. Import the dashboard from the `dashboards/` directory.
3. Start the exporter.
4. Confirm that telemetry data is displayed on the dashboard.

---

# Project Structure

```text
AtmosDemo/
├── src/
│   ├── README.md
│   ├── cmd/
│   │   ├── atmosdemo.py
│   │   └── atmosdemo_ns.py
│   ├── pkg/
│   │   ├── converter/
│   │   ├── infra/
│   │   ├── queue/
│   │   └── usecase/
│   ├── tests/
│   ├── requirements.txt
│   └── docker-compose.yml
└── data/
```