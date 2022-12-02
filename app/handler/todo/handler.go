package todo

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"wb/app/controller/todo"
)

type Handler struct {
	todo.Contract
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

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	todos := h.Contract.Get()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&todos)
}

func (h Handler) Insert(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Desc string `json:"desc,omitempty"`
	}

	var response struct {
		Message string `json:"message,omitempty"`
	}

	response.Message = "Success to Insert"
	json.NewDecoder(r.Body).Decode(&request)
	h.Contract.Insert(request.Desc)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID     int  `json:"id,omitempty"`
		Status bool `json:"status,omitempty"`
	}

	var response struct {
		Message string `json:"message,omitempty"`
	}

	response.Message = "Success to Update"

	json.NewDecoder(r.Body).Decode(&request)
	h.Contract.Update(request.ID, request.Status)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)

}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {

	var response struct {
		Message string `json:"message,omitempty"`
	}

	response.Message = "Success to Delete"

	h.Contract.Delete()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)

}
