# Job Queue and Task Scheduler

A Go HTTP service that accepts jobs, stores them in MongoDB, and processes them with background workers.

## Requirements

- Docker Desktop with Docker Compose, or Go 1.25 and a reachable MongoDB instance

## Run with Docker Compose

From the repository root, start the API and MongoDB:

```powershell
docker compose up --build
```

The API listens on port `8081`. To use a different host port, set `APP_PORT` before starting Compose:

```powershell
$env:APP_PORT = "8090"
docker compose up --build
```

MongoDB data is stored in the `mongo_data` Docker volume. Stop the services with:

```powershell
docker compose down
```

To also delete the stored MongoDB data, run:

```powershell
docker compose down --volumes
```

## Create a job

Send a `POST` request to `/jobs` with a JSON body:

```powershell
$body = @{
    title      = "Print a message"
    job_type   = "PRINT"
    job_status = "PENDING"
    priority   = "HIGH"
    payload    = @{
        message = "Hello from the job queue"
    }
} | ConvertTo-Json

Invoke-RestMethod `
    -Method Post `
    -Uri "http://localhost:8081/jobs" `
    -ContentType "application/json" `
    -Body $body
```

The `PRINT` handler writes the message to the application logs. View the logs with:

```powershell
docker compose logs -f app
```

## Job scheduling

Jobs submitted between 9:00 AM and 6:00 PM in the server's local timezone are queued immediately. Jobs submitted outside that window are scheduled for the next 9:00 AM.

The queue and scheduled jobs are held in memory. MongoDB stores the job record, but restarting the application loses queued or scheduled work that has not yet been processed.

## Configuration

The Compose setup configures the application with:

| Variable | Compose value | Purpose |
| --- | --- | --- |
| `PORT` | `8081` | Port used by the HTTP server inside the container |
| `MONGODB_URI` | `mongodb://mongo:27017` | MongoDB connection URI |
| `DATABASE_NAME` | `job_queue` | MongoDB database used for job records |
| `APP_PORT` | `8081` by default | Host port mapped to the application |

For running the application outside Docker, set `MONGODB_URI`, `DATABASE_NAME`, and optionally `PORT` in the environment before starting it.

## Run tests

```powershell
go test ./...
```
