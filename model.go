/*
Citili, a program for generating CTL formulas for the model checking contest
Copyright (C) 2020  Loïg Jezequel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see https://www.gnu.org/licenses/.
*/

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/loig/pinimili/pnml"
)

type modelType int

const (
	col modelType = iota
	pt
)

type modelInfo struct {
	filePath      string
	directory     string
	modelName     string
	modelType     modelType
	modelInstance string
	twinModel     *modelInfo
	places        []string // ids of places
	transitions   []string // ids of transitions
}

func listModels(inputDir string) []*modelInfo {
	var notDir int
	var noModel int
	var wrongName int
	var duplicateModel int
	models := make([]*modelInfo, 0)
	modelsMap := make(map[string](map[string]*modelInfo))

	inputsInfo, error := ioutil.ReadDir(inputDir)
	if error != nil {
		log.Panic(error)
	}

	for _, fileInfo := range inputsInfo {
		// only directories have to be considered
		if !fileInfo.IsDir() {
			notDir++
			continue
		}

		// directory names have to describe the type of model they contain
		nameOk, error := regexp.MatchString(`\w+(-COL-)|(-PT-)\w+`, fileInfo.Name())
		if error != nil {
			log.Panic(error)
		}
		if !nameOk {
			wrongName++
			continue
		}

		// directories need to contain a file named model.pnml
		directory := filepath.Join(inputDir, fileInfo.Name())
		modelFilePath := filepath.Join(directory, "model.pnml")
		if _, error = os.Stat(modelFilePath); os.IsNotExist(error) {
			noModel++
			continue
		}

		// fill the modelInfo for the current model
		splitName := strings.Split(fileInfo.Name(), "-")
		model := modelInfo{
			filePath:      modelFilePath,
			directory:     directory,
			modelName:     splitName[0],
			modelInstance: splitName[2],
		}
		if splitName[1] == "COL" {
			model.modelType = col
		} else {
			model.modelType = pt
		}

		// check if the model is the col/pt counterpart of an existing pt/col model
		instances, nameExists := modelsMap[model.modelName]
		if nameExists {
			twinModel, instanceExists := modelsMap[model.modelName][model.modelInstance]
			if instanceExists {
				if twinModel.modelType == model.modelType {
					duplicateModel++
					continue
				}
				twinModel.twinModel = &model
				model.twinModel = twinModel
			} else {
				instances[model.modelInstance] = &model
			}
		} else {
			modelsMap[model.modelName] = make(map[string]*modelInfo)
			modelsMap[model.modelName][model.modelInstance] = &model
		}

		// better with model address ?
		models = append(models, &model)

		log.Print(fileInfo.Name())
	}

	log.Print(
		"Warning: ",
		notDir+wrongName+noModel+duplicateModel, " elements were ignored in ", inputDir,
		" (", notDir, " were not directories, ",
		wrongName, " had a non-recognized name, ",
		noModel, " contained no model.pnml file, ",
		duplicateModel, " were duplicates of other models)",
	)

	return models
}

func (m *modelInfo) getids() {
	if m.places == nil || m.transitions == nil {
		m.places, m.transitions = pnml.Getptids(m.filePath)
		log.Print(
			m.modelName, " (", m.modelInstance, ", ", m.modelType, "): ",
			len(m.places), " places and ",
			len(m.transitions), " transitions.",
		)
	}
}
