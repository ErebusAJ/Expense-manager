# Expense Manager Backend

The project folder structure
```
expense-manager/
│
├── cmd/               # Entry points
│   └── main.go
│
├── config/            # Configuration files
│   └── config.go
│
├── internal/
│   ├── handlers/      # Request handlers
│   ├── models/        # Database models
│   ├── repository/    # Database queries
│   ├── services/      # Business logic
│   ├── middleware/    # Middleware functions
│   └── utils/         # Utility functions
│
├── pkg/               # Reusable packages
│
├── migrations/        # Database migration files
│
├── go.mod             # Dependency management
└── go.sum             # Dependency lock file

```

1. First we import/install packages like gin, godotenv 

After that open add key items to .env and .gitignore

Start up the server with gin

2. Now connecting DB

Write the main func in utils