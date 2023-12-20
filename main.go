package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

var addr string
var auth string

func init() {
	flag.StringVar(&addr, "l", ":8888", "listen addr")
	flag.StringVar(&auth, "a", "", "base64_encde(user:password)")
	flag.Parse()

	addrEnv, b := os.LookupEnv("GOPROXY_ADDR")
	if b {
		addr = addrEnv
	}
	authEnv, b := os.LookupEnv("GOPROXY_AUTH")
	if b {
		auth = authEnv
	}
}

func main() {
	if len(auth) == 0 {
		fmt.Printf("listen: %s\n", addr)
	}else {
		fmt.Printf("listen: %s, auth: %s\n", addr, auth)
	}
	handleProxy()
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func basicProxyAuth(proxyAuth string) (auth string, ok bool) {
	if proxyAuth == "" {
		return
	}

	if !strings.HasPrefix(proxyAuth, "Basic ") {
		return
	}

	c := strings.TrimPrefix(proxyAuth, "Basic ")

	return c, true
}

func handleProxy() {
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// lookup ip
			ip,err := net.ResolveIPAddr("ip", r.Host)
			if err != nil {
				http.Error(w,"resolve ip err"+err.Error(),http.StatusServiceUnavailable)
				return
			}

			os.Getenv("HTTP_PROXY")

			dump, _ := httputil.DumpRequest(r, false)
			fmt.Printf("[http] %s -> %s/%s\n%s", r.RemoteAddr, r.Host,ip, string(dump))
			// fmt.Printf("addr:%s,method:%s,host:%s,header:%s\n", r.RemoteAddr, r.Method, r.Host, r.Header)

			if auth != "" {
				proxyAuth, ok := basicProxyAuth(r.Header.Get("Proxy-Authorization"))
				if !ok {
					w.Header().Set("Proxy-Authenticate", `Basic realm=Restricted`)
					http.Error(w, "proxy auth required", http.StatusProxyAuthRequired)
					return
				}

				if proxyAuth != auth {
					http.Error(w, "proxy authentication failed", http.StatusForbidden)
					return

				}
				r.Header.Del("Proxy-Authorization")
			}

			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),

		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Fatal(server.ListenAndServe())
}
