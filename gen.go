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

import "log"

func (m modelInfo) genFormulas(numFormulas, depth int) {

	// CTLFireability
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), generating CTLFireability formulas")
	m.genericGeneration(numFormulas, depth, genCTLFireabilityFormula, CTLFireabilityFileName)

	// CTLCardinality
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), generating CTLCardinality formulas")
	m.genericGeneration(numFormulas, depth, genCTLCardinalityFormula, CTLCardinalityFileName)

}

func (m modelInfo) genericGeneration(numFormulas, depth int, generation func(int, []string) formula, outFileName string) {
	// gen numFormulas formulas
	formulas := make([]formula, numFormulas)
	for i := 0; i < numFormulas; i++ {
		formulas[i] = generation(depth, m.transitions)
	}

	// filter easy formula // TODO

	// write to file
	m.writeFormulas(formulas, outFileName)
}
