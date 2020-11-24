package lib

import (
	"encoding/json"
	"fmt"
	"gitlab.com/c0b/go-ordered-json"
	"reflect"
	"strings"
)

type JsonStr string

func NewJsonStr(rawJsonStr string) (JsonStr, error) {
	jsonStr := JsonStr(rawJsonStr)
	err := jsonStr.validate()
	if err != nil {
		return "", err
	}
	return jsonStr, nil
}

func (jsonStr JsonStr) GetAsGolangString() (string, error) {
	om, err := jsonStr.getAsOrderedMap()
	if err != nil {
		return "", err
	}

	finalString := "type Generated struct {\n"
	iter := om.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		structLine := getAsGolangStructSection(pair.Key, pair.Value, 1)
		finalString += fmt.Sprintf("%s\n", structLine)
	}
	finalString += "}"

	return finalString, nil
}

func (jsonStr JsonStr) getAsOrderedMap() (*ordered.OrderedMap, error) {
	jsonMap := ordered.NewOrderedMap()
	err := json.Unmarshal([]byte(jsonStr), jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}

func (jsonStr JsonStr) validate() error {
	var js map[string]interface{}
	return json.Unmarshal([]byte(jsonStr), &js)
}

func getAsGolangStructSection(key string, value interface{}, indentLevel int) string {
	typeName := getTypeName(value)
	fieldName := strings.Title(key)

	if typeName == "*ordered.OrderedMap" {
		innerStructSection := strings.Repeat(" ", indentLevel * 4)
		innerStructSection += fieldName + " struct {\n"

		om := value.(*ordered.OrderedMap)
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			innerStructLine := getAsGolangStructSection(pair.Key, pair.Value, indentLevel + 1)
			innerStructSection += fmt.Sprintf("%s\n", innerStructLine)
		}

		finalLineFmt := "%s} `json:\"%s\"`"
		innerStructSection += fmt.Sprintf(finalLineFmt, strings.Repeat(" ", indentLevel * 4), key)
		return innerStructSection
	}

	structSection := strings.Repeat(" ", indentLevel * 4)
	structSection += fieldName
	structSection += " " + typeName
	structSection += " " + fmt.Sprintf("`json:\"%s\"`", key)
	return structSection
}

func getTypeName(v interface{}) string {
	if v == nil {
		return "interface{}"
	}
	typeName := reflect.TypeOf(v).String()
	if typeName == "json.Number" {
		if _, err := v.(json.Number).Float64(); err == nil {
			typeName = "float64"
		}
		if _, err := v.(json.Number).Int64(); err == nil {
			typeName = "int"
		}
	}
	return typeName
}
