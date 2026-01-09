package turtlebunny

import (
	"testing"
)

func TestIsUint128(t *testing.T) {
	var ok bool
	// negative integer should fail
	ok = isUint128("-123")
	if ok {
		t.Errorf("expected negative integer to fail")
	}

	// decimal should fail
	ok = isUint128("1.23")
	if ok {
		t.Errorf("expected decimal to fail")
	}

	// leading zeros should fail
	ok = isUint128("0123")
	if ok {
		t.Errorf("expected leading zeros to fail")
	}

	// leading spaces should fail
	ok = isUint128("   123")
	if ok {
		t.Errorf("expected leading spaces to fail")
	}

	// trailing spaces should fail
	ok = isUint128("123   ")
	if ok {
		t.Errorf("expected trailing spaces to fail")
	}

	// non numeric should fail
	ok = isUint128("Hello World")
	if ok {
		t.Errorf("expected non-numeric to fail")
	}

	// overflow should fail
	ok = isUint128("999999999999999999999999999999999999999")
	if ok {
		t.Errorf("expected overflow to fail")
	}
}

func TestIsUint64(t *testing.T) {
	var ok bool
	// negative integer should fail
	ok = isUint64("-123")
	if ok {
		t.Errorf("expected negative integer to fail")
	}

	// decimal should fail
	ok = isUint64("1.23")
	if ok {
		t.Errorf("expected decimal to fail")
	}

	// leading zeros should fail
	ok = isUint64("0123")
	if ok {
		t.Errorf("expected leading zeros to fail")
	}

	// leading spaces should fail
	ok = isUint64("   123")
	if ok {
		t.Errorf("expected leading spaces to fail")
	}

	// trailing spaces should fail
	ok = isUint64("123   ")
	if ok {
		t.Errorf("expected trailing spaces to fail")
	}

	// non numeric should fail
	ok = isUint64("Hello World")
	if ok {
		t.Errorf("expected non-numeric to fail")
	}

	// overflow should fail
	ok = isUint64("99999999999999999999")
	if ok {
		t.Errorf("expected overflow to fail")
	}
}
