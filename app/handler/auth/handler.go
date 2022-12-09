package todo

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"wb/app/controller/auth"
)

type Handler struct {
	auth.Contract
}

func (h Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "register.html")
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "login.html")
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) DoLogin(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	var response struct {
		Message string `json:"message,omitempty"`
	}
	if !ok {
		response.Message = "Format basic auth salah"

		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&response)
		return
	}

	sessionToken := h.Contract.Login(username, password)

	if sessionToken == "" {
		response.Message = "username / password salah"
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&response)
		return
	}

	response.Message = "Selamat Datang"

	http.SetCookie(w, &http.Cookie{
		Name:  "SESSION_TOKEN",
		Value: sessionToken,
		Path:  "/",
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

func (h Handler) DoRegister(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
		Fullname string `json:"fullname,omitempty"`
	}

	var response struct {
		Message string `json:"message,omitempty"`
	}

	json.NewDecoder(r.Body).Decode(&request)
	h.Contract.Register(request.Username, request.Password, request.Fullname)
	response.Message = "Success To Register"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Message string `json:"message,omitempty"`
	}
	response.Message = "Success To Logout"

	cookie, err := r.Cookie("SESSION_TOKEN")
	if err != nil {
		response.Message = "No session"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&response)
		return
	}
	h.Contract.Logout(cookie.Name)

	http.SetCookie(w, &http.Cookie{
		Name:   "SESSION_TOKEN",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}
