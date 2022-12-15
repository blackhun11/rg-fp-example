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

//	@Description	Request body for register
type Request struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Fullname string `json:"fullname,omitempty"`
}

//	@Description	Response for login and register
type Response struct {
	// Message: message response that will be used for alert
	Message string `json:"message,omitempty"`
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

//	@Summary		Login API
//	@Description	Login API using basic auth
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			SESSION_TOKEN	header		string	false	"Session Token from Login"
//  @Security BasicAuth
//	@Success		200				{object}	Response
//	@Failure		401				{object}	Response
//	@Router			/auth/doLogin [post]
func (h Handler) DoLogin(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	var response Response
	if !ok {
		response.Message = "Format basic auth salah"

		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&response)
		return
	}

	sessionToken := h.Contract.Login(username, password)

	if sessionToken == "" {
		response.Message = "username / password salah"
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&response)
		return
	}

	response.Message = "Selamat Datang"

	http.SetCookie(w, &http.Cookie{
		Name:  "SESSION_TOKEN",
		Value: sessionToken,
		Path:  "/",
	})
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

//	@Summary		Register API
//	@Description	Register API
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			SESSION_TOKEN	header		string	false	"Session Token from Login"
//	@Param			data	body		Request	true	"request body"
//	@Success		200				{object}	Response
//	@Router			/auth/doRegister [post]
func (h Handler) DoRegister(w http.ResponseWriter, r *http.Request) {

	var request Request

	var response Response

	json.NewDecoder(r.Body).Decode(&request)
	h.Contract.Register(request.Username, request.Password, request.Fullname)
	response.Message = "Success To Register"

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

//	@Summary		Logout API
//	@Description	Logout API
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			SESSION_TOKEN	header		string	false	"Session Token from Login"
//	@Success		200				{object}	Response
//	@Router			/auth/logout [post]
func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Message string `json:"message,omitempty"`
	}
	response.Message = "Success To Logout"

	cookie, err := r.Cookie("SESSION_TOKEN")
	if err != nil {
		response.Message = "No session"
		w.Header().Add("Content-Type", "application/json")
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
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}
