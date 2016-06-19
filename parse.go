package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/tealeg/xlsx"
)

var XlsxTable xlsxTable

type xlsxTable struct {
	MaxCol  int
	MaxRow  int
	Columns []column
}

var SecretColumns []column

type column struct {
	Name      string
	ShortName string `json:"ShortName,omitempty"`
	Params    params
	Data      []data
}

type params struct {
	Sort     bool   `json:"Sort,omitempty"`
	Hide     bool   `json:"Hide,omitempty"`
	Reg      bool   `json:"Reg,omitempty"`
	Enum     bool   `json:"Enum,omitempty"`
	Filter   bool   `json:"Filter,omitempty"`
	Priority string `json:"Priority,omitempty"`
}

type data struct {
	Value string
	Link  string `json:"Link,omitempty"`
}

func (xt *xlsxTable) ParseXLSX(filename string) error {
	log.Info("Parsing xlsx file")

	xlsx, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}

	// todo данные о размере документа, не всегда соответствуют реальным данным
	xt.MaxCol = xlsx.Sheet[Config.Table.SheetName].MaxCol - 1
	xt.MaxRow = xlsx.Sheet[Config.Table.SheetName].MaxRow - 1

	xt.getColumns(xlsx)

	return nil
}

func (xt *xlsxTable) getColumns(xlsx *xlsx.File) {

	var Columns, SecretColumns []column

	for col := 0; col <= xt.MaxCol; col++ {
		name := xlsx.Sheet[Config.Table.SheetName].Cell(Config.Table.ColumnNames, col).Value
		if name == "" {
			xt.MaxCol = col
			break
		}

		shortName := xlsx.Sheet[Config.Table.SheetName].Cell(Config.Table.ColumnShortNames, col).Value

		Column := column{
			Name:   name,
			Params: xt.getOptions(xlsx, col),
			Data:   xt.getData(xlsx, col),
		}
		if shortName != "" {
			Column.ShortName = shortName
		}

		if !Column.Params.Reg {
			Columns = append(Columns, Column)
		} else {
			SecretColumns = append(SecretColumns, Column)
		}
	}

	xt.Columns = Columns
}

func (xt *xlsxTable) getOptions(xlsx *xlsx.File, col int) (param params) {

	value := xlsx.Sheet[Config.Table.SheetName].Cell(Config.Table.ColumnOptions, col).Value
	o := strings.Split(value, "|")

	for _, property := range o {
		switch property {
		case "hide":
			param.Hide = true
		case "enum":
			param.Enum = true
		case "filter":
			param.Filter = true
		case "reg":
			param.Reg = true
		case "sort":
			param.Sort = true
		case "always":
			param.Priority = "critical"
		case "p1":
			param.Priority = "1"
		case "p2":
			param.Priority = "2"
		case "p3":
			param.Priority = "3"
		case "p4":
			param.Priority = "4"
		case "p5":
			param.Priority = "5"
		case "p6":
			param.Priority = "6"
		}
	}
	return
}

func (xt *xlsxTable) getData(xlsx *xlsx.File, col int) (Data []data) {

	for row := Config.Table.ColumnData; row <= xt.MaxRow; row++ {
		value := xlsx.Sheet[Config.Table.SheetName].Cell(row, col).Value

		if col == 0 && value == "" {
			xt.MaxRow = row - 1
			break
		}

		formula := xlsx.Sheet[Config.Table.SheetName].Cell(row, col).Formula()

		var link string
		if strings.Contains(formula, "HYPERLINK") {
			hl := strings.Split(formula, "\"")
			link = hl[1]
		}
		var d data
		d.Value = value
		if link != "" {
			d.Link = link
		}

		Data = append(Data, d)
	}
	return
}

func (xt *xlsxTable) Save(filename string) error {
	log.Info("Save data json file")

	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fo.Close()

	e := json.NewEncoder(fo)
	if err = e.Encode(xt.Columns); err != nil {
		return err
	}
	return nil
}

func (xt *xlsxTable) Open(filename string) error {
	log.Info("Read data json file in struct")

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(f, &xt.Columns); err != nil {
		return err
	}
	return nil
}
