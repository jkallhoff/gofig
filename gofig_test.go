package gofig

import (
	"encoding/json"
	"regexp"
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

func testJson() string {
	return `{
        "string":"value", "string2":null, "int":34, "float":23.34, "bool":true, "array": [1,2,3,"Test"],
        "obj": {"bool": false,"float": 1.89,"nested": {"wow": "really?", "hah":null}},
        "key_test": {"this_is_a_key":true, "key_2": "value!", "key_3_time": 12.34, "key4here": "checking in"},
        "obj_array": [{"key": "value","int": 10}, {"key": "pair","int": 26}]
    }`
}

func createTestConfig(t *testing.T) (conf *Config) {
	conf = &Config{}
	if rx, e := regexp.Compile(`\_.{1}`); e != nil {
		t.Error("Failed to compile key name regex: %v", e)
		t.FailNow()
	} else {
		conf.keyRx = rx
		if e := json.Unmarshal([]byte(testJson()), &conf.values); e != nil {
			t.Errorf("Failed to parse testing json: %v", e)
			t.FailNow()
		}
	}
	return conf
}

func TestGofigInt(t *testing.T) {
	conf := createTestConfig(t)
	if i, e := conf.Int("int"); e != nil {
		t.Errorf("Int() failed to fetch integer value: %v", e)
	} else if i != 34 {
		t.Errorf("Unexpected value returned from Int(): expected 34, got %d", i)
	}
}

func TestGofigFloat(t *testing.T) {
	conf := createTestConfig(t)
	if f, e := conf.Float("float"); e != nil {
		t.Errorf("Float() failed to fetch float value: %v", e)
	} else if f != 23.34 {
		t.Errorf("Unexpected value returned from Float(): expected 12.34, got %f", f)
	}
}

func TestGofigBool(t *testing.T) {
	conf := createTestConfig(t)
	if b, e := conf.Bool("bool"); e != nil {
		t.Errorf("Bool() failed to fetch boolean value: %v", e)
	} else if !b {
		t.Errorf("Unexpected value returned from Bool(): expected true, got %v", b)
	}
}

func TestGofigArray(t *testing.T) {
	conf := createTestConfig(t)
	if a, e := conf.Array("array"); e != nil {
		t.Errorf("Array() failed to fetch array value: %v", e)
	} else {
		for i := 0; i < 3; i++ {
			if int(a[i].(float64)) != (i + 1) {
				t.Errorf("Unexpected value at element %d, expected %d, got %v.", i, i+1, a[i])
			}
		}

		if a[3].(string) != "Test" {
			t.Errorf("Unexpected value at element 3, expected 'Test', got %v", a[3])
		}
	}
}

func TestGofigValidString(t *testing.T) {
	conf := createTestConfig(t)
	if s, e := conf.Str("string"); e != nil {
		t.Errorf("Str() failed to fetch string value: %v", e)
	} else if s != "value" {
		t.Errorf("Unexpected value returned from Str(): expected 'value', got %s", s)
	}
}

func TestGofigNullString(t *testing.T) {
	conf := createTestConfig(t)
	if s, e := conf.Str("string2"); e != nil {
		t.Errorf("Str() failed to fetch string value: %v", e)
	} else if s != "" {
		t.Errorf("Unexpected value returned from Str(): expected '', got %s", s)
	}
}

func TestGofigNullStringInStruct(t *testing.T) {
	conf := createTestConfig(t)
	obj := &testObj{Nested: &testSubObj{Hah: "should get cleared"}}
	conf.Struct("obj", obj)

	if obj.Nested.Hah != "" {
		t.Errorf("Unexpected value when mapping null strings to struct values.")
	}
}

func TestGofigStructArray(t *testing.T) {
	conf := createTestConfig(t)
	var pairs []testPair
	conf.StructArray("obj_array", &pairs)

	if len(pairs) != 2 {
		t.Errorf("Failed to map all pairs. expected 2, got %d.", len(pairs))
	}

	for _, p := range pairs {
		if p.Key == "" {
			t.Errorf("Pair.Key value should not be blank.")
		}

		if p.Int <= 0 {
			t.Errorf("Pair.Int value should be greater than 0.")
		}
	}
}

func TestGofigKeyNameTranslation(t *testing.T) {
	conf := createTestConfig(t)
	obj := &keyTest{}
	conf.Struct("key_test", obj)

	if obj.ThisIsAKey == false {
		t.Errorf("Key name translation failed (this_is_a_key).")
	}

	if obj.Key2 != "value!" {
		t.Errorf("Key name translation failed (key_2).")
	}

	if obj.Key3Time != 12.34 {
		t.Errorf("Key name translation failed (key_3_time).")
	}

	if obj.Key4here != "checking in" {
		t.Errorf("Key name translation failed (key4here).")
	}
}
