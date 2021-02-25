package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"net/url"

	"github.com/lithammer/shortuuid"

	"constants"
	"handlers/utils"
	"models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var HTTP_CLIENT = http.Client{}
var ACCEPT = "Accept"
var AUTHORIZATION = "Authorization"
var appVersion string
var appBuildDate string

const ERROR_ACCESS_TOKEN_MISSING = "Error: Access token is missing"
const ERROR_GET_USER = "Error: Error occurred while getting user info"
const ERROR_PRIMARY_USER_NOT_FOUND = "Error: Primary user email does not exists"

/*
	Once token request is successful we get access_token in response,
	This lambda is used to fetch userdetails invoking IDP Oauth userinfo api.
*/
func getUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TAG := "[getUser] "
	accessToken := req.Headers[AUTHORIZATION]
	utils.LogInfo(TAG + req.HTTPMethod + " App Version = " + appVersion + " Fetching user ")
	utils.LogDebug(TAG + "Token = " + accessToken)
	var userName string
	var userEmailDetail *models.IDPUserEmailDetails
	var errGetIDPUserEmailDetail error

	if accessToken == constants.EMPTY_STRING {
		utils.LogError(TAG + ERROR_ACCESS_TOKEN_MISSING)
		return utils.GetServErrRespObject(ERROR_ACCESS_TOKEN_MISSING), nil
	}

	//Fetch user details
	user, errGetIDPUser := getIDPUser(accessToken)
	if errGetIDPUser != nil {
		utils.LogError(TAG + "Failed to get identity provider user" + errGetIDPUser.Error())
		return utils.GetServErrRespObject(ERROR_GET_USER), errGetIDPUser
	}

	emailUrl := utils.GetIDPUserEmailsURL()
	// Check if emailUrl is in valid url format
	_, validEmailUrlErr := url.ParseRequestURI(emailUrl)

	if validEmailUrlErr == nil && emailUrl != constants.EMPTY_STRING {
		userEmailDetail, errGetIDPUserEmailDetail = getIDPUserEmailDetails(accessToken)
		if errGetIDPUserEmailDetail != nil {
			utils.LogError(TAG + "Failed to get identity provider user emails" + errGetIDPUserEmailDetail.Error())
			return utils.GetServErrRespObject(ERROR_GET_USER), errGetIDPUserEmailDetail
		}

		if userEmailDetail.Email == constants.EMPTY_STRING || userEmailDetail.Email_Verified == false {
			utils.LogError(TAG + "Primary email not found for user")
			return utils.GetServErrRespObject(ERROR_PRIMARY_USER_NOT_FOUND), nil
		}
		user.Email = userEmailDetail.Email
		userName = user.Name
	} else {
		userName = user.Email[0:strings.Index(user.Email, "@")]
	}

	//Form openIdUser- user to be stored in cognito
	openIdUser := new(models.OpenIDUser)
	if user.Sub != nil {
		openIdUser.Sub = user.Sub
	} else if user.Id != nil {
		openIdUser.Sub = user.Id
	}

	openIdUser.Name = userName
	openIdUser.Preferred_UserName = userName
	openIdUser.Updated_At = user.Updated_At
	openIdUser.Email = user.Email
	openIdUser.Email_Verified = true
	openIdUser.Picture = user.Picture
	openIdUser.IDPName = userName

	//Additional handling in order to avoid overriding user_id
	tenctx, errGetTenantContextByReq := utils.TenantGetContextByReq(req)
	if errGetTenantContextByReq != nil {
		utils.LogError(TAG + errGetTenantContextByReq.Error())
		return utils.GetCliErrRespObject(http.StatusBadRequest, errGetTenantContextByReq.Error()), errGetTenantContextByReq
	}

	userDetails, _ := utils.UserDBGetByName(tenctx, user.Email)

	utils.LogDebug(TAG + "Fetched user, UserId: " + userDetails.UserId)

	if userDetails.UserId == constants.EMPTY_STRING {
		openIdUser.UserId = shortuuid.New()
	} else {
		openIdUser.UserId = userDetails.UserId
	}

	js, _ := json.Marshal(openIdUser)

	utils.LogInfo(TAG + "Successfully have fetched user, OpenId UserId: " + openIdUser.UserId)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       string(js),
	}, nil
}

func makeRequest(requestUrl string, accessToken string) (*http.Response, error) {
	TAG := "[makeRequest] "

	ACCEPT_APPLICATION_JSON := utils.GetAcceptApplicationJson()

	request, errNewRequest := http.NewRequest("GET", requestUrl, nil)
	if errNewRequest != nil {
		utils.LogError(TAG + "Error in creating get access token request" + errNewRequest.Error())
		return nil, errNewRequest
	}
	request.Header.Set(ACCEPT, ACCEPT_APPLICATION_JSON)
	request.Header.Set(AUTHORIZATION, accessToken)
	return HTTP_CLIENT.Do(request)
}

func getIDPUser(accessToken string) (*models.IDPUser, error) {
	TAG := "[getIDPUser] "
	utils.LogDebug(TAG + "Getting identity provider user")

	IDPUserURL := utils.GetIDPUserURL()
	response, errGetResponse := makeRequest(IDPUserURL, accessToken)
	if errGetResponse != nil {
		utils.LogError(TAG + "Error in getting response" + errGetResponse.Error())
		return nil, errGetResponse
	}

	responseBuffer := new(bytes.Buffer)
	responseBuffer.ReadFrom(response.Body)
	responseUser := responseBuffer.String()

	user := new(models.IDPUser)
	errUnmarshall := json.Unmarshal([]byte(responseUser), user)

	if errUnmarshall != nil {
		utils.LogError(TAG + "Failed to unmarshal" + errUnmarshall.Error())
		return nil, errUnmarshall
	}

	if user.Updated_At != constants.EMPTY_STRING {
		updatedTimeInEpoch, errTimeParse := time.Parse(time.RFC3339, user.Updated_At)
		if errTimeParse != nil {
			utils.LogError(TAG + "Failed to parse time" + errTimeParse.Error())
			return nil, errTimeParse
		}
		user.Updated_At = strconv.FormatInt(updatedTimeInEpoch.Unix(), 10)
	}

	utils.LogDebug(TAG + "Successfully fetched IDP user, UserName : " + user.Name)
	return user, nil
}

func getIDPUserEmailDetails(accessToken string) (*models.IDPUserEmailDetails, error) {
	TAG := "[getIDPUserEmailDetails] "
	utils.LogDebug(TAG + "Getting user email details")

	IDPUserEmailsURL := utils.GetIDPUserEmailsURL()
	response, errGetResponse := makeRequest(IDPUserEmailsURL, accessToken)

	if errGetResponse != nil {
		utils.LogError(TAG + "Error in getting response" + errGetResponse.Error())
		return nil, errGetResponse
	}

	responseBuffer := new(bytes.Buffer)
	responseBuffer.ReadFrom(response.Body)
	responseUserEmailDetails := responseBuffer.String()

	userEmailDetails := []models.IDPUserEmailDetails{}
	errUnmarshall := json.Unmarshal([]byte(responseUserEmailDetails), &userEmailDetails)
	if errUnmarshall != nil {
		utils.LogError(TAG + "Failed to unmarshal" + errUnmarshall.Error())
		return nil, errUnmarshall
	}

	utils.LogDebug(TAG + "User Email Details count: " + strconv.Itoa(len(userEmailDetails)))

	userEmailDetail := new(models.IDPUserEmailDetails)

	for _, userEmails := range userEmailDetails {
		if userEmails.Primary == true {
			userEmailDetail.Email = userEmails.Email
			userEmailDetail.Email_Verified = userEmails.Email_Verified
		}
	}

	utils.LogInfo(TAG + "Successfully fetched primary mail info of IDP user :" + userEmailDetail.Email)
	return userEmailDetail, nil
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TAG := "[handleRequest] "
	utils.LogDebug(TAG + "method called req =  " + req.HTTPMethod + " App Version = " + appVersion + ", appBuildDate = " + appBuildDate)
	switch req.HTTPMethod {
	case "GET":
		return getUser(req)
	default:
		return utils.GetCliErrRespObject(http.StatusMethodNotAllowed, constants.ERROR_METHOD_NOT_ALLOWED), nil
	}
}

func main() {
	lambda.Start(handleRequest)
}
