package utils

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"

	"constants"
	"models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const INCORRECT_USERNAME_PASSWD = "Incorrect user name or password"

func DbCreateItem(dbhandle *dynamodb.DynamoDB, table_name string, v interface{}) error {
	av, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		return err
	}
	_, err = dbhandle.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(table_name),
		Item:      av,
	})
	return err
}

/* The entire composite key: Primary + Range Key must be used */
func DbGetItem(dbhandle *dynamodb.DynamoDB, table_name string, query interface{}, out interface{}) error {

	av, err := dynamodbattribute.MarshalMap(query)
	if err != nil {
		return err
	}

	result, err := dbhandle.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table_name),
		Key:       av,
	})
	if err != nil {
		return err
	}
	return dynamodbattribute.UnmarshalMap(result.Item, out)
}

/* Query based on some conditions on the composite keys */
func DbQuery(dbhandle *dynamodb.DynamoDB, tableName, indexName string, limit int64, startKey map[string]*dynamodb.AttributeValue, expr expression.Expression, out interface{}) (error, map[string]*dynamodb.AttributeValue) {
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	if indexName != constants.EMPTY_STRING {
		queryInput.SetIndexName(indexName)
	}

	if limit != 0 {
		queryInput.SetLimit(limit)
	}

	if startKey != nil {
		queryInput.SetExclusiveStartKey(startKey)
	}

	result, err := dbhandle.Query(queryInput)
	if err != nil {
		return err, nil
	}
	return dynamodbattribute.UnmarshalListOfMaps(result.Items, out), result.LastEvaluatedKey
}

func GetDynamoDbConnectionByRegion() *dynamodb.DynamoDB {
	return dynamodb.New(session.New(), aws.NewConfig().WithRegion(GetRegion()))
}

func UserDBGetByName(tenctx models.TenantContext, userName string) (models.UserDetails, error) {
	TAG := "[UserDBGetByName] "
	LogDebug(TAG + "Fetching UserId from userName, customerId: " + tenctx.CustomerId + " userName " + userName)

	var usersTableConnection = GetDynamoDbConnectionByRegion()
	userLists := []models.UserDetails{}
	userDetails := models.UserDetails{}
	errUserNotExists := errors.New(INCORRECT_USERNAME_PASSWD)

	usersTable := TenantGetDBName(tenctx, constants.USERS)

	keyCondition := expression.KeyEqual(expression.Key("user_name"), expression.Value(userName))

	expr, errExpressionBuilder := expression.NewBuilder().
		WithKeyCondition(keyCondition).Build()
	if errExpressionBuilder != nil {
		LogError(TAG + "Error occured while building expression " + errExpressionBuilder.Error())
		return userDetails, errExpressionBuilder
	}

	errGetUserByUserName, _ := DbQuery(usersTableConnection, usersTable, constants.INDEX_USERS_TABLE, 0, nil, expr, &userLists)
	if errGetUserByUserName != nil {
		LogError(TAG + "Error occured while fetching user " + errGetUserByUserName.Error() + " for customer: " + tenctx.CustomerId + " for user: " + userName)
		return userDetails, errGetUserByUserName
	}

	if len(userLists) == 0 {
		LogDebug(TAG + "No users found for customer: " + tenctx.CustomerId)
		return userDetails, errUserNotExists
	}
	userDetails = userLists[0]
	LogDebug(TAG + "Successfully fetched user by name, UserId: " + userDetails.UserId)
	return userDetails, nil
}
