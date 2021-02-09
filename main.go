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

	globalMaxArity = *flag.Int("maxarity", defaultMaxArity, "maximum arity of operators in formulas")
	globalMaxAtomSize = *flag.Int("maxatomsize", defaultMaxAtomSize, "maximum number of transitions/places used in a single atom")
	globalMaxIntegerConstant = *flag.Int("maxintegerconstant", defaultMaxIntegerConstant, "maximum integer constant appearing in integer comparisons in formulas")
	globalMaxFilterTries = *flag.Int("maxfiltertries", defaultMaxFilterTries, "maximum number of times the formulas filter should be called on a given model, each call tries to filter among filtersetsize formulas")
	globalFilterSetSize = *flag.Int("filtersetsize", defaultFilterSetSize, "number of formulas to generate for each round of filtering")
	globalSMCPath = *flag.String("smcpath", defaultSMCPath, "path to SMC, the simple model checker used for filtering formulas")
	globalSMCTmpFileName = *flag.String("smctmpfile", defaultSMCTmpFileName, "path to the file that will be used to store formulas to be given to SMC")
	globalSMCMaxStates = *flag.Int("smcmaxstates", defaultSMCMaxStates, "number of states that SMC should explore before considering that a formula is not easy")
	globalSMClogfile = *flag.String("smclogfile", defaultSMClogfile, "path to the file where SMC log should be stored")

	flag.Parse()

	log.Print(
		"Working with:\n",
		"\t", "models directory: ", *inputDirPtr, "\n",
		"\t", "number of generated formulas per model: ", *numFormulas, "\n",
		"\t", "number of unfolded formulas per COL/PT cuple: ", *numUnfold, "\n",
		"Formulas characteristics:\n",
		"\t", "maximum depth: ", *formulaDepth, "\n",
		"\t", "maximum arity of operator: ", globalMaxArity, "\n",
		"\t", "maximum number of transitions/places per atom: ", globalMaxAtomSize, "\n",
		"\t", "maximum integer constant used in comparisons: ", globalMaxIntegerConstant, "\n",
		"Formulas filtering:\n",
		"\t", "number of filtering rounds per model: ", globalMaxFilterTries, "\n",
		"\t", "number of generated formulas at each filtering round: ", globalFilterSetSize, "\n",
		"\t", "tmp file location: ", globalSMCTmpFileName, "\n",
		"SMC configuration:\n",
		"\t", "path: ", globalSMCPath, "\n",
		"\t", "log file: ", globalSMClogfile, "\n",
		"\t", "maximum number of states to consider: ", globalSMCMaxStates, "\n",
	)

	models := listModels(*inputDirPtr)

	initOperators()

	for pos, m := range models {
		log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), generating formulas")
		if m != nil {
			m.genFormulas(*numFormulas, *formulaDepth, *numUnfold)
			models[pos] = nil
		}
	}

}
