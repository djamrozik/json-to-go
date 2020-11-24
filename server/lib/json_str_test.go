package lib

import "testing"

func TestBasic(t *testing.T) {
	rawJsonStr := "{" +
		"    \"info\": \"some info\"" +
		"}"

	expectedResult := "type Generated struct {\n" +
		"    Info string `json:\"info\"`\n" +
		"}"

	jsonStr, err := NewJsonStr(rawJsonStr)
	if err != nil {
		t.Error("Unexpected error creating new JsonStr", err)
		return
	}
	actualResult, err := jsonStr.GetAsGolangString()
	if err != nil {
		t.Error("Unexpected error returned while converting json to golang", err)
		return
	}
	if actualResult != expectedResult {
		reportIncorrectResults(t, expectedResult, actualResult)
	}
}

func TestPrimitiveTypes(t *testing.T) {
	rawJsonStr := "{" +
		"    \"info\": \"random info here\"," +
		"    \"someInt\": 452," +
		"    \"someFloat\": 1.25," +
		"    \"someNull\": null," +
		"    \"someBoolean\": false" +
		"}"

	expectedResult := "type Generated struct {\n" +
		"    Info string `json:\"info\"`\n" +
		"    SomeInt int `json:\"someInt\"`\n" +
		"    SomeFloat float64 `json:\"someFloat\"`\n" +
		"    SomeNull interface{} `json:\"someNull\"`\n" +
		"    SomeBoolean bool `json:\"someBoolean\"`\n" +
		"}"

	jsonStr, err := NewJsonStr(rawJsonStr)
	if err != nil {
		t.Error("Unexpected error creating new JsonStr", err)
		return
	}
	actualResult, err := jsonStr.GetAsGolangString()
	if err != nil {
		t.Error("Unexpected error returned while converting json to golang", err)
		return
	}
	if actualResult != expectedResult {
		reportIncorrectResults(t, expectedResult, actualResult)
	}
}

// The order of the keys in the json string should match the order of the keys
// in the resulting golang struct.
//
// This test does not 100% guarantee that the key order will stay the same, but successful completion
// means a very high likelyhood that they do.
func TestKeyOrder(t *testing.T) {
	rawJsonStr := "{" +
		"    \"a\": 1," +
		"    \"b\": 2," +
		"    \"c\": 3," +
		"    \"d\": 4," +
		"    \"e\": 5" +
		"}"

	expectedResult := "type Generated struct {\n" +
		"    A int `json:\"a\"`\n" +
		"    B int `json:\"b\"`\n" +
		"    C int `json:\"c\"`\n" +
		"    D int `json:\"d\"`\n" +
		"    E int `json:\"e\"`\n" +
		"}"

	numRuns := 100
	for i := 0; i < numRuns; i++ {
		jsonStr, err := NewJsonStr(rawJsonStr)
		if err != nil {
			t.Error("Unexpected error creating new JsonStr", err)
			return
		}
		actualResult, err := jsonStr.GetAsGolangString()
		if err != nil {
			t.Error("Unexpected error returned while converting json to golang", err)
			return
		}
		if actualResult != expectedResult {
			reportIncorrectResults(t, expectedResult, actualResult)
		}
	}
}

func TestInvalidJSON(t *testing.T) {
	rawJsonStr := "{" +
		"    \"a\": 1," +
		"    \"b\": 2," +
		"    \"c\": 3" +
		"    \"d\": 4," +
		"    \"e\": 5" +
		"}"

	_, err := NewJsonStr(rawJsonStr)
	if err != nil {
		return
	}
	t.Error("Did not get error when there should have been")
}

func TestBasicInnerObject(t *testing.T) {
	rawJsonStr := "{" +
		"    \"info\": \"some info\"," +
		"    \"inner\": {\"innerKey\": true}" +
		"}"

	expectedResult :=
		"type Generated struct {\n" +
		"    Info string `json:\"info\"`\n" +
		"    Inner struct {\n" +
		"        InnerKey bool `json:\"innerKey\"`\n" +
		"    } `json:\"inner\"`\n" +
		"}"

	jsonStr, err := NewJsonStr(rawJsonStr)
	if err != nil {
		t.Error("Unexpected error creating new JsonStr", err)
		return
	}
	actualResult, err := jsonStr.GetAsGolangString()
	if err != nil {
		t.Error("Unexpected error returned while converting json to golang", err)
		return
	}
	if actualResult != expectedResult {
		reportIncorrectResults(t, expectedResult, actualResult)
	}
}

func TestKeyNames(t *testing.T) {
	rawJsonStr := "{" +
		"    \"info\": \"random info here\"," +
		"    \"someInt\": 452," +
		"    \"some_float\": 1.25" +
		"}"

	expectedResult := "type Generated struct {\n" +
		"    Info string `json:\"info\"`\n" +
		"    SomeInt int `json:\"someInt\"`\n" +
		"    SomeFloat float64 `json:\"some_float\"`\n" +
		"}"

	jsonStr, err := NewJsonStr(rawJsonStr)
	if err != nil {
		t.Error("Unexpected error creating new JsonStr", err)
		return
	}
	actualResult, err := jsonStr.GetAsGolangString()
	if err != nil {
		t.Error("Unexpected error returned while converting json to golang", err)
		return
	}
	if actualResult != expectedResult {
		reportIncorrectResults(t, expectedResult, actualResult)
	}
}

func reportIncorrectResults(t *testing.T, expected, actual string) {
	formatStr := "Actual result did not match expected.\n--- Expected --------\n%s\n"
	formatStr += "-------------------\n"
	formatStr += "--- Actual --------\n%s\n-------------------"
	t.Errorf(formatStr, expected, actual)
}
