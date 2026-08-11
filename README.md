# Atmos Demo

A telemetry replay and visualization environment for the Atmos system.

This project replays telemetry data from PKL files, exports the data to InfluxDB, and visualizes it with Grafana.

## Architecture

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
     Grafana
```

## Prerequisites

* Git
* Python 3.12 or later
* Docker
* Docker Compose

### Windows / WSL

Docker Desktop for Windows is required.

Enable **WSL 2 Integration** in Docker Desktop.

Docker Desktop:
https://www.docker.com/products/docker-desktop/

Verify the installation:

```bash
docker --version
docker compose version
```

### Linux

Install Docker Engine and Docker Compose.

```bash
sudo apt update
sudo apt install docker.io docker-compose-plugin
```

Verify:

```bash
docker --version
docker compose version
```

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd atmos-demo
```

### 2. Prepare telemetry data

Place the telemetry PKL files under the `data/` directory.

```text
data/
└── *.pkl
```

The telemetry data is not included in this repository.

### 3. Create the Python environment

```bash
make setup
```

This creates a Python virtual environment and installs the required dependencies.

Activate the environment:

```bash
source .venv/bin/activate
```

### 4. Start InfluxDB and Grafana

```bash
make up
```

Check the running containers:

```bash
docker compose ps
```


### 5. Configure Grafana

Open:

```text
http://localhost:3000

pass: admin
user: admin
```

Configure an InfluxDB data source with:

```text
Query language: Flux
URL:            
Organization:   atmos
Bucket:         telemetry
Token:          atmos-token
```

### 6. Run the Exporter

```bash
make run
```

The exporter reads telemetry data from the `data/` directory and writes it to InfluxDB.

The telemetry data can then be visualized in Grafana.

## Grafana Query Example
set timezone as 2024-01-01 13:46:47 to 2024-01-11 17:46:47


The following Flux query displays the `average` field:

```flux
from(bucket: "telemetry")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) =>
    r._measurement == "telemetry" and
    r._field == "average"
  )
  |> aggregateWindow(
    every: v.windowPeriod,
    fn: mean,
    createEmpty: false
  )
```


For the standard deviation:

```flux
from(bucket: "telemetry")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) =>
    r._measurement == "telemetry" and
    r._field == "std"
  )
  |> aggregateWindow(
    every: v.windowPeriod,
    fn: mean,
    createEmpty: false
  )
```

## Development

### Run tests

```bash
make test
```

### Stop services

```bash
make down
```


## Project Structure

```text
atmos_demo/
├── Makefile
├── README.md
├── docker-compose.yml
├── requirements.txt
├── data/
└── src/
    └── cmd/
        └── atmosdemo.py
```

## Make Commands

| Command        | Description                                                    |
| -------------- | -------------------------------------------------------------- |
| `make setup`   | Create the Python virtual environment and install dependencies |
| `make install` | Install Python dependencies                                    |
| `make up`      | Start InfluxDB and Grafana                                     |
| `make down`    | Stop InfluxDB and Grafana                                      |
| `make run`     | Start the Python Exporter                                      |
| `make test`    | Run unit tests                                                 |
| `make build`   | Build the Analysis Docker image                                |
| `make shell`   | Open a shell in the Analysis Docker image                      |

## Additional Commands

Start Grafana only:

```bash
make up-grafana
```

Start InfluxDB only:

```bash
make up-influxdb
```

Stop the services:

```bash
make down
```

To completely reset the environment, including InfluxDB and Grafana data:

```bash
docker compose down -v
```
