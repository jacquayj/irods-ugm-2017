package main

import (
	"github.com/jjacquay712/GoRODS/msi"
	"strings"
	"testing"
)

func TestBasicExample(t *testing.T) {

	inputParam := msi.NewParam(msi.STR_MS_T).SetString("keytest,valtest")
	outputParam := msi.NewParam(msi.KeyValPair_MS_T)

	if status := BasicExample(UnsafePtrToC(inputParam.Ptr()), UnsafePtrToC(outputParam.Ptr()), nil); status == 0 {

		kvpStr := outputParam.String()

		if strings.Contains(kvpStr, "keytest = valtest") {
			t.Log("Success!")
		} else {
			t.Fatalf("'%v' != '%v'", kvpStr, "keytest = valtest")
		}

	} else {
		t.Fatalf("Status = %v not 0", status)
	}
}
