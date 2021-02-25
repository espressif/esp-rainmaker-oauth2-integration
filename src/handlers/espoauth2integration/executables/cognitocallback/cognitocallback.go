package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"constants"
	"handlers/espoauth2integration"

	"net/http"

	"handlers/utils"
)

var CODE = "code"
var STATE = "state"
var OAUTH_IDPRESPONSE_URL = "/oauth2/idpresponse"
var HTTPS = "https://"
var appVersion string
var appBuildDate string

const ERROR_FETCH_STATE = "Error in fetching cognito state"
const ERROR_FETCH_COGNITO_DOMAIN_NAME = "Error occurred while fetching cognito domain"
const ERROR_FETCH_TENANT_CONTEXT = "Error occurred while fetching tenant context"

func FetchCognitoState(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TAG := "[FetchCognitoState]"
	utils.LogInfo(TAG + request.HTTPMethod + " App Version = " + appVersion + " Fetching cognito state ")

	code := request.QueryStringParameters[CODE]
	key := request.QueryStringParameters[STATE]

	if code == constants.EMPTY_STRING || key == constants.EMPTY_STRING {
		utils.LogError(TAG + "Unable to process further, few parameters are missing")
		return utils.GetServErrRespObject(constants.BAD_REQUEST), nil
	}

	//fetch cognito state
	cognitoState, errGetState := espoauth2integration.GetState(key)
	if errGetState != nil {
		utils.LogError(TAG + ERROR_FETCH_STATE + " error " + errGetState.Error())
		return utils.GetServErrRespObject(ERROR_FETCH_STATE), nil
	}

	//validation
	if cognitoState.State != constants.EMPTY_STRING {
		utils.LogDebug(TAG + "State found and retrieved for key," + key)
	} else {
		utils.LogError(TAG+"No state found for key ", key)
	}

	//get tenant context
	tenctx, errGetTenantContext := utils.TenantGetContext()
	if errGetTenantContext != nil {
		utils.LogError(TAG + "Error TenantContext" + errGetTenantContext.Error())
		return utils.GetServErrRespObject(ERROR_FETCH_TENANT_CONTEXT), nil
	}

	//configuring cognito's idpresponse url to be redirected
	domainName, errGetDomain := utils.GetDomain(tenctx.UserPoolId)
	if errGetDomain != nil {
		utils.LogError(TAG + "Error occured while fetching cognito domain name " + errGetDomain.Error())
		return utils.GetServErrRespObject(ERROR_FETCH_COGNITO_DOMAIN_NAME), nil
	}

	cognitoIdpResponseUrl := HTTPS + domainName + OAUTH_IDPRESPONSE_URL +
		constants.QUESTION + CODE + constants.EQUAL + code +
		constants.AMPERSAND + STATE + constants.EQUAL + cognitoState.State

	utils.LogInfo(TAG + "Successfully redirecting to url " + cognitoIdpResponseUrl)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMovedPermanently,
		Headers: map[string]string{"Location": cognitoIdpResponseUrl, "Access-Control-Allow-Origin": "*",
			"Cache-Control":             "no-cache, no-store, max-age=0, must-revalidate",
			"Pragma":                    "no-cache",
			"Expires":                   "Thu, 01 Jan 1970 00:00:00 UTC",
			"X-Content-Type-Options":    "nosniff",
			"X-Xss-Protection":          "1; mode=block",
			"Strict-Transport-Security": "max-age=31536000 ; includeSubDomains",
			"X-Frame-Options":           "DENY"},
		Body: constants.EMPTY_STRING,
	}, nil
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TAG := "[handleRequest] "
	utils.LogDebug(TAG + " method called req =  " + req.HTTPMethod + " App Version = " + appVersion + ", appBuildDate = " + appBuildDate)
	switch req.HTTPMethod {
	case "GET":
		return FetchCognitoState(req)
	default:
		return utils.GetCliErrRespObject(http.StatusMethodNotAllowed, constants.ERROR_METHOD_NOT_ALLOWED), nil
	}
}

func main() {
	lambda.Start(handleRequest)
}
