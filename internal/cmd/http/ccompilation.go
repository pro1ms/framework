package http

import (
	"bytes"
	"fmt"
	"os/exec"
)

func (s *Scanner) checkCompilation(dir string) error {
	cmd := exec.Command("go", "build", "-o", "/dev/null", dir+"/...")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Gen HTPP: compilation failed:\n%s", stderr.String())
	}

	return nil
}
