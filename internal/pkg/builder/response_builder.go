package builder

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func BuildSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Status: "success",
		Data:   data,
	}
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
