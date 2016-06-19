package main

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

var Config config

type config struct {
	FileId   string
	MimeType string
	Table    table
	Listen   string
	LogFile  bool
}

type table struct {
	SheetName        string
	ColumnOptions    int
	ColumnNames      int
	ColumnShortNames int
	ColumnData       int
}

func (c *config) Init() (err error) {
	f, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Error(err)
		return
	}
	if err = json.Unmarshal(f, &c); err != nil {
		log.Error(err)
		return
	}
	return nil
}
