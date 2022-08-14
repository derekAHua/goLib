package utils

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

// @Author: Derek
// @Description:
// @Date: 2022/8/14 01:06
// @Version 1.0

func RunFmt(directory, fileName string) (err error) {
	_, err = RunCommand(directory, "/bin/bash", "-c", "gofmt -w "+fileName)
	return
}

func RunCommand(path, name string, arg ...string) (msg string, err error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Dir = path

	err = cmd.Run()
	if err != nil {
		msg = fmt.Sprint(err) + ": " + stderr.String()
		err = errors.New(msg)
		log.Println("err", err.Error(), "cmd", cmd.Args)
	}

	log.Println(out.String())
	return
}
