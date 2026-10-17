//go:build !windows

package ytdlp

import (
	"os"
	"os/exec"
)

func interrupt(cmd *exec.Cmd) error {
	return cmd.Process.Signal(os.Interrupt)
}
