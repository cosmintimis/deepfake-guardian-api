package restful

import (
	"net/http"
)

func (app *restfulApi) serverStatus(w http.ResponseWriter, r *http.Request) {
	data := app.healthcheck.Status()
	err := JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
