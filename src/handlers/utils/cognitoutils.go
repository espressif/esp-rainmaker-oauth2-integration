package utils

import (
	"errors"
	"models"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"constants"
)

var cognitoClient = cognitoidentityprovider.New(session.New(), aws.NewConfig().WithRegion(GetRegion()))

var AUTH = ".auth."
var AMAZON_COGNITO = ".amazoncognito.com"

func GetDomain(userPoolId string) (string, error) {
	TAG := "[GetDomain] "
	LogDebug(TAG + "Fetching domain Name for userpoolId: " + userPoolId)

	var customDomainName, cognitoDomainName string

	//fetch userpool
	userPoolDetails, errGetUserPool := GetUserPool(userPoolId)
	if errGetUserPool != nil {
		LogError(TAG + "Error occurred while fetching cognito user pool " + errGetUserPool.Error())
		return constants.EMPTY_STRING, errGetUserPool
	}

	//validation
	if userPoolDetails.UserPool.CustomDomain != nil {
		customDomainName = *userPoolDetails.UserPool.CustomDomain
	}
	if userPoolDetails.UserPool.Domain != nil {
		cognitoDomainName = *userPoolDetails.UserPool.Domain
	}

	LogDebug(TAG + "Domain: " + cognitoDomainName + " Custom Domain: " + customDomainName)

	if customDomainName != constants.EMPTY_STRING {
		LogDebug(TAG + "Custom domain is configured " + customDomainName)
		return customDomainName, nil
	} else if cognitoDomainName != constants.EMPTY_STRING {
		LogDebug(TAG + "Cognito domain is configured " + cognitoDomainName)
		return cognitoDomainName + AUTH + GetRegion() + AMAZON_COGNITO, nil
	} else {
		LogError(TAG + "No Domain is set")
		return constants.EMPTY_STRING, errors.New(constants.INFO_NO_DOMAIN_SET)
	}
}

func GetUserPool(userPoolId string) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	TAG := "[GetUserPoolNameFromId] "
	LogDebug(TAG + "Fetching userpool, Given UserpoolId: " + userPoolId)

	describeUserPoolInput := &cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: aws.String(userPoolId),
	}

	userPoolDetails, err := cognitoClient.DescribeUserPool(describeUserPoolInput)
	if err != nil {
		LogError(TAG + "Error occurred while fetching userpool name, UserPoolId: " + userPoolId + " with error, " + err.Error())
		return nil, err
	}

	LogDebug(TAG + "Successfully fetched userpool ")
	return userPoolDetails, nil
}

func GetCognitoUser(cognitoAccessToken string) (models.AccessTokenHeader, error) {
	TAG := "[GetCognitoUser]"
	LogDebug(TAG + "Fetching cognito user ")

	var accessTokenHeader models.AccessTokenHeader
	var getUserInput cognitoidentityprovider.GetUserInput
	var errParsing error
	getUserInput.AccessToken = &cognitoAccessToken

	cognitoIdPoolClient := cognitoidentityprovider.New(session.New(), aws.NewConfig())

	user, errGetUser := cognitoIdPoolClient.GetUser(&getUserInput)
	if errGetUser != nil {
		LogError(TAG + "Error in getting user " + errGetUser.Error())
		return accessTokenHeader, errGetUser
	}

	for i := 0; i < len(user.UserAttributes); i++ {
		if *user.UserAttributes[i].Name == "custom:user_id" {
			accessTokenHeader.UserId = *user.UserAttributes[i].Value
		}
		if *user.UserAttributes[i].Name == "email" {
			accessTokenHeader.Email = *user.UserAttributes[i].Value
		}
		if *user.UserAttributes[i].Name == constants.CUSTOM_ADMIN {
			accessTokenHeader.IsAdmin, errParsing = strconv.ParseBool(*user.UserAttributes[i].Value)
			if errParsing != nil {
				LogError(TAG + "Error in parsing" + errParsing.Error())
				return accessTokenHeader, errParsing
			}
		}
		if *user.UserAttributes[i].Name == constants.CUSTOM_MAINTAINER {
			accessTokenHeader.IsMaintainer, errParsing = strconv.ParseBool(*user.UserAttributes[i].Value)
			if errParsing != nil {
				LogError(TAG + "Error in parsing" + errParsing.Error())
				return accessTokenHeader, errParsing
			}
		}
	}

	LogDebug(TAG + "Successully fetched Cognito user, UserId: " + accessTokenHeader.UserId)
	return accessTokenHeader, nil
}
