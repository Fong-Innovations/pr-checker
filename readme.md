# Go App

This is a Go-based web application built with the Gin framework. The application can be run in a Dockerized environment or locally with live reload functionality for development.

## Getting Started

### Running with Docker

To build and run the application using Docker, follow these steps:

1. Build the Docker image:
   ```bash
   docker build -t gin-app .
   ```

### Running Locally with Live Reload

For local development, the application uses the [Air](https://github.com/cosmtrek/air) library to enable live reloading. Follow these steps to run the application locally:

1. Install the Air library if you haven't already:

   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. Run the application with Air:
   ```bash
   air
   ```

This will start the application and automatically reload it whenever you make changes to the source code.
