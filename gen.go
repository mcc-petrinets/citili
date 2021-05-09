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
	"log"
)

func (m *modelInfo) genFormulas(numFormulas, depth, numUnfold int, logger *log.Logger, routineNum int) {

	// should never occur, to remove after test
	if m.twinModel != nil {
		if m.modelType != col {
			logger.Print("Found a corresponding COL model, unfolding needed")
			return
		}
	}

	var error error
	m.getids(logger)
	if m.twinModel != nil {
		m.twinModel.getids(logger)
		error = m.twinModel.mapids(logger)
	}

	canUnfold := true
	if error != nil {
		logger.Print("Warning: will not unfold formulas: impossible mapping")
		canUnfold = false
	}

	// CTLFireability
	logger.Print("Generating ", numFormulas, " CTLFireability formulas")
	m.genericGenerationAndWriting(numFormulas, depth, numUnfold, canUnfold, genCTLFireabilityFormula, CTLFireabilityXMLFileName, CTLFireabilityHRFileName, "CTLFireability", logger, routineNum)

	// CTLCardinality
	logger.Print("Generating ", numFormulas, " CTLCardinality formulas")
	m.genericGenerationAndWriting(numFormulas, depth, numUnfold, canUnfold, genCTLCardinalityFormula, CTLCardinalityXMLFileName, CTLCardinalityHRFileName, "CTLCardinality", logger, routineNum)

}

func (m *modelInfo) genericGenerationAndWriting(numFormulas, depth, numUnfold int, canUnfold bool, generation func(int, modelInfo) formula, outXMLFileName, outHRFileName string, formulaType string, logger *log.Logger, routineNum int) {

	modelType := "COL"
	if m.modelType != col {
		modelType = "PT"
	}
	logger.Print("Working on ", modelType, " model")

	// gen numFormulas formulas
	formulas := m.genericGeneration(numFormulas, depth, canUnfold, generation, logger, routineNum)

	// write to file
	logger.Print("Writting formulas")
	m.writexmlFormulas(formulas, outXMLFileName, formulaType, true, logger)
	m.writehrFormulas(formulas, outHRFileName, formulaType, true, logger)

	if m.twinModel == nil {
		return
	}

	// If there is a corresponding PT model
	logger.Print("Found a corresponding PT model, switching to it")
	m = m.twinModel

	// unfolding numUnfold formulas if possible
	if !canUnfold {
		numUnfold = 0
	}
	logger.Print("Unfolding ", numUnfold, " formulas")
	for i := 0; i < numUnfold; i++ {
		formulas[i] = m.unfolding(formulas[i])
	}

	// generating numFormulas - numUnfold formulas
	newFormulas := m.genericGeneration(numFormulas-numUnfold, depth, canUnfold, generation, logger, routineNum)
	for i := numUnfold; i < numFormulas; i++ {
		formulas[i] = newFormulas[i-numUnfold]
	}

	// write to file
	logger.Print("Writting formulas")
	m.writexmlFormulas(formulas, outXMLFileName, formulaType, true, logger)
	m.writehrFormulas(formulas, outHRFileName, formulaType, true, logger)
}

func (m *modelInfo) genericGeneration(numFormulas, depth int, canUnfold bool, generation func(int, modelInfo) formula, logger *log.Logger, routineNum int) []formula {
	numFound := 0
	filterRounds := 0
	formulas := make([]formula, numFormulas)
	for numFound < numFormulas && filterRounds < globalMaxFilterTries {
		// gen numFormulas formulas
		logger.Print("Generating formulas")
		tmpFormulas := make([]formula, globalFilterSetSize)
		for i := 0; i < globalFilterSetSize; i++ {
			tmpFormulas[i] = generation(depth, *m)
		}

		// filter out easy formula
		logger.Print("Filtering formulas")
		toKeep := m.filter(tmpFormulas, numFormulas-numFound, canUnfold, logger, routineNum)
		for i := 0; i < len(toKeep) && numFound < numFormulas; i++ {
			formulas[numFound] = tmpFormulas[toKeep[i]]
			numFound++
		}

		filterRounds++

		// display info on generation
		logger.Print("Round ", filterRounds, ", kept ", len(toKeep), " formulas, ", numFormulas-numFound, " to go")
	}

	// if not enough formulas, complete with completely random ones
	if numFound < numFormulas {
		logger.Print("Found only ", numFound, " formulas, will add random ones to go up to ", numFormulas)
		for ; numFound < numFormulas; numFound++ {
			formulas[numFound] = generation(depth, *m)
		}
	}

	return formulas
}
