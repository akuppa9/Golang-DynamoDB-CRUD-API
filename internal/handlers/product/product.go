package product

import(
	"errors"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
	"time"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/repository/adapter"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/controllers/product"
	EntityProduct "github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/entities/product"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/handlers"
	Rules "github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/rules"
	RulesProduct "github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/rules/product"
	HttpStatus "github.com/akuppa9/Golang-DynamoDB-CRUD-API/utils/http"
)

type Handler struct{
	handler.Interface
	Controller product.Interface
	Rules Rules.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface{
	return &Handler{
		Controller: product.NewController(repository),
		Rules: RulesProduct.NewRules(),
	}
}

func (h *Handler)Get(w http.ResponseWriter, r *http.Request){
	if chi.URLParam(r, "ID") != ""{
		h.getOne(w, r)
	} else{
		h.getAll(w, r)
	}
}

func (h *Handler)getOne(w http.ResponseWriter, r *http.Request){
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil{
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}
	response, err := h.Controller.ListOne(ID)
	if err != nil{
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler)getAll(w http.ResponseWriter, r *http.Request){
	response, err := h.Controller.ListAAll()
	if err != nil{
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request){
	productBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	ID, err := h.Controller.Create(productBody)

	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusOK(w, r, map[string]interface{}{"id": ID.String()})
}

func (h *Handler)Put(w http.ResponseWriter, r *http.Request){
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil{
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}
	productBody, err := h.getBodyAndValidate(r, ID)
	if err != nil{
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}
	if err := h.Controller.Update(ID, productBody); err != nil{
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler)Delete(w http.ResponseWriter, r *http.Request){
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil{
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}
	if err := h.Controller.Remove(ID); err != nil{
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler)Options(w http.ResponseWriter, r *http.Request){
	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID)(*EntityProduct.Product){
	productBody := &EntityProduct.Product{}
	body, err := h.Rules.ConvertIOReaderToStruct(r.Body, productBody)
	if err != nil{
		return &EntityProduct.Product{}, errors.New("body is required")
	}
	productParsed, err := EntityProduct.InterfaceToModel(body)
	if err != nil{
		return &EntityProduct.Product{}, errors.New("error on converting body to model")
	}
	setDefaultValues(productParsed, ID)
	return productParsed, h.Rule.Validate(productParsed)
}

func setDefaultValues(product *EntityProduct.Product, ID uuid.UUID){
	product.UpdatedAt = time.Now()
	if ID == uuid.Nil{
		product.ID = uuid.New()
		product.CreatedAt = time.Now
	}else{
		product.ID = ID
	}
}