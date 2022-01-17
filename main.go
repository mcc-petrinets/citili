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

	/*
		inputDirPtr := flag.String("inputs", defaultInputDir, "directory where the models can be found")
		numFormulas := flag.Int("numformulas", defaultNumFormulas, "number of formulas to generate")
		formulaDepth := flag.Int("formuladepth", defaultFormulaDepth, "max depth of the formulas to generate")
		numUnfold := flag.Int("numunfold", defaultNumUnfold, "number of formulas to unfold from COL to PT when possible")

		tmpMaxArity := flag.Int("maxarity", defaultMaxArity, "maximum arity of operators in formulas")
		tmpMaxAtomSize := flag.Int("maxatomsize", defaultMaxAtomSize, "maximum number of transitions/places used in a single atom")
		tmpMaxIntegerConstant := flag.Int("maxintegerconstant", defaultMaxIntegerConstant, "maximum integer constant appearing in integer comparisons in formulas")
		tmpMaxFilterTries := flag.Int("maxfiltertries", defaultMaxFilterTries, "maximum number of times the formulas filter should be called on a given model, each call tries to filter among filtersetsize formulas")
		tmpFilterSetSize := flag.Int("filtersetsize", defaultFilterSetSize, "number of formulas to generate for each round of filtering")
		tmpSMCPath := flag.String("smcpath", defaultSMCPath, "path to SMC, the simple model checker used for filtering formulas")
		tmpSMCTmpFileName := flag.String("smctmpfile", defaultSMCTmpFileName, "path to the file that will be used to store formulas to be given to SMC")
		tmpSMCMaxStates := flag.Int("smcmaxstates", defaultSMCMaxStates, "number of states that SMC should explore before considering that a formula is not easy")
		tmpSMClogfile := flag.String("smclogfile", defaultSMClogfile, "path to the file where SMC log should be stored")

		numProc := flag.Int("numproc", defaultNumProc, "number of cores available for generation")
	*/

	configFile := flag.String("conf", "config.json", "path to the configuration file")

	flag.Parse()
	getConfig(*configFile) // TODO : get the config file from command line

	log.Print(globalConfiguration)

	/*
		globalMaxArity = *tmpMaxArity
		globalMaxAtomSize = *tmpMaxAtomSize
		globalMaxIntegerConstant = *tmpMaxIntegerConstant
		globalMaxFilterTries = *tmpMaxFilterTries
		globalFilterSetSize = *tmpFilterSetSize
		globalSMCPath = *tmpSMCPath
		globalSMCTmpFileName = *tmpSMCTmpFileName
		globalSMCMaxStates = *tmpSMCMaxStates
		globalSMClogfile = *tmpSMClogfile
	*/

	log.Print(
		"Working with:\n",
		"\t", "cores: ", globalConfiguration.NumProc, "\n",
		"\t", "models directory: ", globalConfiguration.InputDir, "\n",
		"\t", "number of generated formulas per model: ", globalConfiguration.NumFormulas, "\n",
		"\t", "number of unfolded formulas per COL/PT cuple: ", globalConfiguration.NumUnfold, "\n",
		"Formulas characteristics:\n",
		"\t", "maximum depth: ", globalConfiguration.FormulaDepth, "\n",
		"\t", "maximum arity of operator: ", globalConfiguration.MaxArity, "\n",
		"\t", "maximum number of transitions/places per atom: ", globalConfiguration.MaxAtomSize, "\n",
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
