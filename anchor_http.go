package junkboy

import (
	"fmt"
	"net/http"
	"strconv"
)

type anchorService interface {
	AddAnchor(a Anchor) (int, error)
	UpdateAnchor(a Anchor) error
	GetAnchor(id int) (Anchor, error)
	GetAnchors() ([]Anchor, error)
	DeleteAnchor(id int) error
}

type AnchorHTTPHandler struct {
	service anchorService
}

func NewAnchorHTTPHandler(s anchorService) *AnchorHTTPHandler {
	return &AnchorHTTPHandler{
		service: s,
	}
}

func (h *AnchorHTTPHandler) RegisterRoutes(r *Router) {
	r.AddRoute([]string{"POST", "OPTIONS"}, "/anchor", h.addAnchorHandler)
	r.AddRoute([]string{"GET", "OPTIONS"}, "/anchors", h.getAnchorsHandler)
	r.AddRoute([]string{"GET", "OPTIONS"}, "/anchor/([^/]+)", h.getAnchorHandler)
	r.AddRoute([]string{"PUT", "OPTIONS"}, "/anchor", h.updateAnchorHandler)
	r.AddRoute([]string{"DELETE", "OPTIONS"}, "/anchor/([^/]+)", h.deleteAnchorHandler)
}

func (h *AnchorHTTPHandler) addAnchorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	type Response struct {
		ID int `json:"id"`
	}

	if !contentTypeIsValid(w, r, "application/json") {
		return
	}

	var anchor Anchor
	if err := readJSON(w, r, &anchor); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.AddAnchor(anchor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, Response{ID: id})
}

func (h *AnchorHTTPHandler) getAnchorsHandler(w http.ResponseWriter, r *http.Request) {
	anchors, err := h.service.GetAnchors()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, anchors)
}

func (h *AnchorHTTPHandler) getAnchorHandler(w http.ResponseWriter, r *http.Request) {
	idField := getField(r, 0)
	id, err := strconv.Atoi(idField)

	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid anchor id '%s'", idField))
		return
	}

	anchors, err := h.service.GetAnchor(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, anchors)
}

func (h *AnchorHTTPHandler) updateAnchorHandler(w http.ResponseWriter, r *http.Request) {
	if !contentTypeIsValid(w, r, "application/json") {
		return
	}

	var anchor Anchor
	if err := readJSON(w, r, &anchor); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.UpdateAnchor(anchor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnchorHTTPHandler) deleteAnchorHandler(w http.ResponseWriter, r *http.Request) {
	idField := getField(r, 0)
	id, err := strconv.Atoi(idField)

	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid anchor id '%s'", idField))
		return
	}

	err = h.service.DeleteAnchor(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
