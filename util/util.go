package util

import (
	"path/filepath"
	"runtime"
)

func GetCurrentPath() string {
	_, file, _, _ := runtime.Caller(1)

	return filepath.Dir(file)
}

func GetProjectRootPath() string {
	_, b, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(b), "../")
}
