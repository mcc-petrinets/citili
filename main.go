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
	"flag"
	"log"
)

func main() {

	inputDirPtr := flag.String("inputs", defaultInputDir, "directory where the models can be found")
	numFormulas := flag.Int("numformulas", defaultNumFormulas, "number of formulas to generate")
	formulaDepth := flag.Int("depth", defaultFormulaDepth, "max depth of the formulas to generate")
	numUnfold := flag.Int("numunfold", defaultNumUnfold, "number of formulas to unfold from COL to PT when possible")

	flag.Parse()

	log.Print(
		"Working with:\n",
		"\t", "models directory: ", *inputDirPtr, "\n",
		"\t", "number of formulas: ", *numFormulas, "\n",
		"\t", "number of unfolded formulas: ", *numUnfold, "\n",
		"\t", "formula depth: ", *formulaDepth, "\n",
		"\t", "maximum arity: ", globalMaxArity, "\n",
		"\t", "maximum atom size: ", globalMaxAtomSize, "\n",
		"\t", "maximum integer constant: ", globalMaxIntegerConstant, "\n",
	)

	models := listModels(*inputDirPtr)

	for pos, m := range models {
		log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), generating formulas")
		if m != nil {
			m.genFormulas(*numFormulas, *formulaDepth, *numUnfold)
			models[pos] = nil
		}
	}

}
