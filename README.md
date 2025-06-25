# Internal Documentation

This directory contains internal implementation details of the task manager application.

## Overview

The task manager application is built using the Go programming language and utilizes various libraries and frameworks to provide a robust and scalable solution.

## Directory Structure

- `config`: Configuration files and settings for the application.
- `database`: Database-related code, including MongoDB interactions.
- `handlers`: Request handlers for the application's API endpoints.
- `middleware`: Middleware functions for authentication, logging, and other purposes.
- `models`: Data models and structs used throughout the application.
- `repositories`: Data access objects for interacting with the database.
- `utils`: Utility functions for various tasks, such as hashing and token generation.

## Configuration

The application uses a configuration file (`config/config.go`) to load settings from environment variables or a `.env` file.

## Database

The application uses a MongoDB database to store data. The database connection is established using the `database/mongodb.go` file.

## API Endpoints

The application provides several API endpoints for managing tasks, users, and authentication. These endpoints are defined in the `handlers` directory.

## Authentication

The application uses JSON Web Tokens (JWT) for authentication. The `middleware/auth.go` file contains the authentication middleware.
