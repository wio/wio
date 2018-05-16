package log

import "os"

// Special function to be used when printing error logs.
// It terminates the program after printing the logs
func Error(newLine bool, format string, a ...interface{}) {
    if !Norm.status || !Norm.verbose {
        Norm.Red(true, "Turn Verbose mode to see the detailed error")
        os.Exit(3)
    }

    Norm.Red(newLine, "Error Report: ")
    Norm.Write(newLine, format, a...)
    os.Exit(2)
}
