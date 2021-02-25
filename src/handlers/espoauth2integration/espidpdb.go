package espoauth2integration

import (
	"constants"
	"handlers/utils"

	"strconv"
	"time"

	"models"
)

var cognitoStateTableConnection = utils.GetDynamoDbConnectionByRegion()

func StoreState(state models.State) error {
	TAG := "[StoreState] "
	utils.LogDebug(TAG + "Storing state ")

	currentTime := time.Now().Unix()
	timeToLive, errParseInt := strconv.ParseInt(utils.GetTimeToLive(), 10, 64)
	if errParseInt != nil {
		utils.LogError(TAG + "Error parsing time to live " + errParseInt.Error())
		return errParseInt
	}

	expirationTime := currentTime + timeToLive
	state.ExpiresOn = expirationTime

	return utils.DbCreateItem(cognitoStateTableConnection, constants.COGNITO_STATE_TABLE, state)
}

func GetState(key string) (models.State, error) {
	TAG := "[GetState] "
	utils.LogDebug(TAG + "Fetching state ")
	var result models.State
	errGetState := utils.DbGetItem(cognitoStateTableConnection, constants.COGNITO_STATE_TABLE, models.State{Key: key}, &result)
	if errGetState != nil {
		utils.LogError(TAG + "Error occured while fetching state, error " + errGetState.Error())
		return result, errGetState
	}
	utils.LogDebug(TAG + "Successfully fetched state")
	return result, nil
}
