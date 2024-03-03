package velox

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/FernandoJVideira/velox/filesystems"
	"github.com/gabriel-vasile/mimetype"
)

func (v *Velox) UploadFile(r *http.Request, destination, field string, fs filesystems.FS) error {
	fileName, err := v.getFileToUpload(r, field)
	if err != nil {
		v.ErrorLog.Println(err)
		return err
	}

	if fs != nil {
		err = fs.Put(fileName, destination)
		if err != nil {
			v.ErrorLog.Println(err)
			return err
		}
	} else {
		err = os.Rename(fileName, fmt.Sprintf("%s/%s", destination, path.Base(fileName)))
		if err != nil {
			v.ErrorLog.Println(err)
			return err
		}
	}
	defer func() {
		_ = os.Remove(fileName)
	}()

	return nil
}

func (v *Velox) getFileToUpload(r *http.Request, fieldName string) (string, error) {
	_ = r.ParseMultipartForm(v.config.uploads.maxUploadSize)

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		return "", err
	}

	//go back to the beginning of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	if !inSlice(v.config.uploads.allowedMimeTypes, mimeType.String()) {
		v.ErrorLog.Println(v.config.uploads.allowedMimeTypes)
		return "", errors.New("invalid file type")
	}

	dst, err := os.Create(fmt.Sprintf("/tmp/%s", header.Filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/tmp/%s", header.Filename), nil
}

func inSlice(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
