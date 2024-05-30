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

//var randomGenerator *rand.Rand

const (
	year                               string = "2024"
	version                            string = "v2024"
	CTLFireabilityXMLFileName          string = "CTLFireability.xml"
	CTLCardinalityXMLFileName          string = "CTLCardinality.xml"
	ReachabilityFireabilityXMLFileName string = "ReachabilityFireability.xml"
	ReachabilityCardinalityXMLFileName string = "ReachabilityCardinality.xml"
	CTLFireabilityHRFileName           string = "CTLFireability.txt"
	CTLCardinalityHRFileName           string = "CTLCardinality.txt"
	ReachabilityFireabilityHRFileName  string = "ReachabilityFireability.txt"
	ReachabilityCardinalityHRFileName  string = "ReachabilityCardinality.txt"
)

var defaultConfiguration config = config{
	MaxArity:               2,        // max arity for operators
	MaxFireabilityAtomSize: 1,        // max number of transitions in any atom
	MaxCardinalityAtomSize: 1,        // max number of places in any atom
	MinIntegerConstant:     0,        // min constant to appear in integere comparisions in formulas
	MaxIntegerConstant:     100,      // max constant to appear in integer comparisons in formulas
	InputDir:               "INPUTS", // where to find the models
	NumFormulas:            16,       // number of formulas to generate
	NumUnfold:              8,        // number of formulas from COL models to unfold for generating formulas for PT models
	FormulaDepth:           2,        // maximum depth of generated formulas
	MaxFilterTries:         3,        // maximum number of call to SMC per model
	FilterSetSize:          16,       // number of formula to generate for one round of SMC filtering
	SMCPath:                "smc.py",
	SMCTmpFileName:         "tmp",
	SMClogfile:             "smclog",
	SMCMaxStates:           2000,
	NumProc:                1, // number of cores to use for generating formulas
}

/*
var (
	globalMaxArity           int
	globalMaxAtomSize        int
	globalMaxIntegerConstant int
	globalMaxFilterTries     int
	globalFilterSetSize      int
	globalSMCPath            string
	globalSMCTmpFileName     string
	globalSMCMaxStates       int
	globalSMClogfile         string
)
*/
