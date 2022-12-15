package todo

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"wb/app/controller/todo"
)

type Handler struct {
	todo.Contract
}

//	@Description	Request for insert and update
//	@Description	for update, fill id and status
//	@Description	for insert, fill desc
type Request struct {
	// ID: id of todo, for UPDATE
	ID int `json:"id,omitempty" example:"1"`
	// Status: status of todo, for UPDATE
	// * true - Todo is done
	// * false - Todo is not done
	Status bool `json:"status,omitempty" example:"true"`
	// Desc: description of todo, for INSERT
	Desc string `json:"desc,omitempty" example:"my todo"`
}

//	@Description	Response for insert update and delete
type Response struct {
	// Message: message response that will be used for alert
	Message string `json:"message,omitempty"`
}

func (h Handler) Todo(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "todo.html")
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//	@Summary		Get Todo List
//	@Description	Get Todo List by Allowed Session
//	@Tags			Todo
//	@Accept			json
//	@Produce		json
//	@Param			SESSION_TOKEN	header		string	true	"Session Token from Login"
//	@Success		200				{object}	[]model.Todo
//	@Router			/todo/get [get]
func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	userIDstr := r.Context().Value("user_id").(string)
	userID, _ := strconv.ParseInt(userIDstr, 10, 64)
	todos := h.Contract.Get(userID)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&todos)
}

//	@Summary		Insert Todo List
//	@Description	Insert Todo List by Allowed Session
//	@Tags			Todo
//	@Accept			json
//	@Produce		json
//	@Param			SESSION_TOKEN	header		string	true	"Session Token from Login"
//	@Param			data			body		Request	true	"todo data"
//	@Success		200				{object}	Response
//	@Router			/todo/insert [post]
func (h Handler) Insert(w http.ResponseWriter, r *http.Request) {

	var request Request

	var response Response
	userIDstr := r.Context().Value("user_id").(string)
	userID, _ := strconv.ParseInt(userIDstr, 10, 64)
	response.Message = "Success to Insert"
	json.NewDecoder(r.Body).Decode(&request)
	h.Contract.Insert(request.Desc, userID)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

//	@Summary		Update Todo List
//	@Description	Update Todo List by Allowed Session
//	@Tags			Todo
//	@Accept			json
//	@Produce		json
//	@Param			SESSION_TOKEN	header		string	true	"Session Token from Login"
//	@Param			data			body		Request	true	"todo data"
//	@Success		200				{object}	Response
//	@Router			/todo/update [post]
func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	var request Request
	var response Response

	response.Message = "Success to Update"

	json.NewDecoder(r.Body).Decode(&request)
	h.Contract.Update(request.ID, request.Status)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)

}

//	@Summary		Delete Todo List
//	@Description	Soft Delete Todo List that already done by Allowed Session
//	@Tags			Todo
//	@Accept			json
//	@Produce		json
//	@Param			SESSION_TOKEN	header		string	true	"Session Token from Login"
//	@Success		200				{object}	Response
//	@Router			/todo/delete [post]
func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {

	var response Response

	response.Message = "Success to Delete"

	h.Contract.Delete()

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)

}
