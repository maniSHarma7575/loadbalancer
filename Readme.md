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
```

Note: Before building you have to make the required configuration file as defined in

### Configuration

You can manage the configuration for the loadbalancer by either json or yaml.

1. Create the symlink from the `config.json.sample` or `config.yaml.sample`

```
cp internal/config/config.json.sample internal/config/config.json
cp internal/config/config.yaml.sample internal/config/config.yaml
```
```
# Build
go build

# Run Load balancer
go run main.go
```

2. Change the configuration as per you requirement



`Strategy` can have the following values:

- `round-robin`
- `static`
- `traditional-hash`
- `consistent-hash`

## Usage

You can use loadbalancer utility by running command in your terminal:

1. Run the test servers:

```bash
  #change directory
  cd test/httpserver

  #build
  go build

  #run servers
  ./httpserver
```
2. Build and run the loadbalancer

```bash
  go run build && ./loadbalancer
```

3. Configuring TLS

```bash
# update config.yaml or config.json and replace the cert file and key file

tls_enabled: true
tls_cert_file: "/path/on/container/cert.pem"
tls_key_file: "/path/on/container/key.pem"
```

4. Configuring Content Based Routing (CBR)

```bash
# Append to config.yaml
routing:
  rules:
    - conditions:
      - path_prefix: "/api/v1"
        method: "GET"
        headers:
      actions:
        route_to: "app2"
  
# OR

# Append to config.json
"routing": {
  "rules": [
    {
      "conditions": [
        {
          "path_prefix": "/api/v1",
          "method": "GET",
          "headers": {
            "header": "header_value"
          }
        }
      ],
      "actions": {
        "route_to": "app2"
      }
    }
  ]
}
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
