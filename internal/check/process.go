package check

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	if err := c.ValidateParameters([]string{"pidfile", "match"}); err != nil {
		return
	}
	parameterPidFile := c.GenericCheck.Parameters["pidfile"]
	parameterMatch := c.GenericCheck.Parameters["match"]

	pidFileContent, err := os.ReadFile(parameterPidFile)
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

	// try to get process name; fail if not possible
	processName, gatherProcessNameError := processWithPid.Name()
	if gatherProcessNameError != nil {
		c.GenericCheck.State = StateFailed
		c.GenericCheck.Error = gatherProcessNameError
		return
	}

	matcher, matchCompileError := regexp.Compile(parameterMatch)
	if matchCompileError != nil {
		c.GenericCheck.State = StateUnable
		c.GenericCheck.Error = matchCompileError
		return
	}

	// compare process name to match parameter
	if !matcher.MatchString(processName) {
		c.GenericCheck.State = StateWarning
		c.GenericCheck.Error = fmt.Errorf("process name %+q does not match %+q", processName, matcher)
		return
	}
	c.GenericCheck.State = StateOK
}
