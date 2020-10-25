package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	keys = map[string]bool{}
)

type RequestId struct {
	sync.Mutex
	count int
}

func (req *RequestId) next() string {
	req.Lock()
	defer req.Unlock()

	req.count++

	return fmt.Sprintf("%v#%v", time.Now().Format(time.RFC3339), req.count)
}

func main() {
	port := flag.Int("port", 8080, "Server port to listen.")
	response := flag.String("response", "", "Response to send.")
	logs := flag.String("log", "*", "Comma separated list of keys to be logged.\n"+
		"Supported keys: \n"+
		"- *\n"+
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
		"- Body (not in default)\n"+
		"\n"+
		"\"Header.*\" will match all headers.\n"+
		"Alternatively list the explicit headers (e.g. \"Header.Accept\").\n"+
		"\n"+
		"The body can be very large and is therefore not part of the \"*\" group.\n"+
		"To log all keys and the body you need to specify -log '*,body'"+
		"\n"+
		"The default will log everything except the body.\n")
	flag.Parse()

	if *logs != "" {
		for _, key := range strings.Split(*logs, ",") {
			keys[strings.ToLower(key)] = true
		}
	}

	requestId := RequestId{}

	log.Println("Started httplogger", *port, *logs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, *response)
		if err != nil {
			log.Printf("WARN Failed writing response: %v", err)
		}

		req := requestId.next()

		if MatchesKey("*", "Proto") {
			LogString(req, "Proto", html.EscapeString(r.Proto))
		}
		if MatchesKey("*", "Host") {
			LogString(req, "Host", r.Host)
		}
		if MatchesKey("*", "RequestURI") {
			LogString(req, "RequestURI", html.EscapeString(r.RequestURI))
		}
		if MatchesKey("*", "Method") {
			LogString(req, "Method", r.Method)
		}
		if MatchesKey("*", "URL") {
			LogString(req, "URL", html.EscapeString(r.URL.Path))
		}
		if MatchesKey("*", "RemoteAddr") {
			LogString(req, "RemoteAddr", r.RemoteAddr)
		}
		if MatchesKey("*", "TransferEncoding") {
			LogStrings(req, "TransferEncoding", r.TransferEncoding)
		}
		for _, k := range SortedKeys(r.Header) {
			key := fmt.Sprintf("Header.%v", k)
			if MatchesKey("*", "Header.*", key) {
				LogStrings(req, key, r.Header[k])
			}
		}

		if authorizations, ok := r.Header["Authorization"]; ok {
			for _, authorization := range authorizations {
				if strings.HasPrefix(authorization, "Bearer ") {
					token := strings.TrimPrefix(authorization, "Bearer ")
					if MatchesKey("*", "JWT.*", "JWT.Header") {
						if jwtHeader, err := DecodeJWT(token, 0); err == nil {
							LogString(req, "JWT.Header", jwtHeader)
						} else {
							log.Printf("WARN Invalid JWT header: %v", err)
						}
					}
					if MatchesKey("*", "JWT.*", "JWT.Payload") {
						if jwtPayload, err := DecodeJWT(token, 1); err == nil {
							LogString(req, "JWT.Payload", jwtPayload)
						} else {
							log.Printf("WARN Invalid JWT payload: %v", err)
						}
					}
				}
			}
		}

		if MatchesKey("Body") {
			if bodyBytes, err := ioutil.ReadAll(r.Body); err == nil {
				body := EncodeBackslash(string(bodyBytes))
				if body != "" {
					LogString(req, "Body", body)
				}
			} else {
				log.Printf("WARN Failed to read body: %v", err)
			}
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}

func MatchesKey(key ...string) bool {
	for _, k := range key {
		if _, ok := keys[strings.ToLower(k)]; ok {
			return true
		}
	}
	return false
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

func EncodeBackslash(str string) string {
	result := strings.Builder{}

	for _, c := range str {
		switch c {
		case '\n':
			result.WriteString("\\n")
		default:
			result.WriteString(string(c))
		}
	}

	return result.String()
}
