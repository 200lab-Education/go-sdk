package imgprocessing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/200Lab-Education/go-sdk/sdkcm"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (imgproc *imgProcessing) Resize(file *multipart.FileHeader, folder string, longEdge int, quality int) (*sdkcm.Image, error) {
	f, err := file.Open()
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// body to store request body
	var requestBody bytes.Buffer

	// multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	// initial file field
	fileWriter, err := multipartWriter.CreateFormFile("file", file.Filename)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// copy file to form file
	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// add folder field
	folderWriter, err := multipartWriter.CreateFormField("folder")
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}
	_, err = folderWriter.Write([]byte(folder))
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// add long edge field
	longEdgeWriter, err := multipartWriter.CreateFormField("long_edge")
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}
	_, err = longEdgeWriter.Write([]byte(string(longEdge)))
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// add quality field
	qualityWriter, err := multipartWriter.CreateFormField("quality")
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}
	_, err = qualityWriter.Write([]byte(string(quality)))
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// close multipart writer
	if err = multipartWriter.Close(); err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// new request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", imgproc.cfg.host, "resize"), &requestBody)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// set header to multipart
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// do the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}
	defer response.Body.Close()

	var result Response

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	if result.StatusCode == 200 && result.Data != nil {
		err = nil
	} else {
		err = result
	}

	return result.Data, err
}

func (imgproc *imgProcessing) ResizeFile(filePath string, folder string, longEdge int, quality int) (*sdkcm.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}
	defer f.Close()

	// body to store request body
	var requestBody bytes.Buffer

	// multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	// initial file field
	fileWriter, err := multipartWriter.CreateFormFile("file", f.Name())
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// copy file to form file
	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// add folder field
	folderWriter, err := multipartWriter.CreateFormField("folder")
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}
	_, err = folderWriter.Write([]byte(folder))
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// add long edge field
	longEdgeWriter, err := multipartWriter.CreateFormField("long_edge")
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}
	_, err = longEdgeWriter.Write([]byte(strconv.Itoa(longEdge)))
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// add quality field
	qualityWriter, err := multipartWriter.CreateFormField("quality")
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}
	_, err = qualityWriter.Write([]byte(strconv.Itoa(quality)))
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// close multipart writer
	if err = multipartWriter.Close(); err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// new request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", imgproc.cfg.host, "resize"), &requestBody)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	// set header to multipart
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// do the request
	client := &http.Client{Timeout: 3 * time.Minute}
	response, err := client.Do(req)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}
	defer response.Body.Close()

	var result Response

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.ErrCannotProcessImage)
	}

	if result.StatusCode == 200 && result.Data != nil {
		err = nil
	} else {
		err = result.AppError
	}

	return result.Data, err
}
