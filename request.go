package kurt

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

type Request struct {
	*http.Request

	kurtis Kurtis

	errChan chan error

	body     []byte
	bodyErr  error
	bodyOnce *sync.Once
}

func newRequest(r *http.Request, kurt Kurtis) Request {
	return Request{
		Request:  r,
		kurtis:   kurt,
		errChan:  make(chan error, 1),
		bodyOnce: new(sync.Once),
	}
}

func (r Request) Body() []byte {
	r.bodyOnce.Do(func() {
		body, err := ioutil.ReadAll(r.Request.Body)
		if err != nil {
			r.errChan <- errors.Wrap(err, "failed to read body")
		} else {
			r.body = body
		}
	})

	return r.body
}

func (r Request) String() string {
	return string(r.Body())
}

func (r Request) JSON(v interface{}) {
	if err := r.kurtis.UnmarshalJSON(r.Body(), v); err != nil {
		r.errChan <- errors.Wrap(err, "failed to unmarshal json")
	}
}
