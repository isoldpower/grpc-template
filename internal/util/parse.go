package util

import (
	"encoding/json"
	"golang-grpc/internal/log"
	"net/http"
)

func ParseBody(request *http.Request, result any) error {
	decoder := json.NewDecoder(request.Body)
	decodeError := decoder.Decode(result)
	if decodeError != nil {
		return decodeError
	}

	return nil
}

func WriteError(writer http.ResponseWriter, statusCode int, err error) {
	writer.WriteHeader(statusCode)
	writer.Header().Set("Content-Type", "application/json")
	_, writingError := writer.Write(json.RawMessage(`{"message":"` + err.Error() + `"}`))
	if writingError != nil {
		log.PrintError("Failed to marshal the error message json", writingError)
	}
}

func WriteResponse(writer http.ResponseWriter, statusCode int, result interface{}) error {
	writer.WriteHeader(statusCode)
	writer.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(result)

	if err != nil {
		log.Errorln("Failed to marshal the response json")
		return err
	}
	_, writingError := writer.Write(response)
	if writingError != nil {
		log.PrintError("Failed to marshal the response json", writingError)
	}

	return nil
}
