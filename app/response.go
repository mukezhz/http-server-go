package main

type Response struct {
	StatusCode int
	Body       string
}

func NewResponse(statusCode int, body string) *Response {
	return &Response{
		StatusCode: statusCode,
		Body:       body,
	}
}

func (r *Response) String() string {
	if r.StatusCode == 200 {
		return okResponse() + r.Body
	}
	return notFoundResponse()
}

func okResponse() string {
	return "HTTP/1.1 200 OK\r\n\r\n"
}

func notFoundResponse() string {
	return "HTTP/1.1 404 Not Found\r\n\r\n"
}
