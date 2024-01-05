## Technologies Used

- [Golang](https://go.dev/): It is efficient, simple and easy to scale. It can be easily containerized since the go compiler builds a standalone executable file.
- [Echo](https://github.com/labstack/echo): It is a minimalist http router library with some prebuilt HTTP middleware handlers.
- PostgreSQL: It's an open-source, production-ready database with a good feature set. It has several popular index types, including [GIN](https://www.postgresql.org/docs/current/gin-intro.html) index. A GIN index efficiently accelerates searches, which is particularly useful for quickly finding notes based on keywords in your text search functionality.
- [sqlc](https://docs.sqlc.dev/en/stable/index.html): sqlc generates fully type-safe idiomatic Go code from SQL. It is an alternative to using opinionated ORMs.
- [Goose](https://github.com/pressly/goose): Database migration tool that works well with sqlc

## How to use

1. Update the .env file or use the default value and export them
2. Start the Postgresql database using `make db-up`
3. Install the go compiler if not already.
4. Install necessary tools using `make tools`
5. Migrate db using `make migrate-up`
6. Execute the application using `make run`

## How to test

### Unit test
1. Run `make unittest`

### end-to-end test
1. Export the env variables from the .env file
2. Run `make e2e`

NOTE: The written tests are very basic; they don't cover all cases because of time constraints. However, I have experience writing both unit and e2e tests for open-source projects. Please check out [unittest](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/14699) and [e2e test](https://github.com/splunk/splunk-connect-for-kubernetes/pull/707).


### Further Improvements
- The server currently uses an in-memory store for rate-limiting. It makes the application stateful. To scale up the application, [redis](https://redis.io/) can be used to store the data related to rate-limiting. It makes the server stateless and can be easily scaled horizontally.
- The application is executed directly, but it can be easily containerized. It makes the deployment process relatively easy when using container orchestration tools like Kubernetes.
- There is no mechanism for token expiry. We can create a new table to keep track of active tokens and use the caching mechanism to check if the token is valid efficiently. A new endpoint, `/api/auth/logout`, that invalidates the cache, can be introduced.
- The current server serves HTTP traffic. We can enable HTTPS by adding SSL certificates(either self-signed or issued from a known certificate authority)
- We can improve the observability by adding more logs(not done already due to time constraints), performance metrics(like CPU/memory usage, response time, throughput), database metrics(like the performance of SQL queries, connection pool stats, etc) and application-level metrics(like auth, error, etc), and tracing. We can leverage tools like Prometheus and Opentelemetry to collect and export these data.
- Can add more unit and integration tests(not done already due to time constraints)

## Endpoints

### POST /api/auth/signup
```bash
curl --location 'http://localhost:8080/api/auth/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "user1",
    "email": "user1@gmail.com",
    "password": "Hello@123"
}'
```

### POST /api/auth/login
```bash
curl --location 'http://localhost:8080/api/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "user1@email.com",
    "password": "Hello@123"
}'
``` 

### GET /api/notes/
```bash
curl --location 'http://localhost:8080/api/notes/' \
--header 'Authorization: Bearer <TOKEN>'
```

### POST /api/notes/
```bash
curl --location 'http://localhost:8080/api/notes/' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <TOKEN>' \
--data '{
    "title": "title 1",
    "content": "content 1"
}'
```

### PUT /api/notes/:id
```bash
curl --location --request PUT 'http://localhost:8080/api/notes/1' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <TOKEN>' \
--data '{
    "title": "note 1.1",
    "content": "My updated content 1.1"
}'
```

### DELETE /api/notes/1
```bash
curl --location --request DELETE 'http://localhost:8080/api/notes/1' \
--header 'Authorization: Bearer <TOKEN>' \
--data ''
```

### POST /api/notes/1/share
```bash
curl --location 'http://localhost:8080/api/notes/1/share' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "shared_with": "user2@gmail.com"
}'
```

### GET /api/notes/search?q=query
```bash
curl --location 'http://localhost:8080/api/notes/search?q=my%20content' \
--header 'Authorization: Bearer <TOKEN>'
```
