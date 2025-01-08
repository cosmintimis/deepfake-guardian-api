package restful

import (
	"errors"
	"net/http"

	"github.com/cosmintimis/deepfake-guardian-api/pck/business/repositories"
	"github.com/cosmintimis/deepfake-guardian-api/pck/utils"
	"github.com/go-chi/chi/v5"
)

func (app *restfulApi) serverStatus(w http.ResponseWriter, r *http.Request) {
	data := app.healthcheck.Status()
	err := JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *restfulApi) getMediaById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		app.badRequest(w, r, utils.ErrMissingID)
		return
	}
	media, err := app.mediaRepository.GetByID(id)
	if err != nil {
		if errors.Is(err, utils.ErrMediaNotFound) {
			app.notFound(w, r)
			return
		}
		app.somethingWentWrong(w, r)
		return
	}
	err = JSON(w, http.StatusOK, media)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *restfulApi) deleteMediaById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		app.badRequest(w, r, utils.ErrMissingID)
		return
	}
	ok, err := app.mediaRepository.Delete(id)
	if err != nil {
		if errors.Is(err, utils.ErrMediaNotFound) {
			app.notFound(w, r)
			return
		}
		app.somethingWentWrong(w, r)
		return
	}
	err = JSON(w, http.StatusOK, map[string]bool{"deleted": ok})
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *restfulApi) addNewMedia(w http.ResponseWriter, r *http.Request) {
	var payload repositories.MediaPayload
	err := DecodeJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	createdMedia, err := app.mediaRepository.Create(&payload)
	if err != nil {
		app.somethingWentWrong(w, r)
		return
	}
	err = JSON(w, http.StatusCreated, createdMedia)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *restfulApi) updateMedia(w http.ResponseWriter, r *http.Request) {
	var payload repositories.MediaPayload
	err := DecodeJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		app.badRequest(w, r, utils.ErrMissingID)
		return
	}
	updatedMedia, err := app.mediaRepository.Update(id, &payload)
	if err != nil {
		if errors.Is(err, utils.ErrMediaNotFound) {
			app.notFound(w, r)
			return
		}
		app.somethingWentWrong(w, r)
		return
	}
	err = JSON(w, http.StatusOK, updatedMedia)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *restfulApi) getAllMedia(w http.ResponseWriter, r *http.Request) {
	allMedia, err := app.mediaRepository.GetAll()
	if err != nil {
		app.somethingWentWrong(w, r)
		return
	}
	err = JSON(w, http.StatusOK, allMedia)
	if err != nil {
		app.serverError(w, r, err)
	}
}
