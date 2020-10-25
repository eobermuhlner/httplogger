package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
)


func main() {
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	log.Println("Started httplogger")
	count := 0

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count++
		req := fmt.Sprintf("#%v", count)

		LogString(req, "Method", r.Method)
		LogString(req, "URL", html.EscapeString(r.URL.Path))
		LogString(req, "RemoteAddr", r.RemoteAddr)
		for k, v := range r.Header {
			LogStrings(req, fmt.Sprintf("Header[%v]", k), v)
		}

		if authorizations, ok := r.Header["Authorization"]; ok {
			for _, authorization := range authorizations {
				if strings.HasPrefix(authorization, "Bearer ") {
					token := authorization[7:]
					jwt, err := DecodeJWT(token);
					if jwt != "" {
						LogString(req, "JWT", jwt)
					}
					if err != nil {
						log.Printf("Invalid JWT: %v", err)
					}
				}
			}
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}

func DecodeJWT(token string) (string, error) {
	splitToken := strings.Split(token, ".")
	if len(splitToken) >= 2 {
		tokenPayload := splitToken[1]
		jwt, err := base64.RawURLEncoding.DecodeString(tokenPayload)
		if err != nil {
			return string(jwt), err
		}
		return string(jwt), nil
	}
	return "", errors.New("invalid JWT")
}

func LogString(req string, name string, value string) {
	log.Printf("%v %-30v = %v\n", req, name, value)
}

func LogStrings(req string, name string, value []string) {
	log.Printf("%v %-30v = %v\n", req, name, value)
}
