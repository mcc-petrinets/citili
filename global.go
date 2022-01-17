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

const (
	version                            string = "v2022"
	CTLFireabilityXMLFileName          string = "CTLFireability.xml"
	CTLCardinalityXMLFileName          string = "CTLCardinality.xml"
	ReachabilityFireabilityXMLFileName string = "ReachabilityFireability.xml"
	ReachabilityCardinalityXMLFileName string = "ReachabilityCardinality.xml"
	CTLFireabilityHRFileName           string = "CTLFireability.txt"
	CTLCardinalityHRFileName           string = "CTLCardinality.txt"
	ReachabilityFireabilityHRFileName  string = "ReachabilityFireability.txt"
	ReachabilityCardinalityHRFileName  string = "ReachabilityCardinality.txt"
	defaultMaxArity                    int    = 2        // max arity for operators
	defaultMaxAtomSize                 int    = 5        // max number of transitions/places in any atom
	defaultMaxIntegerConstant          int    = 100      // max constant to appear in integer comparisons in formulas
	defaultInputDir                    string = "INPUTS" // where to find the models
	defaultNumFormulas                 int    = 16       // number of formulas to generate
	defaultNumUnfold                   int    = 8        // number of formulas from COL models to unfold for generating formulas for PT models
	defaultFormulaDepth                int    = 2        // maximum depth of generated formulas
	defaultMaxFilterTries              int    = 3        // maximum number of call to SMC per model
	defaultFilterSetSize               int    = 16       // number of formula to generate for one round of SMC filtering
	defaultSMCPath                     string = "smc.py"
	defaultSMCTmpFileName              string = "tmp"
	defaultSMClogfile                  string = "smclog"
	defaultSMCMaxStates                int    = 2000
	defaultNumProc                     int    = 1 // number of cores to use for generating formulas
)

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
