# go build

# Sign in example 
POST http://localhost:8000/token
Content-Type: application/json
{
        "username": "user1",
        "password": "password1"
}


# Test welcome using cookies 
GET http://localhost:8000/welcome

# Refesh token 
POST http://localhost:8000/refresh