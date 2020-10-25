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
	response := flag.String("response", "", "Response")
	logs := flag.String("log", "", "Comma separated list of keys to be logged\n"+
		"Supported keys: \n"+
		"- Method\n"+
		"- URL\n"+
		"- RemoteAddr\n"+
		"- Header.*\n"+
		"- JWT.*\n"+
		"- JWT.PayLoad")
	flag.Parse()

	var keys = map[string]bool{}
	if *logs != "" {
		for _, key := range strings.Split(*logs, ",") {
			keys[key] = true
		}
	}

	log.Println("Started httplogger")
	count := 0

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, *response)
		if err != nil {
			log.Printf("WARN Failed writing response: %v", err)
		}

		count++
		req := fmt.Sprintf("%v#%v", time.Now().Format(time.RFC3339), count)

		if MatchesKeys("Method", keys) {
			LogString(req, "Method", r.Method)
		}
		if MatchesKeys("URL", keys) {
			LogString(req, "URL", html.EscapeString(r.URL.Path))
		}
		if MatchesKeys("RemoteAddr", keys) {
			LogString(req, "RemoteAddr", r.RemoteAddr)
		}
		for _, k := range SortedKeys(r.Header) {
			key := fmt.Sprintf("Header.%v", k)
			if MatchesKeys("Header.*", keys) || MatchesKeys(key, keys) {
				LogStrings(req, key, r.Header[k])
			}
		}

		if authorizations, ok := r.Header["Authorization"]; ok {
			for _, authorization := range authorizations {
				if strings.HasPrefix(authorization, "Bearer ") {
					token := authorization[7:]
					jwtPayload, err := DecodeJWT(token)
					if jwtPayload != "" {
						if MatchesKeys("JWT.*", keys) || MatchesKeys("JWT.Payload", keys) {
							LogString(req, "JWT.Payload", jwtPayload)
						}
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

func MatchesKeys(key string, keys map[string]bool) bool {
	if len(keys) == 0 {
		return true
	}

	_, ok := keys[key]
	return ok
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
