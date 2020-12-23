package lib

import (
	"encoding/json"
	"fmt"
	"gitlab.com/c0b/go-ordered-json"
	"reflect"
	"regexp"
	"strings"
	"time"
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

func combineOrderedMaps(orderedMaps []*ordered.OrderedMap) *ordered.OrderedMap {
	result := ordered.NewOrderedMap()
	for _, om := range orderedMaps {
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			resValue, ok := result.GetValue(pair.Key)
			if !ok {
				result.Set(pair.Key, pair.Value)
			} else if resValue != nil {
				typeInResult := getTypeName(resValue)
				pairType := getTypeName(pair.Value)
				if typeInResult != pairType {
					result.Set(pair.Key, nil)
				}
			}
		}
	}
	return result
}

func convertToOrderedMapList(val interface{}) []*ordered.OrderedMap {
	result := make([]*ordered.OrderedMap, 0)
	valList := val.([]interface{})
	for _, valItem := range valList {
		converted := valItem.(*ordered.OrderedMap)
		result = append(result, converted)
	}
	return result
}

func doesStrListContainStr(strList []string, str string) bool {
	for _, strItem := range strList {
		if strItem == str {
			return true
		}
	}
	return false
}

func getAsCamelCase(str string) string {
	link := regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}

func getAsGolangStructSection(key string, value interface{}, indentLevel int) string {
	typeName := getTypeName(value)
	fieldName := strings.Title(getAsCamelCase(key))

	if typeName == "*ordered.OrderedMap" {
		innerStructSection := strings.Repeat(" ", indentLevel*4)
		innerStructSection += fieldName + " struct {\n"

		om := value.(*ordered.OrderedMap)
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			innerStructLine := getAsGolangStructSection(pair.Key, pair.Value, indentLevel+1)
			innerStructSection += fmt.Sprintf("%s\n", innerStructLine)
		}

		finalLineFmt := "%s} `json:\"%s\"`"
		innerStructSection += fmt.Sprintf(finalLineFmt, strings.Repeat(" ", indentLevel*4), key)
		return innerStructSection
	}

	if typeName == "[]interface {}" {
		listTypes := getListTypeNames(value.([]interface{}))
		if len(listTypes) == 1 {
			typeName = fmt.Sprintf("[]%s", listTypes[0])
		}
	}

	// type name could be set to this from above section
	if typeName == "[]*ordered.OrderedMap" {
		convertedOrderedMapList := convertToOrderedMapList(value)
		combinedMap := combineOrderedMaps(convertedOrderedMapList)
		innerStructSection := strings.Repeat(" ", indentLevel*4)
		innerStructSection += fieldName + " []struct {\n"

		iter := combinedMap.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			innerStructLine := getAsGolangStructSection(pair.Key, pair.Value, indentLevel+1)
			innerStructSection += fmt.Sprintf("%s\n", innerStructLine)
		}

		finalLineFmt := "%s} `json:\"%s\"`"
		innerStructSection += fmt.Sprintf(finalLineFmt, strings.Repeat(" ", indentLevel*4), key)
		return innerStructSection
	}

	structSection := strings.Repeat(" ", indentLevel*4)
	structSection += fieldName
	structSection += " " + typeName
	structSection += " " + fmt.Sprintf("`json:\"%s\"`", key)
	return structSection
}

func getListTypeNames(list []interface{}) []string {
	typeNames := make([]string, 0)
	for _, item := range list {
		typeName := getTypeName(item)
		if typeName == "[]interface {}" {
			listTypes := getListTypeNames(item.([]interface{}))
			if len(listTypes) == 1 {
				typeName = fmt.Sprintf("[]%s", listTypes[0])
			}
		}
		if !doesStrListContainStr(typeNames, typeName) {
			typeNames = append(typeNames, typeName)
		}
	}
	return typeNames
}

func getTypeName(v interface{}) string {
	if v == nil {
		return "interface{}"
	}
	typeName := reflect.TypeOf(v).String()
	if typeName == "string" {
		if isStrInTimeFormat(v.(string)) {
			typeName = "time.Time"
		}
	}
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

func isStrInTimeFormat(str string) bool {
	validFormats := []string{
		time.ANSIC, time.RFC822, time.RFC822Z, time.RFC1123,
		time.RFC1123Z, time.RFC3339,
	}
	for _, validFormat := range validFormats {
		if _, err := time.Parse(validFormat, str); err == nil {
			return true
		}
	}
	return false
}
