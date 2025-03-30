/*
 * Copyright (C) 2025 Ian M. Fink.  All rights reserved.
 *
 * This program is free software:  you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option)
 * any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
 * or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
 * more details.
 *
 * You should have received a copy of the GNU General Public License along
 * with this program.  If not, please see: https://www.gnu.org/licenses.
 *
 * Tabstop:	4
 */

package main

/*
 * Imports
 */

import (
	"fmt"
	"log"
	"io"
	"bufio"
	"strings"
	"time"
	"os/exec"
)

/**********************************************************************/

func outputReceiver(theReadCloser io.ReadCloser, outputChannel chan string,
		outputDoneChannel chan bool) {
	var (
		myReader	*bufio.Reader
		bufString	string
		err			error
	)

	myReader = bufio.NewReader(theReadCloser)

	bufString, err = myReader.ReadString('\n')
	for err != io.EOF {
		bufString = strings.TrimSuffix(bufString, "\n")
		// fmt.Printf("myWorker = '%s'\n", bufString)
		outputChannel <- bufString
		bufString, err = myReader.ReadString('\n')
	}

	outputDoneChannel <- true

} /* outputReceiver */

/**********************************************************************/

func main() {
	var (
		stdoutChannel		chan string
		stderrChannel		chan string
		stdoutDoneChannel	chan bool
		stderrDoneChannel	chan bool
		tmpStdoutString		string
		tmpStderrString		string
		doneStdout			bool
		doneStderr			bool
		err					error
		cmdPtr				*exec.Cmd
		cmdStdout			io.ReadCloser
		cmdStderr			io.ReadCloser
		theTimer			*time.Timer
		theTime				time.Time
		theDuration			time.Duration
	)

	fmt.Println("****************--****************")

	// set the duration of the timer
	theDuration = 500 * time.Millisecond

	// allocate the channels
	stdoutChannel		= make(chan string)
	stderrChannel		= make(chan string)
	stdoutDoneChannel	= make(chan bool)
	stderrDoneChannel	= make(chan bool)

	// use whatever command and its parameters you like
	// cmdPtr = exec.Command("ls", "-las")
	cmdPtr = exec.Command("testStderr") // outputs to stdout & stderr

	cmdStdout, err = cmdPtr.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	cmdStderr, err = cmdPtr.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	go outputReceiver(cmdStdout, stdoutChannel, stdoutDoneChannel)
	go outputReceiver(cmdStderr, stderrChannel, stderrDoneChannel)

	err = cmdPtr.Start()
	if err != nil {
		log.Fatal(err)
	}

	// initialize the conditions
	doneStdout = false
	doneStderr = false

	// start a timer and set its duration
	theTimer = time.NewTimer(theDuration)

	for !doneStdout || !doneStderr {
		select {
			case tmpStdoutString = <- stdoutChannel:
				_ = theTimer.Stop()
				fmt.Printf("stdoutChannel Stdout: '%s'\n", tmpStdoutString)
				theTimer.Reset(theDuration)
				
			case tmpStderrString = <- stderrChannel:
				_ = theTimer.Stop()
				fmt.Printf("stderrChannel Stderr: '%s'\n", tmpStderrString)
				theTimer.Reset(theDuration)

			case theTime = <- theTimer.C:
				_ = theTimer.Stop()
				fmt.Printf("Timer fired at: %02d:%02d:%02d.%03d\n",
					theTime.Hour(), theTime.Minute(), theTime.Second(),
					theTime.Nanosecond() / 1000000) // convert to milliseconds
				theTimer.Reset(theDuration)
				
			case doneStdout = <- stdoutDoneChannel:

			case doneStderr = <- stderrDoneChannel:
		}
	}

	// wait for the command to exit
	err = cmdPtr.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	fmt.Println("****************--****************")

} /* main */

/**********************************************************************/

/*
 * End of file:	main.go
 */


