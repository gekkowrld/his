# Health Information System (HIS)

A simple program designed to simulate a basic Health Information System (HIS) for managing client profiles and health programs/services.

## Features

This application includes the following features:

- **Create Health Programs**: e.g., Tuberculosis (TB), HIV, etc.
- **Register New Clients**: Capture personal and medical information for new clients.
- **Enroll Clients in Programs**: Assign clients to one or more health programs.
- **Search Client Profiles**: Quickly find and view detailed client information.
- **Client Information API**: Expose client data via a basic RESTful API.
- **User Login**: Simple authentication for doctors or users (basic implicit login).

## Prerequisites

Before building and running the system, ensure you have the following installed:

- **Go 1.23+**: Required for building and running the application.
- **C Compiler**: Needed for linking CGO dependencies (used by SQLite3).

If you're new to Go, follow the [official Go tutorial to get started](https://go.dev/doc/tutorial/getting-started).

## Build & Run

### Building the Application

To build and run the client or server, follow the instructions in the official Go tutorial for setting up your environment and compiling the code.

Alternatively, you can use the provided `tools/run.go` file to build and run the application:

```bash
go run tools/run.go
```

This will start both the client and the server.

## Dependencies

All dependencies will be automatically installed using Go's build system. No additional tools are required, apart from the C compiler for SQLite3 support.

## API

The application exposes basic client information through a simple RESTful API. You can interact with the API to retrieve and manage client data.

All endpoints require basic authentication using `user:password`.
If the authentication is missing or incorrect, the request will fail.
Consult your tool or client on how to include authentication headers.

### Available Endpoints

#### POST /client/new

Creates a new client.

**Requires:**

- `first_name` (string)
- `middle_name` (string) (optional)
- `last_name` (string)
- `day` (int)
- `month` (int)
- `year` (int)

**Returns:**

- `id` (string)
- `first_name` (string)
- `middle_name` (string) (optional)
- `last_name` (string)
- `day` (int)
- `month` (int)
- `year` (int)

---

#### POST /program/new

Creates a new health program.

**Requires:**

- `name` (string)
- `description` (string)
- `start_date` (int) [seconds from Unix epoch]
- `end_date` (int) [seconds from Unix epoch] (optional)

**Returns:**

- `id` (string)
- `name` (string)
- `description` (string)
- `start_date` (int) [seconds from Unix epoch]
- `end_date` (int) [seconds from Unix epoch]

---

#### POST /client/{id}

Associates a client with one or more programs.
The `{id}` is the client ID.

**Requires:**

- `programs` [array of program IDs]

**Returns:**

- `client_id` (string)
- `programs` [array of program IDs]

---

#### GET /search/client

Searches for client names in the database.

**Requires:**

- `q` (string) — query term for searching client names.

**Format:** `/search/client?q=query`

**Returns:**

- `id` (string)
- `first_name` (string)
- `last_name` (string)

---

#### GET /client/{id}

Returns detailed information about a client.
The `{id}` is the client ID.

**Returns:**

- `id` (string)
- `first_name` (string)
- `last_name` (string)
- `programs` (array of program dictionaries with `name` and `id`)

---

## License

This project is licensed under the [AFL-3.0 license](LICENSE).

## Contribution

Contribution is not actively encouraged.
However, if you'd like to use the code for your own purposes, feel free to do so.
Please note: Any open issues or pull requests will not be attended to.
