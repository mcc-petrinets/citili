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
	"log"
	"os"
	"path/filepath"
)

const indent string = "   "

// print a set of formulas in a file for a given model
func (m modelInfo) writeFormulas(formulas []formula, fileName string) {

	f, error := os.Create(filepath.Join(m.directory, fileName))
	if error != nil {
		log.Panic(error)
	}

	f.WriteString(
		fmt.Sprint(
			"<?xml version\"1.0\"?>\n",
			"<property-set xmnls=\"http://mcc.lip6.fr/\">\n",
		))

	for i := 0; i < len(formulas); i++ {
		f.WriteString(formulas[i].xmlPrint(m, i))
	}

	f.WriteString(
		fmt.Sprint(
			"</property-set>\n",
		))

	f.Sync()

}

// output one formula as xml
func (f formula) xmlPrint(m modelInfo, num int) (xmlp string) {
	kind := "COL"
	if m.modelType != col {
		kind = "PT"
	}

	fid := fmt.Sprintf(
		"%s-%s-%s-%2.2d",
		m.modelName, kind, m.modelInstance, num,
	)

	xmlf := f.asxml(indent + indent + indent)

	xmlp = fmt.Sprint(
		indent, "<property>\n",
		indent, indent, "<id>", fid, "</id>\n",
		indent, indent, "<description>Automatically generated by Citili ", version, "</description>\n",
		indent, indent, "<formula>\n",
		xmlf,
		indent, indent, "</formula>\n",
		indent, "</property>\n",
	)
	return xmlp
}

func (f formula) asxml(currentIndent string) (xmlf string) {

	switch f.operator.name {
	case "A":
		xmlff := f.operand[0].asxml(currentIndent + indent)
		xmlf = fmt.Sprint(
			currentIndent, "<all-paths>\n",
			xmlff,
			currentIndent, "</all-paths>\n",
		)
	case "E":
		xmlff := f.operand[0].asxml(currentIndent + indent)
		xmlf = fmt.Sprint(
			currentIndent, "<exists-path>\n",
			xmlff,
			currentIndent, "</exists-path>\n",
		)
	case "not":
		xmlff := f.operand[0].asxml(currentIndent + indent)
		xmlf = fmt.Sprint(
			currentIndent, "<negation>\n",
			xmlff,
			currentIndent, "</negation>\n",
		)
	case "and":
		xmlff := f.operand[0].asxml(currentIndent + indent)
		for i := 1; i < len(f.operand); i++ {
			xmlff = xmlff + f.operand[i].asxml(currentIndent+indent)
		}
		xmlf = fmt.Sprint(
			currentIndent, "<conjunction>\n",
			xmlff,
			currentIndent, "</conjunction>\n",
		)
	case "or":
		xmlff := f.operand[0].asxml(currentIndent + indent)
		for i := 1; i < len(f.operand); i++ {
			xmlff = xmlff + f.operand[i].asxml(currentIndent+indent)
		}
		xmlf = fmt.Sprint(
			currentIndent, "<disjunction>\n",
			xmlff,
			currentIndent, "</disjunction>\n",
		)
	case "G":
		xmlff := f.operand[0].asxml(currentIndent + indent)
		xmlf = fmt.Sprint(
			currentIndent, "<globally>\n",
			xmlff,
			currentIndent, "</globally>\n",
		)
	case "F":
		xmlff := f.operand[0].asxml(currentIndent + indent)
		xmlf = fmt.Sprint(
			currentIndent, "<finally>\n",
			xmlff,
			currentIndent, "</finally>\n",
		)
	case "X":
		xmlff := f.operand[0].asxml(currentIndent + indent)
		xmlf = fmt.Sprint(
			currentIndent, "<next>\n",
			xmlff,
			currentIndent, "</next>\n",
		)
	case "U":
		xmlbefore := f.operand[0].asxml(currentIndent + indent + indent)
		xmlafter := f.operand[1].asxml(currentIndent + indent + indent)
		xmlf = fmt.Sprint(
			currentIndent, "<until>\n",
			currentIndent, indent, "<before>\n", xmlbefore,
			currentIndent, indent, "</before>\n",
			currentIndent, indent, "<after>\n", xmlafter,
			currentIndent, indent, "</after>\n",
			currentIndent, "</until>\n",
		)
	case "is-fireable":
		xmlt := f.operand[0].astransition(currentIndent + indent)
		for i := 1; i < len(f.operand); i++ {
			xmlt = xmlt + f.operand[i].astransition(currentIndent+indent)
		}
		xmlf = fmt.Sprint(
			currentIndent, "<is-fireable>\n",
			xmlt,
			currentIndent, "</is-fireable>\n",
		)
	case "leq":
		xmlsmall := f.operand[0].asxml(currentIndent + indent)
		xmlbig := f.operand[1].asxml(currentIndent + indent)
		xmlf = fmt.Sprint(
			currentIndent, "<integer-le>\n",
			xmlsmall, xmlbig,
			currentIndent, "</integer-le>\n",
		)
	case "token-count":
		xmlp := f.operand[0].asplace(currentIndent + indent)
		for i := 1; i < len(f.operand); i++ {
			xmlp = xmlp + f.operand[i].asplace(currentIndent+indent)
		}
		xmlf = fmt.Sprint(
			currentIndent, "<token-count>\n",
			xmlp,
			currentIndent, "</token-count>\n",
		)
	case "integer-constant":
		xmlf = fmt.Sprint(
			currentIndent, "<integer-constant>",
			f.operand[0].operator.name,
			"</integer-constant>\n",
		)
	}

	return xmlf
}

func (f formula) astransition(currentIndent string) (t string) {
	t = fmt.Sprint(
		currentIndent, "<transition>", f.operator.name, "</transition>\n",
	)
	return t
}

func (f formula) asplace(currentIndent string) (p string) {
	p = fmt.Sprint(
		currentIndent, "<place>", f.operator.name, "</place>\n",
	)
	return p
}