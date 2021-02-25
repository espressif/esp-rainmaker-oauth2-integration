package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var lambdaClient = lambda.New(session.New(), aws.NewConfig().WithRegion(GetRegion()))

func GetLambdaEnv(funtionName string) (*lambda.EnvironmentResponse, error) {
	TAG := "[GetLambdaEnv] "
	LogDebug(TAG + "Getting environment for lambda: " + funtionName)

	input := &lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(funtionName),
	}

	result, errGetEnv := lambdaClient.GetFunctionConfiguration(input)
	if errGetEnv != nil {
		LogError(TAG + "Error occurred while fetching environment for lambda: " + funtionName + " with error, " + errGetEnv.Error())
		return nil, errGetEnv
	}

	LogDebug(TAG + "Fetched environment for lambda successfully")
	return result.Environment, nil

}

func UpdateLambdaEnv(funtionName string, environment *lambda.EnvironmentResponse) error {
	TAG := "[UpdateLambdaEnv] "
	LogDebug(TAG + "Updating environment for lambda: " + funtionName)

	env := new(lambda.Environment)

	env.Variables = environment.Variables
	input := &lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(funtionName),
		Environment:  env,
	}

	_, errUpdateEnv := lambdaClient.UpdateFunctionConfiguration(input)
	if errUpdateEnv != nil {
		LogError(TAG + "Error occurred while updating environment for lambda: " + funtionName + " with error, " + errUpdateEnv.Error())
		return errUpdateEnv
	}

	LogDebug(TAG + "Updated environment for lambda successfully")
	return nil
}
