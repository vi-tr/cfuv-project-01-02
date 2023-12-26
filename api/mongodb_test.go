package main

import "testing"
import (
    "os"
    "context"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func TestInit(t *testing.T) {
    client, err := mongo.Connect(context.Background(),
        options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
    if err != nil { t.Log(err); t.FailNow() }
    defer client.Disconnect(context.Background())

    dbHandle := client.Database("tests")
    _, err = NewMongoDB[struct{}](dbHandle, "init")
    if err != nil { t.Log(err); t.FailNow() }
}

func TestAddFind(t *testing.T) {
    client, err := mongo.Connect(context.Background(),
        options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
    if err != nil { t.Log(err); t.FailNow() }
    defer client.Disconnect(context.Background())

    dbHandle := client.Database("tests")
    test, err := NewMongoDB[struct{}](dbHandle, "addfind")
    if err != nil { t.Log(err); t.FailNow() }

    test.Add(context.Background(), struct{}{})
}
