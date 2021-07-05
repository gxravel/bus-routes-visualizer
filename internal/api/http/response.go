package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	"github.com/gxravel/bus-routes-visualizer/internal/logger"
)

const (
	headerContentType   = "Content-Type"
	headerContentLength = "Content-Length"
	mimeApplicationJSON = "application/json"
	mimeImagePNG        = "image/png"
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

	w.Header().Set(headerContentType, mimeApplicationJSON)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("encoding data to respond with json")
	}
}

func RespondPNG(ctx context.Context, w http.ResponseWriter, path string) {
	file, err := os.Open(path)
	if err != nil {
		logger.FromContext(ctx).WithErr(err).Error("unable to open file")
		return
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		logger.FromContext(ctx).WithErr(err).Error("unable to read file info")
		return
	}

	w.Header().Set(headerContentType, mimeImagePNG)
	w.Header().Set(headerContentLength, strconv.FormatInt(fi.Size(), 10))

	_, err = io.Copy(w, file)
	if err != nil {
		logger.FromContext(ctx).WithErr(err).Error("unable to open file")
		return
	}
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
