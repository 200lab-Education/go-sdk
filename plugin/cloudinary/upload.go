package cloudinary

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/200Lab-Education/go-sdk/sdkcm"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func (cd *cloudinary) VideoUpload(ctx context.Context, filePath string, uploadPreset string, folder string, format string) (*VideoResult, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.CustomError("CannotProcessFile", ""))
	}

	// body to store request body
	var requestBody bytes.Buffer

	// multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	signatureValues := map[string]string{
		"folder":        folder,
		"format":        format,
		"timestamp":     fmt.Sprint(time.Now().Unix()),
		"upload_preset": uploadPreset,
	}

	values := map[string]string{
		"cloud_name": cd.config.cloudName,
		"api_key":    cd.config.apiKey,
		"api_secret": cd.config.apiSecret,
		"signature":  generateSignature(signatureValues, cd.config.apiSecret),
	}

	// initial file field
	fileWriter, err := multipartWriter.CreateFormFile("file", f.Name())
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.CustomError("CannotProcessFile", ""))
	}

	// copy file to form file
	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.CustomError("CannotProcessFile", ""))
	}

	// other string field
	if err = writeMultipartField(multipartWriter, values); err != nil {
		return nil, err
	}

	if err = writeMultipartField(multipartWriter, signatureValues); err != nil {
		return nil, err
	}

	// close multipart writer
	if err = multipartWriter.Close(); err != nil {
		return nil, sdkcm.ErrCustom(err, sdkcm.CustomError("CannotProcessFile", ""))
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/video/upload", cd.config.cloudName),
		&requestBody,
	)
	if err != nil {
		return nil, sdkcm.ErrInvalidRequest(err)
	}

	// set header to multipart
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := http.Client{
		Timeout: time.Minute * 5,
	}

	res, err := client.Do(req.WithContext(ctx))

	if err != nil {
		return nil, sdkcm.ErrInvalidRequest(err)
	}

	defer res.Body.Close()

	out, _ := ioutil.ReadAll(res.Body)

	var result VideoResult

	if err := json.Unmarshal(out, &result); err != nil {
		return nil, sdkcm.ErrInvalidRequest(err)
	}

	if res.StatusCode != 200 {
		return nil, sdkcm.NewAppErr(errors.New(result.Error.Message), res.StatusCode, result.Error.Message).WithCode("")
	}

	return &result, nil
}

func writeMultipartField(multipartWriter *multipart.Writer, values map[string]string) error {
	for k, v := range values {
		fw, err := multipartWriter.CreateFormField(k)
		if err != nil {
			return sdkcm.ErrCustom(err, sdkcm.CustomError("CannotProcessFile", ""))
		}

		_, err = fw.Write([]byte(v))
		if err != nil {
			return sdkcm.ErrCustom(err, sdkcm.CustomError("CannotProcessFile", ""))
		}
	}
	return nil
}

func generateSignature(values map[string]string, apiSecret string) string {
	var lstKeyPairs []string
	for k, v := range values {
		lstKeyPairs = append(lstKeyPairs, fmt.Sprintf("%s=%v", k, v))
	}

	sort.Strings(lstKeyPairs)
	payloadToSign := strings.Join(lstKeyPairs, "&")

	h := sha1.New()
	h.Write([]byte(payloadToSign + apiSecret))
	return fmt.Sprintf("%x", h.Sum(nil))
}
