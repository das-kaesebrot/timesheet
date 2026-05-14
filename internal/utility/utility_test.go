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
