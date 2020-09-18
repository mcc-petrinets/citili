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
	globalMaxArity           int = 2
	globalMaxAtomSize        int = 5
	globalMaxIntegerConstant int = 100
)

type operator struct {
	name           string
	minArity       int
	maxArity       int
	isOverBooleans bool
}

var atom operator = operator{"atom", 0, 0, false}
var isfireable operator = operator{name: "is-fireable"}
var tokencount operator = operator{name: "token-count"}
var leqOperator operator = operator{"leq", 2, 2, false}
var integerconstant operator = operator{"integer-constant", 1, 1, false}

var booleanOperators []operator = []operator{
	atom,
	operator{"A", 1, 1, false},
	operator{"E", 1, 1, false},
	operator{"not", 1, 1, true},
	operator{"and", 2, globalMaxArity, true},
	operator{"or", 2, globalMaxArity, true},
}

var pathOperators []operator = []operator{
	operator{"G", 1, 1, false},
	operator{"F", 1, 1, false},
	operator{"X", 1, 1, false},
	operator{"U", 2, 2, false},
}
