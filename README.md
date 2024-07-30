# Library Management System

This project is a Library Management System built using the MVC architecture and MySQL. The system has separate client and admin portals, secure login functionalities, and comprehensive management features.

## Features

### Separate Client and Admin Portals

### Authentication & Authorization
- Secure login functionalities for both clients and admins with role-based access control.

### Admin Features
- **Manage Book Catalog**: Admins can list, update, add, and remove books.
- **Approve/Deny Requests**:
  - Checkout and check-in requests from clients.
  - Requests from users seeking admin privileges.

### Client Features
- **View Books**: Clients can view the list of available books.
- **Request Management**:
  - Request checkout and check-in of books from the admin.
  - View their borrowing history.

### Security Features
- **Secure Password Storage**: Passwords are hashed and salted before being stored in the database.
- **JWT-based Session Management**: Implemented custom session management using JWT tokens.

## Getting Started

### Prerequisites
- Go 1.16+
- MySQL
- golang-migrate

### Installation

1. Clone the Repository:
    ```bash
    git clone https://github.com/Chalhotra/LMS_Go.git
    cd LMS_Go
    ```

2. Build and Run the Application:
    ```bash
    make run
    ```

### Makefile Commands
- **build**: Clean and build the application.
    ```bash
    make build
    ```

- **run**: Build and run the application.
    ```bash
    make run
    ```

- **migrate-up**: Run database migrations up.
    ```bash
    make migrate-up
    ```

- **migrate-down**: Run database migrations down.
    ```bash
    make migrate-down
    ```

- **clean**: Clean the build.
    ```bash
    make clean
    ```

- **test**: Run all tests.
    ```bash
    make test
    ```

- **lint**: Lint the code.
    ```bash
    make lint
    ```

- **deps**: Update dependencies.
    ```bash
    make deps
    ```

- **fmt**: Format the code.
    ```bash
    make fmt
    ```

- **vendor**: Vendor dependencies.
    ```bash
    make vendor
    ```

- **debug**: Debug environment variables.
    ```bash
    make debug
    ```

## Usage
### Login and Register Portals:
- Access via `/login` and `/register`
- Initially when you run ``` make migrate-up```, a super admin user is created by default with credentials, username: admin, password: admin
### Admin Portal:
- Access via `/api/admin/`.
- Manage books and user requests.

### Client Portal:
- Access via `/api/client/`.
- View available books, request checkouts/check-ins, and view borrowing history.
