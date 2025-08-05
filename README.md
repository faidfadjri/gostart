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
## 💻 Required Dependencies
| Dependency                                                            | Description                                                 | Link                                                 |
| --------------------------------------------------------------------- | ----------------------------------------------------------- | ---------------------------------------------------- |
| [go-chi/chi](https://github.com/go-chi/chi)                               | Lightweight, idiomatic web router for Go.                   | [GitHub](https://github.com/go-chi/chi)              |
| [gorm](https://gorm.io/)                                              | The ORM for database interactions in Go.                    | [Website](https://gorm.io/)                          |
| [cosmtrek/air](https://github.com/cosmtrek/air)                                | Live reloading for Go applications.                         | [GitHub](https://github.com/cosmtrek/air)            |
| [godotenv](https://github.com/joho/godotenv)                          | Load environment variables from a `.env` file.              | [GitHub](https://github.com/joho/godotenv)           |
| [MySQL](https://www.mysql.com/)                                       | The primary database used in the project.                   | [Website](https://www.mysql.com/)                    |
| [go-chi/httprate](github.com/go-chi/httprate)              | Go Chi Rate Limiter                                         | [Github](github.com/go-chi/httprate)                 |
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
5. Congrats! Your API live now! Try to access it
   ```
   http://localhost:8000/
   ```

---

## 🛠️ CLI Commands

Use the following commands to generate boilerplate code:

```bash
# Generate sample folder structure
gostart init

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

## 📚 About Dev
I'm Mohamad Faid Fadjri, a dedicated and experienced fullstack developer with 3 years of experience building modern, scalable web apps and backend services.

🌐 Portfolio: https://faidfadjri.github.io/

✍️ Medium: https://medium.com/@faidfadjri

💼 LinkedIn: https://linkedin.com/in/faidfadjri


## 📄 License

This project is licensed under the **MIT License**.