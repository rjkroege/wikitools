package wiki

import (
	"io/ioutil"
	"log"
	"os"
)

func ReplaceStdout() func() {
	logFile, err := ioutil.TempFile("/tmp", "wikipp")
	if err != nil {
		log.Panicf("leap couldn't make a logging file: %v", err)
	}
	log.SetOutput(logFile)
	originalstdout := os.Stdout
	os.Stdout = logFile

	return func() {
		log.SetOutput(os.Stderr)
		os.Stdout = originalstdout
	}
}
