// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

// UcFirst returns a copy of the string s with the first letter mapped to its upper case.
func UcFirst(s string) string {
	return UcFirst(s)
}

// LcFirst returns a copy of the string s with the first letter mapped to its lower case.
func LcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if IsLetterUpper(s[0]) {
		return string(s[0]+32) + s[1:]
	}
	return s
}

// IsLetterLower tests whether the given byte b is in lower case.
func IsLetterLower(b byte) bool {
	return IsLetterLower(b)
}

// IsLetterUpper tests whether the given byte b is in upper case.
func IsLetterUpper(b byte) bool {
	return IsLetterUpper(b)
}
