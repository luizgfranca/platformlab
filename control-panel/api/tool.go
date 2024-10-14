package api

import (
	"encoding/json"
	"log"
	"net/http"
	"platformlab/controlpanel/model"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Tool struct {
	websocketUpgrader websocket.Upgrader
}

func (t *Tool) GetEventRresponseTEST() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input model.ToolEvent
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		if !input.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "invalid request data"})
			return
		}

		helloWorldMock := model.ToolEvent{
			Class:   model.EventClassOperation,
			Type:    model.EventTypeDisplay,
			Project: "x",
			Tool:    "y",
			Display: &model.DisplayDefniition{
				Type: model.DisplayDefniitionTypeView,
				Elements: &[]model.DisplayElement{
					{Type: "heading", Text: "Hello world", Description: "This is a hello world test"},
				},
			},
		}

		json.NewEncoder(w).Encode(helloWorldMock)
	}
}

// simple echo
// after this i should return a mock of an event
func connectionHandler(connection *websocket.Conn) {
	defer connection.Close()
	var event model.ToolEvent

	for {
		msgtype, message, err := connection.ReadMessage()
		if err != nil {
			log.Print("websocket message receiving error: ", err.Error())
			break
		}

		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Print("websocket message payload parsing error: ", err.Error())
			break
		}

		log.Print("EVENT: class ", event.Class)

		err = connection.WriteMessage(msgtype, message)
		if err != nil {
			log.Print("websocket message sending error: ", err.Error())
			break
		}
	}
}

func upgraderAllowAllOrigins(r *http.Request) bool {
	return true
}

func (t *Tool) ToolClientWebsocket() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := t.websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("websocket upgrade error: ", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		go connectionHandler(connection)
	}
}

func ToolRestAPI(db *gorm.DB) *Tool {
	return &Tool{websocket.Upgrader{
		CheckOrigin: upgraderAllowAllOrigins,
	}}
}
