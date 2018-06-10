package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type crashCounter struct {
	Filepath    string
	Filename    string
	Datepattern string
	Textpattern string
}

func (cc *crashCounter) getLogFileReverseContent() (*bufio.Reader, error) {
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

func (cc *crashCounter) getDate() (string, error) {
	monthBytes, err := exec.Command("date", "+%b").Output()

	if err != nil {
		return "", err
	}

	dayBytes, err := exec.Command("date", "+%d").Output()

	if err != nil {
		return "", err
	}

	return dateToManagerLogFormt(string(monthBytes), string(dayBytes)), nil
}

func dateToManagerLogFormt(m, d string) string {
	month := strings.Trim(m, "\n")
	day := strings.Trim(d, "\n")

	var date = month + " " + day
	// the manager.log date has two spaces when day is 1 digit
	if len(day) == 1 {
		date = month + "  " + day
	}
	return date
}

func newCrashCounter() *crashCounter {
	return &crashCounter{
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
