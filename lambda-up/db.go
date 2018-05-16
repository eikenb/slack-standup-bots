package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type dbi interface {
	getOne(where, who string) (standup, error)
	getAll(where string) ([]standup, error)
	putOne(standup) error
}

type dyndbi interface {
	Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
}

type dyndb struct {
	dyndbi
	tablename string
}

func realDb() dbi {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		log.Fatal(err)
	}

	return dyndb{
		dyndbi:    dynamodb.New(sess),
		tablename: "standup",
	}
}

type standup struct {
	Where, Who, When, What string
}

func (s standup) String(err error) string {
	if err != nil || s.Who == "" {
		return ""
	}
	return fmt.Sprintf("%s - %s\n%s", s.Who, s.When, s.What)
}

//
func (db dyndb) getAll(where string) ([]standup, error) {
	where_k := expression.Key("Where").Equal(expression.Value(where))
	exp, err := expression.NewBuilder().WithKeyCondition(where_k).Build()
	if err != nil {
		log.Fatal(err)
	}
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(db.tablename),
		KeyConditionExpression:    exp.KeyCondition(),
		ExpressionAttributeValues: exp.Values(),
		ExpressionAttributeNames:  exp.Names(),
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}

	results, err := db.Query(input)
	if err != nil {
		return nil, err
	}
	standups := []standup{}
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &standups)
	if err != nil {
		return nil, err
	}

	return standups, nil
}

func (db dyndb) getOne(where, who string) (standup, error) {
	key := map[string]*dynamodb.AttributeValue{
		"Where": {S: aws.String(where)},
		"Who":   {S: aws.String(who)},
	}
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.tablename),
		Key:       key,
	}
	if err := input.Validate(); err != nil {
		return standup{}, err
	}
	result, err := db.GetItem(input)
	if err != nil {
		return standup{}, err
	}

	item := standup{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	return item, nil
}

//
func (db dyndb) putOne(item standup) error {

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.tablename),
	}
	_, err = db.PutItem(input)

	return err
}
