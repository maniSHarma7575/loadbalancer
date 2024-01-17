# Loadbalancer

The LoadBalancer tool, built in Golang, offers functionality for distributing loads across multiple servers using various strategies.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)

## Installation

```bash
git clone https://github.com/maniSHarma7575/loadbalancer

# Change directory
cd loadbalancer

# Build
go build

# Run Load balancer
go run main.go
```

## Usage

You can use loadbalancer utility by running command in your terminal:

1. Change the Backend Configuration in `proxy.go`. Default configuration are listed below

```go
	configs := map[string]interface{}{
		"Backends": []map[string]interface{}{
			{"Host": "localhost", "Port": 8085, "HealthStatusUrl": "/health"},
			{"Host": "localhost", "Port": 8086, "HealthStatusUrl": "/health"},
			{"Host": "localhost", "Port": 8087, "HealthStatusUrl": "/health"},
		},
		"Strategy": "consistent_hash",
	}
```

`Strategy` can have the following values:

- `round-robin`
- `static`
- `traditional-hash`
- `consistent-hash`

2. Run the test servers:

```bash
  #change directory
  cd test/httpserver

  #build
  go build

  #run servers
  ./httpserver
```
3. Build and run the loadbalancer

```bash
  go run build && ./loadbalancer
```

### Description

Following Strategy are available:

1. Round-Robin
2. Static balancing strategy
3. Traditional Hash
4. Consistent Hash

## Contributing

Thank you for your interest in contributing to our project! We welcome your suggestions, improvements, or contributions. To get started, follow these steps:

### 1. Fork the Project

Click the "Fork" button on the top-right corner of this repository to create your own copy of the project.

### 2. Create a New Branch

Once you've forked the project, it's a good practice to create a new branch for your changes. This keeps your changes isolated and makes it easier to manage multiple contributions.

```bash
git checkout -b your-new-branch
