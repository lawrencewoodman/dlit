package dlit

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		in        interface{}
		want      *Literal
		wantError error
	}{
		{6, makeLit(6), nil},
		{6.0, makeLit(6.0), nil},
		{6.6, makeLit(6.6), nil},
		{float32(6.6), makeLit(float32(6.6)), nil},
		{int64(922336854775807), makeLit(922336854775807), nil},
		{int64(9223372036854775807), makeLit(9223372036854775807), nil},
		{"98292223372036854775807", makeLit("98292223372036854775807"), nil},
		{complex64(1), makeLit(ErrInvalidKind("complex64")),
			ErrInvalidKind("complex64")},
		{complex128(1), makeLit(ErrInvalidKind("complex128")),
			ErrInvalidKind("complex128")},
		{"6", makeLit("6"), nil},
		{"6.6", makeLit("6.6"), nil},
		{"abc", makeLit("abc"), nil},
		{true, makeLit(true), nil},
		{false, makeLit(false), nil},
		{errors.New("This is an error"), makeLit(errors.New("This is an error")),
			nil},
	}

	for _, c := range cases {
		got, err := New(c.in)
		if !errorMatch(err, c.wantError) {
			t.Errorf("New(%q) - err == %q, wantError == %q", c.in, err, c.wantError)
		}

		if got.String() != c.want.String() {
			t.Errorf("New(%q) - got == %s, want == %s", c.in, got, c.want)
		}
	}
}

func TestMustNew(t *testing.T) {
	cases := []struct {
		in   interface{}
		want *Literal
	}{
		{6, makeLit(6)},
		{6.0, makeLit(6.0)},
		{6.6, makeLit(6.6)},
		{float32(6.6), makeLit(float32(6.6))},
		{int64(922336854775807), makeLit(922336854775807)},
		{int64(9223372036854775807), makeLit(9223372036854775807)},
		{"98292223372036854775807", makeLit("98292223372036854775807")},
		{"6", makeLit("6")},
		{"6.6", makeLit("6.6")},
		{"abc", makeLit("abc")},
		{true, makeLit(true)},
		{false, makeLit(false)},
		{errors.New("This is an error"), makeLit(errors.New("This is an error"))},
	}

	for _, c := range cases {
		got := MustNew(c.in)
		if got.String() != c.want.String() {
			t.Errorf("MustNew(%q) - got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func TestMustNew_panic(t *testing.T) {
	cases := []struct {
		in        interface{}
		wantPanic string
	}{
		{6, ""},
		{complex64(1), ErrInvalidKind("complex64").Error()},
	}

	for _, c := range cases {
		paniced := false
		defer func() {
			if r := recover(); r != nil {
				if r.(string) == c.wantPanic {
					paniced = true
				} else {
					t.Errorf("MustNew(%q) - got panic: %s, wanted: %s",
						c.in, r, c.wantPanic)
				}
			}
		}()
		MustNew(c.in)
		if c.wantPanic != "" && !paniced {
			t.Errorf("MustNew(%q) - failed to panic with: %s", c.in, c.wantPanic)
		}
	}
}

func TestInt(t *testing.T) {
	cases := []struct {
		in        *Literal
		want      int64
		wantIsInt bool
	}{
		{makeLit(6), 6, true},
		{makeLit(6.0), 6, true},
		{makeLit(float32(6.0)), 6, true},
		{makeLit("6"), 6, true},
		{makeLit("6.0"), 6, true},
		{makeLit(fmt.Sprintf("%d", int64(math.MinInt64))),
			int64(math.MinInt64), true},
		{makeLit(fmt.Sprintf("%d", int64(math.MaxInt64))),
			int64(math.MaxInt64), true},
		{makeLit(fmt.Sprintf("-1%d", int64(math.MinInt64))), 0, false},
		{makeLit(fmt.Sprintf("1%d", int64(math.MaxInt64))), 0, false},
		{makeLit("-9223372036854775809"), 0, false},
		{makeLit("9223372036854775808"), 0, false},
		{makeLit(6.6), 0, false},
		{makeLit("6.6"), 0, false},
		{makeLit("abc"), 0, false},
		{makeLit(true), 0, false},
		{makeLit(false), 0, false},
		{makeLit(errors.New("This is an error")), 0, false},
	}

	for _, c := range cases {
		got, gotIsInt := c.in.Int()
		if got != c.want || gotIsInt != c.wantIsInt {
			t.Errorf("Int() with Literal: %q - return: %q, %q - want: %q, %q",
				c.in, got, gotIsInt, c.want, c.wantIsInt)
		}
	}
}

func TestFloat(t *testing.T) {
	cases := []struct {
		in          *Literal
		want        float64
		wantIsFloat bool
	}{
		{makeLit(6), 6.0, true},
		{makeLit(int64(922336854775807)), 922336854775807.0, true},
		{makeLit(fmt.Sprintf("%G", float64(math.SmallestNonzeroFloat64))),
			math.SmallestNonzeroFloat64, true},
		{makeLit(fmt.Sprintf("%f", float64(math.MaxFloat64))),
			float64(math.MaxFloat64), true},
		{makeLit(6.0), 6.0, true},
		{makeLit("6"), 6.0, true},
		{makeLit(6.678934), 6.678934, true},
		{makeLit("6.678394"), 6.678394, true},
		{makeLit("abc"), 0, false},
		{makeLit(true), 0, false},
		{makeLit(false), 0, false},
		{makeLit(errors.New("This is an error")), 0, false},
	}

	for _, c := range cases {
		got, gotIsFloat := c.in.Float()
		if got != c.want || gotIsFloat != c.wantIsFloat {
			t.Errorf("Float() with Literal: %q - return: %q, %q - want: %q, %q",
				c.in, got, gotIsFloat, c.want, c.wantIsFloat)
		}
	}
}

func TestBool(t *testing.T) {
	cases := []struct {
		in         *Literal
		want       bool
		wantIsBool bool
	}{
		{makeLit(1), true, true},
		{makeLit(2), false, false},
		{makeLit(0), false, true},
		{makeLit(1.0), true, true},
		{makeLit(2.0), false, false},
		{makeLit(2.25), false, false},
		{makeLit(0.0), false, true},
		{makeLit(true), true, true},
		{makeLit(false), false, true},
		{makeLit("true"), true, true},
		{makeLit("false"), false, true},
		{makeLit("True"), true, true},
		{makeLit("False"), false, true},
		{makeLit("TRUE"), true, true},
		{makeLit("FALSE"), false, true},
		{makeLit("t"), true, true},
		{makeLit("f"), false, true},
		{makeLit("T"), true, true},
		{makeLit("F"), false, true},
		{makeLit("1"), true, true},
		{makeLit("0"), false, true},
		{makeLit("bob"), false, false},
		{makeLit(errors.New("This is an error")), false, false},
	}

	for _, c := range cases {
		got, gotIsBool := c.in.Bool()
		if got != c.want || gotIsBool != c.wantIsBool {
			t.Errorf("Bool() with Literal: %q - return: %q, %q - want: %q, %q",
				c.in, got, gotIsBool, c.want, c.wantIsBool)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		in   *Literal
		want string
	}{
		{makeLit(124), "124"},
		{makeLit(int64(922336854775807)), "922336854775807"},
		{makeLit(int64(9223372036854775807)), "9223372036854775807"},
		{makeLit("98292223372036854775807"), "98292223372036854775807"},
		{makeLit("Hello my name is fred"), "Hello my name is fred"},
		{makeLit(124.0), "124"},
		{makeLit(124.56728482274629), "124.56728482274629"},
		{makeLit(true), "true"},
		{makeLit(false), "false"},
		{makeLit(errors.New("This is an error")), "This is an error"},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("String() with Literal: %q - return: %q, want: %q",
				c.in, got, c.want)
		}
	}
}

func TestErr(t *testing.T) {
	cases := []struct {
		in        *Literal
		want      error
		wantIsErr bool
	}{
		{makeLit(1), nil, false},
		{makeLit(2), nil, false},
		{makeLit("true"), nil, false},
		{makeLit(2.25), nil, false},
		{makeLit("hello"), nil, false},
		{makeLit(errors.New("This is an error")),
			errors.New("This is an error"), true},
	}

	for _, c := range cases {
		got, gotIsErr := c.in.Err()
		if !errorMatch(c.want, got) || gotIsErr != c.wantIsErr {
			t.Errorf("Err() with Literal: %q - return: %q, %q - want: %q, %q",
				c.in, got, gotIsErr, c.want, c.wantIsErr)
		}
	}
}

/***********************
   Helper functions
************************/
func makeLit(v interface{}) *Literal {
	l, err := New(v)
	if err != nil {
		panic(fmt.Sprintf("MakeLit(%q) gave err: %q", v, err))
	}
	return l
}

func errorMatch(e1 error, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	if e1.Error() == e2.Error() {
		return true
	}
	return false
}
