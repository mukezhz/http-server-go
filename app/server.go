package main

import (
	"flag"
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

type Server struct {
	Listener   net.Listener
	Connection net.Conn
}

func NewServer(l net.Listener) *Server {
	return &Server{
		Listener: l,
	}
}

func (s *Server) handleConnection() {
	c, err := s.Listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	s.Connection = c
	defer s.Connection.Close()
	data := make([]byte, 1024)
	_, err = c.Read(data)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

	httpRequest, err := ParseHTTPRequest(string(data))
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		os.Exit(1)
	}
	resp := NewResponse(200, "")
	s.handleRequest(httpRequest, resp)
}

func (s *Server) Start() {
	for {
		s.handleConnection()
	}
}

var assetsPath string

func init() {
	fmt.Println("Starting server...")
	directory := flag.String("directory", ".", "Specify the directory to use")
	flag.Parse()
	assetsPath = *directory
	fmt.Printf("Using directory: %s\n", *directory)
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!", assetsPath)

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	s := NewServer(l)
	s.Start()
}

func (s *Server) handleRequest(httpRequest *HTTPRequest, resp *Response) {
	if httpRequest.RequestLine.Path == "/" {
		resp.StatusCode = 200
	} else if strings.Contains(httpRequest.RequestLine.Path, "/files/") {
		value := httpRequest.getValueFromDynamicPath("/files/:dynamic")
		content, err := os.ReadFile(assetsPath + value)
		if err != nil {
			resp.StatusCode = 404
		} else {
			resp.Body = string(content)
			resp.Header = map[string]string{
				"Content-Type":   "application/octet-stream",
				"Content-Length": fmt.Sprintf("%d", len(content)),
			}
			resp.StatusCode = 200
		}
	} else if strings.Contains(httpRequest.RequestLine.Path, "/echo/") {
		value := httpRequest.getValueFromDynamicPath("/echo/:dynamic")
		resp.Body = value
		resp.Header = map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(value)),
		}
		resp.StatusCode = 200
	} else if strings.Contains(httpRequest.RequestLine.Path, "/user-agent") {
		resp.Body = httpRequest.Headers["user-agent"]
		resp.Header = map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(httpRequest.Headers["user-agent"])),
		}
	} else {
		resp.StatusCode = 404
	}

	_, err := s.Connection.Write([]byte(resp.String()))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}
