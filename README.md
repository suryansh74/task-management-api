# Task Management API with Redis Caching

A production-ready REST API built with Go, PostgreSQL, and Redis implementing advanced caching patterns and rate limiting.

## Features

- ✅ CRUD operations for tasks with PostgreSQL
- ✅ Redis-based caching (cache-aside pattern)
- ✅ Session-based authentication
- ✅ Rate limiting (100 req/min per user)
- ✅ Clean architecture (Ports & Adapters)
- ✅ Structured logging with Zerolog
- ✅ Docker Compose setup

## Tech Stack

- **Language:** Go 1.23
- **Framework:** Fiber v2
- **Database:** PostgreSQL 16
- **Cache:** Redis 7
- **Validation:** go-playground/validator

## Quick Start

Clone repository
git clone <your-repo>
cd task-management-api-project

Start with Docker Compose
docker-compose up -d --build

View logs
docker-compose logs -f app

API available at http://localhost:8000
text

## API Endpoints

### Authentication

- `POST /register` - Create new account
- `POST /login` - Login and create session
- `POST /logout` - Logout and destroy session

### Tasks (Protected)

- `GET /tasks` - Get all user tasks
- `POST /tasks` - Create new task
- `GET /tasks/:id` - Get task by ID (cached)
- `PUT /tasks/:id` - Update task (invalidates cache)
- `DELETE /tasks/:id` - Delete task (invalidates cache)

### Health

- `GET /check_health` - Service health check

## Architecture Highlights

- **Cache-Aside Pattern:** Read from cache → miss → fetch DB → populate cache
- **Cache Invalidation:** Updates/deletes automatically invalidate cache
- **Rate Limiting:** Redis-based token bucket (100 req/min)
- **Session Storage:** HTTP-only cookies with Redis backend
- **Ownership Policy:** Users can only access their own tasks

## Environment Variables

See `.env.docker` for Docker configuration or `.env` for local development.

## Development

Local development (requires PostgreSQL and Redis)
go run main.go

text

## Project Structure

internal/
├── clients/ # Database & Redis clients
├── handler/ # HTTP handlers
├── service/ # Business logic
├── repository/ # Data access layer
├── models/ # Domain models
├── ports/ # Interface definitions
└── config/ # Configuration management

text

## Learning Outcomes

This project demonstrates:

- Advanced Redis usage (caching, sessions, rate limiting)
- Clean architecture principles
- Production-ready error handling
- Structured logging practices
- Docker containerization
