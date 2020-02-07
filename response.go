package kurt

import "net/http"

const headerContentType = "Content-Type"

type Response struct {
	status int
	header http.Header
	body   []byte
	json   interface{}
}

func Respond() Response {
	return Response{status: http.StatusOK}
}

func (r Response) Status(status int) Response {
	newR := r.copy()
	newR.status = status
	return newR
}

func (r Response) Header(key, value string) Response {
	newR := r.copy()
	newR.initHeader()
	newR.header.Add(key, value)
	return newR
}

func (r Response) Headers(header http.Header) Response {
	newR := r.copy()
	newR.initHeader()

	for key, values := range header {
		for _, value := range values {
			newR.header.Add(key, value)
		}
	}

	return newR
}

func (r Response) Body(body []byte) Response {
	newR := r.copy()
	newR.json = nil
	newR.header.Del(headerContentType)
	newR.body = body
	return newR
}

func (r Response) Bytes(body []byte) Response {
	return r.Body(body)
}

func (r Response) String(body string) Response {
	return r.Body([]byte(body))
}

func (r Response) JSON(body interface{}) Response {
	newR := r.copy()
	newR.json = body
	return newR.Header(headerContentType, "application/json")
}

func (r Response) initHeader() {
	if r.header == nil {
		r.header = make(http.Header)
	}
}

func (r Response) copy() Response {
	newR := r
	newR.header = copyHeader(r.header)
	newR.body = copyBytes(r.body)
	return newR
}

func copyHeader(in http.Header) http.Header {
	out := make(http.Header)
	for key, values := range in {
		for _, value := range values {
			out.Add(key, value)
		}
	}

	return out
}

func copyBytes(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)
	return out
}
