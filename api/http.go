package api

import (
	"encoding/json"
	errcross2 "github.com/ercross/errcross/errcross"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	errs "github.com/pkg/errors"
)

type ErrcrossHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	errcrossService errcross2.ErrcrossService
}

func NewHandler(errcrossService errcross2.ErrcrossService) ErrcrossHandler {
	return &handler{errcrossService}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	e, err := h.errcrossService.Find(key)
	if err != nil {
		if errs.Cause(err) == errcross2.ErrKeyNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, e.URL, http.StatusMovedPermanently)
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	e, err := decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = h.errcrossService.Store(e)
	if err != nil {
		if errs.Cause(err) == errcross2.ErrKeyNotFound {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseBody, err := encode(e)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(w, responseBody, http.StatusCreated)
}

func setupResponse(w http.ResponseWriter, responseBody []byte, statusCode int) {
	w.Header().Set("Content-type", "JSON")
	w.WriteHeader(statusCode)
	_, err := w.Write(responseBody)
	if err != nil {
		log.Println(err)
	}
}

func decode(input []byte) (*errcross2.Errcross, error) {
	e := errcross2.Errcross{}

	if err := json.Unmarshal(input, &e); err != nil {
		return nil, errs.Wrap(err, "errcross.Decode")
	}
	return &e, nil
}

func encode(input *errcross2.Errcross) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errs.Wrap(err, "errcross.Encode")
	}
	return rawMsg, nil
}
