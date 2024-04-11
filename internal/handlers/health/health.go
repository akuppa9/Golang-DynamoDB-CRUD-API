package health

import(
	"errors"
	"net/http"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/repository/adapter"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/handlers"
)

type Handler struct{
	handlers.Interface
	Repository adapter.Interface
}

func NewHandler(repository adapter.Interface){
	return &Handler{
		Repository: repository, 
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request){
	if !h.Repository.Health(){
		HttpStatus.StatusInternalServerError(w, r, errors.New("Relational database not live"))
		return
	}
	HttpStatus.StatusOK(w, r, "Service OK")
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request){
	HttpStatus.StatusMethodNotAllowed(w, r)
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request){
	HttpStatus.StatusMethodNotAllowed(w, r)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request){
	HttpStatus.StatusMethodNotAllowed(w, r)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request){
	HttpStatus.StatusMethodNoContent(w, r)
}