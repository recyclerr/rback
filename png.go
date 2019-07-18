package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/emicklei/dot"
)

func writePng(graph *dot.Graph, filename string) error {
	buf := new(bytes.Buffer)
	cmd := exec.Command("dot", "-Tpng", "-Gsplines=spline", "-Kdot")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(graph.String()), buf, os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute dot. Is Graphviz installed? Error: %v", err)
	}

	if err := ioutil.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing to file: %v\n", err)
	}

	return nil
}
