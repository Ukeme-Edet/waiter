# Waiter

![Waiter Logo](./assets/waiter.gif)
_A lightweight HTTP server for development and testing_

Waiter is a simple HTTP server written in Go that handles various HTTP requests and supports file operations.

## Features

-   Basic HTTP server functionality
-   Echo endpoint that returns provided text
-   User-Agent detection
-   File storage and retrieval
-   Gzip compression support
-   Keep-alive connection handling

## Requirements

-   Go 1.18 or later

## Installation

1. Clone this repository:

    ```bash
    git clone <your-repository-url>
    cd waiter
    ```

2. Build the application:

    ```bash
    make build
    ```

    This will create the executable in the `bin` directory.

## Usage

Run the server:

```bash
./bin/app
```

By default, the server will listen on `0.0.0.0:4221`.

### Specifying a directory for file operations

To specify a directory for file operations:

```bash
./bin/app --directory /path/to/directory
```

## API Endpoints

-   **`/`**: Root endpoint

    -   Method: GET
    -   Response: 200 OK with text/plain content type

-   **`/echo/{text}`**: Echo endpoint

    -   Method: GET
    -   Response: 200 OK with the provided text echoed back

-   **`/user-agent`**: User agent endpoint

    -   Method: GET
    -   Response: 200 OK with the User-Agent header value

-   **`/files/{filename}`**: File operations
    -   Method: GET
        -   Description: Retrieve a file
        -   Response: 200 OK with file contents or 404 Not Found
    -   Method: POST
        -   Description: Create or update a file
        -   Response: 201 Created on success

## Development

Format, lint and build:

```bash
make
```

Clean build artifacts:

```bash
make clean
```

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
