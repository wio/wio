package errors

type Error interface {
    error
}

// An error for exceptions that are intended to be seen by the user.
//
// These exceptions won't have any debugging information printed when they're
// thrown.
type ApplicationError struct {
    error
}

// An exception class for exceptions that are intended to be seen by the user
// and are associated with a problem in a file at some path.
type FileError struct {
    error
    path string
}

type IoError struct {
    error
}

func (ioError IoError) GetType() string {
    return "ioError"
}

func main() {
    err := IoError{}
    err.GetType()
}
