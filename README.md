# Go Process Management Web Server

## Overview

This project is a simple web server written in Go, designed to manage system processes within a containerized environment. It allows users to start, list, and kill processes via a web interface, making it ideal for educational purposes, debugging, or monitoring in a development or testing environment.

## Features
- Start a long-running subprocess
- List all current processes in a tree-like structure
- Kill specific processes or process groups

## Prerequisites
- Go (1.x or newer)
- Docker (optional for containerization)

## Installation

### Clone the Repository
```bash
git clone https://github.com/moabukar/go-process-mgmt.git
cd go-process-mgmt
```

### Build the app

```bash
go build
```

### Run the server

```bash
./go-process-mgmt
```

The server will start on port 8080.

### Usage

Usage
After starting the server, you can interact with it using the following endpoints:
```bash
http://localhost:8080/run: Start a subprocess.
http://localhost:8080/ps: List all processes.
http://localhost:8080/kill: Kill a process. Use ?pid=<process_id> to specify a process.
```

### Containerization

The server can be run in a containerized environment using the included Dockerfile. To build the image, run the following command from the root directory of the project:

```bash
docker build -t go-process-mgmt .
```

To run the container, use the following command:

```bash

docker run -p 8080:8080 go-process-mgmt
```
