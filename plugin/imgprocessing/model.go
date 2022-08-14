package imgprocessing

import "github.com/200Lab-Education/go-sdk/sdkcm"

type Response struct {
	sdkcm.AppError
	Data *sdkcm.Image `json:"data"`
}
