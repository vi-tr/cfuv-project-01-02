package main

import (
    "context"
    "encoding/json"
)

type InMemoryDB[Schema any] struct {
    lastid uint
    data *map[uint]Schema
}

func (self InMemoryDB[Schema]) Find(c context.Context, query any) (result []Schema, ε ε) {
    mq := make(map[string]interface{})
    F1(json.Unmarshal(F2(json.Marshal(query)), &mq))

    outer: for _, v := range *self.data {
        mv := make(map[string]interface{})
        F1(json.Unmarshal(F2(json.Marshal(v)), &mv))
        for k, vv := range mv {
            if vq, ok := mq[k]; ok {
                if vv != vq { continue outer }
            }
        }
        result = append(result, v)
    }
    return
}

func (self InMemoryDB[Schema]) Add(c context.Context, addition Schema) (ε ε) {
    self.lastid += 1
    (*self.data)[self.lastid] = addition
    return
}

func (self InMemoryDB[Schema]) Edit(c context.Context, query any, change Schema) (ε ε) {
    mq := make(map[string]interface{})
    F1(json.Unmarshal(F2(json.Marshal(query)), &mq))

    outer: for i, v := range *self.data {
        mv := make(map[string]interface{})
        F1(json.Unmarshal(F2(json.Marshal(v)), &mv))
        for k, vv := range mv {
            if vq, ok := mq[k]; ok {
                if vv != vq { continue outer }
            }
        }
        (*self.data)[i] = change
        return
    }
    return NothingFoundε
}

func (self InMemoryDB[Schema]) Remove(c context.Context, query any) (ε ε) {
    mq := make(map[string]interface{})
    F1(json.Unmarshal(F2(json.Marshal(query)), &mq))

    outer: for i, v := range *self.data {
        mv := make(map[string]interface{})
        F1(json.Unmarshal(F2(json.Marshal(v)), &mv))
        for k, vv := range mv {
            if vq, ok := mq[k]; ok {
                if vv != vq { continue outer }
            }
        }
        delete(*self.data, i)
    }
    return
}

func NewInMemoryDB[Schema any]() (result InMemoryDB[Schema], ε ε) { return }

var _ DBProvider[struct{}] = InMemoryDB[struct{}]{} // check if InMemoryDB implements DBProvider
