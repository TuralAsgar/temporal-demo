# Temporal Lead Calling Workflow

[YouTube Video](https://youtu.be/qZ3l5aZ_sKM)

[Presentation](./temporal.pptx)

A demo Go application that orchestrates customer call workflows using Temporal. It manages the complete lifecycle of
lead outreach, from initiating calls via 3CX to updating lead records and logging interactions.

## Overview

This project demonstrates a distributed workflow system for lead management with the following capabilities:

- **Workflow Orchestration**: Uses Temporal to manage complex, long-running lead calling processes
- **Activity-based Architecture**: Separates concerns into discrete, reusable activities
- **Signal-driven Communication**: Receives call results asynchronously via webhooks
- **Retry Logic**: Implements exponential backoff with configurable retry policies
- **REST API**: Exposes endpoints to trigger workflows and receive webhook callbacks

## Architecture

### Components

- **Workflow** (`workflows.go`): Orchestrates the lead calling process with retry policies
- **Activities** (`activites.go`): Performs individual tasks (calling, updating records, logging)
- **API** (`api.go`): HTTP endpoints for workflow control and webhook handling
- **Worker** (`main.go`): Temporal worker that executes workflows and activities

### Workflow Flow

```
Start → Pick Lead → Call Customer → Wait for Signal → Update Lead → Create Interaction → End
```

## Prerequisites

- Go 1.25.5 or higher
- Temporal Server v1.39.0 running locally

## Installation

```shell
go mod download
```

```shell
brew install temporal
```

## How to Run the Demo (Live)

### 1️⃣ Start Temporal Cluster

```shell
temporal server start-dev
```

**UI**: http://localhost:8233

### 2️⃣ Start Worker + API

```shell
go run .
```

### 3️⃣ Start Workflow via HTTP

```shell
curl -X POST http://localhost:8080/start-call \
  -H "Content-Type: application/json" \
  -d '{
    "ID": "1",
    "Phone": "+994501234567",
    "Status": "NEW"
  }'
```

### 4️⃣ Simulate 3CX Webhook

```shell
curl -X POST http://localhost:8080/3cx/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "leadId": "1",
    "result": "ANSWERED"
  }'
```

## API Endpoints

### POST `/start-call`

Initiates a new lead calling workflow.

**Request:**

```json
{
  "ID": "lead-123",
  "Phone": "+1234567890",
  "Status": "active"
}
```

**Response:**

```json
{
  "workflowId": "lead-lead-123",
  "runId": "workflow-run-id"
}
```

### POST `/3cx/webhook`

Receives call results from 3CX system and signals the workflow.

**Request:**

```json
{
  "leadId": "lead-123",
  "result": "completed"
}
```

## Configuration

The application connects to Temporal using default settings. Configure the connection in `main.go` if needed:

```go
c, err := client.Dial(client.Options{
// Add custom options if needed
})
```

## Retry Policy

Activities are configured with the following retry behavior:

- **Initial Interval**: 1 second
- **Backoff Coefficient**: 2.0
- **Maximum Interval**: 10 seconds
- **Maximum Attempts**: 3
- **Non-retryable Errors**: `ValidationError`

## Project Structure

```
├── main.go              # Application entry point
├── workflows.go         # Workflow definitions
├── activites.go         # Activity implementations
├── api.go               # HTTP API handlers
├── types.go             # Data type definitions
├── go.mod               # Module dependencies
└── README.md            # This file
```

