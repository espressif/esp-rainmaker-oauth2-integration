package main

import (
	"encoding/json"
	"models"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"constants"
	"handlers/utils"
)

var appVersion string
var appBuildDate string

var availableFunctions = []string{"esp-Oauth2Authorize", "esp-Oauth2Token", "esp-Oauth2UserInfo", "esp-Oauth2CognitoCallback", "esp-Oauth2UpdateEnv", "esp-Oauth2GetEnv"}

func updateEnv(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	TAG := "[updateEnv] "
	utils.LogInfo(TAG + req.HTTPMethod + " App Version = " + appVersion + " Updating environment for lambda")

	rainmakerOauth2Urls := new(models.RainmakerOauth2Urls)

	//fetch request body parameters
	errUnMarshal := json.Unmarshal([]byte(req.Body), rainmakerOauth2Urls)
	if errUnMarshal != nil {
		utils.LogError(TAG + "Error occurred while unmarshalling " + " error\t" + errUnMarshal.Error())
		return utils.GetCliErrRespObject(http.StatusBadRequest, constants.ERROR_INVALID_PAYLOAD)
	}

	// Fetch and update environment (oauth2 urls) for available lambda functions
	for _, functionName := range availableFunctions {

		// Fetch funtion environment
		functionEnv, getEnvErr := utils.GetLambdaEnv(functionName)
		if getEnvErr != nil {
			utils.LogError(TAG + "Error occurred while fetching function environment, function name: " + functionName + " error\t" + getEnvErr.Error())
			return utils.GetServErrRespObject(constants.UPDATE_OAUTH2_URLS_ERROR)
		}

		// Update oauth2 url values as per request body
		if rainmakerOauth2Urls.RainmakerOauth2AuthorizeUrl != constants.EMPTY_STRING {
			functionEnv.Variables[utils.Rainmaker_Oauth2_Authorize_URL] = &rainmakerOauth2Urls.RainmakerOauth2AuthorizeUrl
		}
		if rainmakerOauth2Urls.RainmakerOauth2TokenUrl != constants.EMPTY_STRING {
			functionEnv.Variables[utils.Rainmaker_Oauth2_Token_URL] = &rainmakerOauth2Urls.RainmakerOauth2TokenUrl
		}
		if rainmakerOauth2Urls.RainmakerOauth2UserinfoUrl != constants.EMPTY_STRING {
			functionEnv.Variables[utils.Rainmaker_Oauth2_UserInfo_URL] = &rainmakerOauth2Urls.RainmakerOauth2UserinfoUrl
		}
		functionEnv.Variables[utils.Rainmaker_Oauth2_Email_URL] = &rainmakerOauth2Urls.RainmakerOauth2EmailUrl

		// Update the oauth2 urls in the environment
		updateEnvErr := utils.UpdateLambdaEnv(functionName, functionEnv)
		if updateEnvErr != nil {
			utils.LogError(TAG + "Error occurred while updating function environment, function name: " + functionName + " error\t" + getEnvErr.Error())
			return utils.GetServErrRespObject(constants.UPDATE_OAUTH2_URLS_ERROR)
		}
	}

	apiResponse := models.APIResponse{Status: constants.SUCCESS, Description: constants.UPDATE_OAUTH2_URLS_SUCCESSFUL}

	utils.LogInfo(TAG + "Updated environment for availble lambdas successfully")

	js, _ := json.Marshal(apiResponse)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       string(js),
	}
}

func getEnv(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	TAG := "[getEnv] "
	utils.LogInfo(TAG + req.HTTPMethod + " App Version = " + appVersion + " Getting oauth2 urls from environment")

	rainmakerOauth2Urls := new(models.RainmakerOauth2Urls)

	// Fetch oauth2 urls from environment of any functions from availableFunctions
	functionEnv, getEnvErr := utils.GetLambdaEnv(availableFunctions[0])
	if getEnvErr != nil {
		utils.LogError(TAG + "Error occurred while fetching function environment, function name: " + availableFunctions[0] + " error\t" + getEnvErr.Error())
		return utils.GetServErrRespObject(constants.UPDATE_OAUTH2_URLS_ERROR)
	}

	// get the oauth2 urls
	urls := functionEnv.Variables
	if urls[utils.Rainmaker_Oauth2_Authorize_URL] != nil {
		rainmakerOauth2Urls.RainmakerOauth2AuthorizeUrl = *urls[utils.Rainmaker_Oauth2_Authorize_URL]
	}
	if urls[utils.Rainmaker_Oauth2_Token_URL] != nil {
		rainmakerOauth2Urls.RainmakerOauth2TokenUrl = *urls[utils.Rainmaker_Oauth2_Token_URL]
	}
	if urls[utils.Rainmaker_Oauth2_UserInfo_URL] != nil {
		rainmakerOauth2Urls.RainmakerOauth2UserinfoUrl = *urls[utils.Rainmaker_Oauth2_UserInfo_URL]
	}
	if urls[utils.Rainmaker_Oauth2_Email_URL] != nil {
		rainmakerOauth2Urls.RainmakerOauth2EmailUrl = *urls[utils.Rainmaker_Oauth2_Email_URL]
	}

	utils.LogInfo(TAG + "Fetched oauth2 urls form lambda environment successfully")

	js, _ := json.Marshal(rainmakerOauth2Urls)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       string(js),
	}
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TAG := "[handleRequest] "
	utils.LogDebug(TAG + "method called req =  " + req.HTTPMethod + " App Version = " + appVersion + ", appBuildDate = " + appBuildDate)

	accessToken := req.Headers["Authorization"]
	superAdmin, errCheckSuperAdmin := utils.IsSuperAdmin(accessToken)
	if errCheckSuperAdmin != nil {
		utils.LogError(TAG + "Error occured while cheking if user is super admin, error: " + errCheckSuperAdmin.Error())
		return utils.GetServErrRespObject(constants.NOT_AUTHORIZED), nil
	} else if !superAdmin {
		utils.LogError(TAG + "User is not super admin")
		return utils.GetCliErrRespObject(http.StatusBadRequest, constants.NOT_AUTHORIZED), nil
	}

	switch req.HTTPMethod {
	case "POST":
		return updateEnv(req), nil
	case "GET":
		return getEnv(req), nil
	default:
		return utils.GetCliErrRespObject(http.StatusMethodNotAllowed, constants.ERROR_METHOD_NOT_ALLOWED), nil
	}
}

func main() {
	lambda.Start(handleRequest)
}
