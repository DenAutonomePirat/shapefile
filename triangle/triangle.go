package triangle

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func Triangulate() {
	cmd := exec.Command(
		"triangle",
		"-p",
		"-n",
		"-v",
		//"-j",
		"-Q",
		"temp",
	)
	cmd.Dir = "./triangle"
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(out.String())
}
