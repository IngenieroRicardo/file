// file filejs.go

// +build js

package file

import (
	"encoding/base64"
	"net/http"
	"strings"
	"syscall/js"
	"fmt"
)

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

func GetContentTypeFile(b64Str string) string {
	data, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil || len(data) < 12 {
		return "application/octet-stream"
	}

	contentType := http.DetectContentType(data[:12])
	
	switch {
	case strings.HasPrefix(contentType, "text/plain") && len(data) > 0:
		if isLikelyJSON(data) {
			return "application/json"
		}
		if isLikelyXML(data) {
			return "application/xml"
		}
	case strings.HasPrefix(contentType, "application/octet-stream"):
		if isPDF(data) {
			return "application/pdf"
		}
	}
	return contentType
}

// isDirectoryContent analiza el HTML para determinar si es un listado de directorio
func isDirectoryContent(html string) bool {
    trimmed := strings.TrimSpace(html)
    if !strings.HasPrefix(trimmed, "<pre>") {
        return false
    }
    return strings.Contains(html, "<a href=\"")
}

func IsDirectory(path string) (bool, error) {
    resp, status := Get(path, "", "")
    if status != 200 {
        // No es directorio (no existe o error)
        return false, nil
    }
    return isDirectoryContent(resp), nil
}

// fetchBytes realiza una petición GET y devuelve el cuerpo como []byte.
func fetchBytes(url string) ([]byte, int, error) {
    ch := make(chan struct {
        data   []byte
        status int
        err    error
    }, 1)

    success := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        response := args[0]
        // Verificar status antes de leer el cuerpo (opcional pero recomendado)
        status := response.Get("status").Int()
        if status != 200 {
            ch <- struct {
                data   []byte
                status int
                err    error
            }{nil, status, fmt.Errorf("HTTP error: %d", status)}
            return nil
        }

        // Obtener arrayBuffer
        arrayBufferPromise := response.Call("arrayBuffer")
        arrayBufferPromise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
            arrayBuffer := args[0]
            // Crear Uint8Array desde el ArrayBuffer
            uint8Array := js.Global().Get("Uint8Array").New(arrayBuffer)
            length := uint8Array.Get("length").Int()
            buf := make([]byte, length)
            js.CopyBytesToGo(buf, uint8Array)
            ch <- struct {
                data   []byte
                status int
                err    error
            }{buf, status, nil}
            return nil
        }), js.FuncOf(func(this js.Value, args []js.Value) interface{} {
            ch <- struct {
                data   []byte
                status int
                err    error
            }{nil, status, fmt.Errorf("error al leer arrayBuffer")}
            return nil
        }))
        return nil
    })

    errorFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        ch <- struct {
            data   []byte
            status int
            err    error
        }{nil, 0, fmt.Errorf("error en fetch: %v", args[0].Get("message").String())}
        return nil
    })

    fetchPromise := js.Global().Call("fetch", url)
    fetchPromise.Call("then", success, errorFunc)

    result := <-ch
    success.Release()
    errorFunc.Release()
    return result.data, result.status, result.err
}

// RBFile lee un archivo del servidor (por su ruta) y devuelve su contenido en base64.
func RBFile(inputPath string) string {
    origin := js.Global().Get("location").Get("origin").String()
    url := origin + inputPath
    data, status, err := fetchBytes(url)
    if err != nil || status != 200 {
        return ""
    }
    return base64.StdEncoding.EncodeToString(data)
}

// RTFile lee un archivo de texto del servidor y devuelve su contenido como string.
func RTFile(inputPath string) string {
    origin := js.Global().Get("location").Get("origin").String()
    url := origin + inputPath
    body, status := Get(url, "", "")
    if status != 200 {
        return ""
    }
    return body
}

// PathExists verifica si una ruta en el servidor existe (código 200).
func PathExists(path string) bool {
    origin := js.Global().Get("location").Get("origin").String()
    url := origin + path
    _, status := Head(url, "", "")
    return status == 200
}

// parseDirectoryListing extrae los nombres de archivos (no directorios) del HTML de un listado.
func parseDirectoryListing(html string) []string {
    var files []string
    startTag := "<a href=\""
    for {
        start := strings.Index(html, startTag)
        if start == -1 {
            break
        }
        html = html[start+len(startTag):]
        endQuote := strings.Index(html, "\"")
        if endQuote == -1 {
            break
        }
        href := html[:endQuote]
        closeTag := "</a>"
        closeIdx := strings.Index(html, closeTag)
        if closeIdx == -1 {
            break
        }
        if !strings.HasSuffix(href, "/") {
            files = append(files, href)
        }
        html = html[closeIdx+len(closeTag):]
    }
    return files
}

// ListFiles obtiene el listado de archivos de un directorio en el servidor.
func ListFiles(dirPath string) []string {
    origin := js.Global().Get("location").Get("origin").String()
    if !strings.HasSuffix(dirPath, "/") {
        dirPath += "/"
    }
    url := origin + dirPath
    html, status := Get(url, "", "")
    if status != 200 {
        return nil
    }
    if !isDirectoryContent(html) {
        return nil
    }
    return parseDirectoryListing(html)
}

// WBFile no está soportado en entorno WASM (no se puede escribir en el servidor)
func WBFile(b64Str, outputPath string) error {
    return fmt.Errorf("WBFile no está soportado en el navegador")
}

// WTFile no está soportado en entorno WASM
func WTFile(textStr, outputPath string) error {
    return fmt.Errorf("WTFile no está soportado en el navegador")
}

// CreateDir no está soportado en entorno WASM
func CreateDir(path string) error {
    return fmt.Errorf("CreateDir no está soportado en el navegador")
}

func Get(url, headers, body string) (string, int) {
    return doRequest("GET", url, headers, body)
}

func doRequest(method, url, headersStr, body string) (string, int) {
    // Preparar el objeto de opciones para fetch
    options := make(map[string]interface{})
    options["method"] = method

    // Construir objeto Headers
    headersObj := js.Global().Get("Headers").New()
    if headersStr != "" {
        lines := strings.Split(headersStr, "\n")
        for _, line := range lines {
            line = strings.TrimSpace(line)
            if line == "" {
                continue
            }
            parts := strings.SplitN(line, ":", 2)
            if len(parts) != 2 {
                continue
            }
            key := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])
            headersObj.Call("append", key, value)
        }
    }
    options["headers"] = headersObj

    // Añadir cuerpo si existe y método lo permite
    if body != "" && method != "GET" && method != "HEAD" {
        options["body"] = body
    }

    // Canal para recibir el resultado
    ch := make(chan struct {
        body   string
        status int
    }, 1)

    // Crear callbacks
    success := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        response := args[0]
        // Leer el cuerpo como texto
        textPromise := response.Call("text")
        textPromise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
            bodyText := args[0].String()
            status := response.Get("status").Int()
            ch <- struct {
                body   string
                status int
            }{bodyText, status}
            return nil
        }), js.FuncOf(func(this js.Value, args []js.Value) interface{} {
            // Error al leer el cuerpo
            ch <- struct {
                body   string
                status int
            }{"", 0}
            return nil
        }))
        return nil
    })

    errorFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        ch <- struct {
            body   string
            status int
        }{"", 0}
        return nil
    })

    // Llamar a fetch
    fetchPromise := js.Global().Call("fetch", url, js.ValueOf(options))
    fetchPromise.Call("then", success, errorFunc)

    // Esperar el resultado
    result := <-ch
    success.Release()
    errorFunc.Release()
    return result.body, result.status
}

func Head(url, headers, body string) (string, int) {
    return doRequest("HEAD", url, headers, body)
}
