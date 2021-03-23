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

func (m *modelInfo) filter(formulas []formula, numToFind int, canUnfold bool, logger *log.Logger) []int {

	// model
	modelPath := m.filePath

	// formulas
	if m.modelType == col {
		// if colored with no twin we can do nothing, just keep the formulas
		if m.twinModel == nil {
			logger.Print("COL model without PT equivalent, cannot filter formulas")
			res := make([]int, len(formulas))
			for i := 0; i < len(formulas); i++ {
				res[i] = i
			}
			return res
		}
		// if colored with twin but no correct mapping to PT
		if !canUnfold {
			logger.Print("COL model with PT equivalent but that cannot be unfold (impossible mapping), cannot filter formulas")
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
	m.writexmlFormulas(formulas, globalSMCTmpFileName, "ForFiltering", false, logger)

	// smc run
	logger.Print("running SMC on model ", modelPath, " with formulas file ", globalSMCTmpFileName)

	return runSMC(modelPath, globalSMCTmpFileName, numToFind, logger)
}

func runSMC(model, formulas string, numToFind int, logger *log.Logger) []int {
	tokeep := make([]int, 0)
	smcMaxStates := fmt.Sprint("--max-states=", globalSMCMaxStates)
	smcStopAfter := fmt.Sprint("--mcc15-stop-after=", numToFind)
	smcCommand := exec.Command("python", globalSMCPath, "--use10", smcMaxStates, smcStopAfter, model, formulas)
	logSMCCommand := exec.Command("tee", "-a", globalSMClogfile)
	filterCommand1 := exec.Command("grep", "-v", "^smc:")
	filterCommand2 := exec.Command("grep", "?")
	cutCommand := exec.Command("cut", "-d-", "-f5")
	logReader, smcWriter := io.Pipe()
	command1Reader, logWriter := io.Pipe()
	command2Reader, command1Writer := io.Pipe()
	cutCommandReader, command2Writer := io.Pipe()
	smcCommand.Stdout = smcWriter
	logSMCCommand.Stdin = logReader
	logSMCCommand.Stdout = logWriter
	filterCommand1.Stdin = command1Reader
	filterCommand1.Stdout = command1Writer
	filterCommand2.Stdin = command2Reader
	filterCommand2.Stdout = command2Writer
	cutCommand.Stdin = cutCommandReader
	stdout, err := cutCommand.StdoutPipe()
	if err != nil {
		logger.Print("ERROR: filter, StdoutPipe(): ", err)
		return tokeep
	}
	smcCommandOutput := bufio.NewReader(stdout)
	stderr, err := smcCommand.StderrPipe()
	if err != nil {
		logger.Print("ERROR: filter, StderrPipe(): ", err)
		return tokeep
	}
	smcCommandError := bufio.NewReader(stderr)
	if err := smcCommand.Start(); err != nil {
		logger.Print("ERROR: filter, start SMC: ", err)
		return tokeep
	}
	if err := logSMCCommand.Start(); err != nil {
		logger.Print("ERROR: filter, start tee: ", err)
		return tokeep
	}
	if err := filterCommand1.Start(); err != nil {
		logger.Print("ERROR: filter, start grep 1: ", err)
		return tokeep
	}
	if err := filterCommand2.Start(); err != nil {
		logger.Print("ERROR: filter, start grep 2: ", err)
		return tokeep
	}
	if err := cutCommand.Start(); err != nil {
		logger.Print("ERROR: filter, start cut: ", err)
		return tokeep
	}
	res, err := smcCommandError.ReadString('\n')
	for ; err == nil; res, err = smcCommandError.ReadString('\n') {
		logger.Print("SMC ERROR: ", res)
	}
	if err != io.EOF {
		logger.Print("ERROR: filter, stderr reading error: ", err)
		return tokeep
	}
	if err := smcCommand.Wait(); err != nil {
		logger.Print("ERROR: filter, wait SMC: ", err)
		return tokeep
	}
	smcWriter.Close()
	if err := logSMCCommand.Wait(); err != nil {
		logger.Print("ERROR: filter, wait tee: ", err)
		return tokeep
	}
	logReader.Close()
	logWriter.Close()
	if err := filterCommand1.Wait(); err != nil {
		logger.Print("WARNING: problem during grep while filtering formulas: ", err)
	}
	command1Reader.Close()
	command1Writer.Close()
	if err := filterCommand2.Wait(); err != nil {
		logger.Print("WARNING: problem during grep while filtering formulas (probably, no difficult formula was found): ", err)
	}
	command2Reader.Close()
	command2Writer.Close()
	res, err = smcCommandOutput.ReadString('\n')
	for ; err == nil; res, err = smcCommandOutput.ReadString('\n') {
		v, err := strconv.Atoi(strings.TrimSuffix(res, "\n"))
		if err != nil {
			logger.Print("ERROR: filter, atoi: ", err)
		} else {
			tokeep = append(tokeep, v)
		}
	}
	if err != io.EOF {
		logger.Print("ERROR: filter, stdout reading error: ", err)
		return tokeep
	}
	if err := cutCommand.Wait(); err != nil {
		logger.Print("ERROR: filter, wait cut: ", err)
		return tokeep
	}
	cutCommandReader.Close()
	return tokeep
}
