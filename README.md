# Recipes Server

![GitHub last commit](https://img.shields.io/github/last-commit/StGrozdanov/recipes-v2-server)

A Golang server to support the remake of the recipes website - both the frontend client and the mobile app CMS. You can find the: 

[Frontend here](https://github.com/StGrozdanov/recipes-v2-client)

[Mobile App here](https://github.com/StGrozdanov/recipes-v2-cms-client)

## Project Overview

### Technologies Used

- :whale: Docker
- :cloud: AWS SDK 
- :lock: Bcrypt 
- 🍹Gin 
- :key: JSON Web Token 
- 🏬 PostgreSQL 
- 📅 Moment 
- 👨‍💻 Logrus
- 🔦 Prometheus
- 🪝 pre-commit and pre-push hooks
- Websockets

### Database

The server uses PostgreSQL database with a mixture of JSONB and standard columns.

## Features

- CI/CD pipeline consisting of custom linters, unit tests, integration tests and a single deployment environment
- S3 bucket files download, upload, delete and read
- Authentication with Bcrypt and JWT
- Logging with Logrus
- Tracking by IP and referer
- Analytics
- Metrics by prometheus
- Healths endpoint
- Authenticated, non authenticated, admin groups and owner groups path guards
- Response time measuring on each endpoint
- A full range of CRUD operations for each supported resource
- Websockets
- Block user by IP
