package telestream

import (
	"errors"
	"testing"
)

func Test_serviceToStdOut_prints(t *testing.T) {

	ServiceToStdOut.printError("name", errors.New("some error"))
	ServiceToStdOut.printError("", nil)

	ServiceToStdOut.printInfo("info")
	ServiceToStdOut.printInfo("")

	ServiceToStdOut.printTable(nil, nil)
	ServiceToStdOut.printTable(nil, [][]interface{}{{"val2", "val2"}})
	ServiceToStdOut.printTable([]interface{}{"col1", "col2"}, [][]interface{}{{"val2", "val2"}})
	ServiceToStdOut.printTable(nil, [][]interface{}{{"val2", "val2"}})
	ServiceToStdOut.printTable([]interface{}{"col1", "col2"}, [][]interface{}{{}})
}

func Test_serviceToStdOut_printStructContent(t *testing.T) {

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
		structure interface{}
	}{
		{&TestedStruct1{}},
		{&TestedStruct2{}},
		{&TestedStruct3{}},
		{&TestedStruct4{}},
		{&TestedStruct5{}},
	}

	for _, testEl := range testVector {

		ServiceToStdOut.printStructContent(testEl.structure)
	}
}
