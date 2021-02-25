package utils

import (
	"constants"
	"errors"
)

func IsSuperAdmin(accessToken string) (bool, error) {
	TAG := "[IsSuperAdmin] "
	LogDebug(TAG + "Checking if user is superadmin")

	//TODO, if above todo is complete and mapping addition is done
	//This call be removed when all the user's mapping is migrated to table cognito_user_mapping
	accessTokenDetails, errAccessTokenDetails := GetCognitoUser(accessToken)
	if errAccessTokenDetails != nil {
		LogError(TAG + "Getting details from access token failed " + errAccessTokenDetails.Error())
		return false, errors.New(constants.ERROR_FETCH_USER_DETAILS_FROM_ACCESS_TOKEN)
	}

	// check if super admin
	if !accessTokenDetails.IsAdmin || !accessTokenDetails.IsMaintainer {
		LogDebug(TAG + "User is not superadmin")
		return false, nil
	}
	LogDebug(TAG + "User is superadmin")
	return true, nil
}
