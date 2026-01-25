package testutil

import (
	"reflect"
	"testing"
)

// AssertEqual 断言相等
func AssertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// AssertNotNil 断言不为 nil
func AssertNotNil(t *testing.T, obj interface{}) {
	t.Helper()
	if obj == nil || reflect.ValueOf(obj).IsNil() {
		t.Error("expected non-nil value")
	}
}

// AssertNil 断言为 nil
func AssertNil(t *testing.T, obj interface{}) {
	t.Helper()
	if obj != nil && !reflect.ValueOf(obj).IsNil() {
		t.Errorf("expected nil, got %v", obj)
	}
}

// AssertError 断言有错误
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// AssertNoError 断言无错误
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
