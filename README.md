# httplogger

The `httplogger` is a standalone HTTP server that accepts all request and logs them.

## Features

The `httplogger` logs the following information for every request:
- Proto (for example `HTTP/1.1`)
- Host (for example `localhost:9081`)
- RequestURI (for example `/foo` - see RFC 7230, Section 3.1.1)
- Method (for example `GET`)
- URL (for example `/foo` - use this instead of RequestURI)
- RemoteAddr 
- TransferEncoding
- Header.* (all headers present in the request)
- JWT.* (if the header `Authorization: Bearer` contains a valid JWT)

Every logged request is identified with a unique request id consisting of:
- RFC3339 timestamp
- '#'
- counter since start of `httplogger`

Example for a request id: `2020-10-25T11:17:35+01:00#1`


## Usage

```
Usage of httplogger:
  -log string
        Comma separated list of keys to be logged.
        Supported keys:
        - Proto
        - Host
        - RequestURI
        - TransferEncoding
        - Method
        - URL
        - RemoteAddr
        - Header.*
        - JWT.*
        - JWT.Header
        - JWT.PayLoad
        "Header.*" will match all headers.
        Alternatively list the explicit headers (e.g. "Header.Accept").
  -port int
        Server port to listen. (default 8080)
  -response string
        Response to send.
```

The default port is 8080.

The default response is an empty response.

The default setting will log everything.
The list of log keys is case insensitive.

Example:
```
httplogger -port 9081 -response "Hello, world" -log "method,url,header.user-agent,jwt.*"
```

## Use Cases

Main use cases are:
- Debugging of HTTP requests
- Auditing of HTTP requests, especially in a Kubernetes/Istio cluster

In a Kubernetes/Istio shadow traffic can be used to send a copy of the HTTP request to `httplogger`. 


## Examples

### Example: GET requests

Example for a GET request (actually 2 requests) from a web browser:

```
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Proto                                    = HTTP/1.1
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Host                                     = localhost:9081
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 RequestURI                               = /
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Method                                   = GET
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 URL                                      = /
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 RemoteAddr                               = [::1]:65011
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Accept                            = text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Accept-Encoding                   = gzip, deflate, br
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Accept-Language                   = en-US,en;q=0.9
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Cache-Control                     = max-age=0
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Connection                        = keep-alive
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Cookie                            = Idea-92d69d34=c3336a1a-d9c7-44ae-bc76-2fc0e86e9a18
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Sec-Fetch-Dest                    = document
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Sec-Fetch-Mode                    = navigate
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Sec-Fetch-Site                    = none
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Sec-Fetch-User                    = ?1
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.Upgrade-Insecure-Requests         = 1
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#1 Header.User-Agent                        = Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Proto                                    = HTTP/1.1
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Host                                     = localhost:9081
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 RequestURI                               = /favicon.ico
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Method                                   = GET
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 URL                                      = /favicon.ico
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 RemoteAddr                               = [::1]:65011
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Accept                            = image/avif,image/webp,image/apng,image/*,*/*;q=0.8
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Accept-Encoding                   = gzip, deflate, br
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Accept-Language                   = en-US,en;q=0.9
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Cache-Control                     = no-cache
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Connection                        = keep-alive
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Cookie                            = Idea-92d69d34=c3336a1a-d9c7-44ae-bc76-2fc0e86e9a18
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Pragma                            = no-cache
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Referer                           = http://localhost:9081/
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Sec-Fetch-Dest                    = image
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Sec-Fetch-Mode                    = no-cors
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.Sec-Fetch-Site                    = same-origin
2020/10/25 15:18:09 2020-10-25T15:18:09+01:00#2 Header.User-Agent                        = Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36
```

### Example: GET requests with JWT token

Example with JWT bearer token (see [jwt.io](https://jwt.io/#debugger-io)) using curl:

```shell
curl -s -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c" http://localhost:9081/foo
```

Note that the `Authorization: Bearer` information shows up multiple times,
once as header and once decoded as JWT header and payload.
```
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 Proto                                    = HTTP/1.1
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 Host                                     = localhost:9081
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 RequestURI                               = /foo
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 Method                                   = GET
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 URL                                      = /foo
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 RemoteAddr                               = [::1]:65023
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 Header.Accept                            = */*
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 Header.Authorization                     = Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 Header.User-Agent                        = curl/7.65.3
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 JWT.Header                               = {"alg":"HS256","typ":"JWT"}
2020/10/25 15:19:10 2020-10-25T15:19:10+01:00#3 JWT.Payload                              = {"sub":"1234567890","name":"John Doe","iat":1516239022}
```

### Example: Limiting log keys

Testing the same example with a JWT bearer token but limiting the logged keys:

```
httplogger -port 9081 -response "Hello, world" -log "method,url,header.user-agent,jwt.payload
```

Only the requested log keys show up.
```
2020/10/25 15:20:13 2020-10-25T15:20:13+01:00#1 Method                                   = GET
2020/10/25 15:20:13 2020-10-25T15:20:13+01:00#1 URL                                      = /foo
2020/10/25 15:20:13 2020-10-25T15:20:13+01:00#1 Header.User-Agent                        = curl/7.65.3
2020/10/25 15:20:13 2020-10-25T15:20:13+01:00#1 JWT.Payload                              = {"sub":"1234567890","name":"John Doe","iat":1516239022}
```

### Example: POST request

Example for sending JSON as a POST request using curl:

```shell
curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST http://localhost:9081/data
```

Note that output does not contain the body of the request.
Because the body may be very large it is per default not logged.
```
2020/10/25 15:39:07 2020-10-25T15:39:07+01:00#4 Proto                                    = HTTP/1.1
2020/10/25 15:39:07 2020-10-25T15:39:07+01:00#4 Host                                     = localhost:9081
2020/10/25 15:39:07 2020-10-25T15:39:07+01:00#4 RequestURI                               = /data
2020/10/25 15:39:07 2020-10-25T15:39:07+01:00#4 Method                                   = POST
2020/10/25 15:39:07 2020-10-25T15:39:07+01:00#4 URL                                      = /data
2020/10/25 15:39:07 2020-10-25T15:39:07+01:00#4 RemoteAddr                               = [::1]:65308
2020/10/25 15:39:07 2020-10-25T15:39:07+01:00#4 Header.Accept                            = */*
2020/10/25 15:39:07 2020-10-25T15:39:07+01:00#4 Header.Content-Length                    = 34
2020/10/25 15:39:07 2020-10-25T15:39:07+01:00#4 Header.Content-Type                      = application/json
2020/10/25 15:39:07 2020-10-25T15:39:07+01:00#4 Header.User-Agent                        = curl/7.65.3
```

### Example: POST request with body

To log the body for request you start the server using:
```shell
httplogger -port 9081 -body
```

```shell
curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST http://localhost:9081/data
```

```
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 Proto                                    = HTTP/1.1
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 Host                                     = localhost:9081
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 RequestURI                               = /data
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 Method                                   = POST
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 URL                                      = /data
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 RemoteAddr                               = [::1]:65528
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 Header.Accept                            = */*
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 Header.Content-Length                    = 34
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 Header.Content-Type                      = application/json
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 Header.User-Agent                        = curl/7.65.3
2020/10/25 15:52:10 2020-10-25T15:52:10+01:00#1 Body                                     = {"key1":"value1", "key2":"value2"}
```

or you can explicitely list `Body` key in the `-log` option.
```shell
httplogger -port 9081 -log "method,url,body"
```

```shell
curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST http://localhost:9081/data
```

Only the requested keys including the body are logged.
```
2020/10/25 16:01:47 2020-10-25T16:01:47+01:00#1 Method                                   = POST
2020/10/25 16:01:47 2020-10-25T16:01:47+01:00#1 URL                                      = /data
2020/10/25 16:01:47 2020-10-25T16:01:47+01:00#1 Body                                     = {"key1":"value1", "key2":"value2"}
```

Newlines in the body will be escaped as `\n` to keep the log on a single line:
```shell
curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST http://localhost:9081/data
```

Only the requested keys including the body are logged.
```shell
curl -d '{
  "key1":"value1",
  "key2":"value2"
}' -H "Content-Type: application/json" -X POST http://localhost:9081/data
```

```
2020/10/25 16:56:26 2020-10-25T16:56:26+01:00#1 Method                                   = POST
2020/10/25 16:56:26 2020-10-25T16:56:26+01:00#1 URL                                      = /data
2020/10/25 16:56:26 2020-10-25T16:56:26+01:00#1 Body                                     = {\n  "key1":"value1",\n  "key2":"value2"\n}
```
