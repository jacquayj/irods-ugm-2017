package main

import (
	vision "cloud.google.com/go/vision/apiv1"
	"github.com/jjacquay712/GoRODS/msi"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"unsafe"
)

// #cgo CFLAGS: -I/usr/include/irods
// #cgo LDFLAGS: -lirods_server -lirods_common -lpthread
/*
#include "msParam.h"
#include "re_structs.h"
*/
import "C"

//export ProcessAlbumDirectory
func ProcessAlbumDirectory(imagePath *C.msParam_t, APIAuthPath *C.msParam_t, rei *C.ruleExecInfo_t) int {
	msi.Configure(unsafe.Pointer(rei))

	imageFilePath := msi.ToParam(unsafe.Pointer(imagePath)).String()
	apiAuthFile := msi.ToParam(unsafe.Pointer(APIAuthPath)).String()

	GetImageLabels(imageFilePath, apiAuthFile)

	return 0
}

func GetImageLabels(filepath string, apiAuthFile string) {

	authOption := option.WithServiceAccountFile(apiAuthFile)

	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx, authOption)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	file, err := msi.NewObjReader(filepath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	labels, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	log.Println("Labels:")
	for _, label := range labels {
		log.Println(label.Description)
	}
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C static archive
}
