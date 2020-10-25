# httplogger

The `httplogger` is a standalone HTTP server that accepts all request and logs them.

## Features

The `httplogger` logs:
- Method
- URL
- Headers
- JWT payload


## Usage

```
httplogger [-port PORT]
```

The default port is 8080.


## Use Cases

Main use cases are:
- Debugging of HTTP requests
- Auditing of HTTP requests, especially in a Kubernetes/Istio cluster

In a Kubernetes/Istio shadow traffic can be used to send a copy of the HTTP request to `httplogger`. 

