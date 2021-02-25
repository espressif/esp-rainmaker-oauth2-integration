package utils

import (
	"os"
)

var ACCOUNT_ID = "ACCOUNT_ID"
var REGION = "REGION"
var APIGATEWAYURL = "API_GATEWAY_URL"
var STAGE_NAME = "STAGE_NAME"
var TIME_TO_LIVE = "TIME_TO_LIVE"
var LOG_LEVEL = "LOG_LEVEL"
var BASE_URL = "BASE_URL"
var Rainmaker_Oauth2_Authorize_URL = "Rainmaker_Oauth2_Authorize_URL"
var Rainmaker_Oauth2_Token_URL = "Rainmaker_Oauth2_Token_URL"
var Rainmaker_Oauth2_UserInfo_URL = "Rainmaker_Oauth2_UserInfo_URL"
var Rainmaker_Oauth2_Email_URL = "Rainmaker_Oauth2_Email_URL"
var ACCEPT_APPLICATION_JSON = "ACCEPT_APPLICATION_JSON"

func GetUserPoolId() string {
	var userPoolId string
	userPoolId = os.Getenv("COGNITO_USER_POOL_ID")
	return userPoolId
}

func GetAppClientId() string {

	var appClientId string
	appClientId = os.Getenv("COGNITO_APP_CLIENT_ID")
	return appClientId
}

func GetAccountId() string {
	return os.Getenv(ACCOUNT_ID)
}

func GetRegion() string {
	return os.Getenv(REGION)
}

func GetApiGatewayHostURL() string {
	return os.Getenv(APIGATEWAYURL)
}

func GetStageName() string {
	return os.Getenv(STAGE_NAME)
}

func GetTimeToLive() string {
	return os.Getenv(TIME_TO_LIVE)
}

func GetLogLevel() string {
	return os.Getenv(LOG_LEVEL)
}

func GetBaseApi() string {
	return os.Getenv(BASE_URL)
}

func GetIDPAuthorizeURL() string {
	return os.Getenv(Rainmaker_Oauth2_Authorize_URL)
}

func GetIDPAccessTokenURL() string {
	return os.Getenv(Rainmaker_Oauth2_Token_URL)
}

func GetIDPUserURL() string {
	return os.Getenv(Rainmaker_Oauth2_UserInfo_URL)
}

func GetIDPUserEmailsURL() string {
	return os.Getenv(Rainmaker_Oauth2_Email_URL)
}

func GetAcceptApplicationJson() string {
	return os.Getenv(ACCEPT_APPLICATION_JSON)
}
