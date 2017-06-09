package main

import (
	"cloud.google.com/go/translate"
	vision "cloud.google.com/go/vision/apiv1"
	"compress/gzip"
	"github.com/jjacquay712/GoRODS/msi"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"io"
	"log"
	"path/filepath"
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

const API_AUTH_FILE = "/etc/irods/iRODS-UGM-Demo.json"

//export ExtractImageMetadata
func ExtractImageMetadata(imagePath *C.msParam_t, enableGzipInt *C.msParam_t, rei *C.ruleExecInfo_t) int {

	// Setup GoRODS/msi
	msi.Configure(unsafe.Pointer(rei))

	// Convert input to Golang types
	imageFilePath := msi.ToParam(unsafe.Pointer(imagePath)).String()
	enableGzip := msi.ToParam(unsafe.Pointer(enableGzipInt)).Int() > 0

	// Filter out non-images
	validExtensions := []string{".jpg", ".png", ".gif"}
	ext := strings.ToLower(filepath.Ext(imageFilePath))
	if !Contains(validExtensions, ext) {
		return msi.SUCCESS
	}

	log.Printf("msiextract_image_metadata: Extracting metadata for %v", imageFilePath)

	// Extract image metadata via machine learning API
	labelsKVP := GetImageLabels(imageFilePath, API_AUTH_FILE).ToKVP()

	// Extract exif metadata
	exifKVP := ExtractExifData(imageFilePath)

	if enableGzip {
		gzipDataObjPath := imageFilePath + ".gz"
		gzipDesc := msi.NewParam(msi.INT_MS_T)

		if err := msi.Call("msiDataObjCreate", gzipDataObjPath, "", gzipDesc); err != nil {
			log.Print(err)
			return msi.SYS_INTERNAL_ERR
		}

		gzipDataObj := msi.NewObjReaderFromDesc(gzipDesc)
		defer gzipDataObj.Close()

		gzWriter := gzip.NewWriter(gzipDataObj)
		defer gzWriter.Close()

		origDataObj, err := msi.NewObjReader(imageFilePath)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		defer origDataObj.Close()

		io.Copy(gzWriter, origDataObj)

		imageFilePath = gzipDataObjPath
	}

	// Associate metadata to data object
	if err := msi.Call("msiAssociateKeyValuePairsToObj", labelsKVP, imageFilePath, "-d"); err != nil {
		log.Print(err)
		return msi.SYS_INTERNAL_ERR
	}

	if err := msi.Call("msiAssociateKeyValuePairsToObj", exifKVP, imageFilePath, "-d"); err != nil {
		log.Print(err)
		return msi.SYS_INTERNAL_ERR
	}

	return msi.SUCCESS
}

type ExifWalker struct {
	kvpMap map[string]string
}

func NewExifWalker() *ExifWalker {
	w := new(ExifWalker)
	w.kvpMap = make(map[string]string)
	return w
}

func (w *ExifWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	w.kvpMap["exif_"+string(name)] = tag.String()
	return nil
}

func ExtractExifData(rodsPath string) *msi.Param {
	file, err := msi.NewObjReader(rodsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	exif.RegisterParsers(mknote.All...)
	walker := NewExifWalker()

	exifData, err := exif.Decode(file)
	if err != nil {
		log.Print(err)
		return msi.NewParam(msi.KeyValPair_MS_T).SetKVP(walker.kvpMap)
	}

	if wErr := exifData.Walk(walker); wErr != nil {
		log.Print(wErr)
	}

	return msi.NewParam(msi.KeyValPair_MS_T).SetKVP(walker.kvpMap)
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

func (labels ImgLabels) ToKVP() *msi.Param {

	metaMap := make(map[string]string)
	metaMap["tags_english"] = ""
	metaMap["tags_dutch"] = ""

	for _, label := range labels {
		metaMap["tags_english"] += label.english + ","
		metaMap["tags_dutch"] += label.dutch + ","
	}

	metaMap["tags_english"] = strings.TrimRight(metaMap["tags_english"], ",")
	metaMap["tags_dutch"] = strings.TrimRight(metaMap["tags_dutch"], ",")

	return msi.NewParam(msi.KeyValPair_MS_T).SetKVP(metaMap)
}

func (labels ImgLabels) SetEnglish(enLabels []string) {
	for i, label := range enLabels {
		labels[i].english = label
	}
}

var translateClient *translate.Client
var visionClient *vision.ImageAnnotatorClient

func GetImageLabels(filepath, apiAuthFile string) (imageLabels ImgLabels) {

	authOption := option.WithServiceAccountFile(apiAuthFile)

	ctx := context.Background()

	// Lazy load single client
	if visionClient == nil {

		// Creates a client.
		client, err := vision.NewImageAnnotatorClient(ctx, authOption)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		visionClient = client
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

	labels, err := visionClient.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	enLabels := make([]string, len(labels))
	for i, label := range labels {
		enLabels[i] = label.Description
	}

	imageLabels = make(ImgLabels, len(labels))

	imageLabels.SetEnglish(enLabels)

	// Translate the labels into dutch
	imageLabels.FetchTranslations(apiAuthFile)

	return
}

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

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C static archive
}
