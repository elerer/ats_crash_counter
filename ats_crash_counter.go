package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Crash_counter struct {
	Filepath    string
	Filename    string
	Datepattern string
	Textpattern string
}

func (cc *Crash_counter) getLogFileReverseContent() (*bufio.Reader, error) {
	fullfilename := cc.Filepath + "/" + cc.Filename
	//tac - reads reverse
	cmd := "tac " + fullfilename + "|grep -i \"" + cc.Textpattern + "\""

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(out)

	br := bufio.NewReader(r)

	return br, nil

}

func (cc *Crash_counter) getDate() (string, error) {
	dateBytes, err := exec.Command("date", cc.Datepattern).Output()

	if err != nil {
		return "", err
	}

	date := strings.Trim(string(dateBytes), "\n")
	ds := strings.Split(date, " ")
	// the manager.log date has two spaces when day is 1 digit
	if len(ds[1]) == 1 {
		date = strings.Replace(date, " ", "  ", 1)
	}

	return date, nil

}

func newCrashCounter() *Crash_counter {
	return &Crash_counter{
		Filepath:    "/var/log/trafficserver",
		Filename:    "manager.log",
		Datepattern: "+%b%e",
		Textpattern: "Server Process terminated due to Sig",
	}
}

func main() {

	cc := newCrashCounter()
	fileContent, err := cc.getLogFileReverseContent()

	if err != nil {
		return
	}

	date, err := cc.getDate()
	if err != nil {
		return
	}

	count := 0

	for {
		line, err := fileContent.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		if strings.Contains(line, date) {
			count++
		} else {
			break
		}

	}

	fmt.Print(count)
}
