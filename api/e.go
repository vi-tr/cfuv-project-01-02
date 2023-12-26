package main

import (
    "log"
    "reflect"
    "errors"
)

type ε = error
var Assertε ε = errors.New("Assertion error")
func If[C, T any](c C, a T, b T) T {if reflect.ValueOf(c).IsZero() {return b}; return a}
func Or[T any](a T, b T) T {return If(a,a,b)}
func ə(e *ε) {if r:=recover(); r!=nil {*e = r.(ε)}}
func F1(e ε) {if e!=nil { log.Fatal(e) }}
func P1(e ε) {if e!=nil { panic(e) }}
func L1(e ε) {if e!=nil { log.Println(e) }}
func F2[T any] (x T, e ε) T {if e!=nil { log.Fatal(e) }; return x}
func P2[T any] (x T, e ε) T {if e!=nil { panic(e) }; return x}
func L2[T any] (x T, e ε) T {if e!=nil { log.Println(e) }; return x}
func FJ(e ...ε) {if je:=errors.Join(e...); je!=nil {log.Fatal(je)}}
func PJ(e ...ε) {if je:=errors.Join(e...); je!=nil {panic(je)}}
func LJ(e ...ε) {if je:=errors.Join(e...); je!=nil {log.Println(je)}}
func FA(p bool) {if !p { log.Fatal(Assertε) }}
func PA(p bool) {if !p { panic(Assertε) }}
func LA(p bool) {if !p { log.Println(Assertε) }}
