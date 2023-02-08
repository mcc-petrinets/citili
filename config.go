package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type config struct {
	Seed                   int64
	MaxArity               int
	MaxFireabilityAtomSize int
	MaxCardinalityAtomSize int
	MinIntegerConstant     int
	MaxIntegerConstant     int
	InputDir               string
	NumFormulas            int
	NumUnfold              int
	FormulaDepth           int
	MaxFilterTries         int
	FilterSetSize          int
	SMCPath                string
	SMCTmpFileName         string
	SMClogfile             string
	SMCMaxStates           int
	NumProc                int
}

var globalConfiguration config

func getConfig(fileName string) {

	globalConfiguration = defaultConfiguration

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Error when opening config file: ", err)
	}

	err = json.Unmarshal(content, &globalConfiguration)
	if err != nil {
		log.Fatal("Error when reading config file: ", err)
	}

	// TODO forbidden values (ex: MaxArity = 0)
}
