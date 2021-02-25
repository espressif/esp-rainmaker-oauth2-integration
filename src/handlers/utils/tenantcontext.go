package utils

import (
	"constants"
	"models"

	"github.com/aws/aws-lambda-go/events"
)

func TenantGetContextByReq(req events.APIGatewayProxyRequest) (models.TenantContext, error) {
	TAG := "[TenantContext]"
	LogDebug(TAG + "Getting tenant context ")
	custId := constants.EMPTY_STRING
	userPoolId := GetUserPoolId()
	cognitoAppClientId := GetAppClientId()
	return models.TenantContext{CustomerId: custId, UserPoolId: userPoolId, CognitoAppClientId: cognitoAppClientId}, nil
}

func TenantGetContext() (models.TenantContext, error) {
	TAG := "[TenantContext]"
	LogDebug(TAG + "Getting tenant context ")
	custId := constants.EMPTY_STRING
	userPoolId := GetUserPoolId()
	cognitoAppClientId := GetAppClientId()
	return models.TenantContext{CustomerId: custId, UserPoolId: userPoolId, CognitoAppClientId: cognitoAppClientId}, nil
}

func TenantGetDBName(tenctx models.TenantContext, table string) string {
	return table
}
