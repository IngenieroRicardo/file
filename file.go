package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"os"
	"unsafe"
)

//export WBFile
func WBFile(b64Str *C.char, outputPath *C.char) C.int {
	goB64Str := C.GoString(b64Str)
	goOutputPath := C.GoString(outputPath)
	data, err := base64.StdEncoding.DecodeString(goB64Str)
	if err != nil {
		return -1
	}
	err = ioutil.WriteFile(goOutputPath, data, 0644)
	if err != nil {
		return -2
	}
	return 0
}

//export WTFile
func WTFile(textStr *C.char, outputPath *C.char) C.int {
	goTextStr := C.GoString(textStr)
	goOutputPath := C.GoString(outputPath)
	err := ioutil.WriteFile(goOutputPath, []byte(goTextStr), 0644)
	if err != nil {
		return -1
	}
	return 0
}

//export RBFile
func RBFile(inputPath *C.char) *C.char {
	goInputPath := C.GoString(inputPath)
	data, err := ioutil.ReadFile(goInputPath)
	if err != nil {
		return nil
	}
	b64Str := base64.StdEncoding.EncodeToString(data)
	return C.CString(b64Str)
}

//export RTFile
func RTFile(inputPath *C.char) *C.char {
	goInputPath := C.GoString(inputPath)
	data, err := ioutil.ReadFile(goInputPath)
	if err != nil {
		return nil
	}
	return C.CString(string(data))
}


//export CreateDir
func CreateDir(path *C.char) C.int {
	goPath := C.GoString(path)
	err := os.MkdirAll(goPath, 0755)
	if err != nil {
		return -1
	}
	
	return 0
}

//export PathExists
func PathExists(path *C.char) C.int {
	goPath := C.GoString(path)
	_, err := os.Stat(goPath)
	if err == nil {
		return 1 // Existe
	}
	if os.IsNotExist(err) {
		return 0 // No existe
	}
	return -1 // Error al verificar
}

//export ListFiles
func ListFiles(dirPath *C.char) **C.char {
	goDirPath := C.GoString(dirPath)
	files, err := ioutil.ReadDir(goDirPath)
	if err != nil {
		return nil
	}
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	cArray := C.malloc(C.size_t(len(fileNames)+1) * C.size_t(unsafe.Sizeof(uintptr(0))))
	a := (*[1<<30 - 1]*C.char)(unsafe.Pointer(cArray))
	for i, name := range fileNames {
		a[i] = C.CString(name)
	}
	a[len(fileNames)] = nil // Terminador NULL
	return (**C.char)(cArray)
}

//export FreeListFiles
func FreeListFiles(arr **C.char) {
	if arr == nil {
		return
	}
	for i := 0; ; i++ {
		p := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(arr)) + uintptr(i)*unsafe.Sizeof(*arr)))
		if p == nil {
			break
		}
		C.free(unsafe.Pointer(p))
	}
	C.free(unsafe.Pointer(arr))
}

//export GetContentTypeFile
func GetContentTypeFile(b64Str *C.char) *C.char {
	goB64Str := C.GoString(b64Str)
	// Decodificar los primeros bytes del base64
	data, err := base64.StdEncoding.DecodeString(goB64Str)
	if err != nil || len(data) < 12 { // Necesitamos al menos 12 bytes para una buena detección
		return C.CString("application/octet-stream")
	}
	// Usar http.DetectContentType que analiza magic numbers
	contentType := http.DetectContentType(data[:12])
	// Algunos ajustes para tipos comunes
	switch {
	case strings.HasPrefix(contentType, "text/plain") && len(data) > 0:
		if isLikelyJSON(data) {
			return C.CString("application/json")
		}
		if isLikelyXML(data) {
			return C.CString("application/xml")
		}
	case strings.HasPrefix(contentType, "application/octet-stream"):
		if isPDF(data) {
			return C.CString("application/pdf")
		}
	}
	return C.CString(contentType)
}

// Funciones auxiliares para detección más precisa
func isPDF(data []byte) bool {
	return len(data) > 4 && string(data[:4]) == "%PDF"
}

func isLikelyJSON(data []byte) bool {
	firstChar := strings.TrimSpace(string(data[:1]))
	return firstChar == "{" || firstChar == "["
}

func isLikelyXML(data []byte) bool {
	str := strings.TrimSpace(string(data[:32]))
	return strings.HasPrefix(str, "<?xml") || 
	       strings.HasPrefix(str, "<html") || 
	       strings.HasPrefix(str, "<!DOCTYPE html")
}

func main() {}