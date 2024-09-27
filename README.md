
# Golang Web Application

This project is a Golang HTTP server built using the GoFiber framework. The main goal is to implement `auth` and `user` modules. The entities within these modules should be stored in a local `.db` file using the `tidwall/buntdb` library.

## Validation
All request bodies or parameters must be validated using a validator. The validation rules are up to the developer's discretion.

## Code Structure
All features and package utilization should be encapsulated and use abstraction to separate business logic.

## Auth Module
The `auth` module will handle all endpoints related to user authentication. It will use JWT tokens for user authentication. The encryption method to be used is `SHA256` with a TTL of 10 minutes. Additionally, a refresh token mechanism must be implemented to issue new access tokens if the current token expires. When the JWT token expires, a new `AccessToken` and `RefreshToken` should be issued using the refresh token.

### Endpoints

#### 1. Register [POST] `/auth/register`
Registers a new user with the following request parameters:

**Request**:
```json
{
  "username": "<username of the user>",
  "email": "<email of the user>",
  "password": "<password of the user>",
  "name": "<name of the user>",
  "lastname": "<lastname of the user>"
}
```

**Response**:
```json
{
  "username": "<username of the user>",
  "email": "<email of the user>",
  "name": "<name of the user>",
  "lastname": "<lastname of the user>",
  "access_token": "<access token of the user (10 min TTL)>",
  "refresh_token": "<refresh token of the user>"
}
```

#### 2. Login [POST] `/auth/login`
Authenticates a user using their `username` or `email` along with the password.

**Request**:
```json
{
  "identifier": "<username or email of the user>",
  "password": "<password of the user>"
}
```

**Response**:
```json
{
  "username": "<username of the user>",
  "email": "<email of the user>",
  "name": "<name of the user>",
  "lastname": "<lastname of the user>",
  "access_token": "<access token of the user (10 min TTL)>",
  "refresh_token": "<refresh token of the user>"
}
```

#### 3. Logout [POST] `/auth/logout`
Logs out the user and deletes the refresh token associated with the session.

**Request**:
- The access token is sent in the header for user identification.

#### 4. Refresh [POST] `/auth/refresh`
Issues a new `AccessToken` and `RefreshToken` using the refresh token.

**Request**:
```json
{
  "identifier": "<identifier of the user [email or username]>",
  "refresh_token": "<refresh token of the user>"
}
```

**Response**:
```json
{
  "username": "<username of the user>",
  "email": "<email of the user>",
  "name": "<name of the user>",
  "lastname": "<lastname of the user>",
  "access_token": "<access token of the user (10 min TTL)>",
  "refresh_token": "<refresh token of the user>"
}
```

## User Module

### Endpoints

#### 1. GetProfile [GET] `/user/profile`
Retrieves the profile details of the user.

**Response**:
```json
{
  "username": "<username of the user>",
  "timezone": "<location of the user e.g., Europe/Istanbul>",
  "language": "<language of the user e.g., TR | EN>",
  "sport_branches": ["<Sport branches of the user>", "This is an array"]
}
```

#### 2. UpdateProfile [PATCH] `/user/profile`
Updates the user's profile details. **Note:** This is a PATCH method, so the handler logic should fit the requirements of a PATCH operation.

**Request**:
```json
{
  "username": "<username of the user>",
  "timezone": "<location of the user e.g., Europe/Istanbul>",
  "language": "<language of the user e.g., TR | EN>",
  "sport_branches": ["<Sport branches of the user>", "This is an array"]
}
```

**Response**:
```json
{
  "username": "<username of the user>",
  "timezone": "<location of the user e.g., Europe/Istanbul>",
  "language": "<language of the user e.g., TR | EN>",
  "sport_branches": ["<Sport branches of the user>", "This is an array"]
}
```

#### 3. Delete [DELETE] `/user`
Deletes the current user from the local database.
