#!/bin/bash
cd infra && docker-compose -f docker-compose.demo.yaml up -d
sleep 15
cd ../orchestrator/
go run ./cmd/migrator/postgres  --migrations-path=./migrations

