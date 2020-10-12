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
	m.genericGeneration(numFormulas, depth, numUnfold, genCTLFireabilityFormula, CTLFireabilityFileName)

	// CTLCardinality
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), generating ", numFormulas, " CTLCardinality formulas")
	m.genericGeneration(numFormulas, depth, numUnfold, genCTLCardinalityFormula, CTLCardinalityFileName)

}

func (m *modelInfo) genericGeneration(numFormulas, depth, numUnfold int, generation func(int, modelInfo) formula, outFileName string) {
	// gen numFormulas formulas
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), generating formulas")
	formulas := make([]formula, numFormulas)
	for i := 0; i < numFormulas; i++ {
		formulas[i] = generation(depth, *m)
	}

	// filter easy formula // TODO
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), filtering formulas")

	// write to file
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), writting formulas")
	m.writeFormulas(formulas, outFileName)

	if m.twinModel == nil {
		return
	}

	// If there is a corresponding PT model
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), Found a corresponding PT model, switching to it")
	m = m.twinModel
	//m.getids()
	//m.mapids()

	// unfolding numUnfold formulas
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), unfolding ", numUnfold, " formulas")
	for i := 0; i < numUnfold; i++ {
		formulas[i] = m.unfolding(formulas[i])
	}

	// generating numFormulas - numUnfold formulas
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), generating ", numFormulas-numUnfold, " new formulas")
	for i := numUnfold; i < numFormulas; i++ {
		formulas[i] = generation(depth, *m)
	}

	// filter easy formula // TODO
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), filtering formulas")

	// write to file
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), writting formulas")
	m.writeFormulas(formulas, outFileName)
}
