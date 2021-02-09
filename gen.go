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

func (m *modelInfo) genFormulas(numFormulas, depth, numUnfold int) {

	// should never occur, to remove after test
	if m.twinModel != nil {
		if m.modelType != col {
			log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), Found a corresponding COL model, unfolding needed")
			return
		}
	}

	var error error
	m.getids()
	if m.twinModel != nil {
		m.twinModel.getids()
		error = m.twinModel.mapids()
	}

	if error != nil {
		log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), Warning: will not unfold formulas: impossible mapping")
		numUnfold = 0
	}

	// CTLFireability
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), generating ", numFormulas, " CTLFireability formulas")
	m.genericGenerationAndWriting(numFormulas, depth, numUnfold, genCTLFireabilityFormula, CTLFireabilityFileName)

	// CTLCardinality
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), generating ", numFormulas, " CTLCardinality formulas")
	m.genericGenerationAndWriting(numFormulas, depth, numUnfold, genCTLCardinalityFormula, CTLCardinalityFileName)

}

func (m *modelInfo) genericGenerationAndWriting(numFormulas, depth, numUnfold int, generation func(int, modelInfo) formula, outFileName string) {

	// gen numFormulas formulas
	formulas := m.genericGeneration(numFormulas, depth, generation)

	// write to file
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), writting formulas")
	m.writeFormulas(formulas, outFileName, true)

	if m.twinModel == nil {
		return
	}

	// If there is a corresponding PT model
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), Found a corresponding PT model, switching to it")
	m = m.twinModel

	// unfolding numUnfold formulas
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), unfolding ", numUnfold, " formulas")
	for i := 0; i < numUnfold; i++ {
		formulas[i] = m.unfolding(formulas[i])
	}

	// generating numFormulas - numUnfold formulas
	newFormulas := m.genericGeneration(numFormulas-numUnfold, depth, generation)
	for i := numUnfold; i < numFormulas; i++ {
		formulas[i] = newFormulas[i-numUnfold]
	}

	// write to file
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), writting formulas")
	m.writeFormulas(formulas, outFileName, true)
}

func (m *modelInfo) genericGeneration(numFormulas, depth int, generation func(int, modelInfo) formula) []formula {
	numFound := 0
	filterRounds := 0
	formulas := make([]formula, numFormulas)
	for numFound < numFormulas && filterRounds < globalMaxFilterTries {
		// gen numFormulas formulas
		log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), generating formulas")
		tmpFormulas := make([]formula, globalFilterSetSize)
		for i := 0; i < globalFilterSetSize; i++ {
			tmpFormulas[i] = generation(depth, *m)
		}

		// filter out easy formula
		log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), filtering formulas")
		toKeep := m.filter(tmpFormulas, numFormulas-numFound)
		for i := 0; i < len(toKeep) && numFound < numFormulas; i++ {
			formulas[numFound] = tmpFormulas[toKeep[i]]
			numFound++
		}

		filterRounds++

		// display info on generation
		log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), round ", filterRounds, ", kept ", len(toKeep), " formulas, ", numFormulas-numFound, " to go")
	}

	// if not enough formulas, complete with completely random ones
	if numFound < numFormulas {
		log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), found only ", numFound, " formulas, will add random ones to go up to ", numFormulas)
		for ; numFound < numFormulas; numFound++ {
			formulas[numFound] = generation(depth, *m)
		}
	}

	return formulas
}
