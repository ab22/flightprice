# Flight Price

Go microservice that retrieves flight prices from local mocked APIs.


## Quickstart

```sh
# Create our .env file.
cp .env.example .env

# This will build our backend and launch it with docker compose.
make build && make up

# Perform a quick test
curl -i http://localhost:8080/ping

# or make use of helper recipes already on backend/Makefile.
make ping

# perform a request to /flights/search without a valid JWT token. This should
# show an http.StatusNotFound (404) instead of a forbidden. We show a 404 and not a
# 403 for security purposes.
make flights-invalid

# create a sample JWT token by calling our API.
make login

# Copy the JWT token into our `.env` file inside of the API_TOKEN variable. Now
# we should be able to hit our flights endpoint.
make flights

# Test the websocket server. The current recipe sends a path parameter
# so that it receives an update every 3 seconds. This can be changed to
# /subscribe/x where x is the number of seconds on which you will receive the
# next message.
make ws
``
