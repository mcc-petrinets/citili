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

func unfolding(f formula, mapping map[string][]string) formula {

	switch f.operator.name {
	case "A", "E", "not", "and", "or", "G", "F", "X", "U", "leq", "integer-constant":
		for i := 0; i < len(f.operand); i++ {
			f.operand[i] = unfolding(f.operand[i], mapping)
		}
	case "is-fireable":
		unfOperand := make([]formula, 0)
		for _, t := range f.operand {
			for _, ut := range mapping[t.operator.name] {
				unfOperand = append(
					unfOperand,
					formula{operator: operator{name: ut}},
				)
			}
		}
		f.operand = unfOperand
	case "token-count":
		unfOperand := make([]formula, 0)
		for _, p := range f.operand {
			for _, up := range mapping[p.operator.name] {
				unfOperand = append(
					unfOperand,
					formula{operator: operator{name: up}},
				)
			}
		}
		f.operand = unfOperand
	}

	return f
}
