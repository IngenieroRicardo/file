# FILE

Una biblioteca ligera para manipular archivos en C.  
Compilada usando: `go build -o file.dll -buildmode=c-shared file.go` (Windows) o `go build -o file.so -buildmode=c-shared file.go` (Linux).

---

## 📥 Descargar la librería

| Linux | Windows |
| --- | --- |
| `wget https://github.com/IngenieroRicardo/file/releases/download/2.0/file.so` | `Invoke-WebRequest https://github.com/IngenieroRicardo/file/releases/download/2.0/file.dll -OutFile ./file.dll` |
| `wget https://github.com/IngenieroRicardo/file/releases/download/2.0/file.h` | `Invoke-WebRequest https://github.com/IngenieroRicardo/file/releases/download/2.0/file.h -OutFile ./file.h` |

---

## 🛠️ Compilar

| Linux | Windows |
| --- | --- |
| `gcc -o main.bin main.c ./file.so` | `gcc -o main.exe main.c ./file.dll` |
| `x86_64-w64-mingw32-gcc -o main.exe main.c ./file.dll` (cross) |  |

---

## 🧪 Ejemplos de uso

### Escritura y lectura de archivos

```c
#include <stdio.h>
#include <stdlib.h>
#include "file.h"

int main() {
    // 1. Escritura binaria desde base64
    char* base64Data = "SGVsbG8gV29ybGQh"; // "Hello World!" en base64
    char* binaryPath = "./salida.bin";

    int wb = WBFile(base64Data, binaryPath);
    if (wb == 0) {
        printf("Archivo binario creado: %s\n", binaryPath);
    } else {
        printf("Error al crear binario (código %d)\n", wb);
    }

    // 2. Escritura de texto plano
    char* textData = "Este es un texto de ejemplo\nSegunda línea";
    char* textPath = "./salida.txt";

    int wt = WTFile(textData, textPath);
    if (wt == 0) {
        printf("Archivo de texto creado: %s\n", textPath);
    }

    // 3. Lectura binaria (a base64)
    char* base64Result = RBFile(binaryPath);
    if (base64Result != NULL) {
        printf("Base64 del archivo binario: %s\n", base64Result);
        free(base64Result);
    }

    // 4. Lectura de texto plano
    char* textResult = RTFile(textPath);
    if (textResult != NULL) {
        printf("Contenido del archivo de texto:\n%s\n", textResult);
        free(textResult);
    }

    return 0;
}
```

### Detección de tipo MIME (Content-Type)

```c
#include <stdio.h>
#include <stdlib.h>
#include "file.h"

int main() {
    // PNG 1x1 en base64
    char* imageBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HgAGgwJ/lK3Q6wAAAABJRU5ErkJggg==";
    char* contentType = GetContentTypeFile(imageBase64);
    printf("Content-Type (PNG): %s\n", contentType);
    free(contentType);

    // JSON en base64
    char* jsonBase64 = "ewogICJuYW1lIjogIkpvaG4gRG9lIiwKICAiYWdlIjogMzAKfQ=="; // {"name": "John Doe", "age": 30}
    contentType = GetContentTypeFile(jsonBase64);
    printf("Content-Type (JSON): %s\n", contentType);
    free(contentType);

    return 0;
}
```

### Operaciones con directorios

```c
#include <stdio.h>
#include <stdlib.h>
#include "file.h"

int main() {
    char* dirPath = ".";
    char* newDir = "./nuevo_directorio";

    // Verificar si existe
    if (PathExists(dirPath) == 1) {
        printf("El directorio '%s' existe\n", dirPath);
    }

    // Crear directorio
    if (CreateDir(newDir) == 0) {
        printf("Directorio '%s' creado\n", newDir);
    }

    // Listar archivos (no directorios)
    char** files = ListFiles(dirPath);
    if (files != NULL) {
        printf("Archivos en '%s':\n", dirPath);
        for (int i = 0; files[i] != NULL; i++) {
            printf("- %s\n", files[i]);
        }
        FreeListFiles(files);
    } else {
        printf("Error al listar archivos o directorio vacío\n");
    }

    // Verificar si una ruta es directorio
    int isDir = IsDirectory(dirPath);
    if (isDir == 1) {
        printf("'%s' es un directorio\n", dirPath);
    } else if (isDir == 0) {
        printf("'%s' es un archivo\n", dirPath);
    } else {
        printf("Error al acceder a '%s'\n", dirPath);
    }

    return 0;
}
```

---

## 📚 Documentación de la API

### Funciones para archivos binarios (Base64)

| Función | Descripción | Retorno |
|---------|-------------|---------|
| `char* RBFile(char* inputPath)` | Lee un archivo binario y lo devuelve codificado en Base64. | Puntero a cadena con el Base64, o `NULL` si hay error. Debe liberarse con `free()`. |
| `int WBFile(char* b64Str, char* outputPath)` | Escribe un archivo binario a partir de una cadena Base64. | `0` si éxito, `-1` si falla la decodificación Base64, `-2` si falla la escritura del archivo. |
| `char* GetContentTypeFile(char* b64Str)` | Detecta el tipo MIME de un archivo a partir de su contenido en Base64 (lee primeros bytes). | Puntero a cadena con el Content-Type (ej: `"image/png"`, `"application/json"`). Debe liberarse con `free()`. |

### Funciones para archivos de texto plano

| Función | Descripción | Retorno |
|---------|-------------|---------|
| `char* RTFile(char* inputPath)` | Lee un archivo de texto y lo devuelve como cadena. | Puntero a cadena con el contenido, o `NULL` si error. Debe liberarse con `free()`. |
| `int WTFile(char* textStr, char* outputPath)` | Escribe una cadena de texto en un archivo. | `0` si éxito, `-1` si falla la escritura. |

### Funciones para directorios

| Función | Descripción | Retorno |
|---------|-------------|---------|
| `int CreateDir(char* path)` | Crea un directorio (y los intermedios si es necesario). | `0` si éxito, `-1` si error. |
| `int PathExists(char* path)` | Verifica si una ruta (archivo o directorio) existe. | `1` si existe, `0` si no existe, `-1` si error. |
| `int IsDirectory(char* path)` | Determina si una ruta es un directorio. | `1` si es directorio, `0` si es archivo, `-1` si error. |
| `char** ListFiles(char* dirPath)` | Devuelve un arreglo de cadenas con los nombres de los archivos (no directorios) en la ruta dada. El arreglo termina con `NULL`. | Puntero a arreglo, o `NULL` si error o directorio vacío. Debe liberarse con `FreeListFiles()`. |
| `void FreeListFiles(char** files)` | Libera la memoria asignada por `ListFiles()`. | — |

---

## ⚠️ Notas importantes

- Todas las cadenas devueltas por la librería (excepto `ListFiles`) deben liberarse con `free()`.
- `ListFiles` solo devuelve nombres de archivos (no directorios). Para listar todo, puedes usar `IsDirectory` adicionalmente.
- `GetContentTypeFile` analiza los primeros 12 bytes del archivo decodificado para detectar el tipo MIME. Para archivos pequeños puede ser suficiente, pero para una detección más precisa se recomienda usar la cabecera completa.
- Los códigos de error negativos se reservan para indicar fallos específicos.



Este README ahora cubre todas las funciones exportadas, incluye ejemplos actualizados y documentación clara. Además, he corregido algunos detalles como la inclusión de `IsDirectory` y los códigos de retorno específicos de `WBFile`.
