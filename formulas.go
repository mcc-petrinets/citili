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
	"fmt"
	"math/rand"
)

type formula struct {
	operator operator
	operand  []formula
}

// Determine maxDepth of a formula from previous maxDepth
func getDepth(oldDepth int, randomChoose bool) (newDepth int) {
	if oldDepth > 1 && randomChoose {
		newDepth = rand.Intn(oldDepth-1) + 1
	} else {
		newDepth = oldDepth - 1
	}
	return
}

// Generation of a boolean formula
func genBooleanFormula(maxDepth int) (f formula) {
	if maxDepth <= 1 {
		f = formula{operator: booleanOperators[0]}
		return f
	}

	// choose operator
	opNum := rand.Intn(len(booleanOperators))
	f = formula{operator: booleanOperators[opNum]}

	// generate subformulas
	arity := rand.Intn(f.operator.maxArity+1-f.operator.minArity) + f.operator.minArity
	f.operand = make([]formula, arity)
	if f.operator.isOverBooleans {
		for i := 0; i < arity; i++ {
			f.operand[i] = genBooleanFormula(getDepth(maxDepth, i > 0))
		}
	} else {
		for i := 0; i < arity; i++ {
			f.operand[i] = genPathFormula(getDepth(maxDepth, i > 0))
		}
	}

	return f
}

// Generation of a path formula
func genPathFormula(maxDepth int) (f formula) {
	// choose operator
	opNum := rand.Intn(len(pathOperators))
	f = formula{operator: pathOperators[opNum]}

	// generate subformulas
	arity := rand.Intn(f.operator.maxArity+1-f.operator.minArity) + f.operator.minArity
	f.operand = make([]formula, arity)
	for i := 0; i < arity; i++ {
		f.operand[i] = genBooleanFormula(getDepth(maxDepth, i > 0))
	}

	return f
}

// Generation of a generic CTL formula
func genCTLFormula(maxDepth int) formula {
	f := genBooleanFormula(maxDepth)
	for isInOtherCategory(f) {
		f = genBooleanFormula(maxDepth)
	}
	return f
}

// Checks if a CTL formula also belongs to another category
// LTL: A xxx (or just xxx) with xxx containing no A or E operator
// Reachability EF xxx or AG xxx with xxx containing no A or E operator
func isInOtherCategory(f formula) bool {
	if f.operator.name == "E" {
		if f.operand[0].operator.name != "F" {
			return false
		}
	}
	for _, operand := range f.operand {
		if containsCTLOperator(operand) {
			return false
		}
	}
	return true
}

// Checks if a formula contains at least one CTL operator
func containsCTLOperator(f formula) bool {
	if f.operator.name == "E" || f.operator.name == "A" {
		return true
	}
	for _, operand := range f.operand {
		if containsCTLOperator(operand) {
			return true
		}
	}
	return false
}

// Generation of a CTLFireability formula
func genCTLFireabilityFormula(maxDepth int, m modelInfo) (f formula) {
	f = genCTLFormula(maxDepth)
	f.CTLFireabilitySubstituteAtoms(m.transitions)
	return f
}

func (f *formula) CTLFireabilitySubstituteAtoms(transitions []string) {
	if f.operator == atom {
		*f = genCTLFireabilityAtom(transitions)
		return
	}
	for opNum := 0; opNum < len(f.operand); opNum++ {
		f.operand[opNum].CTLFireabilitySubstituteAtoms(transitions)
	}
}

func genCTLFireabilityAtom(transitions []string) (f formula) {
	f = formula{operator: isfireable}
	f.operand = make([]formula, 0)
	maxTransitions := len(transitions)
	if maxTransitions > globalMaxAtomSize {
		maxTransitions = globalMaxAtomSize
	}
	numTransitions := rand.Intn(maxTransitions) + 1
	rand.Shuffle(len(transitions), func(i, j int) { transitions[i], transitions[j] = transitions[j], transitions[i] })
	for i := 0; i < numTransitions; i++ {
		ff := formula{operator: operator{name: transitions[i]}}
		f.operand = append(f.operand, ff)
	}
	return f
}

// Generation of a CTLCardinality formula
func genCTLCardinalityFormula(maxDepth int, m modelInfo) (f formula) {
	f = genCTLFormula(maxDepth)
	f.CTLCardinalitySubstituteAtoms(m.places)
	return f
}

func (f *formula) CTLCardinalitySubstituteAtoms(places []string) {
	if f.operator == atom {
		*f = genCTLCardinalityAtom(places)
		return
	}
	for opNum := 0; opNum < len(f.operand); opNum++ {
		f.operand[opNum].CTLCardinalitySubstituteAtoms(places)
	}
}

func genCTLCardinalityAtom(places []string) (f formula) {
	f = formula{operator: leqOperator}
	f.operand = make([]formula, 2)
	tokencountChoice := rand.Intn(3) // 0 : tokencount on the left, 1: tokencount on the right, 2: tokencount on both sides
	switch tokencountChoice {
	case 0:
		f.operand[0] = genTokencount(places)
		f.operand[1] = genIntconstant()
	case 1:
		f.operand[0] = genIntconstant()
		f.operand[1] = genTokencount(places)
	case 2:
		f.operand[0] = genTokencount(places)
		f.operand[1] = genTokencount(places)
	}
	return f
}

func genTokencount(places []string) (f formula) {
	f = formula{operator: tokencount}
	f.operand = make([]formula, 0)
	maxPlaces := len(places)
	if maxPlaces > globalMaxAtomSize {
		maxPlaces = globalMaxAtomSize
	}
	numPlaces := rand.Intn(maxPlaces) + 1
	rand.Shuffle(len(places), func(i, j int) { places[i], places[j] = places[j], places[i] })
	for i := 0; i < numPlaces; i++ {
		ff := formula{operator: operator{name: places[i]}}
		f.operand = append(f.operand, ff)
	}
	return f
}

func genIntconstant() (f formula) {
	f = formula{operator: integerconstant}
	f.operand = make([]formula, 1)
	val := fmt.Sprint(rand.Intn(globalMaxIntegerConstant) + 1)
	f.operand[0] = formula{operator: operator{name: val}}
	return f
}
