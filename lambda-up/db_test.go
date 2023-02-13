package main

import (
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T) {
	// db := newDb()
	db := fakeDb()
	for _, sup := range sups {
		// PutOne
		err := db.putOne(sup)
		assert.NoError(t, err)
	}
	sup := sups[0]
	// GetOne
	item, err := db.getOne(sup.Where, sup.Who)
	assert.NoError(t, err)
	assert.Equal(t, sup, item)
	// GetAll
	s, err := db.getAll(sup.Where)
	if assert.NoError(t, err) {
		assert.Equal(t, sups, s)
	}
	// negative tests
	_, err = db.getOne("doesn't", "exist")
	assert.Error(t, err)
}

////////////////////////////////////////////////////////////////////////
var sups = []standup{
	standup{Who: "testuser", When: "now", Where: "foo", What: "standup"},
	standup{Who: "user2", When: "now", Where: "foo", What: "standup"},
	standup{Who: "user3", When: "now", Where: "foo", What: "standup"},
}

func fakeDb() dbi {
	return dyndb{
		dyndbi:    fakedb{items: make(map[string]attrVals)},
		tablename: "standup",
	}
}

func populatedFakeDb() dbi {
	db := fakedb{items: make(map[string]attrVals)}
	for _, sup := range sups {
		av, _ := dynamodbattribute.MarshalMap(sup)
		db.items[sup.Where+sup.Who] = av
	}
	return dyndb{
		dyndbi:    db,
		tablename: "standup",
	}
}

type attrVals = map[string]*dynamodb.AttributeValue
type fakedb struct {
	items map[string]attrVals
}

func (db fakedb) Query(in *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	key := aws.StringValue(in.ExpressionAttributeValues[":0"].S)
	results := make([]attrVals, 0, 2)
	for k, v := range db.items {
		if strings.HasPrefix(k, key) {
			results = append(results, v)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		w1 := aws.StringValue(results[i]["Who"].S)
		w2 := aws.StringValue(results[j]["Who"].S)
		return w1 < w2
	})
	return &dynamodb.QueryOutput{Items: results}, nil
}

func (db fakedb) GetItem(in *dynamodb.GetItemInput,
) (*dynamodb.GetItemOutput, error) {
	where := aws.StringValue(in.Key["Where"].S)
	who := aws.StringValue(in.Key["Who"].S)
	item, ok := db.items[where+who]
	if ok {
		return &dynamodb.GetItemOutput{Item: item}, nil
	}
	return nil, awserr.New(dynamodb.ErrCodeResourceNotFoundException, "", nil)
}

func (db fakedb) PutItem(in *dynamodb.PutItemInput,
) (*dynamodb.PutItemOutput, error) {
	where := aws.StringValue(in.Item["Where"].S)
	who := aws.StringValue(in.Item["Who"].S)
	db.items[where+who] = in.Item
	return nil, nil
}
