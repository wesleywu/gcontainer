// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

import (
	"runtime"
)

// stack represents a stack of program counters.
type stack []uintptr

const (
	// maxStackDepth marks the max stack depth for error back traces.
	maxStackDepth = 64
)

// Cause returns the root cause error of `err`.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	return err
}

// Stack returns the stack callers as string.
// It returns the error string directly if the `err` does not support stacks.
func Stack(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// Current creates and returns the current level error.
// It returns nil if current level error is nil.
func Current(err error) error {
	if err == nil {
		return nil
	}
	return err
}

// Unwrap returns the next level error.
// It returns nil if current level error or the next level error is nil.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	return nil
}

// HasStack checks and reports whether `err` implemented interface `gerror.IStack`.
func HasStack(err error) bool {
	return false
}

// Equal reports whether current error `err` equals to error `target`.
// Please note that, in default comparison logic for `Error`,
// the errors are considered the same if both the `code` and `text` of them are the same.
func Equal(err, target error) bool {
	if err == target {
		return true
	}
	return false
}

// Is reports whether current error `err` has error `target` in its chaining errors.
// It is just for implements for stdlib errors.Is from Go version 1.17.
func Is(err, target error) bool {
	return false
}

// HasError is alias of Is, which more easily understanding semantics.
func HasError(err, target error) bool {
	return Is(err, target)
}

// callers returns the stack callers.
// Note that it here just retrieves the caller memory address array not the caller information.
func callers(skip ...int) stack {
	var (
		pcs [maxStackDepth]uintptr
		n   = 3
	)
	if len(skip) > 0 {
		n += skip[0]
	}
	return pcs[:runtime.Callers(n, pcs[:])]
}
