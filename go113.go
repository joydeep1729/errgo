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

// UnwrapWithOuter decomposes a wrapped error into the immediate wrapped inner error
// and the outer context/message added at the current level.
//
// If err is nil, it returns (nil, nil). If err does not wrap another error, it
// returns (nil, err). Otherwise, it returns the inner error and a new error
// representing only the annotation message added at this level.
func UnwrapWithOuter(err error) (inner, outer error) {
	if err == nil {
		return nil, nil
	}

	// Try Go 1.13+ Unwrap first, then fallback to legacy pkg/errors causer.
	inner = stderrors.Unwrap(err)
	if inner == nil {
		type causer interface {
			Cause() error
		}
		if c, ok := err.(causer); ok {
			inner = c.Cause()
		}
	}

	if inner == nil {
		return nil, err
	}

	msg := err.Error()
	innerMsg := inner.Error()

	// Extract the annotation message added at this specific level by stripping
	// the inner error's message from the outer error's message.
	var outerMsg string
	if len(msg) > len(innerMsg) {
		suffix := ": " + innerMsg
		if len(msg) >= len(suffix) && msg[len(msg)-len(suffix):] == suffix {
			outerMsg = msg[:len(msg)-len(suffix)]
		} else if len(msg) >= len(innerMsg) && msg[len(msg)-len(innerMsg):] == innerMsg {
			outerMsg = msg[:len(msg)-len(innerMsg)]
		} else {
			outerMsg = msg
		}
	} else if msg == innerMsg {
		outerMsg = ""
	} else {
		outerMsg = msg
	}

	return inner, stderrors.New(outerMsg)
}

// UnwrapToCauseWithOuter unwraps the error chain completely down to the root cause.
//
// If err is nil, it returns (nil, nil). If err does not wrap another error, it
// returns (nil, err). Otherwise, it returns the root cause and a new error
// representing the consolidated messages of all wrapping layers.
func UnwrapToCauseWithOuter(err error) (cause, outer error) {
	if err == nil {
		return nil, nil
	}

	current := err
	for {
		next := stderrors.Unwrap(current)
		if next == nil {
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
		return nil, err
	}

	msg := err.Error()
	innerMsg := current.Error()

	// Extract the combined wrapper messages by stripping the root cause.
	var outerMsg string
	if len(msg) > len(innerMsg) {
		suffix := ": " + innerMsg
		if len(msg) >= len(suffix) && msg[len(msg)-len(suffix):] == suffix {
			outerMsg = msg[:len(msg)-len(suffix)]
		} else if len(msg) >= len(innerMsg) && msg[len(msg)-len(innerMsg):] == innerMsg {
			outerMsg = msg[:len(msg)-len(innerMsg)]
		} else {
			outerMsg = msg
		}
	} else if msg == innerMsg {
		outerMsg = ""
	} else {
		outerMsg = msg
	}

	return current, stderrors.New(outerMsg)
}
