package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

func main() {
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	log.Println("Started httplogger")
	count := 0

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count++
		req := fmt.Sprintf("%v#%v", time.Now().Format(time.RFC3339), count)

		LogString(req, "Method", r.Method)
		LogString(req, "URL", html.EscapeString(r.URL.Path))
		LogString(req, "RemoteAddr", r.RemoteAddr)
		for _, k := range SortedKeys(r.Header) {
			LogStrings(req, fmt.Sprintf("Header.%v", k), r.Header[k])
		}

		if authorizations, ok := r.Header["Authorization"]; ok {
			for _, authorization := range authorizations {
				if strings.HasPrefix(authorization, "Bearer ") {
					token := authorization[7:]
					jwt, err := DecodeJWT(token)
					if jwt != "" {
						LogString(req, "JWT.Payload", jwt)
					}
					if err != nil {
						log.Printf("WARN Invalid JWT: %v", err)
					}
				}
			}
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}

func SortedKeys(m map[string][]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
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
	log.Printf("%v %-40v = %v\n", req, name, value)
}

func LogStrings(req string, name string, value []string) {
	log.Printf("%v %-40v = %v\n", req, name, value)
}
