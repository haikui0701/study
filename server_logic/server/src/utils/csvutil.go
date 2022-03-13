package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type CsvUtilMgr struct {
}

var csvUtilMgr *CsvUtilMgr = nil

func GetCsvUtilMgr() *CsvUtilMgr {
	if csvUtilMgr == nil {
		csvUtilMgr = new(CsvUtilMgr)
	}
	return csvUtilMgr
}

func (self *CsvUtilMgr) getFieldMap(config interface{}) map[string][]string {
	t := reflect.TypeOf(config)
	fieldMap := make(map[string][]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		data := make([]string, 2, 2)
		data[0] = field.Name
		data[1] = fmt.Sprintf("%v", field.Type)
		fieldMap[tag] = data
	}
	return fieldMap
}

func (self *CsvUtilMgr) readCsv(fileName string) [][]string {
	csvReadFile, err := os.Open(fileName)
	defer csvReadFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		return [][]string{}
	}

	r := csv.NewReader(bufio.NewReader(csvReadFile))
	var records [][]string
	var recordIndex int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		// remove bom head
		if recordIndex == 0 {
			for index := range record {
				if index == 0 && strings.Contains(record[index], "\ufeff") {
					record[index] = strings.Replace(record[index], "\ufeff", "", 1)
					break
				}
			}
			recordIndex += 1
		}
		records = append(records, record)
	}
	return records
}

func (self *CsvUtilMgr) getTagMap(fields []string) map[int]string {
	tagMap := make(map[int]string, len(fields))
	for index, v := range fields {
		tagMap[index] = v
	}

	return tagMap
}

func (self *CsvUtilMgr) getValueType(SlicePtr interface{}) reflect.Type {
	value := reflect.ValueOf(SlicePtr)
	if value.Kind() == reflect.Ptr {
		// 返回指针指向的值,这里是slice
		value = value.Elem()
	}

	// 返回当前值对应的类型
	outType := value.Type()
	// 返回这些类型包含的一个值
	outInnerType := outType.Elem()

	return outInnerType
}

// 指针只能取interface{}, 值得取Elem,否则会蹦
// tagMap: csv第一行
// filedMap: key: tag
func (self *CsvUtilMgr) genConfig(dataPtr interface{}, csvData [][]string, tagMap map[int]string,
	fieldMap map[string][]string, trimFlag map[int]bool, fileName string, keyTag string) (err error) {
	dataVal := reflect.Indirect(reflect.ValueOf(dataPtr))
	outInnerType := self.getValueType(dataPtr)

	//if fileName == "expspeedup" {
	//litter.Dump(tagMap)
	//litter.Dump(fieldMap)
	//}

	for r := 1; r < len(csvData); r++ {
		data := reflect.New(outInnerType.Elem())

		key := 0
		for c := 0; c < len(csvData[r]); c++ {
			tag := tagMap[c]
			if _, ok := trimFlag[c]; ok {
				tag = self.trimNumber(tag)
			}

			fieldInfo, ok := fieldMap[tag]
			if !ok {
				continue
			}

			if len(fieldInfo) != 2 {
				continue
			}
			fieldName := fieldInfo[0]
			filedType := fieldInfo[1]
			cellValue := csvData[r][c]
			switch filedType {
			case "int":
				v, err := strconv.Atoi(cellValue)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", index:", r)
					break
				}
				reflect.Indirect(data).FieldByName(fieldName).SetInt(int64(v))
				if tag == keyTag && key == 0 {
					key = v
					//fmt.Println("fileName:", fileName, ", tag:", tag)
				}
			case "int64":
				v, err := strconv.ParseInt(cellValue, 10, 64)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName)
					break
				}
				reflect.Indirect(data).FieldByName(fieldName).SetInt(v)
			case "string":
				reflect.Indirect(data).FieldByName(fieldName).SetString(cellValue)
			case "[]int":
				v, err := strconv.Atoi(cellValue)
				if err != nil {
					fv, err := strconv.ParseFloat(cellValue, 64)
					if err != nil {
						fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", row:", r, ", col:", c, ", ", strings.Join(csvData[r], ","))
					} else {
						c := reflect.Indirect(data).FieldByName(fieldName)
						newSlice := reflect.Append(c, reflect.ValueOf(int(fv*10)))
						reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
					}
					break
				} else {
					c := reflect.Indirect(data).FieldByName(fieldName)
					newSlice := reflect.Append(c, reflect.ValueOf(v))
					reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
				}
			case "[]int64":
				v, err := strconv.ParseInt(cellValue, 10, 0)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName)
					break
				}
				c := reflect.Indirect(data).FieldByName(fieldName)
				newSlice := reflect.Append(c, reflect.ValueOf(v))
				reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
			case "[]string":
				c := reflect.Indirect(data).FieldByName(fieldName)
				newSlice := reflect.Append(c, reflect.ValueOf(cellValue))
				reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
			case "float32":
				fv, err := strconv.ParseFloat(cellValue, 32)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", row:", r, ", col:", c, ", ", strings.Join(csvData[r], ","))
				} else {
					reflect.Indirect(data).FieldByName(fieldName).SetFloat(fv)
				}
			case "float64":
				fv, err := strconv.ParseFloat(cellValue, 64)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", row:", r, ", col:", c, ", ", strings.Join(csvData[r], ","))
				} else {
					reflect.Indirect(data).FieldByName(fieldName).SetFloat(fv)
				}
			case "[]float32":
				fv, err := strconv.ParseFloat(cellValue, 32)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", row:", r, ", col:", c, ", ", strings.Join(csvData[r], ","))
				} else {
					c := reflect.Indirect(data).FieldByName(fieldName)
					newSlice := reflect.Append(c, reflect.ValueOf(float32(fv)))
					reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
				}
			case "[]float64":
				fv, err := strconv.ParseFloat(cellValue, 64)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", row:", r, ", col:", c, ", ", strings.Join(csvData[r], ","))
				} else {
					c := reflect.Indirect(data).FieldByName(fieldName)
					newSlice := reflect.Append(c, reflect.ValueOf(fv))
					reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
				}
			}
		}

		kind := reflect.TypeOf(dataVal.Interface()).Kind()
		if kind == reflect.Slice {
			dataVal.Set(reflect.Append(dataVal, data))
		} else if kind == reflect.Map {
			dataVal.SetMapIndex(reflect.ValueOf(key), data)
		}
	}

	return nil
}

func (self *CsvUtilMgr) trimNumber(s string) string {
	subString := strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsNumber(r)
	})

	return subString

}

func (self *CsvUtilMgr) LoadCsv(fileName string, SlicePtr interface{}) {
	csvFile := "csv/" + fileName + ".csv"
	csvData := self.readCsv(csvFile)
	if len(csvData) <= 1 {
		fmt.Println("len(csvData) <= 1, filename:", fileName)
		os.Exit(1)
		return
	}

	//LogDebug("Read File:", fileName, len(csvData))
	err := self.ParseDataSimple(csvData, SlicePtr, fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (self *CsvUtilMgr) LoadEventsCsv(fileName string, SlicePtr interface{}) {
	csvFile := "csv/ChapterMap/" + fileName + ".csv"
	csvData := self.readCsv(csvFile)
	if len(csvData) <= 1 {
		fmt.Println("len(csvData/ChapterMap) <= 1, fileName:", fileName)
		os.Exit(1)
		return
	}

	err := self.ParseDataSimple(csvData, SlicePtr, fileName)
	if err != nil {
		return
	}
}

func (self *CsvUtilMgr) getFieldMapSimple(config interface{}, tagMap map[int]string) (map[string][]string, map[int]bool, string) {
	t := reflect.TypeOf(config)
	fieldMap := make(map[string][]string, t.NumField())
	trimFlag := make(map[int]bool)
	keyTag := ""
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		trim := field.Tag.Get("trim")

		if trim == "1" || trim == "" {
			tag = self.trimNumber(tag)
		}

		for key, v := range tagMap {
			if (trim == "" || trim == "1") && tag == self.trimNumber(v) {
				trimFlag[key] = true
			}
		}

		// 设置表数据结构
		for key, v := range tagMap {
			if (trim == "" || trim == "1") && tag == self.trimNumber(v) {
				tagMap[key] = tag
				break
			}
		}

		data := make([]string, 2, 2)
		data[0] = field.Name
		data[1] = fmt.Sprintf("%v", field.Type)
		fieldMap[tag] = data
		if i == 0 {
			keyTag = tag
		}
	}
	return fieldMap, trimFlag, keyTag
}

func (self *CsvUtilMgr) ParseDataSimple(csvData [][]string, dataPtr interface{}, fileName string) (err error) {
	outInnerType := self.getValueType(dataPtr)
	data := reflect.New(outInnerType.Elem())

	value := reflect.Indirect(data)
	tagMap := self.getTagMap(csvData[0])
	fieldMap, trimFlag, keyTag := self.getFieldMapSimple(value.Interface(), tagMap)
	err = self.genConfig(dataPtr, csvData, tagMap, fieldMap, trimFlag, fileName, keyTag)

	return
}
