package helpers

type Response map[string]any

func Success(data any, message string) Response {
	msg := "Your request has been successfully processed"
	if message != "" {
		msg = message
	}
	return Response{
		"success": true,
		"message": msg,
		"data":    data,
	}
}

func Error(errorMsg any) Response {
	if _, ok := errorMsg.(string); ok {
		return Response{
			"errors":  make(map[string][]string),
			"success": false,
			"message": errorMsg,
		}
	}

	if _, ok := errorMsg.(map[string][]string); ok {
		return Response{
			"success": false,
			"errors":  errorMsg,
			"message": "Your request has been failed to process",
		}
	}

	if errHttp, ok := errorMsg.(*CustomError); ok {
		return Response{
			"errors":  errHttp.Errors,
			"success": false,
			"message": errHttp.Msg,
		}
	}

	if err, ok := errorMsg.(error); ok {
		return Response{
			"errors":  make(map[string][]string),
			"success": false,
			"message": err.Error(),
		}
	}

	return Response{
		"success": false,
		"message": "Your request has been failed to process",
	}
}
