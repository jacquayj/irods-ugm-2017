package main

import (
	"cloud.google.com/go/translate"
	vision "cloud.google.com/go/vision/apiv1"
	"github.com/jjacquay712/GoRODS/msi"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
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

//export ExtractImageMetadata
func ExtractImageMetadata(imagePath *C.msParam_t, APIAuthPath *C.msParam_t, rei *C.ruleExecInfo_t) int {
	msi.Configure(unsafe.Pointer(rei))

	// Convert input to Golang strings
	imageFilePath := msi.ToParam(unsafe.Pointer(imagePath)).String()
	apiAuthFile := msi.ToParam(unsafe.Pointer(APIAuthPath)).String()

	// Extract image labels via machine learning API
	enLabels := GetImageLabels(imageFilePath, apiAuthFile)

	// Translate the labels into dutch
	labels := make(ImgLabels, len(enLabels))
	for i, label := range enLabels {
		labels[i].english = label
	}

	labels.FetchTranslations(apiAuthFile)

	metaKVPs := msi.NewParam(msi.KeyValPair_MS_T)

	metaMap := make(map[string]string)
	metaMap["tags_english"] = ""
	metaMap["tags_dutch"] = ""

	for _, label := range labels {
		metaMap["tags_english"] += label.english + ","
		metaMap["tags_dutch"] += label.dutch + ","
	}

	metaMap["tags_english"] = strings.TrimRight(metaMap["tags_english"], ",")
	metaMap["tags_dutch"] = strings.TrimRight(metaMap["tags_dutch"], ",")

	metaKVPs.SetKVP(metaMap)

	if err := msi.Call("msiAssociateKeyValuePairsToObj", metaKVPs, imageFilePath, "-d"); err != nil {
		log.Print(err)
		return -1
	}

	return 0
}

type ImgLabel struct {
	english string
	dutch   string
}

type ImgLabels []ImgLabel

func (labels ImgLabels) FetchTranslations(apiAuthFile string) {
	words := make([]string, len(labels))
	for i, label := range labels {
		words[i] = label.english
	}

	translations := TranslateString("nl", words, apiAuthFile)
	for i, translation := range translations {
		labels[i].dutch = translation
	}
}

func GetImageLabels(filepath, apiAuthFile string) (imageLabels []string) {

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

	for _, label := range labels {
		imageLabels = append(imageLabels, label.Description)
	}

	return
}

var translateClient *translate.Client

func TranslateString(targetLang string, words []string, apiAuthFile string) []string {

	authOption := option.WithServiceAccountFile(apiAuthFile)

	ctx := context.Background()

	// Lazy load single client
	if translateClient == nil {

		// Creates a client.
		client, err := translate.NewClient(ctx, authOption)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		translateClient = client
	}

	// Sets the target language.
	target, err := language.Parse(targetLang)
	if err != nil {
		log.Fatalf("Failed to parse target language: %v", err)
	}

	translations, err := translateClient.Translate(ctx, words, target, nil)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	var responseTranslations []string

	for _, translation := range translations {
		responseTranslations = append(responseTranslations, translation.Text)
	}

	return responseTranslations
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C static archive
}
