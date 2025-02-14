# Future Fit Take Home Project

## Overview

This project is designed to create a simple appointment scheduling API using Go and the Gin web framework. It allows users to manage appointments with trainers, check availability, and retrieve existing appointments.

## Project Structure

The project is organized as follows:

- `README.md`: This file provides an overview and documentation for the project.
- `.air.toml`: Configuration file for [Air](https://github.com/cosmtrek/air), a live reloading tool for Go applications. It specifies build and watch configurations.
- `.editorconfig`: Defines coding styles to maintain consistency across different editors and IDEs.
- `.gitignore`: Specifies files and directories that should be ignored by Git, such as temporary files and build artifacts.
- `bin/test_api.sh`: A shell script for testing the API endpoints using curl commands.
- `data/appointments.json`: Sample data file containing a list of appointments used to pre-load the application.
- `docker-compose.yml`: Defines the service configuration for Docker to run the application in a containerized environment.
- `go.mod` and `go.sum`: Go modules files that manage dependencies required by the project.
- `main.go`: The entry point of the application where the HTTP server is set up and routes are registered.
- `models/appointments.go`: Contains the `Appointment` struct and validation logic for appointments.
- `routes/routes.go`: Defines the API routes and their corresponding handler functions for managing appointments.
- `store/memory_store.go`: In-memory data store for appointments, handling loading, retrieving, and adding appointments.
- `utils/utils.go`: Contains utility functions for validating times.

## Getting Started

### Prerequisites

- Go (1.23 or higher) installed on your machine.
- Docker (optional) for running in a containerized environment.

### Running the Application

1. **Clone the Repository**:

   ```sh
   git clone <repository-url>
   cd <repository-directory>
   ```

2. **Install Dependencies**:
   Navigate to the project directory and run:

   ```sh
   go mod download
   ```

3. **Using Docker**:
   If you prefer running the application with Docker, execute:

   ```sh
   docker-compose up --build
   ```

   This command will build the Docker image and start the container.

4. **Run the Application**:
   If running locally without Docker, you can start the application with:
   ```sh
   go run main.go
   ```
   The application will start and listen on `http://localhost:8080`.

### API Endpoints

The following API endpoints are available:

- **GET** `/appointments`: Retrieve all appointments.
- **GET** `/appointments/:id`: Retrieve a specific appointment by ID.
- **GET** `/trainers/:id/availability?starts_at=<timestamp>&ends_at=<timestamp>`: Get availability for a trainer within the specified time range.
- **GET** `/trainers/:id/appointments`: Get all appointments for a specific trainer.
- **POST** `/trainers/:id/appointments`: Create a new appointment for a specific trainer (requires JSON body).

### Testing the API

You can test the API using the provided `bin/test_api.sh` script.

1. **Make it executable**:

   ```sh
   chmod +x bin/test_api.sh
   ```

2. **Run the test script**:
   ```sh
   ./bin/test_api.sh
   ```

The script will perform various GET and POST requests against the API and display the responses.

### Logging

Errors and build messages are logged in `build-errors.log`.

### Notes

- Ensure the timestamps used in the API requests follow the format specified in the API documentation (RFC3339).
- The application uses an in-memory store for appointments, which means data will not persist across application restarts.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more information.
