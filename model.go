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
	"errors"
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
	filePath                 string
	directory                string
	modelName                string
	modelType                modelType
	modelInstance            string
	modelInstanceSeparators  int // HACK: for dealing with models with - in the instance, MCC2021 surprise
	twinModel                *modelInfo
	pnml                     *pnml.Pnml
	places                   []string            // ids of places to use for generation
	unmappedPlaces           []string            // ids of places that will not be used for generation
	transitions              []string            // ids of transitions to use for generation
	unmappedTransitions      []string            // ids of transitions that will not be used for generation
	placesMapping            map[string][]string // mapping of ids of places to ids of the twin model
	transitionsMapping       map[string][]string // mapping of ids of transitions
	maxConstantInMarking     int
	maxConstantInTransitions int
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
		instanceName := strings.Join(splitName[2:], "-")
		separators := len(splitName[2:]) - 1
		model := modelInfo{
			filePath:                modelFilePath,
			directory:               directory,
			modelName:               splitName[0],
			modelInstance:           instanceName,
			modelInstanceSeparators: separators,
		}
		if splitName[1] == "COL" {
			model.modelType = col
		} else {
			model.modelType = pt
		}

		// check if the model is the col/pt counterpart of an existing pt/col model
		_, nameExists := modelsMap[model.modelName]
		if nameExists {
			twinModel, instanceExists := modelsMap[model.modelName][model.modelInstance]
			if instanceExists {
				if twinModel.modelType == model.modelType {
					duplicateModel++
					continue
				}
				twinModel.twinModel = &model
				model.twinModel = twinModel
				if model.modelType == col {
					modelsMap[model.modelName][model.modelInstance] = &model
				}
			} else {
				modelsMap[model.modelName][model.modelInstance] = &model
			}
		} else {
			modelsMap[model.modelName] = make(map[string]*modelInfo)
			modelsMap[model.modelName][model.modelInstance] = &model
		}

		//models = append(models, &model)
		log.Print(fileInfo.Name())
	}

	// all the COL models and all the PT models with no twin are added to the set of models
	for _, instances := range modelsMap {
		for _, modelPtr := range instances {
			models = append(models, modelPtr)
		}
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

func (m *modelInfo) getpnml(logger *log.Logger) {
	if m.pnml == nil {
		m.pnml = pnml.GetPnml(m.filePath, false)
		logger.Print(
			"Pnml parsed",
		)
	}
}

func (m *modelInfo) getids(logger *log.Logger) {
	if m.places == nil || m.transitions == nil {
		m.places, m.transitions = m.pnml.Getptids()
		logger.Print(
			len(m.places), " places and ",
			len(m.transitions), " transitions.",
		)
	}
}

func (m *modelInfo) getMaxConstants(logger *log.Logger) {
	m.maxConstantInMarking = -1
	if m.modelType == pt {
		m.maxConstantInMarking = int(m.pnml.GetMaxConstantInMarking())
	}
	logger.Print(
		"maximum constant appearing in marking: ", m.maxConstantInMarking,
	)
}

// checks if a place/transition of a PT model was unfolded from a place/transition
// of the twin COL model
func isUnfolding(ptNode, colNode string, colNodes []string, logger *log.Logger) bool {
	if strings.HasPrefix(ptNode, colNode) {
		logger.Print("Nodes mapping: ", ptNode, " could have been obtained from ", colNode)
		if ptNode == colNode {
			logger.Print("Nodes mapping: yes, it was (equality)")
			return true
		}
		for _, n := range colNodes {
			if ptNode == n {
				logger.Print("Nodes mapping: no, it was not, ", n, " exists in COL")
				return false
			}
		}
		logger.Print("Nodes mapping: yes, it was (prefix)")
		return true
	}
	return false
}

func (m *modelInfo) mapids(logger *log.Logger) error {
	// when this function is called, m should always be the PT model

	if m.placesMapping == nil || m.transitionsMapping == nil {

		m.placesMapping = make(map[string][]string)
		m.transitionsMapping = make(map[string][]string)

		checkPlaces := make([]bool, len(m.places))
		unmappedPlaces := make([]string, 0)
		mappedPlaces := make([]string, 0)
		for _, p := range m.twinModel.places {
			// will not work if a place of the COL model has an id which is a prefix of another place id of this model
			for i, pp := range m.places {
				if isUnfolding(pp, p, m.twinModel.places, logger) {
					m.placesMapping[p] = append(m.placesMapping[p], pp)
					checkPlaces[i] = true
				}
			}
			// check that p was unfolded into something
			if len(m.placesMapping[p]) == 0 {
				logger.Print(
					"Warning, colored model has a place not mapped to a PT place: ",
					p,
				)
				unmappedPlaces = append(unmappedPlaces, p)
			} else {
				mappedPlaces = append(mappedPlaces, p)
			}
		}
		// check that every PT place is the unfolding of something
		for i, v := range checkPlaces {
			if !v {
				logger.Print(
					"Warning, PT model has a place not unfolded from a COL place: ",
					m.places[i],
				)
			}
		}
		// check that the set of places of the COL net that were
		// unfolded into places of the PT net is not empty
		if len(mappedPlaces) == 0 {
			logger.Print(
				"Warning, colored model has an empty set of mapped places",
			)
			return errors.New("empty set of places")
		}
		m.twinModel.places = mappedPlaces
		m.twinModel.unmappedPlaces = unmappedPlaces

		checkTransitions := make([]bool, len(m.transitions))
		unmappedTransitions := make([]string, 0)
		mappedTransitions := make([]string, 0)
		for _, t := range m.twinModel.transitions {
			// will not work if a transition of the COL model has an id which is a prefix of another transition id of this model
			for i, tt := range m.transitions {
				if isUnfolding(tt, t, m.twinModel.transitions, logger) {
					m.transitionsMapping[t] = append(m.transitionsMapping[t], tt)
					checkTransitions[i] = true
				}
			}
			// check that t was unfolded into something
			if len(m.transitionsMapping[t]) == 0 {
				logger.Print(
					"Warning, colored model has a transition not mapped to a PT transition: ",
					t,
				)
				unmappedTransitions = append(unmappedTransitions, t)
			} else {
				mappedTransitions = append(mappedTransitions, t)
			}
		}
		// check that every PT transition is the unfolding of something
		for i, v := range checkTransitions {
			if !v {
				logger.Print(
					"Warning, PT model has a transition not unfolded from a COL transition: ",
					m.transitions[i],
				)
			}
		}
		// check that the set of transitions of the COL net that were
		// unfolded into transitions of the PT net is not empty
		if len(mappedTransitions) == 0 {
			logger.Print(
				"Warning, colored model has an empty set of mapped transitions",
			)
			return errors.New("empty set of transitions")
		}
		m.twinModel.transitions = mappedTransitions
		m.twinModel.unmappedTransitions = unmappedTransitions
	}

	return nil
}
