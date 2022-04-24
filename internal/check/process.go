package check

import (
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/v3/process"
)

func init() {
	RegistryInstance.Register("process", func() Check {
		return &ProcessCheck{}
	})
}

type ProcessCheck struct {
	GenericCheck
}

func (c *ProcessCheck) Run() {
	// validate config parameters
	if err := c.ValidateParameters([]string{"pidfile"}); err != nil {
		return
	}
	pidFile := c.GenericCheck.Parameters["pidfile"]
	pidFileContent, err := os.ReadFile(pidFile)
	if err != nil {
		c.GenericCheck.State = StateFailed
		c.GenericCheck.Error = err
		return
	}
	// format pid file content to ensure correct int32 parsing
	pidFileContentString := strings.TrimSuffix(string(pidFileContent), "\n")

	// parse pid string to int32; fail if not possible
	pid, pidConvertErr := strconv.ParseInt(pidFileContentString, 10, 32)
	if pidConvertErr != nil {
		c.GenericCheck.State = StateFailed
		c.GenericCheck.Error = pidConvertErr
		return
	}
	// create process instance; fail if not possible
	processWithPid, processWithPidError := process.NewProcess(int32(pid))
	if processWithPidError != nil {
		c.GenericCheck.State = StateFailed
		c.GenericCheck.Error = processWithPidError
		return
	}

	// try sending SIG-NULL signal
	processSignalError := processWithPid.SendSignal(syscall.Signal(0))
	if processSignalError != nil {
		c.GenericCheck.State = StateFailed
		c.GenericCheck.Error = processSignalError
		return
	}
	c.GenericCheck.State = StateOK
}
