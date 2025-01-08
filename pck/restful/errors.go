package restful

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

func (app *restfulApi) reportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	requestAttrs := slog.Group("request", "method", method, "url", url)
	app.logger.Error(message, requestAttrs, "trace", trace)

}

func (app *restfulApi) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		app.reportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *restfulApi) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.reportServerError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (app *restfulApi) notFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorMessage(w, r, http.StatusNotFound, message, nil)
}

func (app *restfulApi) somethingWentWrong(w http.ResponseWriter, r *http.Request) {
	message := "Something went wrong"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (app *restfulApi) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (app *restfulApi) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

func (app *restfulApi) failedValidation(w http.ResponseWriter, r *http.Request, errors []error) {
	jsonErrors := map[string][]error{"errors": errors}

	err := JSON(w, http.StatusUnprocessableEntity, jsonErrors)
	if err != nil {
		app.serverError(w, r, err)
	}
}
