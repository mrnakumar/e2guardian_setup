package pkg

import (
	"fmt"
	"github.com/ricochet2200/go-disk-usage/du"
)

func Size(folderPath string) (uint64, error) {
	usage := du.NewDiskUsage(folderPath)
	if usage == nil {
		return 0, fmt.Errorf("invalid path '%s'", folderPath)
	}
	mb := uint64(1024 * 1024)
	return usage.Size() / mb, nil
}
