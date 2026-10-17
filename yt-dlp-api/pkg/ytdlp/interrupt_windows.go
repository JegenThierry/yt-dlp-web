package ytdlp

import (
	"os/exec"
	"strconv"
)

func interrupt(cmd *exec.Cmd) error {
	return exec.Command(string(TASKKILL), string(TASKKILL_TREE), string(TASKKILL_FORCE), string(TASKKILL_PID), strconv.Itoa(cmd.Process.Pid)).Run()
}
