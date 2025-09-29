package util

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentPath(t *testing.T) {
	// Test that GetCurrentPath returns a valid directory path
	path := GetCurrentPath()

	assert.NotEmpty(t, path)
	assert.False(t, strings.HasSuffix(path, "/"))

	// Verify it's a directory (not a file)
	assert.DirExists(t, path)

	// The path should contain "util" since we're in the util package
	assert.Contains(t, path, "util")
}

func TestGetProjectRootPath(t *testing.T) {
	// Test that GetProjectRootPath returns a valid directory path
	path := GetProjectRootPath()

	assert.NotEmpty(t, path)
	assert.False(t, strings.HasSuffix(path, "/"))

	// Verify it's a directory (not a file)
	assert.DirExists(t, path)

	// The path should be the project root, containing go.mod
	goModPath := filepath.Join(path, "go.mod")
	assert.FileExists(t, goModPath)
}

func TestPathConsistency(t *testing.T) {
	// Test that the paths are consistent and related
	currentPath := GetCurrentPath()
	projectRootPath := GetProjectRootPath()

	// Current path should be a subdirectory of project root
	relPath, err := filepath.Rel(projectRootPath, currentPath)
	assert.NoError(t, err)
	assert.NotEqual(t, ".", relPath) // Should not be the same directory

	// The relative path should contain "util"
	assert.Contains(t, relPath, "util")
}

func TestGetCurrentPath_CallerInfo(t *testing.T) {
	// Test that GetCurrentPath correctly uses runtime.Caller
	// This test verifies the function works as expected with the caller parameter

	// Get the path using the function
	funcPath := GetCurrentPath()

	// Get the path manually using runtime.Caller(0) - this gives the current test file's directory
	_, file, _, ok := runtime.Caller(0)
	assert.True(t, ok)
	testFilePath := filepath.Dir(file)

	// The paths should be the same because GetCurrentPath uses Caller(1)
	// and we're calling it from this test function (caller level 1)
	assert.Equal(t, testFilePath, funcPath)
}

func TestGetProjectRootPath_FromUtil(t *testing.T) {
	// Test that GetProjectRootPath correctly navigates from util directory
	path := GetProjectRootPath()

	// The path should contain the project name
	assert.Contains(t, path, "go-hexagonal")

	// Verify it contains the expected directory structure
	assert.FileExists(t, filepath.Join(path, "go.mod"))
	assert.DirExists(t, filepath.Join(path, "util"))
	assert.DirExists(t, filepath.Join(path, "domain"))
	assert.DirExists(t, filepath.Join(path, "application"))
}

func TestPathValidity(t *testing.T) {
	// Test that both functions return valid, absolute paths
	currentPath := GetCurrentPath()
	projectRootPath := GetProjectRootPath()

	// Both should be absolute paths
	assert.True(t, filepath.IsAbs(currentPath))
	assert.True(t, filepath.IsAbs(projectRootPath))

	// Neither should contain relative components
	assert.False(t, strings.Contains(currentPath, "./"))
	assert.False(t, strings.Contains(currentPath, "../"))
	assert.False(t, strings.Contains(projectRootPath, "./"))
	assert.False(t, strings.Contains(projectRootPath, "../"))
}

func TestMultipleCallsConsistency(t *testing.T) {
	// Test that multiple calls to the same function return consistent results
	path1 := GetCurrentPath()
	path2 := GetCurrentPath()

	assert.Equal(t, path1, path2)

	rootPath1 := GetProjectRootPath()
	rootPath2 := GetProjectRootPath()

	assert.Equal(t, rootPath1, rootPath2)
}

func TestCrossPlatformPathHandling(t *testing.T) {
	// Test that paths are handled correctly regardless of platform
	currentPath := GetCurrentPath()
	projectRootPath := GetProjectRootPath()

	// Paths should use the correct separator for the current platform
	if runtime.GOOS == "windows" {
		assert.Contains(t, currentPath, "\\")
		assert.Contains(t, projectRootPath, "\\")
	} else {
		assert.Contains(t, currentPath, "/")
		assert.Contains(t, projectRootPath, "/")
	}
}

func TestPathNavigation(t *testing.T) {
	// Test that we can navigate from current path to project root and vice versa
	currentPath := GetCurrentPath()
	projectRootPath := GetProjectRootPath()

	// Navigate from project root to util directory
	utilPathFromRoot := filepath.Join(projectRootPath, "util")
	assert.DirExists(t, utilPathFromRoot)

	// Navigate from util directory to project root
	relPath, err := filepath.Rel(currentPath, projectRootPath)
	assert.NoError(t, err)

	// The relative path should go up one level
	assert.True(t, strings.HasPrefix(relPath, ".."))
}

func TestFilepathOperations(t *testing.T) {
	// Test that the returned paths work with standard filepath operations
	currentPath := GetCurrentPath()
	projectRootPath := GetProjectRootPath()

	// Test Join operation
	joinedPath := filepath.Join(currentPath, "test_file.txt")
	assert.Contains(t, joinedPath, "util")
	assert.Contains(t, joinedPath, "test_file.txt")

	// Test Base operation
	base := filepath.Base(currentPath)
	assert.Equal(t, "util", base)

	// Test Dir operation
	dir := filepath.Dir(currentPath)
	assert.Equal(t, projectRootPath, dir)
}

func TestRuntimeCallerBehavior(t *testing.T) {
	// Test the behavior of runtime.Caller to understand how GetCurrentPath works

	// Caller(0) returns information about the current function
	pc, file, line, ok := runtime.Caller(0)
	assert.True(t, ok)
	assert.NotNil(t, pc)
	assert.NotEmpty(t, file)
	assert.Greater(t, line, 0)

	// Caller(1) would return information about the caller of this function
	_, file1, line1, ok1 := runtime.Caller(1)
	assert.True(t, ok1)
	assert.NotEmpty(t, file1)
	assert.Greater(t, line1, 0)
}

func TestErrorHandling(t *testing.T) {
	// Test that the functions handle edge cases gracefully
	// These functions don't return errors, but we can test their behavior

	// Multiple calls should not panic
	for i := 0; i < 10; i++ {
		path := GetCurrentPath()
		assert.NotEmpty(t, path)

		rootPath := GetProjectRootPath()
		assert.NotEmpty(t, rootPath)
	}
}

// Benchmark tests
func BenchmarkGetCurrentPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetCurrentPath()
	}
}

func BenchmarkGetProjectRootPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetProjectRootPath()
	}
}

// Test helper function to verify path structure
func verifyPathStructure(t *testing.T, path string, expectedComponents []string) {
	for _, component := range expectedComponents {
		assert.Contains(t, path, component)
	}
}

func TestPathStructure(t *testing.T) {
	currentPath := GetCurrentPath()
	projectRootPath := GetProjectRootPath()

	// Verify current path structure
	verifyPathStructure(t, currentPath, []string{"go-hexagonal", "util"})

	// Verify project root path structure
	verifyPathStructure(t, projectRootPath, []string{"go-hexagonal"})

	// Project root should not contain "util" in the path
	assert.False(t, strings.Contains(projectRootPath, "/util/"))
	assert.False(t, strings.Contains(projectRootPath, "\\util\\"))
}

func TestConcurrentAccess(t *testing.T) {
	// Test that the functions can be called concurrently
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			path := GetCurrentPath()
			assert.NotEmpty(t, path)

			rootPath := GetProjectRootPath()
			assert.NotEmpty(t, rootPath)

			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
