package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func (m *modelInfo) filter(formulas []formula, numToFind int) []int {

	// model
	modelPath := m.filePath

	// formulas
	if m.modelType == col {
		// if colored with no twin we can do nothing, just keep the formulas
		if m.twinModel == nil {
			log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), COL model without PT equivalent, cannot filter formulas")
			res := make([]int, len(formulas))
			for i := 0; i < len(formulas); i++ {
				res[i] = i
			}
			return res
		}
		// if colored and has twin, unfold for using SMC
		unfoldedFormulas := make([]formula, len(formulas))
		for i := 0; i < len(formulas); i++ {
			unfoldedFormulas[i] = m.twinModel.unfolding(formulas[i])
		}
		formulas = unfoldedFormulas
		// change the model accordingly
		modelPath = m.twinModel.filePath
	}
	m.writeFormulas(formulas, globalSMCTmpFileName, false)

	// smc run
	log.Print(m.modelName, " (", m.modelInstance, ", ", m.modelType, "), running SMC on model ", modelPath, " with formulas file ", globalSMCTmpFileName)

	return runSMC(modelPath, globalSMCTmpFileName, numToFind)
}

func runSMC(model, formulas string, numToFind int) []int {
	tokeep := make([]int, 0)
	smcMaxStates := fmt.Sprint("--max-states=", globalSMCMaxStates)
	smcStopAfter := fmt.Sprint("--mcc15-stop-after=", numToFind)
	smcCommand := exec.Command("python", globalSMCPath, "--use10", smcMaxStates, smcStopAfter, model, formulas)
	filterCommand1 := exec.Command("grep", "-v", "^smc:")
	filterCommand2 := exec.Command("grep", "?")
	cutCommand := exec.Command("cut", "-d-", "-f5")
	command1Reader, smcWriter := io.Pipe()
	command2Reader, command1Writer := io.Pipe()
	cutCommandReader, command2Writer := io.Pipe()
	smcCommand.Stdout = smcWriter
	filterCommand1.Stdin = command1Reader
	filterCommand1.Stdout = command1Writer
	filterCommand2.Stdin = command2Reader
	filterCommand2.Stdout = command2Writer
	cutCommand.Stdin = cutCommandReader
	stdout, err := cutCommand.StdoutPipe()
	if err != nil {
		log.Fatal("filter, StdoutPipe(): ", err)
	}
	smcCommandOutput := bufio.NewReader(stdout)
	stderr, err := smcCommand.StderrPipe()
	if err != nil {
		log.Fatal("filter, StderrPipe(): ", err)
	}
	smcCommandError := bufio.NewReader(stderr)
	if err := smcCommand.Start(); err != nil {
		log.Fatal("filter, start SMC: ", err)
	}
	if err := filterCommand1.Start(); err != nil {
		log.Fatal("filter, start grep 1: ", err)
	}
	if err := filterCommand2.Start(); err != nil {
		log.Fatal("filter, start grep 2: ", err)
	}
	if err := cutCommand.Start(); err != nil {
		log.Fatal("filter, start cut: ", err)
	}
	res, err := smcCommandError.ReadString('\n')
	for ; err == nil; res, err = smcCommandError.ReadString('\n') {
		log.Print("SMC error: ", res)
	}
	if err != io.EOF {
		log.Fatal("filter, stderr reading error: ", err)
	}
	if err := smcCommand.Wait(); err != nil {
		log.Fatal("filter, wait SMC: ", err)
	}
	smcWriter.Close()
	if err := filterCommand1.Wait(); err != nil {
		log.Print("WARNING: problem during grep while filtering formulas: ", err)
	}
	command1Reader.Close()
	command1Writer.Close()
	if err := filterCommand2.Wait(); err != nil {
		log.Print("WARNING: problem during grep while filtering formulas: ", err)
	}
	command2Reader.Close()
	command2Writer.Close()
	res, err = smcCommandOutput.ReadString('\n')
	for ; err == nil; res, err = smcCommandOutput.ReadString('\n') {
		v, err := strconv.Atoi(strings.TrimSuffix(res, "\n"))
		if err != nil {
			log.Fatal("filter, atoi: ", err)
		}
		tokeep = append(tokeep, v)
	}
	if err != io.EOF {
		log.Fatal("filter, stdout reading error: ", err)
	}
	if err := cutCommand.Wait(); err != nil {
		log.Fatal("filter, wait cut: ", err)
	}
	cutCommandReader.Close()
	return tokeep
}
