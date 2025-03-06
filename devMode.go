package main

import (
	"bufio"
	"os/exec"
	"strings"
)

func startDevMode(reload chan bool) {
	cmd := exec.Command("tailwindcss", "-i", "./assets/input.css", "-o", "./assets/output.css", "--watch=always")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	go func() {
		reader := bufio.NewReader(stderr)
		line, err := reader.ReadString('\n')
		for err == nil {
			if strings.Contains(line, "Done in ") {
				reload <- true
			}
			line, err = reader.ReadString('\n')
		}
	}()

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}
