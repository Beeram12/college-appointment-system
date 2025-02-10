# College Appointment System

The **College Appointment System** is a web application that enables students to book appointments with professors based on their availability. The system is built using Go (Golang) as the backend language, MongoDB as the database, and JWT (JSON Web Token) for secure authentication.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Features](#features)
3. [Technologies Used](#technologies-used)
4. [Setup Instructions](#setup-instructions)
5. [Endpoints](#endpoints)
6. [Authentication](#authentication)
7. [Testing](#testing)
8. [Contributing](#contributing)
9. [License](#license)

---

## Project Overview

The **College Appointment System** simplifies the process of scheduling appointments between students and professors. Professors can set their availability, while students can view and book available slots. The system ensures secure access and proper role-based permissions using JWT-based authentication.

---

## Features

- **User Registration and Login**:
  - Students and professors can register and log in to the system.
- **Set Availability**:
  - Professors can set their available time slots.
- **Book Appointments**:
  - Students can book appointments with professors based on their availability.
- **View Appointments**:
  - Students can view their booked appointments.
- **Cancel Appointments**:
  - Students can cancel their appointments.
- **Role-Based Access Control**:
  - Separate access levels for `student` and `professor` roles.

---

## Technologies Used

- **Backend**: Go (Golang)
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Token)
- **API Framework**: Gorilla Mux
- **Utilities**:
  - `golang.org/x/crypto` for password hashing
  - `github.com/golang-jwt/jwt/v5` for JWT generation and validation
- **Build Tool**: Makefile
- **Testing**: Postman

---

## Setup Instructions

### Prerequisites

1. **Go**:
   - Install Go from [https://golang.org/dl/](https://golang.org/dl/).
   - Ensure `go version` returns a compatible version (e.g., `1.20+`).

2. **MongoDB**:
   - Install MongoDB locally or use a cloud-hosted instance like MongoDB Atlas.
   - Update the connection string in `.env`.

3. **Postman**:
   - Install Postman for API testing: [https://www.postman.com/downloads/](https://www.postman.com/downloads/).

4. **Make**:
   - Install `make` for running build and test commands.

---

### Steps to Set Up

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Beeram12/college-appointment-system.git
   cd college-appointment-system
   ```

2. **Install Dependencies**:
   Make sure you have `go mod` initialized and install the required dependencies:
   ```bash
   make deps
   ```
   This will execute `go mod tidy`, downloading and organizing all dependencies           specified in the `go.mod` file.

3. **Configure Environment Variables**:
   Create a `.env` file in the root directory and add the following configurations:
   ```env
   MONGO_URI=mongodb://localhost:27017
   MONGO_DB="Your DB name"
   JWT_SECRET=your_secret_key
   PORT=8000
   ```
    - Replace `mongodb://localhost:27017` with your MongoDB connection string (e.g., for MongoDB Atlas).
    - Replace your`_jwt_secret_key_`here with a strong secret key for signing JWT tokens.
    - Ensure the `.env` file is added to `.gitignore` to prevent sensitive information from being committed:


4. **Run the Application**:
   Start the server using the following command:
   ```bash
   make run
   ```
   The application will run at `http://localhost:8000`.

---

## Endpoints

### Authentication
| Method | Endpoint                  | Description                |Request Body                     |
|--------|---------------------------|----------------------------|---------------------------------|
| POST   | `/auth/register`          | Register a user            |`{ "username": "user1", "password": "password123", "role": "student" }`|
| POST   | `/auth/login`             | Login to get a JWT token   |`{ "username": "user1", "password": "password123" }` |

### Professor Endpoints
| Method | Endpoint                  | Description                       |Request Body                |
|--------|---------------------------|-----------------------------------|----------------------------|
| POST   | `/professor/availability` | Set available time slots          |`{ "time_slot": "09:00 AM"}`|
| GET    | `/professor/appointments` | View all appointments             |`None`                      |

### Student Endpoints
| Method | Endpoint                     | Description                    |Request Body|
|--------|------------------------------|--------------------------------|------------|
| GET    | `/appointments/{professor_id}`              | View available slots          |`None`       |
| POST   | `/appointments/book`         | Book an appointment           |`{ "professor_id": "67a62b8f1dae6c55fdddf2d9", "time_slot": "09:00 AM" }`|
| DELETE | `/appointments/cancel/{appointment_id}`   | Cancel a booked appointment   |`None`|
| GET    | `/student/appointments`      | View booked appointments      |`None`|

---

## Authentication
- The system uses JWT(JSON Web Token) for secure authentication.
- Users must include their token in the `Authorization` header for protected endpoints:
  ```bash
  Authorization: Bearer <your_token>
  ```

---

## Testing

### Manual Testing
- Use **Postman** to test all endpoints individually.
- Include the JWT token in headers for authenticated routes.

