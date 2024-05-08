# University Canteen Management System

## Overview
This project is structured as a full-stack Go application, integrating frontend and backend services, and designed to be containerized using Docker. It employs a modular architecture with separation of concerns among user interface, API layer, business logic, and utility functions.

## Project Structure

### cmd
- **myapp**: Contains the entry point for the Go application. This directory typically includes the main function initializing and starting the server.

### frontend
- **public**: Houses static files such as images, JavaScript, and CSS that are required by the web application.
- **templates**: Contains HTML templates which are rendered by the backend. This supports a dynamic user interface that interacts with backend services.

### internal
- **api**: Defines the API endpoints and associated handlers. This module acts as the interface through which the frontend communicates with the backend.
- **service**: Implements the business logic of the application. It processes data received from the API layer and interacts with models to perform CRUD operations.
- **structs**: Defines various data structures and models used across the application, ensuring data integrity and providing a blueprint for data manipulation.

### pkg
- **server**: Sets up and configures the HTTP server, including routes, middleware, and API endpoints.
- **utils**: Includes utility functions like data validators, loggers, and other helper functions that are reused across the project.

## Setup Instructions

### Requirements
- Docker
- Docker Compose
- Go (version specified in `go.mod`)

### Running the Application
1. Clone the repository: 
  ```bash
  git clone https://github.com/Arshynbeek/GO-Final.git
  ```

2. Navigate to the project directory: 
  ```bash
  cd GO-Final
  ```

3. Build and run the Docker containers: 
  ```bash
  docker-compose up --build
  ```

This command builds the application as defined in the `Dockerfile` and `docker-compose.yml`, and runs it along with its dependencies.

## Usage
After launching the application, you can access the web interface through your browser at `http://localhost:[port]`. Replace `[port]` with the actual port number configured in the Docker or Go settings (in this case: http://localhost:2024).