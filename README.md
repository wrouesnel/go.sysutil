[![Build Status](https://travis-ci.org/wrouesnel/go.sysutil.svg?branch=master)](https://travis-ci.org/wrouesnel/go.sysutil)
[![Coverage Status](https://coveralls.io/repos/github/wrouesnel/go.sysutil/badge.svg?branch=master)](https://coveralls.io/github/wrouesnel/go.sysutil?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/wrouesnel/go.sysutil)](https://goreportcard.com/report/github.com/wrouesnel/go.sysutil)

# System programming utility libraries for Golang

This repository is a collection of useful functions I find myself
commonly re-implementing in my Go projects, which I've broken out here
to be more reusable.

## logutil
Primarily impements the LogWriter, which allows dumping the output of
io.Writer's to Golang logging library compatible loggers.

## fsutil
Implements many python-like file functions. If you're doing a lot of
operations on files, then these are not the most efficient way to do it
but they do provide a very easy way to get through some trivial
ones.

## executil
Subprocess execution utility functions, focused on providing an easy
shell-like experience in Go.
