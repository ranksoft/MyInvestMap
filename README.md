# MyInvestMap

MyInvestMap is a web application for investment portfolio management that enables users to track stocks, buy and sell them, and analyze their investments.

## Key Features

- **Registration and Authentication**: Allows users to create an account and log in to access their portfolios.
- **Asset Management**: Users can add, update, sell, and delete assets from their portfolio.
- **Portfolio Viewing**: Ability to view the current state of the investment portfolio, including the total investment value and profit/loss.
- **Automatic Data Update**: Automatic updating of stock prices and other asset information through an API.
- **API Key Configuration**: Users can enter their personal API key for integration with external financial services.

## Technologies

- **Frontend**: React, Bootstrap
- **Backend**: Go (Golang) using the Gorilla Mux library for routing and JWT for authentication.
- **Database**: SQLite
- **Additional**: Docker for deployment and development simplification.

## Installation and Running

### Requirements

For the application to work, Docker and Docker Compose must be installed.

### Installation Steps

1. **Cloning the Repository**

Copy code
```bash
git clone <repository-url>
cd myinvestmap
```
2. **Running with Docker**

Copy code
```bash
docker-compose up --build
```

After executing these commands, the application will be accessible at `http://localhost:3000`.

### Usage
After starting, go to `http://localhost:3000` and register in the system. After registration and logging in, you will be able to manage your portfolio.

### Contributing
Contributions are welcome. To propose changes, create a pull request with a description of your changes.