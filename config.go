package gofig

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Str returns the string value of a key, or an empty string if null.
func (c *Config) Str(key string) (string, error) {
	var v string
	err := c.getValue(key, &v)
	return v, err
}

// Int returns the integer value of the given key, truncating any floating point values.
func (c *Config) Int(key string) (int, error) {
	f, err := c.Float(key)
	return int(f), err
}

// Float returns the float64 value of the given key.
func (c *Config) Float(key string) (float64, error) {
	var f float64
	err := c.getValue(key, &f)
	return f, err
}

// Bool returns the boolean value of the given key.
func (c *Config) Bool(key string) (bool, error) {
	var b bool
	err := c.getValue(key, &b)
	return b, err
}

// Array returns an array of interface{} for the given key, as JSON arrays can contain any datatype.
func (c *Config) Array(key string) ([]interface{}, error) {
	var a []interface{}
	err := c.getValue(key, &a)
	return a, err
}

// StructArray populates a typed array for the given key. Pass the address of an array and
// if there is no error, the array will be populated from the value.
func (c *Config) StructArray(key string, array interface{}) error {
	aType := reflect.TypeOf(array).Elem()
	if aType.Kind() != reflect.Array && aType.Kind() != reflect.Slice {
		return errors.New("array parameter is not an array or slice.")
	}
	if a, err := c.Array(key); err != nil {
		return err
	} else {
		aVal := reflect.ValueOf(array).Elem()
		for _, v := range a {
			obj := reflect.New(aType.Elem()).Interface()
			c.mapToStruct(v.(map[string]interface{}), obj)
			r := reflect.Append(aVal, reflect.ValueOf(obj).Elem())
			aVal.Set(r)
		}
		return nil
	}
}

// Map returns a map[string]interface{} for the given key.
func (c *Config) Map(key string) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := c.getValue(key, &m)
	return m, err
}

// Struct inflates an interface{} based on the config value of the given key.
// Pass the address of the interface{} and any exported fields found in the object
// will be populated from the value.
func (c *Config) Struct(key string, s interface{}) error {
	if m, err := c.Map(key); err != nil {
		return err
	} else {
		c.mapToStruct(m, s)
		return nil
	}
}

func (c *Config) getValue(key string, result interface{}) (err error) {
	if v, ok := c.values[key]; ok {
		rType := reflect.TypeOf(result)
		rVal := reflect.ValueOf(result).Elem()

		if v == nil {
			rVal.Set(reflect.New(rType.Elem()).Elem())
		} else {
			vType := reflect.PtrTo(reflect.TypeOf(v))
			if vType.AssignableTo(rType) {
				vVal := reflect.ValueOf(v)
				rVal.Set(vVal)
			} else {
				err = errors.New("Cannot assign config value to expected type.")
			}
		}

	} else {
		err = errors.New(fmt.Sprintf("Key '%s' does not exist.", key))
	}
	return
}

func (c *Config) mapToStruct(m map[string]interface{}, s interface{}) {
	pVal := reflect.ValueOf(s).Elem()
	for k, v := range m {
		vType := reflect.TypeOf(v)
		vVal := reflect.ValueOf(v)
		f := pVal.FieldByName(strings.Title(k))

		if f.Kind() != reflect.Invalid && f.Type().AssignableTo(vType) {
			f.Set(vVal)
		} else if vType.Kind() == reflect.Map && vType.Key().Kind() == reflect.String && vType.Elem().Kind() == reflect.Interface {
			var obj interface{}
			var objVal reflect.Value

			if f.Kind() == reflect.Ptr {
				obj = reflect.New(f.Type().Elem()).Interface()
				objVal = reflect.ValueOf(obj)
			} else {
				obj = reflect.New(f.Type()).Interface()
				objVal = reflect.ValueOf(obj).Elem()
			}

			c.mapToStruct(v.(map[string]interface{}), obj)
			f.Set(objVal)
		}
	}
}
