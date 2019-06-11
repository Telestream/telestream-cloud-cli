package telestream

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

var ServiceToStdOut = NewServiceToStdOut(table.StyleLight)

type serviceToStdOut struct {
	tableStyle table.Style
}

func NewServiceToStdOut(style table.Style) *serviceToStdOut {

	printer := new(serviceToStdOut)
	printer.tableStyle = style

	return printer
}

func (printer *serviceToStdOut) printError(fName string, err error) {

	if err != nil {
		fmt.Println(fName + err.Error())
	} else {
		fmt.Println(fName)
	}
}

func (printer *serviceToStdOut) printStructContent(j interface{}) {

	fmt.Println()
	e := reflect.ValueOf(j).Elem()

	for i := 0; i < e.NumField(); i++ {

		varField := e.Type().Field(i)
		varName := varField.Name
		varValue := e.Field(i).Interface()

		if jsonTag := varField.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {

				varName = jsonTag[:commaIdx]
			} else {

				varName = jsonTag
			}
		}

		fmt.Printf("%v: %v\n", varName, varValue)
	}
	fmt.Println()
}

func (printer *serviceToStdOut) printInfo(info string) {

	fmt.Println(info)
}

func (printer *serviceToStdOut) printTable(colNames []interface{}, rows [][]interface{}) {

	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(os.Stdout)
	tableWriter.SetStyle(printer.tableStyle)

	tableWriter.AppendHeader(colNames)

	for _, el := range rows {
		tableWriter.AppendRow(el)
	}

	tableWriter.Render()
}
