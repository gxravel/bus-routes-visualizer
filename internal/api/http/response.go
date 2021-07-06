package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
)

type MIME string

func (m MIME) String() string { return string(m) }

const (
	MIMEImagePNG        MIME = "image/png"
	MIMEApplicationJSON MIME = "application/json"
)

const (
	headerContentType   = "Content-Type"
	headerContentLength = "Content-Length"
)

type RangeItemsResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

// Response describes http response for api v1.
type Response struct {
	Data  interface{}    `json:"data,omitempty"`
	Error *ierr.APIError `json:"error,omitempty"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	if data == nil {
		w.WriteHeader(code)
		return
	}

	w.Header().Set(headerContentType, MIMEApplicationJSON.String())

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		log.
			FromContext(ctx).
			WithErr(err).
			Error("encode data to respond with json")
	}
}

func RespondBytes(ctx context.Context, w http.ResponseWriter, code int, mime MIME, size int64, data []byte) {
	w.Header().Set(headerContentType, mime.String())
	w.Header().Set(headerContentLength, strconv.FormatInt(size, 10))

	w.WriteHeader(code)

	if _, err := w.Write(data); err != nil {
		log.
			FromContext(ctx).
			WithErr(err).
			Error("write data")
	}
}

func RespondImageOK(ctx context.Context, w http.ResponseWriter, size int64, data []byte) {
	RespondBytes(ctx, w, http.StatusAccepted, MIMEImagePNG, size, data)
}

func RespondEmpty(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func RespondCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// RespondDataOK responds with 200 status code and JSON in format: {"data": <val>}.
func RespondDataOK(ctx context.Context, w http.ResponseWriter, val interface{}) {
	RespondData(ctx, w, http.StatusOK, val)
}

// RespondEmptyItems responds with empty items and 200 status code.
func RespondEmptyItems(ctx context.Context, w http.ResponseWriter) {
	RespondData(ctx, w, http.StatusOK, RangeItemsResponse{})
}

// RespondData responds with custom status code and JSON in format: {"data": <val>}.
func RespondData(ctx context.Context, w http.ResponseWriter, code int, val interface{}) {
	RespondJSON(ctx, w, code, &Response{
		Data: val,
	})
}

// RespondError converts error to Reason, resolves http status code and responds with APIError.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) {
	reason := ierr.ConvertToReason(err)
	code := ierr.ResolveStatusCode(ierr.Cause(reason.Err))

	RespondJSON(ctx, w, code, &Response{
		Error: &ierr.APIError{
			Reason: &ierr.APIReason{
				RType:   reason.RType,
				Err:     reason.Error(),
				Message: reason.Message,
			},
		},
	})
}
