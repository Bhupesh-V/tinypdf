package shared

import (
	"os"
	"os/exec"
)

func FileSizeBytes(path string) int64 {
	info, err := os.Stat(path)
	if err != nil || info.Size() < 0 {
		return 0
	}
	return info.Size()
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
