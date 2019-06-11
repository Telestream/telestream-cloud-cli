package telestream

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"tcs-cli/cli"
)

func Test_structToProperties(t *testing.T) {

	type TestedStruct1 struct {
		Field1 string `json:"field_1"`
		Field2 string `json:"field_2"`
	}

	type TestedStruct2 struct {
		Field1 string
		Field2 string `json:""`
	}

	type TestedStruct3 struct {
		Field1 string `json:"field_1"`
		Field2 string `json:"-"`
	}

	type TestedStruct4 struct {
		Field1 string `json:"field_1,omitempty"`
		Field2 int16
	}

	type TestedStruct5 struct {
		Field1 bool
		Field2 int32
	}

	var testVector = []struct {
		structure  interface{}
		properties []string
	}{
		{&TestedStruct1{}, []string{"field_1", "field_2"}},
		{&TestedStruct2{}, []string{"Field1", "Field2"}},
		{&TestedStruct3{}, []string{"field_1", "Field2"}},
		{&TestedStruct4{}, []string{"field_1"}},
		{&TestedStruct5{}, []string{"Field1", "Field2"}},
	}

	for _, testEl := range testVector {

		properties := structToProperties(testEl.structure)

		if !reflect.DeepEqual(properties, testEl.properties) {

			t.Errorf("%v vs %v", properties, testEl.properties)
		}
	}
}

func Test_propertiesToStruct(t *testing.T) {

	type TestedStruct1 struct {
		Field1 string `json:"field_1,omitempty"`
		Field2 string `json:"field_2"`
	}

	type TestedStruct2 struct {
		Field1 string `json:"field_1"`
		Field2 string
	}

	type TestedStruct3 struct {
		Field1 bool
		Field2 int32
	}

	field1Value := "field_1_value"
	field2Value := "field_2_value"

	fieldBool := "true"
	fieldint32 := "32"

	var testVector = []struct {
		structure    interface{}
		flagMap      cli.FlagMap
		outputStruct interface{}
	}{
		{&TestedStruct1{}, cli.FlagMap{"field_1": cli.FlagProperties{Value: &field1Value, IsRequired: true},
			"field_2": cli.FlagProperties{Value: &field2Value, IsRequired: true}},
			&TestedStruct1{"field_1_value", "field_2_value"}},
		{&TestedStruct2{}, cli.FlagMap{"field_1": cli.FlagProperties{Value: &field1Value, IsRequired: true},
			"Field2": cli.FlagProperties{Value: &field2Value, IsRequired: true}},
			&TestedStruct2{"field_1_value", "field_2_value"}},
		{&TestedStruct3{}, cli.FlagMap{"Field1": cli.FlagProperties{Value: &fieldBool, IsRequired: true},
			"Field2": cli.FlagProperties{Value: &fieldint32, IsRequired: true}},
			&TestedStruct3{true, 32}},
	}

	for _, testEl := range testVector {

		propertiesToStruct(testEl.structure, testEl.flagMap)

		if !reflect.DeepEqual(testEl.structure, testEl.outputStruct) {

			t.Errorf("%v vs %v", testEl.structure, testEl.outputStruct)
		}
	}
}

func Test_PageOpt(t *testing.T) {

	flags := map[string]bool{}
	addPageOpt(flags)
	assert.True(t, reflect.DeepEqual(flags, map[string]bool{"page": false, "per_page": false}))

	emptyString := ""
	one := "1"

	argsMap := cli.FlagMap{"page": {Value: &one, IsRequired: false}, "per_page": {Value: &one, IsRequired: false},
		"some_flag": {Value: &emptyString, IsRequired: false}}

	opts, err := getPageOpt(&argsMap)
	assert.Equal(t, err, nil)
	assert.True(t, reflect.DeepEqual(opts, map[string]interface{}{"page": int32(1), "perPage": int32(1)}))
	fmt.Print("map: ", argsMap)
	assert.True(t, reflect.DeepEqual(argsMap, cli.FlagMap{"some_flag": cli.FlagProperties{Value: &emptyString, IsRequired: false}}))

	argsMap["page"] = cli.FlagProperties{Value: &emptyString, IsRequired: false}
	opts, err = getPageOpt(&argsMap)
	assert.Equal(t, err, nil)
	assert.True(t, reflect.DeepEqual(opts, map[string]interface{}{}))
	assert.True(t, reflect.DeepEqual(argsMap, cli.FlagMap{"some_flag": cli.FlagProperties{Value: &emptyString, IsRequired: false}}))
}
