# Chirpy

A twitter(ish) clone built as a boot.dev project
Built using:
- go 1.24.1
- PostgresSQL 16.9

## Installation

`go get github.com/zdelk/chirpy`

## User Guide

Requires a .env file with the following:
- DB_URL
- PLATFORM
- JWT_SECRET
- POLKA_KEY

### Functions
POST /api/users - Creates User  
json: 
```
{
    "email":"<user@email.com>",  
    "password":"<password>"
}  
```

POST /api/login - Logs in to Account
```
{
    "email":"<user@email.com>",  
    "password":"<password>"
}  
```

PUT /api/users - Updates User info. Must be logged in
```
{
    "email":"<new@email.com>",  
    "password":"<new_password>"
}  
```
POST /api/polka/webhooks - Upgrades a user to premium

POST /api/chirps - Create Chirp. Must be fewer than 140 characters    
json :
```
{
    "body":"text for the chirp"
}
```

GET /api/chirps  - Gets All chirps  
Optional Sorting:  
GET /api/chirps?sort=asc - Ascending by creation timestamp (Default)  
GET /api/chirps?sort=desc - Descending by creation timestamp  

GET /api/chirps/{chirpID} - Return single chirp  

DELETE /api/chirps/{chirpID} - Deletes chirp. Must be author  



POST /api/refresh - Refreshes a users privileges
POST /api/revoke - Revokes a users privileges 



GET /admin/metrics - Show how many times the page has been accessed


POST /admin/reset - Resets metrics and server
