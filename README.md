# GoStart 🚀

**GoStart** is a lightweight code generator tool to help you quickly build a clean and maintainable Go project structure following the **Hexagonal Architecture** pattern.

---

## 📁 Project Structure

```
src/
├── app/
│   ├── controllers/        # Business logic controllers
│   ├── usecases/           # Application use cases
│   └── config/             # Application configuration
├── infrastructure/
│   ├── middlewares/        # HTTP middlewares
│   ├── databases/          # Database connections and models
│   │   └── models/         # Database models
│   ├── repositories/       # Data access layer
│   └── services/           # External/internal services
└── interface/
    ├── handlers/           # HTTP handlers
    ├── request/            # Request DTOs and parsers
    └── response/           # Response DTOs and formatters
```

---

## ⚙️ Getting Started

1. Install Package
   ```bash
   go install github.com/faidfadjri/gostart@latest
   ```
2. Copy the `.env.example` file to `.env` and update the environment variables as needed.
3. Install the dependencies:
   ```bash
   go mod tidy
   ```
4. Run the application:
   ```bash
   air
   ```

---

## 🛠️ CLI Commands

Use the following commands to generate boilerplate code:

```bash
# Generate a new usecase
gostart create usecase <name>

# Generate a new repository
gostart create repository <name>

# Generate a new handler
gostart create handler <name>

# Generate a new feature it will generate: repository, usecase, handler
gostart create feature <name>
```

Replace `<name>` with your feature name (e.g., `user`, `task`, `auth`, etc).

---

## 📄 License

This project is licensed under the **MIT License**.
