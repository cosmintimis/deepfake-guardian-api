package restful

import (
	"errors"
	"net/http"

	"github.com/cosmintimis/deepfake-guardian-api/pck/business/repositories"
	"github.com/cosmintimis/deepfake-guardian-api/pck/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *restfulApi) wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer conn.Close()

	// Register the connection
	app.connLock.Lock()
	app.connections[conn] = struct{}{}
	app.connLock.Unlock()

	defer func() {
		// Unregister the connection when it is closed
		app.connLock.Lock()
		delete(app.connections, conn)
		app.connLock.Unlock()
	}()

	// Only read messages from the client and log them
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			app.logger.Error("WebSocket: Error reading message", "error", err)
			return
		}
		app.logger.Info("WebSocket: Received message -> ", "message", string(p))
		// Echo the message back to the client
		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	app.logger.Error("WebSocket: Error writing message", "error", err)
		// 	return
		// }
	}

}

func (app *restfulApi) broadcastMessage(messageType int, message []byte) {
	app.connLock.Lock()
	defer app.connLock.Unlock()

	for conn := range app.connections {
		if err := conn.WriteMessage(messageType, message); err != nil {
			app.logger.Error("WebSocket: Error broadcasting message", "error", err)
			conn.Close()
			delete(app.connections, conn) // Remove broken connections
		}
	}
}

func (app *restfulApi) serverStatus(w http.ResponseWriter, r *http.Request) {
	app.broadcastMessage(websocket.TextMessage, []byte("Hello, WebSocket clients!"))
	data := app.healthcheck.Status()
	err := JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *restfulApi) getMediaById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
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
