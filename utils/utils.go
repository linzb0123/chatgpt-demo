package utils

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Dump(a any) []byte {
	b, _ := json.Marshal(a)
	return b
}

func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}

	return b
}
