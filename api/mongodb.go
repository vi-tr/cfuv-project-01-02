package main

import (
    "context"

    "go.mongodb.org/mongo-driver/mongo"
)

type MongoDB[Schema any] struct {
    collectionHandle *mongo.Collection
}

func (self MongoDB[Schema]) Find(c context.Context, query any) (result []Schema, ε ε) {
    defer ə(&ε)
    cursor := P2(self.collectionHandle.Find(c, query))
    P1(cursor.All(c, &result))
    return
}

func (self MongoDB[Schema]) Add(c context.Context, addition Schema) (ε ε) {
    _, ε = self.collectionHandle.InsertOne(c, addition); return
}

func (self MongoDB[Schema]) Edit(c context.Context, query any, change Schema) (ε ε) {
    _, ε = self.collectionHandle.UpdateOne(c, query, change); return
}

func (self MongoDB[Schema]) Remove(c context.Context, query any) (ε ε) {
    _, ε = self.collectionHandle.DeleteMany(c, query); return
}

func NewMongoDB[Schema any](db *mongo.Database, collection string) (result MongoDB[Schema], ε ε) {
    result.collectionHandle = db.Collection(collection); return
}

var _ DBProvider[struct{}] = MongoDB[struct{}]{} // check if MongoDB implements DBProvider
