package utils

import (
	"bytes"
	"constants"
	"encoding/json"
	"models"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GetRespObject(status int, v interface{}) events.APIGatewayProxyResponse {
	body, err := json.Marshal(v)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "JSON Marshaling Failure",
		}
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)
	return events.APIGatewayProxyResponse{
		StatusCode:      status,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers:         map[string]string{"Access-Control-Allow-Origin": "*"},
	}
}

func GetServErrRespObject(errorText string) events.APIGatewayProxyResponse {
	return GetRespObject(http.StatusInternalServerError,
		models.APIErrResponse{Status: constants.FAILURE,
			Description: errorText})
}

func GetCliErrRespObject(status int, errorText string) events.APIGatewayProxyResponse {
	return GetRespObject(status,
		models.APIErrResponse{Status: constants.FAILURE,
			Description: errorText})
}
