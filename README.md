# POS Microservices System

A production-ready Point of Sale/Venue management system built with Go microservices architecture, showcasing modern backend development practices and distributed system design.

## Architecture Overview

**Multi-tenant microservices** with JWT-based authentication and NGINX gateway routing. Each service maintains its own PostgreSQL database and communicates via HTTP APIs (There were plans on refactoring service-service communication to use gRPC later on, but it wasn't a priority)

<img width="1251" height="1011" alt="System Architecture" src="https://github.com/user-attachments/assets/fe287c61-ea2f-49f4-b248-c806e293b8c8" />

## Core Services

### 🔐 Auth Service
- **JWT Authentication** - 3-minute access tokens + 7-day refresh tokens (configurable)
- **Multi-tenant Claims** - Tenant ID embedded in JWT for data isolation
- **Token Management** - Automated cleanup job for expired tokens (verification/password-reset etc)
- **Email Integration** - Password reset & verification flows
- **Secure Hashing** - bcrypt for password storage

### 🏢 Venue Service  
- **Venue Management** - CRUD operations with soft delete/archive
- **Item Management** - Product catalog with full lifecycle support
- **Tenant Isolation** - Data segregation via JWT middleware
- **RESTful Design** - Clean API endpoints with proper HTTP verbs

### 📧 Communication Service
- **Email Service** - SMTP integration with gomail
- **Template System** - Verification, password reset, welcome emails
- **Service Decoupling** - Async communication pattern

## Technical Implementation

### **Go Architecture Patterns**
- **Repository Pattern** - Clean data layer abstraction
- **Service Layer** - Business logic separation
- **Dependency Injection** - Constructor-based DI
- **Middleware Chain** - Authentication, logging, recovery
- **Context Timeouts** - Graceful request handling

### **Database Design**
- **Connection Pooling** - pgx/v5 with configurable pool sizes
- **ACID Compliant** - although no transactions used
- **Migration Ready** - Structured schema management
- **Multi-database** - Service-per-database pattern

### **Security & Reliability**
- **JWT Middleware** - Request validation with tenant extraction  
- **Input Validation** - go-playground/validator integration
- **Password Security** - bcrypt hashing with salt
- **Timeout Management** - Context-based request timeouts
- **Error Handling** - Structured error responses

### **DevOps & Deployment**
- **Docker Compose** - Multi-service orchestration
- **NGINX Gateway** - Load balancing and routing
- **Environment Config** - 12-factor app compliance
- **Health Checks** - Service monitoring endpoints
- **Hot Reload** - Development-optimized containers

## Tech Stack

**Backend:** Go 1.24, Chi Router, pgx PostgreSQL driver  
**Database:** PostgreSQL with connection pooling, golang migrate cli tool for migrations  
**Authentication:** JWT with HMAC-SHA256 signing  
**Communication:** HTTP REST APIs, SMTP email, JSON
**Infrastructure:** Docker, NGINX, Docker Compose  
**Libraries:** validator, gomail, bcrypt, envconfig

## Quick Start

```bash
# Clone and start all services
git clone <repository>
cd pos-microservices

# for email communication, setup example is in .env.example in the communications service
# you can get env details at your own gmail google account, its free
# make sure to have `communication-service/.env` available otherwise you won't get email communications

docker-compose up --build

# Services available at by default:
# Gateway: http://localhost:8080
# Auth: http://localhost:8082  
# Venue: http://localhost:8081
# Communication: http://localhost:8083
```

## API Endpoints

You can find these in the `routes.go` file in each respective service

**Authentication:** Bearer token required for all `/api/*` endpoints (except auth) - NOTE: they expire very fast, so if you want to test it out, I recommend changing some configs in docker compose to make access token last longer 
