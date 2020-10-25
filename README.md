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

## Example

Example for a GET request from a web browser:

```
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Method                                   = GET
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 URL                                      = /
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 RemoteAddr                               = [::1]:49213
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Accept                            = [text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Accept-Encoding                   = [gzip, deflate, br]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Accept-Language                   = [en-US,en;q=0.9]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Cache-Control                     = [max-age=0]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Connection                        = [keep-alive]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Cookie                            = [Idea-92d69d34=c3336a1a-d9c7-44ae-bc76-2fc0e86e9a18]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Sec-Fetch-Dest                    = [document]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Sec-Fetch-Mode                    = [navigate]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Sec-Fetch-Site                    = [none]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Sec-Fetch-User                    = [?1]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.Upgrade-Insecure-Requests         = [1]
2020/10/25 11:17:35 2020-10-25T11:17:35+01:00#1 Header.User-Agent                        = [Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36]
```