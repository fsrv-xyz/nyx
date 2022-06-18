package check

import (
	"fmt"
	"net"
	"time"
)

func init() {
	RegistryInstance.Register("port", func() Check {
		return &PortCheck{}
	})
}

type PortCheck struct {
	GenericCheck
}

func (c *PortCheck) Run() {
	// validate config parameters
	if err := c.ValidateParameters([]string{"port"}); err != nil {
		return
	}
	parameterAddress, parameterAddressProvided := c.GenericCheck.Parameters["address"]
	parameterTimeout, parameterTimeoutProvided := c.GenericCheck.Parameters["timeout"]
	parameterPort := c.GenericCheck.Parameters["port"]

	if !parameterTimeoutProvided {
		parameterTimeout = "10ms"
	}

	// if address is not provided, use localhost
	if !parameterAddressProvided {
		parameterAddress = "localhost"
	}

	timeout, parseTimeoutError := time.ParseDuration(parameterTimeout)
	if parseTimeoutError != nil {
		c.GenericCheck.State = StateFailed
		c.GenericCheck.Error = fmt.Errorf("unable to parse timeout parameter: %+q\n", parseTimeoutError.Error())
		return
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(parameterAddress, parameterPort), timeout)
	if err != nil {
		c.GenericCheck.State = StateFailed
		c.GenericCheck.Error = err
		return
	}
	defer conn.Close()
	c.State = StateOK
}
