package helpers

import "net/http"

func CheckAuth(w http.ResponseWriter, r *http.Request, username string, password string) (ok bool) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Authentication required"`)
	rUsername, rPassword, rOk := r.BasicAuth()
	if rOk && rUsername == username && rPassword == password {
		return true
	} else {
		http.Error(w, "Not authorized", 401)
		return false
	}
}
