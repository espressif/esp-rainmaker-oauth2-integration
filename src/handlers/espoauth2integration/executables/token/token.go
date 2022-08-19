package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"constants"
	"handlers/utils"
	"models"
)

var CLIENT_ID = "client_id"
var CLIENT_SECRET = "client_secret"
var GRANT_TYPE = "grant_type"
var AUTHORIZATION_CODE = "authorization_code"
var CODE = "code"
var appVersion string
var appBuildDate string

const ERROR_ACCESS_TOKEN = "Error: Error occurred while fetching access token"
const ERROR_PARSING = "Error occurred while parsing"

/*
	Once authorize is successfully we get authorization code in response,
	This lambda is used to fetch access token by invoking IDP Oauth token api.
*/
func postToken(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TAG := "[postToken] "
	utils.LogInfo(TAG + req.HTTPMethod + " App Version = " + appVersion + " Posting Token ")

	client := http.Client{}

	//Fetch body from request
	receivedData, errParse := url.ParseQuery(req.Body)
	if errParse != nil {
		utils.LogError(TAG + "Error occurred while parsing " + errParse.Error())
		return utils.GetCliErrRespObject(http.StatusBadRequest, ERROR_PARSING), errParse
	}

	//client_id is IDP app client_id
	clientId := receivedData[CLIENT_ID]

	//client_secret is IDP app client_secret
	clientSecret := receivedData[CLIENT_SECRET]

	//code is authorization code received after successful authorization
	code := receivedData[CODE]

	utils.LogDebug(TAG + "Client ID:" + fmt.Sprint(clientId) + "client_secret" + fmt.Sprint(clientSecret) + "Code:" + fmt.Sprint(code))

	if len(clientId) == 0 || len(clientSecret) == 0 || len(code) == 0 {
		utils.LogError(TAG + "Unable to make request, few parameters are missing")
		return utils.GetServErrRespObject(constants.BAD_REQUEST), nil
	}

	//Form token request
	IDPAccessTokenURL := utils.GetIDPAccessTokenURL()

	//Form request body, previously these were query params, to support razer changed it to request body params
	data := url.Values{}
	data.Set(GRANT_TYPE, AUTHORIZATION_CODE)
	data.Set(CODE, code[0])
	data.Set(CLIENT_SECRET, clientSecret[0])
	data.Set(CLIENT_ID, clientId[0])

	utils.LogDebug(TAG + "Fetching access token uri is " + IDPAccessTokenURL)

	//Form request
	request, errNewRequest := http.NewRequest("POST", IDPAccessTokenURL, strings.NewReader(data.Encode()))
	if errNewRequest != nil {
		utils.LogError(TAG + "Error in creating get access token request" + errNewRequest.Error())
		return utils.GetCliErrRespObject(http.StatusBadRequest, constants.ERROR_INVALID_PAYLOAD), errNewRequest
	}

	//Set headers
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")

	//Make request
	response, errGetAccessToken := client.Do(request)
	if errGetAccessToken != nil {
		utils.LogError(TAG + "Error in fetching response" + errGetAccessToken.Error())
		return utils.GetServErrRespObject(ERROR_ACCESS_TOKEN), errGetAccessToken
	}

	//Fetch access token from response
	responseBuffer := new(bytes.Buffer)
	responseBuffer.ReadFrom(response.Body)
	responseAccessToken := responseBuffer.String()

	utils.LogDebug(fmt.Sprint(responseBuffer))
	utils.LogDebug(TAG + "Response " + fmt.Sprint(responseAccessToken))

	accessToken := new(models.AccessToken)
	errUnmarshall := json.Unmarshal([]byte(responseAccessToken), accessToken)

	if errUnmarshall != nil {
		utils.LogError(TAG + "Failed to unmarshal" + errUnmarshall.Error())
		return utils.GetCliErrRespObject(http.StatusBadRequest, constants.ERROR_INVALID_PAYLOAD), errUnmarshall
	}

	if accessToken.Access_token != constants.EMPTY_STRING {
		utils.LogInfo(TAG + "Successfully fetched access token")
	} else {
		utils.LogError(TAG + "Failed to fetch access token from identity provider")
	}

	js, _ := json.Marshal(accessToken)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       string(js),
	}, nil
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TAG := "[handleRequest] "
	utils.LogDebug(TAG + "method called req =  " + req.HTTPMethod + " App Version = " + appVersion + ", appBuildDate = " + appBuildDate)
	switch req.HTTPMethod {
	case "POST":
		return postToken(req)
	default:
		return utils.GetCliErrRespObject(http.StatusMethodNotAllowed, constants.ERROR_METHOD_NOT_ALLOWED), nil
	}
}

func main() {
	lambda.Start(handleRequest)
}
