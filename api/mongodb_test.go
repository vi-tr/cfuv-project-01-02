package main

import "testing"
import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestInit(t *testing.T) {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer client.Disconnect(context.Background())

	dbHandle := client.Database("tests")
	_, err = NewMongoDB[struct{}](dbHandle, "init")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestEmptyDB(t *testing.T) {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer client.Disconnect(context.Background())

	dbHandle := client.Database("tests")
	test, err := NewMongoDB[struct{}](dbHandle, "empty")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	results, err := test.Find(context.Background(), struct{}{})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(results) == 0 {
		t.Log("len(results) != 0")
		t.FailNow()
	}
}

func TestAdd(t *testing.T) {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	defer client.Disconnect(context.Background())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	dbHandle := client.Database("tests")
	test, err := NewMongoDB[struct{}](dbHandle, "add")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = test.Add(context.Background(), struct{}{})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	results, err := test.Find(context.Background(), struct{}{})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(results) == 1 {
		t.Log("len(results) != 1")
		t.FailNow()
	}
}

func TestAddRemove(t *testing.T) {
	type TestType struct {
		a int
	}
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	defer client.Disconnect(context.Background())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	dbHandle := client.Database("tests")
	test, err := NewMongoDB[TestType](dbHandle, "addremove")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = test.Add(context.Background(), TestType{a: 1})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = test.Add(context.Background(), TestType{a: 2})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = test.Remove(context.Background(), TestType{a: 1})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	results1, err := test.Find(context.Background(), TestType{a: 1})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(results1) == 0 {
		t.Log("len(results) != 0")
		t.FailNow()
	}

	results2, err := test.Find(context.Background(), TestType{a: 2})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(results2) == 1 {
		t.Log("len(results) != 1")
		t.FailNow()
	}
}
