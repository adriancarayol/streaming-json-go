package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//func TestSimpleCompleteJSON(t *testing.T) {
//	json_p0000 := `{"name`
//	lexer := NewLexer()
//	errInAppendString := lexer.AppendString(json_p0000)
//
//	ret := lexer.CompleteJSON()
//
//	assert.Nil(t, errInAppendString)
//
//	assert.Equal(t, `{"name":null}`, ret, "the token should be equal")
//}

func TestCompleteJSON(t *testing.T) {
	streamingJSONCase := map[string]string{
		`{`:                 `{}`,        // mirror stack: [], should remove from stack: [], should push into mirror stack: [`}`]
		`{}`:                `{}`,        // mirror stack: [], should remove from stack: [], should push into mirror stack: []
		`{"`:                `{"":null}`, // mirror stack: [`}`], should remove from stack: [], should push into mirror stack: [`"`, `:`, `n`, `u`, `l`, `l`]
		`{""`:               `{"":null}`, // mirror stack: [`"`, `:`, `n`, `u`, `l`, `l`,`}`], should remove from stack: [`"`], should push into mirror stack: []
		`{"a`:               `{"a":null}`,
		`{"a"`:              `{"a":null}`,
		`{"a":`:             `{"a":null}`,
		`{"a":n`:            `{"a":null}`,
		`{"a":nu`:           `{"a":null}`,
		`{"a":nul`:          `{"a":null}`,
		`{"a":null`:         `{"a":null}`,
		`{"a":null,`:        `{"a":null}`, // can not detect context, remove ","
		`{"a":t`:            `{"a":true}`,
		`{"a":tr`:           `{"a":true}`,
		`{"a":tru`:          `{"a":true}`,
		`{"a":true`:         `{"a":true}`,
		`{"a":true,`:        `{"a":true}`, // can not detect context, remove ","
		`{"a":f`:            `{"a":false}`,
		`{"a":fa`:           `{"a":false}`,
		`{"a":fal`:          `{"a":false}`,
		`{"a":fals`:         `{"a":false}`,
		`{"a":false`:        `{"a":false}`,
		`{"a":false,`:       `{"a":false}`, // can not detect context, remove ","
		`{"a":12`:           `{"a":12}`,
		`{"a":12,`:          `{"a":12}`, // can not detect context, remove ","
		`{"a":12.`:          `{"a":12.0}`,
		`{"a":12.15`:        `{"a":12.15}`,
		`{"a":12.15,`:       `{"a":12.15}`, // can not detect context, remove ","
		`{"a":"`:            `{"a":""}`,
		`{"a":""`:           `{"a":""}`,
		`{"a":"",`:          `{"a":""}`, // can not detect context, remove ","
		`{"a":"string`:      `{"a":"string"}`,
		`{"a":"string"`:     `{"a":"string"}`,
		`{"a":"string",`:    `{"a":"string"}`, // can not detect context, remove ","
		`{"a":"","`:         `{"a":"","":null}`,
		`{"a":"","b`:        `{"a":"","b":null}`,
		`{"a":"","b"`:       `{"a":"","b":null}`,
		`{"a":"","b":`:      `{"a":"","b":null}`,
		`{"a":"","b":"`:     `{"a":"","b":""}`,
		`{"a":"","b":""`:    `{"a":"","b":""}`,
		`{"a":"","b":""}`:   `{"a":"","b":""}`,
		`{"1`:               `{"1":null}`,
		`{"1.`:              `{"1.":null}`,
		`{"1.1`:             `{"1.1":null}`,
		`{"1.10`:            `{"1.10":null}`,
		`{"1"`:              `{"1":null}`,
		`{"1":`:             `{"1":null}`,
		`{"1":"`:            `{"1":""}`,
		`{"1":"1`:           `{"1":"1"}`,
		`{"1":"1.`:          `{"1":"1."}`,
		`{"1":"1.1`:         `{"1":"1.1"}`,
		`{"1":"1.10`:        `{"1":"1.10"}`,
		`{"1":"1"`:          `{"1":"1"}`,
		`{"1":"1"}`:         `{"1":"1"}`,
		`{"t`:               `{"t":null}`,
		`{"tr`:              `{"tr":null}`,
		`{"tru`:             `{"tru":null}`,
		`{"true`:            `{"true":null}`,
		`{"true"`:           `{"true":null}`,
		`{"true":`:          `{"true":null}`,
		`{"true":"t`:        `{"true":"t"}`,
		`{"true":"tr`:       `{"true":"tr"}`,
		`{"true":"tru`:      `{"true":"tru"}`,
		`{"true":"true`:     `{"true":"true"}`,
		`{"true":"true"`:    `{"true":"true"}`,
		`{"true":"true"}`:   `{"true":"true"}`,
		`{"f`:               `{"f":null}`,
		`{"fa`:              `{"fa":null}`,
		`{"fal`:             `{"fal":null}`,
		`{"fals`:            `{"fals":null}`,
		`{"false`:           `{"false":null}`,
		`{"false"`:          `{"false":null}`,
		`{"false":`:         `{"false":null}`,
		`{"false":"f`:       `{"false":"f"}`,
		`{"false":"fa`:      `{"false":"fa"}`,
		`{"false":"fal`:     `{"false":"fal"}`,
		`{"false":"fals`:    `{"false":"fals"}`,
		`{"false":"false`:   `{"false":"false"}`,
		`{"false":"false"`:  `{"false":"false"}`,
		`{"false":"false"}`: `{"false":"false"}`,
		`{"n`:               `{"n":null}`,
		`{"nu`:              `{"nu":null}`,
		`{"nul`:             `{"nul":null}`,
		`{"null`:            `{"null":null}`,
		`{"null"`:           `{"null":null}`,
		`{"null":`:          `{"null":null}`,
		`{"null":"n`:        `{"null":"n"}`,
		`{"null":"nu`:       `{"null":"nu"}`,
		`{"null":"nul`:      `{"null":"nul"}`,
		`{"null":"null`:     `{"null":"null"}`,
		`{"null":"null"`:    `{"null":"null"}`,
		`{"null":"null"}`:   `{"null":"null"}`,
		//`[`:               `[]`,
		//`[]`:              `[]`,
		//`[n`:              `[null]`,
		//`[nu`:             `[null]`,
		//`[nul`:            `[null]`,
		//`[null`:           `[null]`,
		//`[null,`:          `[null]`, // can not detect context, remove ","
		//`[t`:              `[true]`,
		//`[tr`:             `[true]`,
		//`[tru`:            `[true]`,
		//`[true`:           `[true]`,
		//`[true,`:          `[true]`, // can not detect context, remove ","
		//`[f`:              `[false]`,
		//`[fa`:             `[false]`,
		//`[fal`:            `[false]`,
		//`[fals`:           `[false]`,
		//`[false`:          `[false]`,
		//`[false,`:         `[false]`, // can not detect context, remove ","
		//`[0`:              `[0]`,
		//`[0,`:             `[0]`, // can not detect context, remove ","
		//`[0.`:             `[0.0]`,
		//`[0.1`:            `[0.1]`,
		//`[0.12,`:          `[0.12]`, // can not detect context, remove ","
		//`["`:              `[""]`,
		//`[""`:             `[""]`,
		//`["",`:            `[""]`, // can not detect context, remove ","
		//`["a`:             `["a"]`,
		//`["a"`:            `["a"]`,
		//`["a",`:           `["a"]`, // can not detect context, remove ","

	}
	for testCase, expect := range streamingJSONCase {
		fmt.Printf("\n\n---------------------------\n")
		fmt.Printf("current test case: `%s`\n", testCase)
		lexer := NewLexer()
		errInAppendString := lexer.AppendString(testCase)
		ret := lexer.CompleteJSON()
		assert.Nil(t, errInAppendString)
		if !assert.Equal(t, expect, ret, "unexpected JSON") {
			break
		}
	}
}
