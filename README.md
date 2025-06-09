# FILE
Una biblioteca ligera para manipular File en C.  
Compilada usando: `go build -o file.dll -buildmode=c-shared file.go`

---

### 📥 Descargar la librería

| Linux | Windows |
| --- | --- |
| `wget https://github.com/IngenieroRicardo/file/releases/download/1.0/file.so` | `Invoke-WebRequest https://github.com/IngenieroRicardo/file/releases/download/1.0/file.dll -OutFile ./file.dll` |
| `wget https://github.com/IngenieroRicardo/file/releases/download/1.0/file.h` | `Invoke-WebRequest https://github.com/IngenieroRicardo/file/releases/download/1.0/file.h -OutFile ./file.h` |

---

### 🛠️ Compilar

| Linux | Windows |
| --- | --- |
| `gcc -o main.bin main.c ./file.so` | `gcc -o main.exe main.c ./file.dll` |
| `x86_64-w64-mingw32-gcc -o main.exe main.c ./file.dll` |  |

---

### 🧪 Ejemplo de escritura y lectura

```C
#include <stdio.h>
#include <stdlib.h>
#include "file.h"

int main() {
    // 1. Ejemplo de escritura binaria desde base64
    char* base64Data = "SGVsbG8gV29ybGQh"; // "Hello World!" en base64
    char* binaryPath = "./salida.bin";
    
    if (WBFile(base64Data, binaryPath) == 0) {
        printf("Archivo binario creado: %s\n", binaryPath);
    }

    // 2. Ejemplo de escritura de texto
    char* textData = "Este es un texto de ejemplo\nSegunda línea";
    char* textPath = "./salida.txt";
    
    if (WTFile(textData, textPath) == 0) {
        printf("Archivo de texto creado: %s\n", textPath);
    }

    // 3. Ejemplo de lectura binaria (a base64)
    char* base64Result = RBFile(binaryPath);
    if (base64Result != NULL) {
        printf("Base64 del archivo binario: %s\n", base64Result);
        free(base64Result);
    }

    // 4. Ejemplo de lectura de texto
    char* textResult = RTFile(textPath);
    if (textResult != NULL) {
        printf("Contenido del archivo de texto:\n%s\n", textResult);
        free(textResult);
    }

    return 0;
}
```

---

### 🧪 Ejemplo de obtención de content-type

```C
#include <stdio.h>
#include <stdlib.h>
#include "file.h"

int main() {
    // Ejemplo de creación de directorio
    char* dirPath = "./nueva_carpeta/subcarpeta";
    
    if (CreateDir(dirPath) == 0) {
        printf("Directorio creado exitosamente: %s\n", dirPath);
        
        // Ahora podemos usar las otras funciones con este directorio
        char* filePath = "./nueva_carpeta/subcarpeta/archivo.txt";
        char* textData = "Este archivo fue creado en un nuevo directorio";
        
        if (WTFile(textData, filePath) == 0) {
            printf("Archivo creado en el nuevo directorio\n");
        }
    } else {
        printf("Error al crear directorio\n");
    }
    
    return 0;
}
```

---

### 🧪 Ejemplo de directorio

```C
#include <stdio.h>
#include <stdlib.h>
#include "file.h"

int main() {
    char* dirPath = "."; // Directorio actual
    
    // Obtener lista de archivos
    char** files = ListFiles(dirPath);
    
    if (files != NULL) {
        printf("Archivos en el directorio '%s':\n", dirPath);
        
        // Iterar hasta encontrar el terminador NULL
        for (int i = 0; files[i] != NULL; i++) {
            printf("- %s\n", files[i]);
        }
        
        // Liberar memoria
        FreeListFiles(files);
    } else {
        printf("Error al leer el directorio o directorio vacío\n");
    }
    
    return 0;
}

```

---


## 📚 Documentación de la API

#### Manejo de file binario
- `char* RBFile(char* inputPath)`: Retorna el Base64 del archivo leído.
- `int WBFile(char* b64Str, char* outputPath)`: Retorna 0 cuando el archivo se crea correctamente.
- `char* GetContentTypeFile(char* b64Str)`: Retorna el content-type de un base64.

#### Manejo de file de texto
- `char* RTFile(char* inputPath)`: Retorna el texto del archivo leído.
- `int WTFile(char* textStr, char* outputPath)`: Retorna 0 cuando el archivo se crea correctamente.

#### Manejo de dir
- `int CreateDir(char* path)`: Retorna 0 cuando el directorio se crea correctamente.
- `int PathExists(char* path)`: Retorna 1 cuando el directorio o archivo existe.
- `char** ListFiles(char* dirPath)`: Retorna la lista de archivos en la ruta.

#### Utilidades
- `void FreeListFiles(char** files)`: Libera memoria de resultados.
