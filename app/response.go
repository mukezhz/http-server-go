package main

type Response struct {
	StatusCode int
	Body       string
	Header     map[string]string
}

func NewResponse(statusCode int, body string) *Response {
	return &Response{
		StatusCode: statusCode,
		Body:       body,
	}
}

func (r *Response) String() string {
	headers := ""
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			headers += k + ": " + v + "\r\n"
		}
	}
	if r.StatusCode == 200 {
		return okResponse() + headers + "\r\n" + r.Body
	}
	return notFoundResponse()
}

func okResponse() string {
	return "HTTP/1.1 200 OK\r\n"
}

func notFoundResponse() string {
	return "HTTP/1.1 404 Not Found\r\n\r\n"
}
