# RDEV DDGEN (Domain Directory Generate)

[![Go Report Card](https://goreportcard.com/badge/github.com/XaiPhyr/rdev-go-ddgen)](https://goreportcard.com/report/github.com/XaiPhyr/rdev-go-ddgen)
[![GitHub release (latest by SemVer)](https://img.shields.io/github/v/release/XaiPhyr/rdev-go-ddgen?logo=github&color=blue)](https://github.com/XaiPhyr/rdev-go-ddgen/releases)
[![Build Status](https://github.com/XaiPhyr/rdev-go-ddgen/actions/workflows/go.yml/badge.svg)](https://github.com/XaiPhyr/drdev-go-dgen/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`ddgen` is a lightweight, zero-dependency Go command-line tool designed to eliminate the tedious copy-pasting required when scaffolding a Domain-Driven, Clean Architecture backend. It enforces structural consistency across your team, letting you skip the boilerplate and dive straight into writing business logic.

## Why `ddgen`?

* **Zero Boilerplate:** Initialize an entire Go API layout with a single command.
* **Domain-Driven Generation:** Spin up modular, self-contained domain features (`handler`, `service`, `repository`) instantly.
* **Go Idiomatic Architecture:** Enforces separation of concerns by cleanly isolating your transport, business logic, and database access layers.

---

## Installation

Install the binary globally to your computer using the Go toolchain:

```bash
go install github.com/XaiPhyr/rdev-go-ddgen@latest
```

> *Note: Ensure your environment's `$GOPATH/bin` is added to your system's `PATH` variable to execute the command from any directory.*

---

## Commands & Usage

### 1. Initialize a New Project Structure

Bootstrap the entire core foundation layout for your new backend repository in your current working directory.

```bash
rdev-go-ddgen init
```

### 2. Generate a New Domain Feature Layer

Scaffold a complete, isolated business domain folder containing pre-configured boilerplate layers.

```bash
rdev-go-ddgen -d orders
```

*Replace `orders` with any domain concept (e.g., `users`, `products`, `payments`).*

---

## Generated Folder Structure

Running `ddgen init` and adding domains creates a predictable, production-ready Go project layout:

```text
.
├── cmd/                      # Main entry point for building and running the binary
├── internal/
│   ├── config/               # Configuration management; parses environment variables
│   ├── db/                   # Database instance connection pools
│   │   └── migrations/       # SQL files tracking schema history over time
│   ├── middleware/           # Cross-cutting HTTP middleware (Auth, Logging, CORS)
│   ├── shared/               # Global utilities shared across multiple features
│   │   ├── dto/              # System-wide Data Transfer Objects
│   │   ├── helpers/          # Common utility helper functions
│   │   └── models/           # Global database/domain models
│   ├── orders/               # Self-contained Business Domain (Example output from: ddgen -d orders)
│   │   ├── handler.go        # HTTP routers / Controllers mapping endpoints to business logic
│   │   ├── service.go        # Pure core business rules and use-case logic validation
│   │   ├── repository.go     # Direct database access queries and infrastructure storage drivers
│   │   └── types.go          # Domain-specific request/response struct shapes
│   └── templates/            # Static HTML components (e.g., transactional email layouts)
└── go.mod

```

---

## Architecture Flow

The generated code layout follows a strict **dependency rule**: control moves inward from the client to your database layers.

```text
[Client Request] ──> [handler.go] ──> [service.go] ──> [repository.go] ──> [Database]
```

* **Handler:** Handles HTTP parsing, reads input DTOs, and forwards requests to the service.
* **Service:** The brain of your application. Does not care if you use HTTP or gRPC; it only handles core logic rules.
* **Repository:** The boundary layer. Reads and writes data to your storage infrastructure (PostgreSQL, MySQL, etc.).

---
