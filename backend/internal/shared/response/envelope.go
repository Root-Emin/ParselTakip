package response

import "net/http"

// Envelope is the standard success response wrapper used across business APIs.
// MessageKey is a stable, client-facing identifier (see module constants) that
// clients can map to localized text; Data carries the payload.
type Envelope struct {
	Success    bool        `json:"success"`
	MessageKey string      `json:"message_key"`
	Data       interface{} `json:"data,omitempty"`
}

// Success writes a 200 OK enveloped response.
func Success(w http.ResponseWriter, messageKey string, data interface{}) {
	JSON(w, http.StatusOK, Envelope{Success: true, MessageKey: messageKey, Data: data})
}

// CreatedEnvelope writes a 201 Created enveloped response.
func CreatedEnvelope(w http.ResponseWriter, messageKey string, data interface{}) {
	JSON(w, http.StatusCreated, Envelope{Success: true, MessageKey: messageKey, Data: data})
}
