package sdkcm

import (
	"github.com/btcsuite/btcutil/base58"
	"net/http"
)

// Response helpers
var (
	SimpleSuccessResponse = func(data interface{}) Response {
		return newResponse(http.StatusOK, data, nil, nil)
	}

	ResponseWithPaging = func(data, param interface{}, other interface{}) Response {
		if v, ok := other.(Paging); ok {
			if v.NextCursor != "" {
				if !v.CursorIsUID {
					v.NextCursor = base58.Encode([]byte(v.NextCursor))
				}
			}
			return newResponse(http.StatusOK, data, param, v)
		}
		return newResponse(http.StatusOK, data, param, other)
	}
)

type Response struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
	Param  interface{} `json:"param,omitempty"`
	Paging interface{} `json:"paging,omitempty"`
}

func newResponse(code int, data, param, other interface{}) Response {
	return Response{
		Code:   code,
		Data:   data,
		Param:  param,
		Paging: other,
	}
}
