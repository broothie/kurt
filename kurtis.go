package kurt

import "net/http"

type (
	Handler func(r Request) Response

	ErrorHandler func(err error) (status int, body string)

	MarshalJSON   func(v interface{}) (data []byte, err error)
	UnmarshalJSON func(data []byte, v interface{}) error
)

type Kurtis struct {
	MarshalJSON
	UnmarshalJSON
	RequestErrorHandler  ErrorHandler
	ResponseErrorHandler ErrorHandler
}

func (k Kurtis) Handle(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := newRequest(r, k)

		responseChan := make(chan Response)
		doneChan := make(chan struct{})
		go func() {
			defer close(doneChan)

			for {
				select {
				case response := <-responseChan:
					k.writeResponse(w, response)
					return

				case err := <-request.errChan:
					k.requestError(w, err)
					return
				}
			}
		}()

		go func() { responseChan <- h(request) }()
		<-doneChan
	}
}

func (k Kurtis) writeResponse(w http.ResponseWriter, r Response) {
	for key, values := range r.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	if r.status == 0 {
		r.status = http.StatusOK
	}

	w.WriteHeader(r.status)
	if r.body == nil && r.json == nil {
		return
	}

	body := r.body
	if r.json != nil {
		var err error
		if body, err = k.MarshalJSON(r.json); err != nil {
			k.responseError(w, err)
			return
		}
	}

	if _, err := w.Write(body); err != nil {
		k.responseError(w, err)
	}
}

func (k Kurtis) requestError(w http.ResponseWriter, err error) {
	status, body := k.RequestErrorHandler(err)
	http.Error(w, body, status)
}

func (k Kurtis) responseError(w http.ResponseWriter, err error) {
	status, body := k.ResponseErrorHandler(err)
	http.Error(w, body, status)
}
