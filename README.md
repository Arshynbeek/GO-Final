# University Canteen Management System

## Overview
The project entails the development of an extensive university cafeteria management system utilizing the Go programming language. This system will explore multiple facets of software development in Go, including an introduction to Go's syntax, Docker, interaction with Go databases, and API integration, the creation of web applications using Go

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

## Documentation
### API Endpoints
The API provides a range of endpoints for managing users, products, and reports:
  - `/api/v1/signup/`: Register a new user.
  - `/api/v1/signin/`: Authenticate a user.
  - `/api/v1/signout/`: Log out a user.
  - `/api/v1/profile/`: Retrieve and update user profiles.
  - `/api/v1/new/`: Create new product.
  - `/api/v1/edit-product/`: Update existing product.
  - `/api/v1/delete/`: Delete product.
  - `/api/v1/buy/`, `/api/v1/add/`, `/api/v1/remove/`: Manage shopping cart and purchase items.
  - `/api/v1/feedback/`: Submit feedback for products.
  - `/api/v1/search/`: Search for products by name.
  - `/api/v1/reports/`: Generate sales, inventory, preferences, and revenue reports.

### Data Models
Defines the structure of data for categories, users, orders, feedback, and products:
  - **User**: Includes fields like name, username, email, password, and admin status.
  - **Food**: Represents food items with details like name, description, category, price, and quantity.
  - **Order**: Details of user orders including status and quantity.
  - **Feedback**: User feedback with ratings and comments linked to food items.
  - **Category**: Food categories to organize menu items.

### External Dependencies
  - **Gin**: HTTP web framework used for handling requests.
  - **GORM**: ORM library for Go, used for database interactions.
  - **JSON Web Token (JWT)**: Used to transmit data for authentication in client-server applications.
  - **Crypto**: For secure password hashing.
  - **PostgreSQL**: An advanced, enterprise-class open-source relational database.
  - **Docker**: Container platform used to encapsulate application environment.

## Setup Instructions

### Requirements
- Docker
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
After launching the application, you can access the web interface through your browser at `http://localhost:[port]`. Replace `[port]` with the actual port number configured in the Docker or Go settings.