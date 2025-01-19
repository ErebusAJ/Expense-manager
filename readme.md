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

## ***5. Authenticated Admin handlers***
    
We need endpoint routes for admin access for that we need to assign access levels, so we alter user table to add **access_level** field which will have user access by default.

> I also had to change the JWT generater function to embed user role field in token

One approach is we can use if conditions in our user handlers e.g if user.access == admin then we call admin function but it's not secure and can expose all users data to single user so we use seprate auth and handlers

* **migrations/schema**:
    - Write appropriate goose up/down query to alter the users table and add column **access_level**

* **internal/middleware/middleware_admin**:
    - Here we write the middleware to authenticate admin users 
    almost similar to auth middleware just we pass two token claims to the handler from middleware **(userId, userRole)**

* **internal/handlers/handler_admin.go**:
    - **adminLogin**: similar to user login binds request body to a self defined struct **(email,password)** just we before comparing password hash we check if userID == adminID if not return error else go to check password. With all validated generate a authToken
    - **adminGetAllUsers**: takes params from middleware userId, userRole and verify them from env which stores adminID and adminRole if validated return all users data

* **internal/handlers/routes.go**:
    - Add a new group **adminProtected** with path "/admin" add the middleware and then the routes we want to user here /users which calls adminGetAllUsers

> A small optimization I did is everytime a methods err is not nill we need to return the error to client via c.IndentedJSON passing error code and obj and then a log for server which explains about the error breifly. So I creted a helper function in **internal/utils/err.go** which takes input ginContext, error code, client string, server string and error

## ***6. Password reset***
Now we need a handler to reset user password if they forget to help them recover theri account.

* **schema/passwordReset.sql**:
    - We need to create a table **password_tokens** so that we can store the token generated and use it for verification while the user updates his password. Add fields like, id, user_id(foreign key), token, craeted_at, expires_at

* **queries/table_password.sql**:
    - Write a few queries like insert, update, delete, select

* **handlers/handleruser.go**:
    - **resetPasswordRequest**: it verfies the email being received as a request is appropriate and does a user with the email exists or not. After verifyinfg we generate a token and it's expiry time i.e 1hr and insert it into **password_token** table. Then we call sendEmail function passing arguments email received in request and url/password reset link

    - **resetPasswordConfirm**: here we get a newPassword in request and also a param **token** we verify token if it exists in DB. If it does and if it's valid not expired we hash the password and update users table with the new password. After updation is done we have no use of the token so we delete it

* **utils/email.go** : 
    - to send an email I have used gomail. First retrieve your email and app pass form env. Create a **NewMessage()** instance of gomail set values to headers like from, to, subject, and the body where link goes

    - Create a **NewDialer()** of gomail set host argument as **"smtp.gmail.com"** for gmail and port 587 then your email and pass


## ***Users Endpoints***


| Endpoint | HTTP Method | Purpose | Authentication | 
| ----------- |----------- |----------- |----------- |
| `/v1/register` | `POST` | Register a new user | No |
| `/v1/login` | `POST` | Authenticates a user, returns JWT | No |
| `/auth/user` | `GET` | Get logged in user details | Yes (JWT) |
| `/auth/user` | `PUT` | Update user details | Yes (JWT) |
| `/auth/user` | `DELETE` | Delete authenticated user | Yes (JWT) |
| `/v1/user/password-reset` | `POST` | Sends a pass reset request to user email | No |
| `/v1/user/password-reset` | `POST` | Update user's password | Reset Token |
| `/admin/users` | `GET` | Get all users details | Yes (JWT) |
| `/admin/user/:id` | `GET` | Get a user's details by ID | Yes (JWT) |
| `/admin/user/:id` | `DELETE` | Delete a user by its ID | Yes (JWT) |


## ***7. Expenses DB/Handlers**
Now that users endpoints are done. We move on to next table expenses which stores expenses of a user.

As done for users we creat table schema, queries and handlers 

| Endpoint | HTTP Method | Purpose | Authentication | 
| ----------- |----------- |----------- |----------- |
| `/expenses` | `GET` | Gets a user's all expenses | Yes (JWT) |
| `/expense/:id` | `GET` | Gets a user expense of specified ID | Yes (JWT) | 
| `/expense/` | `POST` | Adds a user expense | Yes (JWT) | 
| `/expense/:id` | `PUT` | Updates a user expense of specified ID | Yes (JWT) | 
| `/expense/:id` | `GET` | Deletes a user expense of specified ID | Yes (JWT) | 












