package fsutil_test

import (
	"github.com/wrouesnel/go.sysutil/fsutil"
	. "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type PathSuite struct{}

var _ = Suite(&PathSuite{})

func (s *PathSuite) TestAll(c *C) {
	fsutil.MustLookupPaths("/bin/sh")
	fsutil.MustPathExist("/bin/sh")
	fsutil.MustPathNotExist("/notarealpath")

	exeFolder := fsutil.MustExecutableFolder()
	c.Assert(exeFolder, Not(Equals), "")

	_, err := fsutil.GetFilePerms("/bin/sh")
	c.Assert(err, IsNil)

	fsize := fsutil.MustGetFileSize("/bin/sh")
	c.Assert(fsize, Not(Equals), 0)
}
