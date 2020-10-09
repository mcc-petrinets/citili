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

func (m modelInfo) unfolding(f formula) formula {

	var ff formula
	ff.operator = f.operator

	switch f.operator.name {
	case "A", "E", "not", "and", "or", "G", "F", "X", "U", "leq", "integer-constant":
		ff.operand = make([]formula, len(f.operand))
		for i := 0; i < len(f.operand); i++ {
			ff.operand[i] = m.unfolding(f.operand[i])
		}
	case "is-fireable":
		unfOperand := make([]formula, 0)
		for _, t := range f.operand {
			for _, ut := range m.transitionsMapping[t.operator.name] {
				unfOperand = append(
					unfOperand,
					formula{operator: operator{name: ut}},
				)
			}
		}
		ff.operand = unfOperand
	case "token-count":
		unfOperand := make([]formula, 0)
		for _, p := range f.operand {
			for _, up := range m.placesMapping[p.operator.name] {
				unfOperand = append(
					unfOperand,
					formula{operator: operator{name: up}},
				)
			}
		}
		ff.operand = unfOperand
	}

	return ff
}
