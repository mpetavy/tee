package main

import (
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"os/signal"

	"io"
	"os"
)

var (
	input  = flag.String("i", "", "file to read data, use OS.STDIN when omitted")
	output = flag.String("o", "", "file to save data additionally to OS.STDOUT")
	add    = flag.Bool("a", false, "append the output to the output file")
	ignore = flag.Bool("s", false, "ignore the SIGINT signal")
)

func init() {
	common.Init("1.0.2", "2017", "Passthrough STDIN/file to STDOUT and/or file (optional)", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, run, 0)
}

func run() error {
	outputFlag := os.O_WRONLY | os.O_CREATE
	if *add {
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
		outputFile, err = os.OpenFile(*output, os.O_RDWR|os.O_CREATE|os.O_APPEND, common.DefaultFileMode)

		if err != nil {
			return err
		}

		defer func() {
			common.Error(outputFile.Close())
		}()
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

	*common.FlagNoBanner = true

	common.Run(nil)
}
