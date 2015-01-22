package dispatch

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

const (
	ContentTypeText = "text/plain"
	ContentTypeJson = "application/json"
	ContentTypeXml  = "application/xml"
)

type DispatchResponse struct{ Writer http.ResponseWriter }

func Response(w http.ResponseWriter) *DispatchResponse {
	return &DispatchResponse{Writer: w}
}

func (d *DispatchResponse) Abort(code int) {
	d.Writer.WriteHeader(code)
}

func writeHeaders(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
}

func (d *DispatchResponse) String(code int, format string, data ...interface{}) error {
	writeHeaders(d.Writer, code, ContentTypeText)

	var err error
	if len(data) > 0 {
		_, err = d.Writer.Write([]byte(fmt.Sprintf(format, data...)))
	} else {
		_, err = d.Writer.Write([]byte(format))
	}
	return err
}

func (d *DispatchResponse) JSON(code int, data interface{}) error {
	writeHeaders(d.Writer, code, ContentTypeJson)
	encoder := json.NewEncoder(d.Writer)
	return encoder.Encode(data)
}

func (d *DispatchResponse) XML(code int, data interface{}) error {
	writeHeaders(d.Writer, code, ContentTypeXml)
	encoder := xml.NewEncoder(d.Writer)
	return encoder.Encode(data)
}

func (d *DispatchResponse) Redirect(code int, location string) error {
	d.Writer.Header().Set("Location", location)
	d.Writer.WriteHeader(code)
	return nil
}

func (d *DispatchResponse) NotFound() {
	d.Abort(http.StatusNotFound)
}
