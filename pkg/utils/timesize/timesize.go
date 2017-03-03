// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package timesize

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/CodisLabs/codis/pkg/utils/errors"
	"github.com/CodisLabs/codis/pkg/utils/log"
)

type Duration time.Duration

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d Duration) Int64() int64 {
	return int64(d)
}

func (d Duration) MarshalText() ([]byte, error) {
	if d == 0 {
		return []byte("0"), nil
	}
	var abs = time.Duration(d)
	if abs < 0 {
		abs = -abs
	}
	var val = time.Duration(d)
	switch {
	case abs%time.Hour == 0:
		return []byte(fmt.Sprintf("%dh", int64(val/time.Hour))), nil
	case abs%time.Minute == 0:
		return []byte(fmt.Sprintf("%dm", int64(val/time.Minute))), nil
	case abs%time.Second == 0:
		return []byte(fmt.Sprintf("%ds", int64(val/time.Second))), nil
	case abs%time.Millisecond == 0:
		return []byte(fmt.Sprintf("%dms", int64(val/time.Millisecond))), nil
	case abs%time.Microsecond == 0:
		return []byte(fmt.Sprintf("%dus", int64(val/time.Microsecond))), nil
	default:
		return []byte(fmt.Sprintf("%s", val)), nil
	}
}

func (p *Duration) Set(t time.Duration) {
	*p = Duration(t)
}

func (p *Duration) UnmarshalText(text []byte) error {
	n, err := Parse(string(text))
	if err != nil {
		return err
	}
	*p = Duration(n)
	return nil
}

var (
	fullRegexp = regexp.MustCompile(`^\s*(\-?[\d\.]+)\s*([a-z]+|)\s*$`)
	digitsOnly = regexp.MustCompile(`^\-?\d+$`)
)

var ErrBadTimeSize = errors.New("invalid timesize")

func Parse(s string) (time.Duration, error) {
	if !fullRegexp.MatchString(s) {
		return 0, errors.Trace(ErrBadTimeSize)
	}

	subs := fullRegexp.FindStringSubmatch(s)
	if len(subs) != 3 {
		return 0, errors.Trace(ErrBadTimeSize)
	}

	text := subs[1]
	unit := subs[2]

	switch {
	case unit != "":
		return time.ParseDuration(text + unit)
	case digitsOnly.MatchString(text):
		n, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return 0, errors.Trace(ErrBadTimeSize)
		}
		n *= int64(time.Second)
		return time.Duration(n), nil
	default:
		n, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return 0, errors.Trace(ErrBadTimeSize)
		}
		n *= float64(time.Second)
		return time.Duration(n), nil
	}
}

func MustParse(s string) time.Duration {
	v, err := Parse(s)
	if err != nil {
		log.PanicError(err, "parse timesize failed")
	}
	return v
}
