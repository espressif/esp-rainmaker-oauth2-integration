package main

import (
	"constants"
	"github.com/lithammer/shortuuid"
	"handlers/utils"
	"models"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"handlers/espoauth2integration"
)

var CLIENT_ID = "client_id"
var SCOPE = "scope"
var RESPONSE_TYPE = "response_type"
var STATE = "state"
var COLON = constants.COLON
var AMPERSAND = constants.AMPERSAND
var EQUAL = constants.EQUAL
var appVersion string
var appBuildDate string

const ERROR_STORE_STATE = "Error in storing cognito state"

/*
	This lambda is used to form the authorize request and redirect to third party login to authorize the calling user.
 */
func authorize(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TAG := "[authorize] "
	utils.LogInfo(TAG + request.HTTPMethod + " App Version = " + appVersion + " Authorizing User")

	//client_id is IDP app client_id
	var client_id = request.QueryStringParameters[CLIENT_ID]

	//scope is authorize request scope
	var scope = request.QueryStringParameters[SCOPE]

	//response_type is auauthorize request response_type
	var response_type = request.QueryStringParameters[RESPONSE_TYPE]

	var state = request.QueryStringParameters[STATE]

	utils.LogDebug(TAG + "Request parameters: " + CLIENT_ID + COLON + client_id +
		SCOPE + COLON + scope +
		RESPONSE_TYPE + COLON + response_type +
		STATE + COLON + state)

	if client_id == constants.EMPTY_STRING || scope == constants.EMPTY_STRING || response_type == constants.EMPTY_STRING || state == constants.EMPTY_STRING {
		utils.LogError(TAG + "Unable to make request, few parameters are missing")
		return utils.GetServErrRespObject(constants.BAD_REQUEST), nil
	}

	//This change is made as some IDPs like Github don't handle longer state variable.
	//generate random uuid as state to be passed to IDP
	uuidState := shortuuid.New()

	//Store state and key[generated uuid] to cognito_state tbl
	errStoreState := espoauth2integration.StoreState(models.State{
		State: state,
		Key:   uuidState,
	})
	if errStoreState != nil {
		utils.LogError(TAG + ERROR_STORE_STATE + " error " + errStoreState.Error())
		return utils.GetServErrRespObject(ERROR_STORE_STATE), nil
	}
	utils.LogDebug(TAG + "Successfully stored state ")

	//Form authorize request
	IDPAuthorizeURL := utils.GetIDPAuthorizeURL()
	getAuthorizeCodeUrl := IDPAuthorizeURL + constants.QUESTION + CLIENT_ID + EQUAL + client_id +
		AMPERSAND + STATE + EQUAL + uuidState +
		AMPERSAND + SCOPE + EQUAL + scope +
		AMPERSAND + RESPONSE_TYPE + EQUAL + response_type

	utils.LogDebug(TAG + "Get Authorized code URL: " + getAuthorizeCodeUrl)
	utils.LogInfo(TAG + "Successfully redirecting to IDP for authorization ")

	//Redirecting to IDP authorize url
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusFound,
		Headers:    map[string]string{"Location": getAuthorizeCodeUrl, "Access-Control-Allow-Origin": "*"},
		Body:       constants.EMPTY_STRING,
	}, nil
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TAG := "[handleRequest] "
	utils.LogDebug(TAG + " method called req =  " + req.HTTPMethod + " App Version = " + appVersion + ", appBuildDate = " + appBuildDate)
	switch req.HTTPMethod {
	case "GET":
		return authorize(req)
	default:
		return utils.GetCliErrRespObject(http.StatusMethodNotAllowed, constants.ERROR_METHOD_NOT_ALLOWED), nil
	}
}

func main() {
	lambda.Start(handleRequest)
}
