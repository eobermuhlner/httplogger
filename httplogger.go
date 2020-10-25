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
	port := flag.Int("port", 8080, "Server port to listen.")
	response := flag.String("response", "", "Response to send.")
	logs := flag.String("log", "", "Comma separated list of keys to be logged.\n"+
		"Supported keys: \n"+
		"- Proto\n"+
		"- Host\n"+
		"- RequestURI\n"+
		"- Method\n"+
		"- URL\n"+
		"- RemoteAddr\n"+
		"- TransferEncoding\n"+
		"- Header.*\n"+
		"- JWT.*\n"+
		"- JWT.PayLoad\n"+
		"- JWT.Header\n"+
		"\"Header.*\" will match all headers.\n"+
		"Alternatively list the explicit headers (e.g. \"Header.Accept\").\n"+
		"The default will log everything.")
	flag.Parse()

	var keys = map[string]bool{}
	if *logs != "" {
		for _, key := range strings.Split(*logs, ",") {
			keys[strings.ToLower(key)] = true
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

		if MatchesKey("Proto", keys) {
			LogString(req, "Proto", html.EscapeString(r.Proto))
		}
		if MatchesKey("Host", keys) {
			LogString(req, "Host", r.Host)
		}
		if MatchesKey("RequestURI", keys) {
			LogString(req, "RequestURI", html.EscapeString(r.RequestURI))
		}
		if MatchesKey("Method", keys) {
			LogString(req, "Method", r.Method)
		}
		if MatchesKey("URL", keys) {
			LogString(req, "URL", html.EscapeString(r.URL.Path))
		}
		if MatchesKey("RemoteAddr", keys) {
			LogString(req, "RemoteAddr", r.RemoteAddr)
		}
		if MatchesKey("TransferEncoding", keys) {
			LogStrings(req, "TransferEncoding", r.TransferEncoding)
		}
		for _, k := range SortedKeys(r.Header) {
			key := fmt.Sprintf("Header.%v", k)
			if MatchesKey("Header.*", keys) || MatchesKey(key, keys) {
				LogStrings(req, key, r.Header[k])
			}
		}

		if authorizations, ok := r.Header["Authorization"]; ok {
			for _, authorization := range authorizations {
				if strings.HasPrefix(authorization, "Bearer ") {
					token := strings.TrimPrefix(authorization, "Bearer ")
					if MatchesKey("JWT.*", keys) || MatchesKey("JWT.Header", keys) {
						if jwtHeader, err := DecodeJWT(token, 0); err == nil {
							LogString(req, "JWT.Header", jwtHeader)
						} else {
							log.Printf("WARN Invalid JWT header: %v", err)
						}
					}
					if MatchesKey("JWT.*", keys) || MatchesKey("JWT.Payload", keys) {
						if jwtPayload, err := DecodeJWT(token, 1); err == nil {
							LogString(req, "JWT.Payload", jwtPayload)
						} else {
							log.Printf("WARN Invalid JWT payload: %v", err)
						}
					}
				}
			}
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}

func MatchesKey(key string, keys map[string]bool) bool {
	if len(keys) == 0 {
		return true
	}

	_, ok := keys[strings.ToLower(key)]
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

func DecodeJWT(token string, index int) (string, error) {
	splitToken := strings.Split(token, ".")
	if len(splitToken) >= index+1 {
		jwt, err := base64.RawURLEncoding.DecodeString(splitToken[index])
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

func LogStrings(req string, name string, values []string) {
	for _, value := range values {
		log.Printf("%v %-40v = %v\n", req, name, value)
	}
}
