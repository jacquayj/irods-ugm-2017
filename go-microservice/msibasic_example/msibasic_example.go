package main

import (
	"encoding/csv"
	"github.com/jjacquay712/GoRODS/msi"
	"io"
	"log"
	"strings"
	"unsafe"
)

// #cgo CFLAGS: -I/usr/include/irods
// #cgo LDFLAGS: -lirods_server -lirods_common -lpthread
/*
#include "msParam.h"
#include "re_structs.h"
*/
import "C"

//export BasicExample
func BasicExample(inputParam *C.msParam_t, outputParam *C.msParam_t, rei *C.ruleExecInfo_t) int {

	// Setup GoRODS/msi
	msi.Configure(unsafe.Pointer(rei))

	// Convert *C.msParam_t to golang types
	inputCSV := msi.ToParam(unsafe.Pointer(inputParam)).String()
	outputKVP := msi.ToParam(unsafe.Pointer(outputParam)).ConvertTo(msi.KeyValPair_MS_T)

	// Set output KVP
	outputKVP.SetKVP(GetKVPMap(inputCSV))

	return msi.SUCCESS
}

// Parse CSV and return a map of key value pairs
func GetKVPMap(csvStr string) map[string]string {
	kvpMap := make(map[string]string)

	reader := csv.NewReader(strings.NewReader(csvStr))
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		kvpMap[record[0]] = record[1]
	}

	return kvpMap
}

// Used in testing, since C package can't be imported to test files
func UnsafePtrToC(ptr unsafe.Pointer) *C.msParam_t {
	return (*C.msParam_t)(ptr)
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C static archive
}
