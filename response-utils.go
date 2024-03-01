package velox

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"

	"errors"
)

func (v *Velox) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}
	return nil
}

func (v *Velox) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for k, value := range headers[0] {
			w.Header()[k] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (v *Velox) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for k, value := range headers[0] {
			w.Header()[k] = value
		}
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (v *Velox) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile, fileName string) error {
	fp := path.Join(pathToFile, fileName)
	fileToServe := filepath.Clean(fp)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)
	return nil
}

func (v *Velox) Error404(w http.ResponseWriter, r *http.Request) {
	v.ErrorStatus(w, http.StatusNotFound)
}
func (v *Velox) Error500(w http.ResponseWriter, r *http.Request) {
	v.ErrorStatus(w, http.StatusInternalServerError)
}
func (v *Velox) ErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	v.ErrorStatus(w, http.StatusUnauthorized)
}
func (v *Velox) ErrorForbiden(w http.ResponseWriter, r *http.Request) {
	v.ErrorStatus(w, http.StatusForbidden)
}
func (v *Velox) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)

}
