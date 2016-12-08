package fsutil

import (
	"github.com/kardianos/osext"
	"os"
	"os/exec"
)

// ErrPath is returned as the panic() context when a Must function is invoked
// as a wrapped error..
type ErrPath struct {
	message string
	path string
}

func newErrPath(message string, path string) *ErrPath {
	return &ErrPath{
		message: message,
		path: path,
	}
}

func (this ErrPath) WrappedErrors() []error {
	return []error{}
}

func (this ErrPath) Error() string {
	return this.message
}

func (this ErrPath) Path() string {
	return this.path
}

// Exit with a panic if paths do not exist as command executables
func MustLookupPaths(paths ...string) {
	for _, path := range paths {
		_, err := exec.LookPath(path)
		if err != nil {
			panic(newErrPath("Could not find path", path))
		}
	}

}

// Exit Panicly if path does not exist
func MustPathExist(paths ...string) {
	for _, path := range paths {
		if !PathExists(path) {
			panic(newErrPath("Cannot continue", nil))
		}
	}
}

// Exit Panicly if path exists
func MustPathNotExist(paths ...string) {
	for _, path := range paths {
		if PathExists(path) {
			panic(newErrPath("Cannot continue", nil))
		}
	}
}

// Path does not exist
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Path exists
func PathNotExist(path string) bool {
	if !PathExists(path) {
		return true
	}
	return false
}

// Get the current executable's folder or fail
func MustExecutableFolder() string {
	folder, err := osext.ExecutableFolder()
	if err != nil {
		panic(newErrPath("Could not get executable folder", err))
	}
	return folder
}

func GetFilePerms(filename string) (os.FileMode, error) {
	st, err := os.Stat(filename)
	if err != nil {
		return os.FileMode(0777), err
	}
	return st.Mode(), nil
}

func MustGetFileSize(filename string) int64 {
	size, err := GetFileSize(filename)
	if err != nil {
		panic(newErrPath("Could not get file size", err))
	}
	return size
}

func GetFileSize(filename string) (int64, error) {
	st, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return st.Size(), nil
}