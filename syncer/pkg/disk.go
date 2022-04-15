package pkg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func Size(folderPath string) (uint64, error) {
	cmd := exec.Command("du", "-shm", folderPath)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return 0, fmt.Errorf("failed to check size for path '%s'. caused by: '%v'", folderPath, err)
	}

	output := strings.TrimSpace(out.String())
	parts := strings.Split(output, "\t")
	if len(parts) > 0 {
		return strconv.ParseUint(parts[0], 10, 64)
	}
	return 0, fmt.Errorf("failed to check size for path '%s'. unexpected output '%s'", folderPath, output)
}
