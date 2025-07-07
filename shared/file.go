package shared

import (
	"os"
	"os/exec"

	"github.com/dustin/go-humanize"
)

func FileSize(path string) string {
	info, err := os.Stat(path)
	if err != nil || info.Size() < 0 {
		return "0 B"
	}
	return humanize.Bytes(uint64(info.Size()))
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func IsBinaryAvailable(binFile string) bool {
	_, err := exec.LookPath(binFile)
	return err == nil
}
