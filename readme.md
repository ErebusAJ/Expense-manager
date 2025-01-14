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

## ***1. Import/Install Packages***

After that open add key items to .env and .gitignore

Start up the server with gin

## ***2. DB Connection***
* **utils/db.go**
    - we write the function to connect db like sql.Open and testing the connection via ping and return the connection or error if any

* **cmd/main.go**
    - call the db.go connect function

* **handling schema/queries**
    - to handle creation of schema and queries we use sqlc and goose
* **migrations/**
    - **schema/** : create users table and the run goose up to migrate
    - **queries/** : write the queries like insert/create user retreive user and run sqlc generate

* **handler/handler_users.go**
    - to handle the queries we write handlers for which first open up db connection in **routes.go** and create a apiConfig which store the new DB connection
    - now write the add user function here we have also created hashPassword which hash incoming password from json object and create a user

## ***3. AuthMiddleware***
* **migrations/schema**: 
    - First we need to alter our users table to be able to store token_version, last_logged in so we create a auth schema to alter table and add colums via goose
* **middleware/middleware_auth**:
    - returns a **gin.HandlerFunc** we parse the Headers get **"Authorization"** which includes Bearer token
    - After parsing check for Bearer prefix and then parse the JWT token
    - Create a jwt map claim to send userid key value pair to next func as context
    - parse userid to UUID

## ***4. Authenticated handlers***
* **handlers/handler_users.go**:
    - **loginUser** We write a function to login user where we verify the email and password_hash associated to it if successfull then generate a JWT token.

    - **getAuthUser** This fucntion gets the userID from the middleware which got it from JWT token with secret signed key then we call the sqlc query function **GetUserByID** and pass on the userID and send it back to client


* **utils/jwt.go**:
    - Here we have a fucntion **GenerateJWT** to generate a JWT auth token with a secret signed key once the user is verified in **loginUser** handler
    - Create a **jwt.MapClaims** assign the userID to userID recieved as a parameter and two other keys expiration set to 24hrs and issued_at to time.Now()
    - Then once this is done we create a new instance of this map claim with signingmethodhs256 and load the **godotenv** to get the secretKey 
    - create a signedstring of token and return it