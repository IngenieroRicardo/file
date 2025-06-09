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

### 🧪 Ejemplo para leer file

```C

```

---

### 🧪 Ejemplo para escribir, editar y eliminar file

```C

```

---


## 📚 Documentación de la API

#### Manejo de file binario
- `char* RBFile(char* inputPath)`: Retorna el Base64 del archivo leido.
- `int WBFile(char* b64Str, char* outputPath)`: Retorna 0 cuando el archivo se crea correctamente.
- `char* GetContentTypeFile(char* b64Str)`: Retorna el content-type de un base64.

#### Manejo de file de texto
- `char* RTFile(char* inputPath)`: Retorna el texto del archivo leido.
- `int WTFile(char* textStr, char* outputPath)`: Retorna 0 cuando el archivo se crea correctamente.

#### Manejo de dir
- `int CreateDir(char* path)`: Retorna 0 cuando el directorio se crea correctamente.
- `int PathExists(char* path)`: Retorna 1 cuando el dir/file existe.
- `char** ListFiles(char* dirPath)`: Retorna la lista de archivos en la ruta.

#### Utilidades
- `void FreeListFiles(char** files)`: Libera memoria de resultados.
