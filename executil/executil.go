package executil

import (
	"os/exec"
	"bytes"
	"io"
)

// If the program dies, we need to shutdown gracefully. This channel is
// closed by the signal handler when we need to do a panic and die in the
// main go-routine (which triggers our deferred handlers).
// Every function here should react to this sensibly.
var interruptCh = make(chan interface{})

// Calling this closes the interrupt channel, causes all executing processes to
// die. This let's deferred clean up handlers run in the owning
// goroutine.
func Interrupt() {
	if interruptCh != nil {
		close(interruptCh)
		interruptCh = nil
	}
}

// ErrProcess is returned when there is an error during execution, or as part
// of a panic handler.
type ErrProcess struct {
	message string

	command string
	commandLine []string

	innerError error
}

func newErrProcess(message string, command string, commandLine []string, innerError error) *ErrProcess {
	return &ErrProcess{
		message: message,
		command: command,
		commandLine: commandLine,
		innerError: innerError,
	}
}

func (this ErrProcess) WrappedErrors() []error {
	return []error{this.innerError}
}

func (this ErrProcess) Error() string {
	return this.message
}

// Check for successful execution and return stdout and stderr as strings
func CheckExecWithOutput(command string, commandLine ...string) (string, string, error) {
	//log.Debugln("Executing Command:", command, commandLine)
	cmd := exec.Command(command, commandLine...)

	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)

	//cmd.Stdout = io.MultiWriter(stdoutBuffer,
	//	NewLogWriter(log.With("pipe", "stdout").With("cmd", command).Debugln))
	//cmd.Stderr = io.MultiWriter(stderrBuffer,
	//	NewLogWriter(log.With("pipe", "stderr").With("cmd", command).Debugln))

	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	// Wait on a go-routine for the process to exit
	doneCh := make(chan error)
	go func() {
		doneCh <- cmd.Wait()
	}()

	// Wait for process exit or global interrupt
	select {
	case err := <-doneCh:
		if err != nil {
			return stdoutBuffer.String(), stderrBuffer.String(), err
		}
	case <-interruptCh:
		cmd.Process.Kill()
		return stdoutBuffer.String(), stderrBuffer.String(), newErrProcess("Interrupted by external request", command, commandLine, nil)
	}

	return stdoutBuffer.String(), stderrBuffer.String(), nil
}

func CheckExecWithEnv(env []string, command string, commandLine ...string) error {
	//log.Debugln("Executing Command:", command, commandLine)
	cmd := exec.Command(command, commandLine...)

	cmd.Env = env
	//cmd.Stdout = NewLogWriter(log.With("pipe", "stdout").With("cmd", command).Debugln)
	//cmd.Stderr = NewLogWriter(log.With("pipe", "stderr").With("cmd", command).Debugln)

	if err := cmd.Start(); err != nil {
		return err
	}

	// Wait on a go-routine for the process to exit
	doneCh := make(chan error)
	go func() {
		doneCh <- cmd.Wait()
	}()

	// Wait for process exit or global interrupt
	select {
	case err := <-doneCh:
		if err != nil {
			return err
		}
	case <-interruptCh:
		cmd.Process.Kill()
		return newErrProcess("Interrupted by external request", command, commandLine, nil)
	}

	return nil
}

// Checks for successful execution. Logs all output at default level.
func CheckExec(command string, commandLine ...string) error {
	//log.Debugln("Executing Command:", command, commandLine)
	cmd := exec.Command(command, commandLine...)

	//cmd.Stdout = NewLogWriter(log.With("pipe", "stdout").With("cmd", command).Debugln)
	//cmd.Stderr = NewLogWriter(log.With("pipe", "stderr").With("cmd", command).Debugln)

	if err := cmd.Start(); err != nil {
		return err
	}

	// Wait on a go-routine for the process to exit
	doneCh := make(chan error)
	go func() {
		doneCh <- cmd.Wait()
	}()

	// Wait for process exit or global interrupt
	select {
	case err := <-doneCh:
		if err != nil {
			return err
		}
	case <-interruptCh:
		cmd.Process.Kill()
		return newErrProcess("Interrupted by external request", command, commandLine, nil)
	}

	return nil
}

// Checks for successful execution. Logs all output at default level. Writes
// the supplied string to stdin of the process.
func CheckExecWithInput(input string, command string, commandLine ...string) error {
	//log.Debugln("Executing Command:", command, commandLine)
	cmd := exec.Command(command, commandLine...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	//cmd.Stdout = NewLogWriter(log.With("pipe", "stdout").With("cmd", command).Debugln)
	//cmd.Stderr = NewLogWriter(log.With("pipe", "stderr").With("cmd", command).Debugln)

	if err := cmd.Start(); err != nil {
		return err
	}

	// Wait on a go-routine for the process to exit
	doneCh := make(chan error)
	go func() {
		doneCh <- cmd.Wait()
	}()

	stdinWriteCompleteCh := make(chan error)

	go func() {
		_, err := io.WriteString(stdin, input)
		stdinWriteCompleteCh <- err
		close(stdinWriteCompleteCh)
	}()

	// Wait for process exit or global interrupt
	select {
	case err := <-doneCh:
		if err != nil {
			return err
		}
		// Read the results of writing stdin
		if werr := <-stdinWriteCompleteCh ; werr != nil {
			return err
		}
	case <-interruptCh:
		cmd.Process.Kill()
		return newErrProcess("Interrupted by external request", command, commandLine, nil)
	}

	return nil
}

func CheckExecWithInputAndOutput(input string, command string, commandLine ...string) (string, string, error) {
	//log.Debugln("Executing Command:", command, commandLine)
	cmd := exec.Command(command, commandLine...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", "", err
	}
	defer stdin.Close()

	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)

	//cmd.Stdout = io.MultiWriter(stdoutBuffer,
	//	NewLogWriter(log.With("pipe", "stdout").With("cmd", command).Debugln))
	//cmd.Stderr = io.MultiWriter(stderrBuffer,
	//	NewLogWriter(log.With("pipe", "stderr").With("cmd", command).Debugln))

	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	// Wait on a go-routine for the process to exit
	doneCh := make(chan error)
	go func() {
		doneCh <- cmd.Wait()
	}()

	stdinWriteCompleteCh := make(chan error)

	go func() {
		_, err := io.WriteString(stdin, input)
		stdinWriteCompleteCh <- err
		close(stdinWriteCompleteCh)
	}()

	// Wait for process exit or global interrupt
	select {
	case err := <-doneCh:
		if err != nil {
			return stdoutBuffer.String(), stderrBuffer.String(), err
		}
		// Read the results of writing stdin
		if werr := <-stdinWriteCompleteCh ; werr != nil {
			return stdoutBuffer.String(), stderrBuffer.String(), werr
		}
	case <-interruptCh:
		cmd.Process.Kill()
		return stdoutBuffer.String(), stderrBuffer.String(), newErrProcess("Interrupted by external request", command, commandLine, nil)
	}

	return stdoutBuffer.String(), stderrBuffer.String(), nil
}

func MustExecWithOutput(command string, commandLine ...string) (string, string) {
	stdout, stderr, err := CheckExecWithOutput(command, commandLine...)
	if err != nil {
		panic(newErrProcess("Cannot continue - command failed:", command, commandLine, nil))
	}
	return stdout, stderr
}

func MustExecWithEnv(env []string, command string, commandLine ...string) {
	err := CheckExecWithEnv(env, command, commandLine...)
	if err != nil {
		panic(newErrProcess("Cannot continue - command failed:", command, commandLine, nil))
	}
}

// Exit program if execution is not successful
func MustExec(command string, commandLine ...string) {
	err := CheckExec(command, commandLine...)
	if err != nil {
		panic(newErrProcess("Cannot continue - command failed:", command, commandLine, nil))
	}
}
