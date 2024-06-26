package handlers

import (
	"UnlockEdv2/src/models"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (srv *Server) registerProviderPlatformRoutes() {
	srv.Mux.Handle("GET /api/provider-platforms", srv.applyMiddleware(http.HandlerFunc(srv.HandleIndexProviders)))
	srv.Mux.Handle("GET /api/provider-platforms/{id}", srv.applyMiddleware(http.HandlerFunc(srv.HandleShowProvider)))
	srv.Mux.Handle("POST /api/provider-platforms", srv.applyMiddleware(http.HandlerFunc(srv.HandleCreateProvider)))
	srv.Mux.Handle("PATCH /api/provider-platforms/{id}", srv.applyMiddleware(http.HandlerFunc(srv.HandleUpdateProvider)))
	srv.Mux.Handle("DELETE /api/provider-platforms/{id}", srv.applyMiddleware(http.HandlerFunc(srv.HandleDeleteProvider)))
}

func (srv *Server) HandleIndexProviders(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling provider index request")
	page, perPage := srv.GetPaginationInfo(r)
	total, platforms, err := srv.Db.GetAllProviderPlatforms(page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	paginationData := models.NewPaginationInfo(page, perPage, total)
	response := models.PaginatedResource[models.ProviderPlatform]{
		Data: platforms,
		Meta: paginationData,
	}
	log.Info("Found "+strconv.Itoa(int(total)), " provider platforms")
	if err = srv.WriteResponse(w, http.StatusOK, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (srv *Server) HandleShowProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Error("GET Provider handler Error: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	platform, err := srv.Db.GetProviderPlatformByID(id)
	if err != nil {
		log.Error("Error: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.Resource[models.ProviderPlatform]{
		Data: make([]models.ProviderPlatform, 0),
	}
	response.Data = append(response.Data, *platform)
	if err = srv.WriteResponse(w, http.StatusOK, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (srv *Server) HandleCreateProvider(w http.ResponseWriter, r *http.Request) {
	var platform models.ProviderPlatform
	err := json.NewDecoder(r.Body).Decode(&platform)
	if err != nil {
		log.Error("Error decoding request body: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()
	newProv, err := srv.Db.CreateProviderPlatform(&platform)
	if err != nil {
		log.Error("Error creating provider platform: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response := models.Resource[models.ProviderPlatform]{
		Data:    make([]models.ProviderPlatform, 0),
		Message: "Provider platform created successfully",
	}
	response.Data = append(response.Data, *newProv)
	if err = srv.WriteResponse(w, http.StatusOK, &response); err != nil {
		log.Error("Error writing response: ", err.Error())
		srv.ErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func (srv *Server) HandleUpdateProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Error("PATCH Provider handler Error:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var platform models.ProviderPlatform
	err = json.NewDecoder(r.Body).Decode(&platform)
	if err != nil {
		log.Error("Error decoding request body: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()
	updated, err := srv.Db.UpdateProviderPlatform(&platform, uint(id))
	if err != nil {
		log.Error("Error updating provider platform: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response := models.Resource[models.ProviderPlatform]{
		Data: make([]models.ProviderPlatform, 0),
	}
	response.Data = append(response.Data, *updated)
	if err = srv.WriteResponse(w, http.StatusOK, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (srv *Server) HandleDeleteProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Error("DELETE Provider handler Error: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err = srv.Db.DeleteProviderPlatform(id); err != nil {
		log.Error("Error deleting provider platform: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
