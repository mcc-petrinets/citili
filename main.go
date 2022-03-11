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
	"fmt"
	"log"
	"os"
	"runtime"
)

func main() {

	configFile := flag.String("conf", "config.json", "path to the configuration file")

	flag.Parse()
	getConfig(*configFile)

	log.Print(
		"Working with:\n",
		"\t", "cores: ", globalConfiguration.NumProc, "\n",
		"\t", "models directory: ", globalConfiguration.InputDir, "\n",
		"\t", "number of generated formulas per model: ", globalConfiguration.NumFormulas, "\n",
		"\t", "number of unfolded formulas per COL/PT cuple: ", globalConfiguration.NumUnfold, "\n",
		"Formulas characteristics:\n",
		"\t", "maximum depth: ", globalConfiguration.FormulaDepth, "\n",
		"\t", "maximum arity of operator: ", globalConfiguration.MaxArity, "\n",
		"\t", "maximum number of transitions per atom: ", globalConfiguration.MaxFireabilityAtomSize, "\n",
		"\t", "maximum number of places per atom: ", globalConfiguration.MaxCardinalityAtomSize, "\n",
		"\t", "maximum integer constant used in comparisons: ", globalConfiguration.MaxIntegerConstant, "\n",
		"Formulas filtering:\n",
		"\t", "number of filtering rounds per model: ", globalConfiguration.MaxFilterTries, "\n",
		"\t", "number of generated formulas at each filtering round: ", globalConfiguration.FilterSetSize, "\n",
		"\t", "tmp file location: ", globalConfiguration.SMCTmpFileName, "\n",
		"SMC configuration:\n",
		"\t", "path: ", globalConfiguration.SMCPath, "\n",
		"\t", "log file: ", globalConfiguration.SMClogfile, "\n",
		"\t", "maximum number of states to consider: ", globalConfiguration.SMCMaxStates, "\n",
	)

	// set the number of cores to use
	oldNumProc := runtime.GOMAXPROCS(globalConfiguration.NumProc)
	log.Print("Switching from ", oldNumProc, " cores (default) to ", globalConfiguration.NumProc, " cores")

	models := listModels(globalConfiguration.InputDir)

	initBooleanOperators() // for CTL only
	initStateOperators()   // for reachability only

	routineNum := 0
	doneChan := make(chan int, globalConfiguration.NumProc)
	for pos := range models {
		availableRoutineNum := routineNum
		if routineNum >= globalConfiguration.NumProc {
			availableRoutineNum = <-doneChan
		} else {
			routineNum++
		}
		go handleModel(
			pos, models, globalConfiguration.NumFormulas, globalConfiguration.FormulaDepth,
			globalConfiguration.NumUnfold, availableRoutineNum, doneChan)
	}

	for routineNum > 0 {
		<-doneChan
		routineNum--
	}

}

func handleModel(pos int, models []*modelInfo, numFormulas, formulaDepth, numUnfold int, routineNum int, doneChan chan int) {
	m := models[pos]

	logger := log.New(
		os.Stderr,
		fmt.Sprint("[goroutine-", routineNum, "] (", m.modelName, "-", m.modelInstance, ") "),
		log.LstdFlags,
	)

	logger.Print("Starting goroutine")

	logger.Print("generating formulas")
	m.genFormulas(numFormulas, formulaDepth, numUnfold, logger, routineNum)
	models[pos] = nil
	logger.Print("Ending goroutine")
	doneChan <- routineNum
}
