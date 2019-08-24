package main

import (
	"flag"
	"github.com/mpetavy/common"
	"os/signal"

	"io"
	"os"
)

var (
	input  = flag.String("i", "", "file to read data, use OS.STDIN when omitted")
	output = flag.String("o", "", "file to save data additionally to OS.STDOUT")
	append = flag.Bool("a", false, "append the output to the output file")
	ignore = flag.Bool("s", false, "ignore the SIGINT signal")
)

func init() {
	common.Init("tee", "1.0.2", "2017", "Passthrough STDIN/file to STDOUT and/or file (optional)", "mpetavy", common.APACHE, "https://github.com/mpetavy/tee", false, nil, nil, run, 0)
}

func run() error {
	outputFlag := os.O_WRONLY | os.O_CREATE
	if *append {
		outputFlag |= os.O_APPEND
	}

	if *ignore {
		signal.Ignore(os.Interrupt)
	}

	var inputFile = os.Stdin
	var outputFile *os.File
	var err error

	if *input != "" {
		inputFile, err = os.Open(*input)
		if err != nil {
			return err
		}
	}
	if *output != "" {
		outputFile, err = os.OpenFile(*output, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)

		if err != nil {
			return err
		}

		defer outputFile.Close()
	}

	b := make([]byte, 8192)
	c := 0

	for {
		n, err := inputFile.Read(b)

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if *output != "" {
			outputFile.Write(b[:n])
		}

		_, err = os.Stdout.Write(b[:n])
		if err != nil {
			return err
		}

		c += n
	}

	if *output != "" {
		err := outputFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	defer common.Done()

	common.NoBanner = true

	common.Run(nil)
}
