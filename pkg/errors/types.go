package errors

// causer is an interface implemented by an error that supports wrapping. Use
// it to test if an error wraps another error instance.
type causer interface {
	Cause() error
}

// unpacker is an interface implemented by an error that supports unpack function.
type unpacker interface {
	Unpack() []error
}
