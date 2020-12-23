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
		"    \"someBoolean\": false," +
		"    \"someTime\": \"2019-08-22T18:10:37Z\"" +
		"}"

	expectedResult := "type Generated struct {\n" +
		"    Info string `json:\"info\"`\n" +
		"    SomeInt int `json:\"someInt\"`\n" +
		"    SomeFloat float64 `json:\"someFloat\"`\n" +
		"    SomeNull interface{} `json:\"someNull\"`\n" +
		"    SomeBoolean bool `json:\"someBoolean\"`\n" +
		"    SomeTime time.Time `json:\"someTime\"`\n" +
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
		"    \"some_float\": 1.25," +
		"    \"some_weird_BooleanHere\": true" +
		"}"

	expectedResult := "type Generated struct {\n" +
		"    Info string `json:\"info\"`\n" +
		"    SomeInt int `json:\"someInt\"`\n" +
		"    SomeFloat float64 `json:\"some_float\"`\n" +
		"    SomeWeirdBooleanHere bool `json:\"some_weird_BooleanHere\"`\n" +
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

func TestArrays(t *testing.T) {
	rawJsonStr := "{" +
		"    \"info\": \"random info here\"," +
		"    \"simpleIntArray\": [1, 2, 3]," +
		"    \"simpleCombinedArray\": [1, true, \"a\"]," +
		"    \"twoDIntArray\": [[1,2,3], [1,2,3]]" +
		"}"

	expectedResult := "type Generated struct {\n" +
		"    Info string `json:\"info\"`\n" +
		"    SimpleIntArray []int `json:\"simpleIntArray\"`\n" +
		"    SimpleCombinedArray []interface {} `json:\"simpleCombinedArray\"`\n" +
		"    TwoDIntArray [][]int `json:\"twoDIntArray\"`\n" +
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

func TestArrayOfObjects(t *testing.T) {
	rawJsonStr := "{" +
		"    \"info\": \"random info here\"," +
		"    \"simpleIntArray\": [1, 2, 3]," +
		"    \"arrayOfObjects\": [" +
		"        {\"a\": 1, \"b\": true}," +
		"        {\"a\": 5, \"b\": 0.4}," +
		"        {\"a\": 5, \"b\": 0.6, \"c\": \"strr\"}" +
		"	 ]" +
		"}"

	expectedResult := "type Generated struct {\n" +
		"    Info string `json:\"info\"`\n" +
		"    SimpleIntArray []int `json:\"simpleIntArray\"`\n" +
		"    ArrayOfObjects []struct {\n" +
		"        A int `json:\"a\"`\n" +
		"        B interface{} `json:\"b\"`\n" +
		"        C string `json:\"c\"`\n" +
		"    } `json:\"arrayOfObjects\"`\n" +
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
