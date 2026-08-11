.PHONY: setup install up down run test

VENV := .venv
PYTHON := $(VENV)/bin/python
PIP := $(VENV)/bin/pip
DATA_PATH ?= ../../data/data/raw_data_jan

setup:
	python3 -m venv $(VENV)
	$(PIP) install --upgrade pip
	$(PIP) install -r requirements.txt

install:
	$(PIP) install -r requirements.txt

up:
	docker compose up -d

down:
	docker compose down

run:
	DATA_PATH=$(DATA_PATH) PYTHONPATH=./exporter_app/src/ $(PYTHON) ./exporter_app/src/cmd/atmosdemo.py

test:
	PYTHONPATH=./exporter_app/src/ $(PYTHON) -m pytest

up-grafana:
	docker compose up -d grafana

up-influxdb:
	docker compose up -d influxdb