package telestream

import (
	"reflect"
	"strconv"
	"strings"

	"tcs-cli/cli"
)

type ServiceOutput interface {
	printError(fName string, err error)
	printStructContent(j interface{})
	printTable(colNames []interface{}, rows [][]interface{})
	printInfo(info string)
}

// Convert all structure field names to string slice
func structToProperties(j interface{}) []string {

	propertiesList := []string{}

	e := reflect.ValueOf(j).Elem()

	for i := 0; i < e.NumField(); i++ {

		varField := e.Type().Field(i)
		varName := varField.Name
		varType := e.Field(i).Type().String()

		if jsonTag := varField.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {

			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {

				varName = jsonTag[:commaIdx]
			} else {

				varName = jsonTag
			}
		}

		if varType == "string" || varType == "int32" || varType == "bool" {

			propertiesList = append(propertiesList, varName)
		}
	}

	return propertiesList
}

// set structure field basing on map -> k -> field name, *v.Value -> field value
func propertiesToStruct(j interface{}, argsMap cli.FlagMap) {

	e := reflect.ValueOf(j).Elem()

	for i := 0; i < e.NumField(); i++ {

		varField := e.Type().Field(i)
		varName := varField.Name
		varType := e.Field(i).Type().String()

		if jsonTag := varField.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {

			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {

				varName = jsonTag[:commaIdx]
			} else {

				varName = jsonTag
			}
		}

		if val, ok := argsMap[varName]; ok {

			switch varType {

			case "string":
				e.Field(i).SetString(*val.Value)

			case "bool":
				b, _ := strconv.ParseBool(*val.Value)
				e.Field(i).SetBool(b)

			case "int32":
				in, _ := strconv.ParseInt(*val.Value, 10, 32)
				e.Field(i).SetInt(in)

			}
		}
	}
}

func addPageOpt(flags map[string]bool) {

	flags["page"] = false
	flags["per_page"] = false
}

func getPageOpt(argsMap *cli.FlagMap) (map[string]interface{}, error) {

	resMap := map[string]interface{}{}
	pageOpt := map[string]string{"page": "page", "per_page": "perPage"}

	for key, val := range pageOpt {
		if flagVal, ok := (*argsMap)[key]; ok {

			delete(*argsMap, key)
			if *flagVal.Value != "" {
				i, err := strconv.ParseInt(*flagVal.Value, 10, 32)
				if err != nil {
					return resMap, err
				}
				resMap[val] = int32(i)
			}
		}
	}

	// set page to 1 if parameter per_page passed  and page not
	if _, perPageOk := resMap["perPage"]; perPageOk {

		if _, pageOk := resMap["page"]; !pageOk {

			resMap["page"] = int32(1)
		}
	}

	return resMap, nil
}
