# Go Project

This project was generated using gostart.

## Project Structure

```
src/
├── app/
│   ├── controllers/         # Business logic controllers
│   ├── usecases/           # Application use cases
│   └── config/             # Application configuration
├── infrastructure/
│   ├── middlewares/        # HTTP middlewares
│   ├── databases/          # Database connections and models
│   │   └── models/         # Database models
│   ├── repositories/       # Data access layer
│   └── services/          # External services
└── interface/
    ├── handlers/           # HTTP request handlers
    ├── request/           # Request DTOs and parsers
    └── response/          # Response DTOs and helpers
```

## Getting Started

1. Copy `.env.example` to `.env` and configure your environment variables
2. Install dependencies: `go mod tidy`
3. Run the application: `go run src/main.go`

## Commands

Generate new components using gostart:

```bash
# Generate usecase
gostart create usecase user

# Generate repository  
gostart create repository user
```

## License

This project is licensed under the MIT License.
