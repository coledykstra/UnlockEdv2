package handlers

import (
	"UnlockEdv2/src/models"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (srv *Server) registerProgramsRoutes() {
	srv.Mux.Handle("GET /api/programs", srv.applyMiddleware(http.HandlerFunc(srv.HandleIndexPrograms)))
	srv.Mux.Handle("GET /api/programs/{id}", srv.applyMiddleware(http.HandlerFunc(srv.HandleShowProgram)))
	srv.Mux.Handle("POST /api/programs", srv.applyMiddleware(http.HandlerFunc(srv.HandleCreateProgram)))
	srv.Mux.Handle("DELETE /api/programs/{id}", srv.applyMiddleware(http.HandlerFunc(srv.HandleDeleteProgram)))
	srv.Mux.Handle("PATCH /api/programs/{id}", srv.applyMiddleware(http.HandlerFunc(srv.HandleUpdateProgram)))
	srv.Mux.Handle("PUT /api/programs/{id}/save", srv.applyMiddleware(http.HandlerFunc(srv.HandleFavoriteProgram)))
}

/*
* @Query Params:
* ?page=: page
* ?perPage=: perPage
* ?sort=: sort
* ?filter=: filter
* ?search=: search
* ?searchFields=: searchFields
 */
func (srv *Server) HandleIndexPrograms(w http.ResponseWriter, r *http.Request) {
	page, perPage := srv.GetPaginationInfo(r)
	search := r.URL.Query().Get("search")
	total, programs, err := srv.Db.GetProgram(page, perPage, search)
	if err != nil {
		log.Debug("IndexPrograms Database Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	last := srv.CalculateLast(total, perPage)
	paginationData := models.PaginationMeta{
		PerPage:     perPage,
		LastPage:    int(last),
		CurrentPage: page,
		Total:       total,
	}
	response := models.PaginatedResource[models.Program]{
		Meta: paginationData,
		Data: programs,
	}
	if err = srv.WriteResponse(w, http.StatusOK, response); err != nil {
		log.Error("Error writing response: " + err.Error())
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (srv *Server) HandleShowProgram(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Debug("GET Program handler Error: ", err)
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	program, err := srv.Db.GetProgramByID(id)
	if err != nil {
		log.Debug("GET Program handler Error: ", err)
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err = srv.WriteResponse(w, http.StatusOK, program); err != nil {
		log.Error("Error writing response: " + err.Error())
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (srv *Server) HandleCreateProgram(w http.ResponseWriter, r *http.Request) {
	var program models.Program
	err := json.NewDecoder(r.Body).Decode(&program)
	defer r.Body.Close()
	if err != nil {
		log.Error("CreateProgram Error:" + err.Error())
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = srv.Db.CreateProgram(&program)
	if err != nil {
		log.Error("Error creating program:" + err.Error())
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (srv *Server) HandleUpdateProgram(w http.ResponseWriter, r *http.Request) {
	var program models.Program
	err := json.NewDecoder(r.Body).Decode(&program)
	defer r.Body.Close()
	if err != nil {
		log.Error("UpdateProgram Error:" + err.Error())
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Debug("GET Program handler Error: ", err)
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
	}
	toUpdate, err := srv.Db.GetProgramByID(id)
	if err != nil {
		log.Error("Error getting program:" + err.Error())
	}
	models.UpdateStruct(&toUpdate, &program)
	updated, updateErr := srv.Db.UpdateProgram(toUpdate)
	if updateErr != nil {
		log.Error("Error updating program:" + err.Error())
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := srv.WriteResponse(w, http.StatusOK, updated); err != nil {
		log.Error("Error writing response: " + err.Error())
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (srv *Server) HandleDeleteProgram(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Error("DELETE Program handler Error: " + err.Error())
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err = srv.Db.DeleteProgram(id); err != nil {
		log.Error("Error deleting program:" + err.Error())
		srv.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	log.Info("Program deleted")
	w.WriteHeader(http.StatusNoContent)
}

func (srv *Server) HandleFavoriteProgram(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Error("Favorite Program handler Error: " + err.Error())
		srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	user_id := srv.GetUserID(r)
	if rows := srv.Db.Conn.Delete(&models.UserFavorite{}, "program_id = ? AND user_id = ?", id, user_id).RowsAffected; rows < 1 {
		log.Info("favoriting program")
		favorite := models.UserFavorite{
			ProgramID: uint(id),
			UserID:    user_id,
		}
		if err := srv.Db.Conn.Create(&favorite).Error; err != nil {
			log.Error("Error creating favorite:" + err.Error())
			srv.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := srv.WriteResponse(w, http.StatusCreated, "program favorited"); err != nil {
			log.Errorf("Error writing response: %s", err.Error())
			srv.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := srv.WriteResponse(w, http.StatusNoContent, "program unfavorited"); err != nil {
			log.Errorf("Error writing response: %s", err.Error())
			srv.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
