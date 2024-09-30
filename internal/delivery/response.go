package delivery

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/s3nn1k/ef-mob-task/internal/models"
)

// Statuses that will return to user with response in json body
const (
	statusOk  = "Ok"
	statusErr = "Error"
)

// type Response represents json body of response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"error,omitempty"`
	Result  any    `json:"result,omitempty"`
}

// AsLogValue represents Response struct as slog.Value
// Used for logging
func (r *Response) AsLogValue() slog.Value {
	var logValues []slog.Value

	songs, ok := r.Result.([]models.Song)
	if !ok {
		verses, ok := r.Result.([]string)
		if !ok {
			logValues = append(logValues, slog.AnyValue(r.Result))
		} else {
			for _, verse := range verses {
				logValues = append(logValues, slog.StringValue(verse))
			}
		}
	} else {
		for _, song := range songs {
			logValues = append(logValues, song.AsLogValue())
		}
	}

	return slog.GroupValue(
		slog.String("status", r.Status),
		slog.String("message", r.Message),
		slog.Any("result", logValues),
	)
}

// Ok is an alias func to create success response
func Ok(result any) Response {
	return Response{
		Status: statusOk,
		Result: result,
	}
}

// Error is an alias func to create Error response
func Error(msg string) Response {
	return Response{
		Status:  statusErr,
		Message: msg,
	}
}

// response send's response and log's it
func (h *Handler) response(w http.ResponseWriter, r Response, status int) {
	data, err := json.Marshal(r)
	if err != nil {
		msg := "Can't marshal response"

		h.log.Error(msg+": "+err.Error(), "input", r.AsLogValue())

		r = Error(msg)
		status = http.StatusInternalServerError
	}

	h.log.Info("Response", slog.Any("response", r.AsLogValue()), slog.Int("status", status))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
