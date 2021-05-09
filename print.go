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

// print a set of formulas as xml in a file for a given model
func (m modelInfo) writexmlFormulas(formulas []formula, fileName string, formulaType string, inModelDirectory bool, logger *log.Logger) {

	filePath := fileName
	if inModelDirectory {
		filePath = filepath.Join(m.directory, fileName)
	}

	f, error := os.Create(filePath)
	if error != nil {
		logger.Print("ERROR: cannot create file ", filePath)
		return
	}

	_, error = f.WriteString(
		fmt.Sprint(
			"<?xml version=\"1.0\"?>\n",
			"<property-set xmlns=\"http://mcc.lip6.fr/\">\n",
		))
	if error != nil {
		logger.Print("ERROR: cannot write to file ", filePath)
		return
	}

	for i := 0; i < len(formulas); i++ {
		_, error = f.WriteString(formulas[i].xmlPrint(m, i, formulaType))
		if error != nil {
			logger.Print("ERROR: cannot write to file ", filePath)
			return
		}
	}

	_, error = f.WriteString(
		fmt.Sprint(
			"</property-set>\n",
		))
	if error != nil {
		logger.Print("ERROR: cannot write to file ", filePath)
		return
	}

	error = f.Sync()
	if error != nil {
		logger.Print("ERROR: cannot sync file ", filePath)
		return
	}

	error = f.Close()
	if error != nil {
		logger.Print("ERROR: cannot close file ", filePath)
		return
	}
}

// output one formula as xml
func (f formula) xmlPrint(m modelInfo, num int, formulaType string) (xmlp string) {
	kind := "COL"
	if m.modelType != col {
		kind = "PT"
	}

	fid := fmt.Sprintf(
		"%s-%s-%s-%s-%2.2d",
		m.modelName, kind, m.modelInstance, formulaType, num,
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
			currentIndent, indent, "<reach>\n", xmlafter,
			currentIndent, indent, "</reach>\n",
			currentIndent, "</until>\n",
		)
	case "is-fireable":
		xmlt := f.operand[0].asxmltransition(currentIndent + indent)
		for i := 1; i < len(f.operand); i++ {
			xmlt = xmlt + f.operand[i].asxmltransition(currentIndent+indent)
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
	case "tokens-count":
		xmlp := f.operand[0].asxmlplace(currentIndent + indent)
		for i := 1; i < len(f.operand); i++ {
			xmlp = xmlp + f.operand[i].asxmlplace(currentIndent+indent)
		}
		xmlf = fmt.Sprint(
			currentIndent, "<tokens-count>\n",
			xmlp,
			currentIndent, "</tokens-count>\n",
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

func (f formula) asxmltransition(currentIndent string) (t string) {
	t = fmt.Sprint(
		currentIndent, "<transition>", f.operator.name, "</transition>\n",
	)
	return t
}

func (f formula) asxmlplace(currentIndent string) (p string) {
	p = fmt.Sprint(
		currentIndent, "<place>", f.operator.name, "</place>\n",
	)
	return p
}

// print a set of formulas as a human-readable format in a file for a given model
func (m modelInfo) writehrFormulas(formulas []formula, fileName string, formulaType string, inModelDirectory bool, logger *log.Logger) {

	filePath := fileName
	if inModelDirectory {
		filePath = filepath.Join(m.directory, fileName)
	}

	f, error := os.Create(filePath)
	if error != nil {
		logger.Print("ERROR: cannot create file ", filePath)
		return
	}

	for i := 0; i < len(formulas); i++ {
		_, error = f.WriteString(formulas[i].hrPrint(m, i, formulaType))
		if error != nil {
			logger.Print("ERROR: cannot write to file ", filePath)
			return
		}
	}

	error = f.Sync()
	if error != nil {
		logger.Print("ERROR: cannot sync file ", filePath)
		return
	}

	error = f.Close()
	if error != nil {
		logger.Print("ERROR: cannot close file ", filePath)
		return
	}
}

// output one formula in a human-readable format
func (f formula) hrPrint(m modelInfo, num int, formulaType string) (hrp string) {
	kind := "COL"
	if m.modelType != col {
		kind = "PT"
	}

	fid := fmt.Sprintf(
		"Property %s-%s-%s-%s-%2.2d",
		m.modelName, kind, m.modelInstance, formulaType, num,
	)

	hrf := f.ashr()

	hrp = fmt.Sprint(
		fid, "\n",
		indent, "\"Automatically generated by Citili ", version, "\"\n",
		indent, "is:\n",
		indent, indent, hrf, "\n",
		indent, "end.\n",
	)
	return hrp
}

func (f formula) ashr() (hrf string) {

	switch f.operator.name {
	case "A":
		hrf = "A (" + f.operand[0].ashr() + ")"
	case "E":
		hrf = "E (" + f.operand[0].ashr() + ")"
	case "not":
		hrf = "! (" + f.operand[0].ashr() + ")"
	case "and":
		hrf = "(" + f.operand[0].ashr() + ")"
		for i := 1; i < len(f.operand); i++ {
			hrf = hrf + " & (" + f.operand[i].ashr() + ")"
		}
	case "or":
		hrf = "(" + f.operand[0].ashr() + ")"
		for i := 1; i < len(f.operand); i++ {
			hrf = hrf + " | (" + f.operand[i].ashr() + ")"
		}
	case "G":
		hrf = "G (" + f.operand[0].ashr() + ")"
	case "F":
		hrf = "F (" + f.operand[0].ashr() + ")"
	case "X":
		hrf = "X (" + f.operand[0].ashr() + ")"
	case "U":
		hrf = "(" + f.operand[0].ashr() + ") U (" + f.operand[1].ashr() + ")"
	case "is-fireable":
		hrff := f.operand[0].ashrtransition()
		for i := 1; i < len(f.operand); i++ {
			hrff = hrff + ", " + f.operand[i].ashrtransition()
		}
		hrf = "is-fireable(" + hrff + ")"
	case "leq":
		hrf = f.operand[0].ashr() + " <= " + f.operand[1].ashr()
	case "tokens-count":
		hrff := f.operand[0].ashrplace()
		for i := 1; i < len(f.operand); i++ {
			hrff = hrff + ", " + f.operand[i].ashrplace()
		}
		hrf = "tokens-count(" + hrff + ")"
	case "integer-constant":
		hrf = f.operand[0].operator.name
	}

	return hrf
}

func (f formula) ashrtransition() string {
	return "\"" + f.operator.name + "\""
}

func (f formula) ashrplace() string {
	return "\"" + f.operator.name + "\""
}
