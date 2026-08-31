package utility

import (
	"errors"
	"testing"
)

func TestAssertReturnsValue(t *testing.T) {
	result := assert(42, nil)
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestAssertReturnsString(t *testing.T) {
	result := assert("hello", nil)
	if result != "hello" {
		t.Errorf("expected 'hello', got '%s'", result)
	}
}

func TestAssertPanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("assert did not panic when given an error")
		}
	}()
	assert(0, errors.New("test error"))
}

func TestAssertPanicsWithCorrectError(t *testing.T) {
	errMsg := "expected error"
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("assert did not panic")
		}
		err, ok := r.(error)
		if !ok {
			t.Fatal("panic value is not an error")
		}
		if err.Error() != errMsg {
			t.Errorf("expected panic with '%s', got '%v'", errMsg, err)
		}
	}()
	assert(0, errors.New(errMsg))
}

func TestAssertWithStruct(t *testing.T) {
	type testStruct struct {
		A int
		B string
	}
	input := testStruct{A: 1, B: "test"}
	result := assert(input, nil)
	if result != input {
		t.Errorf("expected %+v, got %+v", input, result)
	}
}

func TestDivmod(t *testing.T) {
	tests := []struct {
		name        string
		numerator   int64
		denominator int64
		wantQ       int64
		wantR       int64
	}{
		{name: "positive divisor divisible", numerator: 10, denominator: 2, wantQ: 5, wantR: 0},
		{name: "positive divisor remainder", numerator: 10, denominator: 3, wantQ: 3, wantR: 1},
		{name: "zero numerator", numerator: 0, denominator: 5, wantQ: 0, wantR: 0},
		{name: "negative numerator", numerator: -10, denominator: 3, wantQ: -4, wantR: 2},
		{name: "negative numerator exact", numerator: -10, denominator: 5, wantQ: -2, wantR: 0},
		{name: "negative denominator", numerator: 10, denominator: -3, wantQ: -4, wantR: -2},
		{name: "negative numerator negative denominator", numerator: -10, denominator: -3, wantQ: 3, wantR: -1},
		{name: "negative denominator exact", numerator: 10, denominator: -5, wantQ: -2, wantR: 0},
		{name: "large values", numerator: 1_000_000_000_000, denominator: 7, wantQ: 142857142857, wantR: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, r := Divmod(tt.numerator, tt.denominator)
			if q != tt.wantQ {
				t.Errorf("Divmod(%d, %d) quotient = %d, want %d", tt.numerator, tt.denominator, q, tt.wantQ)
			}
			if r != tt.wantR {
				t.Errorf("Divmod(%d, %d) remainder = %d, want %d", tt.numerator, tt.denominator, r, tt.wantR)
			}
			if tt.denominator > 0 && (r < 0 || r >= tt.denominator) {
				t.Errorf("Divmod(%d, %d) remainder %d out of range [0, %d)", tt.numerator, tt.denominator, r, tt.denominator)
			}
			if tt.denominator < 0 && (r > 0 || r <= tt.denominator) {
				t.Errorf("Divmod(%d, %d) remainder %d out of range (%d, 0]", tt.numerator, tt.denominator, r, tt.denominator)
			}
			if got := q*tt.denominator + r; got != tt.numerator {
				t.Errorf("invariant violated: %d*%d + %d = %d, want %d", q, tt.denominator, r, got, tt.numerator)
			}
		})
	}
}
