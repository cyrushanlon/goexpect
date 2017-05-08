package expect

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type Process struct {
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	stdout      bytes.Buffer
	stdoutQueue []string
	stderr      bytes.Buffer
	stderrQueue []string
	done        bool

	Timeout int //seconds for expect to wait
}

//listenToPipe
//TODO
//combine these functions so there isnt so much duplicated code
func (p *Process) listenToOut(pipe *bytes.Buffer) {
	for !p.done {
		//fmt.Print("|")
		//declare them once here so we dont waste time inside when reading
		var b byte
		var err error
		line := []byte{}
		for pipe.Len() > 0 {
			b, err = pipe.ReadByte()
			if err == nil {
				line = append(line, b)
			}
		}
		if len(line) > 0 {
			fmt.Print(string(line))
			p.stdoutQueue = append(p.stdoutQueue, string(line))
		}

		time.Sleep(10 * time.Millisecond)
	}
	fmt.Print(string(pipe.Bytes()))
}
func (p *Process) listenToErr(pipe *bytes.Buffer) {
	for !p.done {
		//declare them once here so we dont waste time inside when reading
		var b byte
		var err error
		line := []byte{}
		for pipe.Len() > 0 {
			b, err = pipe.ReadByte()
			if err == nil {
				line = append(line, b)
			}
		}
		if len(line) > 0 {
			p.stdoutQueue = append(p.stderrQueue, string(line))
		}
	}
	fmt.Print(string(pipe.Bytes()))
}

//SendInput passes a string to the subprocesses stdin
func (p *Process) SendInput(in string) {
	in += "\r\n"
	io.WriteString(p.stdin, in)
}

//Start spawns the sub process with the given arguments
func (p *Process) Start(args ...string) error {

	p.cmd = exec.Command(args[0], args[1:]...)

	p.stdin, _ = p.cmd.StdinPipe()
	p.cmd.Stdout = &p.stdout
	p.cmd.Stderr = &p.stderr

	err := p.cmd.Start()
	if err != nil {
		p.done = true
		return err
	}

	p.done = false

	go p.listenToOut(&p.stdout)
	//we dont do anything with the erorr stream so dont bother listening to it
	//go p.listenToErr(&p.stderr)

	return nil
}

//Expect waits for a given string to show up in the input
func (p *Process) Expect(compare string, nocase bool) bool {

	startTime := time.Now()
	for time.Now().Unix()-startTime.Unix() < int64(p.Timeout) {

		if len(p.stdoutQueue) > 0 {

			curLine := p.stdoutQueue[len(p.stdoutQueue)-1]
			if nocase {
				compare = strings.ToLower(compare)
				curLine = strings.ToLower(curLine)
			}

			if strings.Contains(curLine, compare) {
				return true
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
	return false
}

//Close ends the subprocess
func (p *Process) Close() {

	if p.cmd != nil {
		p.cmd.Process.Kill()
	}
	p.done = true
}
