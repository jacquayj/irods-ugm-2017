package main

import (
	"strings"
	"testing"
)

func TestExtractImageMetadata(t *testing.T) {

	inputParam := getMsParam()
	outputParam := getMsParam()

	setMsParamStr(inputParam, "test1234")

	if status := ExtractImageMetadata(inputParam, outputParam, nil); status == 0 {
		output := getMsParamStr(outputParam) // NEED TO FIX

		if strings.Contains(output, "test5678") {
			t.Fatalf("Not expecting string")
		} else {
			t.Log(output)
		}
	} else {
		t.Fatalf("Status = %v not 0", status)
	}
}
