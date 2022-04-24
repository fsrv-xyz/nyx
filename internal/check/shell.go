package check

import (
    "context"
    "fmt"
    "os/exec"
    "time"
)

func init() {
    RegistryInstance.Register("shell", func() Check {
        return &ShellCheck{}
    })
}

type ShellCheck struct {
    GenericCheck
}

func (c *ShellCheck) Run() {
    // validate config parameters
    if err := c.ValidateParameters([]string{"command"}); err != nil {
        return
    }
    command := c.GenericCheck.Parameters["command"]

    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()
    args := append([]string{}, "-c")
    args = append(args, fmt.Sprintf("%s", command))
    cmd := exec.CommandContext(ctx, "/bin/sh", args...)

    err := cmd.Run()
    c.GenericCheck.Error = err
    if err != nil {
        c.State = StateFailed
        return
    }
    c.State = StateOK
}
