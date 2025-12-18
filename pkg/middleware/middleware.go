package middleware

import (
	"net/http"

	"fmt"
)

func InternalOnlyEndpointMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			alertMessage := "Access Forbidden: This endpoint is for internal use only."
			homepagePath := "/"
			fmt.Fprintf(w,
				`<!DOCTYPE html>
					<html>
					<head>
					<title>Forbidden</title>
					</head>
					<body>
					<script>
					alert("%s");
					window.location.href = "%s";
					</script>
					</body>
					</html>`,
				alertMessage,
				homepagePath,
			)

			return
		}

		next(w, r)
	}
}
