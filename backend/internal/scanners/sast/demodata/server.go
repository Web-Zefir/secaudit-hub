package democorp

import "net/http"

func Start(mux http.Handler) error {
	return http.ListenAndServe(":8080", mux)
}
