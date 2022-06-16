package check

import (
	"crypto/tls"
	"fmt"
	"time"
)

func init() {
	RegistryInstance.Register("webssl", func() Check {
		return &WebSSLCheck{}
	})
}

type WebSSLCheck struct {
	GenericCheck
}

func (c *WebSSLCheck) Run() {
	// validate config parameters
	if err := c.ValidateParameters([]string{"url"}); err != nil {
		return
	}
	url := c.GenericCheck.Parameters["url"]
	warningTime, warningTimeProvided := c.GenericCheck.Parameters["warning_time"]

	if !warningTimeProvided {
		warningTime = "720h"
	}

	warningTimeParsed, warningTimeParseError := time.ParseDuration(warningTime)
	if warningTimeParseError != nil {
		c.GenericCheck.Error = warningTimeParseError
		c.State = StateUnable
		return
	}

	conn, err := tls.Dial("tcp", url, nil)
	if err != nil {
		c.GenericCheck.Error = err
		c.State = StateFailed
		return
	}
	defer conn.Close()

	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
	cn := conn.ConnectionState().PeerCertificates[0].Subject.CommonName
	difference := expiry.Sub(time.Now()).Hours() / 24

	if expiry.Before(time.Now().Add(warningTimeParsed)) {
		c.GenericCheck.Error = fmt.Errorf("certificate %+q expires on %s (%2.0f days)", cn, expiry, difference)
		c.State = StateWarning
		return
	}

	c.GenericCheck.Help = fmt.Sprintf("certificate expires on %s", expiry)
	c.State = StateOK
}
