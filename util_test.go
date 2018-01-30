package main

import (
	"testing"
)

type StructName struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
	Ack string `json:"ack"`
}

func TestReflectJsonTags(t *testing.T) {
	result := lookupStructTags(StructName{})
	if len(result) != 3 {
		t.Errorf("result is %d not 3", len(result))
	}
}

func TestReflectStructName(t *testing.T) {
	result := lookupStructName(StructName{})
	if result != "StructName" {
		t.Errorf("expected StructName but got %s", result)
	}
}
