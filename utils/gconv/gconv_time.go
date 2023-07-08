// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wesleywu/gcontainer/internal/utils"
	"github.com/wesleywu/gcontainer/utils/gerror"
	"github.com/wesleywu/gcontainer/utils/gregex"
)

// Time converts `any` to time.Time.
func Time(any interface{}, format ...string) time.Time {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(time.Time); ok {
			return v
		}
	}
	return time.Time{}
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func Duration(any interface{}) time.Duration {
	// It's already this type.
	if v, ok := any.(time.Duration); ok {
		return v
	}
	s := String(any)
	if !utils.IsNumeric(s) {
		d, _ := ParseDuration(s)
		return d
	}
	return time.Duration(Int64(any))
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h", "1d" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h", "d".
//
// Very note that it supports unit "d" more than function time.ParseDuration.
func ParseDuration(s string) (duration time.Duration, err error) {
	var (
		num int64
	)
	if utils.IsNumeric(s) {
		num, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			err = gerror.Wrapf(err, `strconv.ParseInt failed for string "%s"`, s)
			return 0, err
		}
		return time.Duration(num), nil
	}
	match, err := gregex.MatchString(`^([\-\d]+)[dD](.*)$`, s)
	if err != nil {
		return 0, err
	}
	if len(match) == 3 {
		num, err = strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			err = gerror.Wrapf(err, `strconv.ParseInt failed for string "%s"`, match[1])
			return 0, err
		}
		s = fmt.Sprintf(`%dh%s`, num*24, match[2])
		duration, err = time.ParseDuration(s)
		if err != nil {
			err = gerror.Wrapf(err, `time.ParseDuration failed for string "%s"`, s)
		}
		return
	}
	duration, err = time.ParseDuration(s)
	err = gerror.Wrapf(err, `time.ParseDuration failed for string "%s"`, s)
	return
}
