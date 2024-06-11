package main

import (
	"errors"
	"strings"
)

var (
	ErrHTTPInvalidRequest      = errors.New("Invalid http request")
	ErrHTTPInvalidaRequestLine = errors.New("Invalid http request line")
	ErrHTTPInvalidHeader       = errors.New("Invalid http header")
)

type RequestLine struct {
	Method  string
	Path    string
	Version string
}

type HTTPRequest struct {
	RequestLine *RequestLine
	Headers     map[string]string
	Body        string
}

func (r *HTTPRequest) String() string {
	return r.RequestLine.Method + " " + r.RequestLine.Path + " " + r.RequestLine.Version + "\r\n" + r.Body
}

func parseHTTPRequestLine(requestLine string) (*RequestLine, error) {
	requestParts := strings.Split(requestLine, " ")
	if len(requestParts) < 3 {
		return nil, ErrHTTPInvalidaRequestLine
	}

	rl := RequestLine{
		Method:  requestParts[0],
		Path:    requestParts[1],
		Version: requestParts[2],
	}
	return &rl, nil
}

func ParseHTTPRequest(data string) (*HTTPRequest, error) {
	requestData := strings.Split(string(data), "\r\n")
	if len(requestData) < 1 {
		return nil, ErrHTTPInvalidRequest
	}

	requestLine := requestData[0]
	rl, err := parseHTTPRequestLine(requestLine)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	bodySeperator := 0

	for i, line := range requestData[1:] {
		if line == "" {
			bodySeperator = i
			break
		}
		headerParts := strings.Split(line, ": ")
		if len(headerParts) < 2 {
			return nil, ErrHTTPInvalidHeader
		}
		headers[headerParts[0]] = headerParts[1]
	}

	httpRequest := HTTPRequest{
		RequestLine: rl,
		Headers:     headers,
		Body:        requestData[bodySeperator+1],
	}
	return &httpRequest, nil

}

func (r *HTTPRequest) getValueFromDynamicPath(dynamicPath string) string {
	splittedPath := strings.Split(dynamicPath, "/")
	actualSplittedPath := strings.Split(r.RequestLine.Path, "/")
	for idx, path := range splittedPath {
		if path == ":dynamic" {
			return actualSplittedPath[idx]
		}
	}
	return ""
}
