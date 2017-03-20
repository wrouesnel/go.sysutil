package fsutil

import (
	"errors"
	"github.com/hashicorp/errwrap"
	"github.com/kardianos/osext"
	"os"
	"os/exec"
)

var (
	errCouldNotGetExecutableFolder = errors.New("could not get executable folder")
)

type ErrPath struct {
	message string
	path    string
}

func newErrPath(message string, path string) error {
	return error(&ErrPath{
		message: message,
		path:    path,
	})
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
			panic(newErrPath("path is not present in executable PATH", path))
		}
	}

}

// Exit Panicly if path does not exist
func MustPathExist(paths ...string) {
	for _, path := range paths {
		if !PathExists(path) {
			panic(newErrPath("path does not exist but must", path))
		}
	}
}

// Exit Panicly if path exists
func MustPathNotExist(paths ...string) {
	for _, path := range paths {
		if PathExists(path) {
			panic(newErrPath("path exists but musn't", path))
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
		panic(errCouldNotGetExecutableFolder)
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
		panic(err)
	}
	return size
}

func GetFileSize(filename string) (int64, error) {
	st, err := os.Stat(filename)
	if err != nil {
		return 0, errwrap.Wrap(newErrPath("could not get file size", filename), err)
	}
	return st.Size(), nil
}
