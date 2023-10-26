package utilsResponse

import (
	"net/http"

	"github.com/gin-gonic/gin"

	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
)

var ApiResponseConfigData = map[int]modelUtils.ApiResponseConfigStruct{
	http.StatusOK: {
		ErrorCode: "0",
		Error:     nil,
	},
	http.StatusBadRequest: {
		ErrorCode: "DF400",
		Error:     "Bad Request",
	},
	http.StatusUnauthorized: {
		ErrorCode: "DF401",
		Error:     "Unauthorized",
	},
	http.StatusForbidden: {
		ErrorCode: "DF403",
		Error:     "Forbidden",
	},
	http.StatusNotFound: {
		ErrorCode: "DF404",
		Error:     "Not Found",
	},
	http.StatusConflict: {
		ErrorCode: "DF409",
		Error:     "Conflict",
	},
	http.StatusInternalServerError: {
		ErrorCode: "DF500",
		Error:     "Internal Server Error",
	},
}

func ApiResponse(c *gin.Context, param modelUtils.ApiResponseStruct) {
	if param.ResponseCode == 0 {
		param.ResponseCode = 500
	}

	var success_status bool
	if param.ResponseCode >= 200 && param.ResponseCode <= 299 {
		success_status = true
	} else {
		success_status = false
	}

	if param.ErrorCode == "" {
		param.ErrorCode = ApiResponseConfigData[param.ResponseCode].ErrorCode
	}

	if param.Error == nil {
		param.Error = ApiResponseConfigData[param.ResponseCode].Error
	}

	JsonResponse(c, modelUtils.JsonResponseStruct{
		ResponseCode: param.ResponseCode,
		Detail: modelUtils.JsonResponseStructDetail{
			Success:   success_status,
			ErrorCode: param.ErrorCode,
			Data:      param.Data,
			Error:     param.Error,
		},
	})
}
