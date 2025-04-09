# 📘 Article Account Service

A microservice designed to handle account management including login, registration, email verification, OTP, and user lookup features.

---

## 🚀 Features

- User registration & login
- JWT-based authentication
- Email verification via OTP
- Send OTP to email
- Account lookup & filtering
- Built with Go and RESTful API design

---

## 📂 API Endpoints

### 🔐 `POST /v1/account/login`

Login a registered user.

**Request Body**
```json
{
  "email": "user@example.com",
  "password": "yourPassword"
}
📝 POST /v1/account/register

Register a new user account.

Request Body

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securePassword"
}

Response

{
  "message": "Registration successful. Please verify your email."
}

🔍 GET /v1/account/id

Search for an account using filters like id, username, or email.

Query Parameters

    id (optional): User ID

    username (optional): Filter by username

    email (optional): Filter by email

Example

GET /v1/account/id?id=123
GET /v1/account/id?username=johndoe

Response

{
  "id": "123",
  "username": "johndoe",
  "email": "johndoe@example.com"
}

✅ POST /v1/account/verify

Verify a user's email with the OTP code sent to them.

Request Body

{
  "email": "user@example.com",
  "otp": "123456"
}

Response

{
  "message": "Email verified successfully"
}

✉️ POST /v1/account/send-otp

Send an OTP to the provided email address for verification.

Request Body

{
  "email": "user@example.com"
}

Response

{
  "message": "OTP sent to email"
}

⚙️ Getting Started

    Clone the repository

git clone https://github.com/your-username/article-account-service.git
cd article-account-service

    Install dependencies

go mod tidy

    Set up environment variables
    Create a .env file and configure it with:

DATABASE_URL=
JWT_SECRET=
SMTP_HOST=
SMTP_PORT=
SMTP_USERNAME=
SMTP_PASSWORD=

    Run the service

go run main.go

    Access the API via:
    http://localhost:<your-port>/v1/account/...

🧪 Testing

You can use tools like Postman or cURL to test the endpoints.
📬 Contribution

Feel free to fork this repo and submit PRs. Suggestions, issues, and feedback are always welcome!
📄 License

This project is licensed under MIT License.

---

Let me know if you want a version that includes Docker setup, OpenAPI (Swagger), or Postman Collection too!
