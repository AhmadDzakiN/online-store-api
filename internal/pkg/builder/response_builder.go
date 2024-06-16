package builder

import "reflect"

type SuccessResponse struct {
	Status string      `json:"status"`
	Next   *string     `json:"next,omitempty"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func BuildSuccessResponse(data interface{}, pageToken *string) SuccessResponse {
	resp := SuccessResponse{}

	if reflect.ValueOf(data).Kind() == reflect.Slice {
		resp.Data = data
	} else {
		temp := []interface{}{}
		if reflect.ValueOf(data).IsValid() {
			temp = []interface{}{
				data,
			}
		}
		resp.Data = temp
	}

	resp.Status = "success"
	resp.Next = pageToken

	return resp
}

func BuildErrorResponse(err string) ErrorResponse {
	if err == "" {
		err = "Oops, something wrong with the server. Please try again later"
	}

	return ErrorResponse{
		Status: "error",
		Error:  err,
	}
}
