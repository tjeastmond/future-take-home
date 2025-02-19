# Future Fit Take Home Project

## Overview

This project is designed to create a simple appointment scheduling API using Go and the Gin web framework. It allows users to manage appointments with trainers, check availability, and retrieve existing appointments.

## Requirements

You must have a modern version of Docker installed on your machine to run this project. You can download Docker [here](https://www.docker.com/products/docker-desktop).

Personally, I use OrbStache to manage my Docker containers. You can download OrbStache [here](https://orbstache.io).

## Getting Started - Docker

To get started with this project, you must first clone the repository to your local machine. You can do this by running the following command in your terminal:

```sh
docker-compose up -d
```

When you are done with the project, you can stop and remove the containers by running the following command:

```sh
docker-compose down -v
```

The data should persist if you decide to start or stop the containers.

```sh
docker-compose start
docker-compose stop
```

## API Endpoints

The following endpoints are available for use:

- `GET /trainers/1` - will return a list of appointmenrts for trainer 1
- `GET /trainers/1/availability` - will return a list of available times for trainer 1
- `POST /trainers/1` - will create a new appointment for trainer 1

Below are some examples of how to use these endpoints:

```sh
# Get a list of appointments for trainer 1
curl -X GET 'localhost:8080/trainers/1' \
  --url-query 'trainer_id=1'

# Get a list of available times for trainer 1
curl -X GET 'localhost:8080/trainers/1/availability' \
  --url-query 'starts_at=2025-02-19T08:00:00-08:00' \
  --url-query 'ends_at=2025-02-20T08:00:00-08:00'

# Create a new appointment for trainer 1
curl -X POST 'localhost:8080/trainers/1' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "trainer_id": 1,
  "user_id": 1,
  "starts_at": "2025-02-19T14:00:00-08:00",
  "ends_at": "2025-02-19T14:30:00-08:00"
}'
```

## Testing and Notes

- I didn't have time to write tests for this project
- I took a little longer than the suggested 3 or so hours to complete this project after the request to use Postgres
- I didn't use an ORM because I was nearly done when the Posrgres request came in and I'm a little rust with Gorn
- I didn't create a table or check for valid trainer IDs
