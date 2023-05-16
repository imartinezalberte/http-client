package httpclient

import (
	"strings"

	ofuscateStruct "github.com/imartinezalberte/ofuscate-struct"
	"golang.org/x/exp/constraints"
)

const (
	ofuscation = "XXX"
	comma      = ","
)

type (
	OfuscateParamsFunc func(map[string][]string) map[string]string

	OfuscateBodyFunc func(interface{}) map[string]interface{}
)

func ofuscate(params map[string]any) OfuscateParamsFunc {
	return func(input map[string][]string) map[string]string {
		output := make(map[string]string)
		for key, val := range input {
			if _, ok := params[key]; ok {
				output[key] = ofuscation
			} else {
				output[key] = strings.Join(val, ",")
			}
		}
		return output
	}
}

func ofuscateBody(params []string) OfuscateBodyFunc {
	return func(input interface{}) map[string]interface{} {
		if input == nil {
			return nil
		}

		m := input.(map[string]interface{})
		for _, param := range params {
			ofuscateStruct.DoOfuscate(m, param)
		}

		return m
	}
}

func ArrToMap[T comparable](input []T) map[T]any {
	output := make(map[T]any)
	for _, v := range input {
		output[v] = 0
	}
	return output
}

func CheckValueInRange[T constraints.Integer](min, max, bydefault T) func(T) T {
	return func(input T) T {
		if input < min || input > max {
			input = bydefault
		}
		return input
	}
}
