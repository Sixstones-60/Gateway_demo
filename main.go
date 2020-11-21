package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type handle struct {
	host string
	port string
	path string
}

type Service struct {
	hello *handle
	bey   *handle
}

func NewMultipleHostsReverseProxy1(targets []*url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		target := targets[rand.Int()%len(targets)]
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return &httputil.ReverseProxy{Director: director}
}

func (this *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var remote *url.URL

	if strings.Contains(r.RequestURI, "api/hello") {
		remote, _ = url.Parse("http://" + this.hello.host + ":" + this.hello.port)
		proxy := NewMultipleHostsReverseProxy1([]*url.URL{
			{
				Scheme: "http",
				Host:   "localhost:9091",
				Path:   "/hello",
			},
		})
		proxy.ServeHTTP(w, r)
	} else if strings.Contains(r.RequestURI, "api/bey") {
		remote, _ = url.Parse("http://" + this.bey.host + ":" + this.bey.port + this.bey.path)
		fmt.Println(remote)
		proxy := NewMultipleHostsReverseProxy1([]*url.URL{
			{
				Scheme: "http",
				Host:   "localhost:9092",
				Path:   "/bey",
			},
		})
		proxy.ServeHTTP(w, r)

	} else {
		fmt.Fprintf(w, "404 Not Found")
		return

	}
	//proxy := httputil.NewSingleHostReverseProxy(remote)
	//proxy.ServeHTTP(w, r)
}

func startServer() {
	// 注册被代理的服务器 (host， port)
	service := &Service{
		hello: &handle{host: "127.0.0.1", port: "9091"},
		bey:   &handle{host: "127.0.0.1", port: "9092", path: "/bey"},
	}

	err := http.ListenAndServe(":8888", service)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func main() {
	fmt.Println("test")
	startServer()
}
