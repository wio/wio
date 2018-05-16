package commands

import (
    "github.com/go-errors/errors"
    "os"
    "wio/cmd/wio/utils/io/log"
    "github.com/urfave/cli"
)

type Command interface {
    GetContext() (*cli.Context)
    Execute()
}

func RecordError(err error, message string) {
    if err != nil {
        if message != "" {
            log.Norm.Red(true, message)
        }

        log.Error(true, err.(*errors.Error).ErrorStack())
        os.Exit(2)
    }
}
