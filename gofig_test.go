package gofig

import (
	"encoding/json"
	"reflect"
	"testing"
)

type keyTest struct {
	ThisIsAKey bool
	Key2       string
	Key3Time   float64
	Key4here   string
}

type testObj struct {
	Bool   bool
	Float  float64
	Nested *testSubObj
}

type testSubObj struct {
	Wow string
	Hah string
}

type testPair struct {
	Key string
	Int float64
}

type nullTest struct {
	Key   string
	Float float64
	Bool  bool
}

func testJson() string {
	return `{
        "string":"value", "string2":null, "int":34, "float":23.34, "bool":true, "array": [1,2,3,"Test"],
        "obj": {"bool": false,"float": 1.89,"nested": {"wow": "really?", "hah":null}},
        "key_test": {"this_is_a_key":true, "key_2": "value!", "key_3_time": 12.34, "key4here": "checking in"},
        "obj_array": [{"key": "value","int": 10}, {"key": "pair","int": 26}],
        "map": {"name":"john doe", "age":43, "active":true},
        "nullsoft": {"key": null, "float":null, "bool":null},
        "null_obj": null
    }`
}

func createTestConfig(t *testing.T) (conf *Config) {
	conf, err := initConfig()
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
	if e := json.Unmarshal([]byte(testJson()), &conf.values); e != nil {
		t.Errorf("Failed to parse testing json: %v", e)
		t.FailNow()
	}
	return
}

func TestGofigInt(t *testing.T) {
	conf := createTestConfig(t)
	expected := 34
	if i, e := conf.Int("int"); e != nil {
		t.Errorf("Int() failed to fetch integer value: %v", e)
	} else if i != expected {
		t.Errorf("Unexpected value returned from Int(): expected %d, got %d", expected, i)
	}
}

func TestGofigFloat(t *testing.T) {
	conf := createTestConfig(t)
	expected := 23.34
	if f, e := conf.Float("float"); e != nil {
		t.Errorf("Float() failed to fetch float value: %v", e)
	} else if f != expected {
		t.Errorf("Unexpected value returned from Float(): expected %f, got %f", expected, f)
	}
}

func TestGofigBool(t *testing.T) {
	conf := createTestConfig(t)
	expected := true
	if b, e := conf.Bool("bool"); e != nil {
		t.Errorf("Bool() failed to fetch boolean value: %v", e)
	} else if !expected {
		t.Errorf("Unexpected value returned from Bool(): expected %v, got %v", expected, b)
	}
}

func TestGofigMap(t *testing.T) {
	conf := createTestConfig(t)
	expected := map[string]interface{}{
		"name":   "john doe",
		"age":    43.0,
		"active": true,
	}

	if m, e := conf.Map("map"); e != nil {
		t.Errorf("Map() failed to fetch map value: %v", e)
	} else if !reflect.DeepEqual(m, expected) {
		t.Errorf("Unexpected value returned from Map(): expected %v, got %v", expected, m)
	}
}

func TestGofigArray(t *testing.T) {
	conf := createTestConfig(t)
	expected := []interface{}{1.0, 2.0, 3.0, "Test"}
	if a, e := conf.Array("array"); e != nil {
		t.Errorf("Array() failed to fetch array value: %v", e)
	} else if !reflect.DeepEqual(a, expected) {
		t.Errorf("Unexpected value returned from Array(): expected %v, got %v.", expected, a)
	}
}

func TestGofigString(t *testing.T) {
	conf := createTestConfig(t)
	expected := "value"
	if s, e := conf.Str("string"); e != nil {
		t.Errorf("Str() failed to fetch string value: %v", e)
	} else if s != expected {
		t.Errorf("Unexpected value returned from Str(): expected '%s', got '%s'", expected, s)
	}
}

func TestGofigNullString(t *testing.T) {
	conf := createTestConfig(t)
	expected := ""
	if s, e := conf.Str("string2"); e != nil {
		t.Errorf("Str() failed to fetch string value: %v", e)
	} else if s != expected {
		t.Errorf("Unexpected value returned from Str(): expected '%s', got '%s'", expected, s)
	}
}

func TestGofigNestedStruct(t *testing.T) {
	conf := createTestConfig(t)
	obj := &testObj{}

	if err := conf.Struct("obj", obj); err != nil {
		t.Errorf("Struct() failed: %v", err)
		t.FailNow()
	}

	if obj.Nested == nil {
		t.Error("Struct() failed to map nested struct.")
	}
}

func TestGofigNullStringInStruct(t *testing.T) {
	conf := createTestConfig(t)
	expected := ""
	obj := &testObj{Nested: &testSubObj{Hah: "should get cleared"}}

	if err := conf.Struct("obj", obj); err != nil {
		t.Errorf("Struct() failed: %v", err)
		t.FailNow()
	}

	if obj.Nested.Hah != expected {
		t.Errorf("Unexpected value when mapping null strings to struct values: expected '%s', got '%s'", expected, obj.Nested.Hah)
	}
}

func TestGofigNullStructValues(t *testing.T) {
	conf := createTestConfig(t)
	obj := &nullTest{}
	expected := &nullTest{
		Key:   "",
		Float: 0.0,
		Bool:  false,
	}

	if err := conf.Struct("nullsoft", obj); err != nil {
		t.Errorf("Struct() failed: %v", err)
		t.FailNow()
	}

	if !reflect.DeepEqual(obj, expected) {
		t.Error("Struct() failed to set default values when json contains null.")
	}
}

func TestGofigNullStruct(t *testing.T) {
	conf := createTestConfig(t)
	obj := &testObj{}
	expected := &testObj{}

	if err := conf.Struct("null_obj", obj); err != nil {
		t.Errorf("Struct() failed: %v", err)
		t.FailNow()
	}

	if !reflect.DeepEqual(obj, expected) {
		t.Error("Struct() failed to set empty struct when object is null.")
	}
}

func TestGofigStructArray(t *testing.T) {
	conf := createTestConfig(t)
	var pairs []testPair

	if err := conf.StructArray("obj_array", &pairs); err != nil {
		t.Errorf("StructArray() failed: %v", err)
		t.FailNow()
	}

	expected := []testPair{
		testPair{Key: "value", Int: 10},
		testPair{Key: "pair", Int: 26},
	}

	if len(pairs) != len(expected) {
		t.Errorf("Failed to map all pairs. expected %d, got %d.", len(expected), len(pairs))
	}

	if !reflect.DeepEqual(pairs, expected) {
		t.Errorf("Unexpected result from Array(): expected %v, got %v", expected, pairs)
	}
}

func TestGofigKeyNameTranslation(t *testing.T) {
	conf := createTestConfig(t)
	obj := &keyTest{}
	expected := &keyTest{
		ThisIsAKey: true,
		Key2:       "value!",
		Key3Time:   12.34,
		Key4here:   "checking in",
	}

	if err := conf.Struct("key_test", obj); err != nil {
		t.Error("Struct() failed: %v", err)
		t.FailNow()
	}

	if !reflect.DeepEqual(obj, expected) {
		t.Errorf("Struct() failed to map key names correctly. Expected %v, got %v", expected, obj)
	}
}
