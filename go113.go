// +build go1.13

package errors

import (
	stderrors "errors"
)

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool { return stderrors.As(err, target) }

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// UnwrapWithOuter returns the inner and outer errors.
// If err wraps another error, it returns the wrapped error as inner
// and the original err as outer.
func UnwrapWithOuter(err error) (inner, outer error) {
	if err == nil {
		return nil, nil
	}
	
	current := err
	for {
		// Try standard library unwrap
		next := stderrors.Unwrap(current)
		if next == nil {
			// Fallback to older causer interface used by pkg/errors types
			// if they don't implement Go 1.13 Unwrap()
			type causer interface {
				Cause() error
			}
			if c, ok := current.(causer); ok {
				next = c.Cause()
			}
		}
		if next == nil {
			break
		}
		current = next
	}
	
	if current == err {
		// The error does not wrap anything
		return nil, err
	}

	msg := err.Error()
	innerMsg := current.Error()

	if len(msg) > len(innerMsg) {
		suffix := ": " + innerMsg
		if len(msg) >= len(suffix) && msg[len(msg)-len(suffix):] == suffix {
			msg = msg[:len(msg)-len(suffix)]
		} else if len(msg) >= len(innerMsg) && msg[len(msg)-len(innerMsg):] == innerMsg {
			msg = msg[:len(msg)-len(innerMsg)]
		}
	} else if msg == innerMsg {
		msg = ""
	}

	return current, stderrors.New(msg)
}
