# Golang Web Application
This project is a web application built with Golang using the GoFiber framework. It includes user management and authentication (auth) modules. The primary focus of this application is user registration, authentication (with JWT tokens), and managing user profiles. User data is stored in a local .db file using the Tidwall/buntdb library, which is a lightweight and embeddable key/value store.

## Features
User registration with validation.
User login with JWT token-based authentication.
Profile management (get and update user information).
Email-based user search.
JWT-based access and refresh tokens.
Local database storage using Tidwall/buntdb.

## Technologies Used
GoFiber: Web framework for building fast and scalable web applications.
BuntDB: An embeddable, in-memory key/value store with support for persistency.
JWT (JSON Web Tokens): For user authentication and session management.
SHA256: For password hashing.
Fiber middleware: For error handling, logging, and routing.

## Project Structure
![image](https://github.com/user-attachments/assets/3f5bae63-0416-4975-97d4-b0eda4c30d6a)

## Requirements
Go version 1.17 or higher
Git
A local .env file to specify the environment variables (e.g., LOCAL_DB_PATH).

## Installation and Setup
Clone the Repository:
```
git clone https://github.com/alisahin10/Golang-Web-App
cd golang-web-app

```



Install Dependencies: Run the following command to install the necessary Go dependencies:
```
go mod tidy
```

Configure the .env file: Create a .env file in the root directory and configure the environment variables:
```
LOCAL_DB_PATH=./data/buntdb.db
```

Run the Application: Start the application with:
```
go run main.go
```

## Endpoints
### 1. Auth Module (/auth)
Handles all endpoints related to user authentication and token management. It uses JWT tokens for access and refresh token mechanisms.

**Register [POST] /auth/register
Registers a new user.**

Request Body:
```
{
  "username": "string",
  "email": "string",
  "password": "string",
  "name": "string",
  "lastname": "string"
}

```
Response Body:
```
{
  "username": "string",
  "email": "string",
  "name": "string",
  "lastname": "string",
  "access_token": "string (10 min TTL)",
  "refresh_token": "string"
}
```

**Login [POST] /auth/login
Authenticates a user and provides access and refresh tokens.**

Request Body:
```
{
  "identifier": "string (username or email)",
  "password": "string"
}

```
Response Body:
```
{
  "username": "string",
  "email": "string",
  "name": "string",
  "lastname": "string",
  "access_token": "string (10 min TTL)",
  "refresh_token": "string"
}

```
**Logout [POST] /auth/logout
Logs out the user and invalidates the refresh token.**

**Request Header:** Authorization: Bearer <access_token>

**Refresh [POST] /auth/refresh
Issues a new access token and refresh token using the refresh token.**

Request Body:
```
{
  "identifier": "string (username or email)",
  "password": "string"
}

```
Response Body:
```
{
  "username": "string",
  "email": "string",
  "name": "string",
  "lastname": "string",
  "access_token": "string (10 min TTL)",
  "refresh_token": "string"
}
```

### 2. User Module (/user)
Handles user profile management including getting user profiles, updating user information, and deleting users.

**Get User Profile [GET] /user/:id
Retrieves the details of a specific user by their ID.**
Response Body:
```
{
  "id": "string",
  "username": "string",
  "email": "string",
  "name": "string",
  "lastname": "string",
  "age": "integer"
}


```
**Update User Profile [PATCH] /user/:id
Updates the user's profile details.**
Request Body:
```
{
  "email": "string",
  "name": "string",
  "lastname": "string",
  "age": "integer"
}


```
Response Body:
```
{
  "message": "User updated successfully"
}
```
**Delete User [DELETE] /user/:id
Deletes a specific user by their ID.**
Response Body:
```
{
  "message": "User deleted successfully"
}
```
**Search User by Email [GET] /user/search
Searches for a user by their email address.**

Request Body:
```
{
  "email": "string"
}
```
Response Body:
```
{
  "id": "string",
  "username": "string",
  "email": "string",
  "name": "string",
  "lastname": "string",
  "age": "integer"
}
```
