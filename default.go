package kurt

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DefaultRequestErrorHandler(err error) (int, string) {
	return http.StatusBadRequest, fmt.Sprintf(`{"error": "%v"}`, err)
}

func DefaultResponseErrorHandler(err error) (int, string) {
	return http.StatusInternalServerError, fmt.Sprintf(`{"error": "%v"}`, err)
}

func Default() Kurtis {
	return Kurtis{
		MarshalJSON:          json.Marshal,
		UnmarshalJSON:        json.Unmarshal,
		RequestErrorHandler:  DefaultRequestErrorHandler,
		ResponseErrorHandler: DefaultResponseErrorHandler,
	}
}

var DefaultKurtis = Default()

func Handle(h Handler) http.HandlerFunc {
	return DefaultKurtis.Handle(h)
}
