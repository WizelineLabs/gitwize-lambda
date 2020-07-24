package gitnative

import (
	"os/exec"
	"strconv"
	"strings"
)

// git show --unified=0 COMMIT_HASH | grep "^@@" -c

// GetInsertionPoint get number of block of code added/removed
func GetInsertionPoint(repoDir, commit string) (int, error) {
	command := "git show --unified=0 " + commit + " | grep ^@@ -c"
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	o := strings.Trim(string(output), "\n")
	return strconv.Atoi(o)
}
