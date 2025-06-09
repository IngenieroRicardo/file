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
#include <stdio.h>
#include "file.h"

int main() {
    char* file = "{\"nombre\":\"Juan\", \"edad\":30, \"direccion\": {\"pais\":\"Villa Lactea\",\"departamento\":\"Tierra\"}, \"documentos\": [\"B00000001\",\"00000000-1\"], \"foto\":\"iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAArSURBVBhXY/iPA0AlGBgwGFAKlwQmAKrAIgcVRZODCsI5cAAVgVDo4P9/AHe4m2U/OJCWAAAAAElFTkSuQmCC\" }";
    
    // Analizar file
    fileResult resultado = Parsefile(file);
    
    if (resultado.is_valid) {
        printf("file válido: %s\n", resultado.value);
    } else {
        printf("Error: %s\n", resultado.error);
        FreefileResult(resultado);
        return 1;
    }
    
    // Obtener valores
    fileResult nombre = GetfileValue(file, "nombre");
    fileResult pais = GetfileValueByPath(file, "direccion.pais");
    fileResult documento1 = GetfileValueByPath(file, "documentos.0");
    
    // Mostrar valores sin comillas
    printf("Nombre: %s\n", nombre.value);
    printf("País: %s\n", pais.value);
    printf("Primer Documento: %s\n", documento1.value);
    
    // Liberar memoria
    FreefileResult(resultado);
    FreefileResult(nombre);
    FreefileResult(pais);
    FreefileResult(documento1);
    
    return 0;
}
```

---

### 🧪 Ejemplo para escribir, editar y eliminar file

```C
#include <stdio.h>
#include "file.h"

int main() {
    // 1. Crear un objeto file vacío
    fileResult file_vacio = CreateEmptyfile();
    printf("file vacío: %s\n", file_vacio.value);
    FreefileResult(file_vacio);

    // 2. Crear un objeto file con datos básicos de persona
    fileResult persona = CreateEmptyfile();
    persona = AddStringTofile(persona.value, "nombre", "Juan Pérez");
    persona = AddNumberTofile(persona.value, "edad", 30);
    persona = AddBooleanTofile(persona.value, "es_estudiante", 0); // 0 = falso
    
    printf("\nPersona básica:\n%s\n", persona.value);

    // 3. Crear una dirección como file y añadirla a la persona
    fileResult direccion = CreateEmptyfile();
    direccion = AddStringTofile(direccion.value, "calle", "Calle Principal 123");
    direccion = AddStringTofile(direccion.value, "ciudad", "Ciudad Ejemplo");
    direccion = AddStringTofile(direccion.value, "pais", "España");
    
    persona = AddfileTofile(persona.value, "direccion", direccion.value);
    FreefileResult(direccion);

    // 4. Crear un array de pasatiempos y añadirlo
    fileResult pasatiempos = CreateEmptyArray();
    pasatiempos = AddItemToArray(pasatiempos.value, "\"fútbol\"");
    pasatiempos = AddItemToArray(pasatiempos.value, "\"lectura\"");
    pasatiempos = AddItemToArray(pasatiempos.value, "\"programación\"");
    
    persona = AddfileTofile(persona.value, "pasatiempos", pasatiempos.value);
    FreefileResult(pasatiempos);

    // 5. Modificar el file existente
    persona = AddNumberTofile(persona.value, "edad", 31); // Actualizar edad
    persona = AddStringTofile(persona.value, "correo", "juan@ejemplo.com");
    
    printf("\nPersona actualizada:\n%s\n", persona.value);

    // 6. Eliminar una propiedad
    persona = RemoveKeyFromfile(persona.value, "es_estudiante");
    printf("\nPersona sin 'es_estudiante':\n%s\n", persona.value);

    // 7. Crear otro file con información laboral
    fileResult info_laboral = CreateEmptyfile();
    info_laboral = AddStringTofile(info_laboral.value, "empresa", "Soluciones Tecnológicas");
    info_laboral = AddStringTofile(info_laboral.value, "puesto", "Desarrollador");
    
    // Combinar con el file de persona
    persona = Mergefile(persona.value, info_laboral.value);
    printf("\nPersona con información laboral:\n%s\n", persona.value);
    FreefileResult(info_laboral);

    // 8. Verificar si el file es válido
    int es_valido = IsValidfile(persona.value);
    printf("\n¿file válido? %s\n", es_valido ? "Sí" : "No");

    // Liberar memoria
    FreefileResult(persona);

    return 0;
}
```

---


## 📚 Documentación de la API

#### Manejo Básico de file
- `fileResult Parsefile(char* fileStr)`: Analiza una cadena file
- `int IsValidfile(char* file_str)`: Verifica si una cadena es file válido

#### Obtención de Valores
- `fileResult GetfileValue(char* file_str, char* key)`: Obtiene valor por clave
- `fileResult GetfileValueByPath(char* file_str, char* path)`: Obtiene valor por ruta
- `fileResult GetArrayLength(char* file_str)`: Obtiene longitud de array
- `fileResult GetArrayItem(char* file_str, int index)`: Obtiene elemento de array

#### Construcción/Modificación
- `fileResult CreateEmptyfile()`: Crea objeto file vacío
- `fileResult CreateEmptyArray()`: Crea array file vacío
- `fileResult AddStringTofile(char* file_str, char* key, char* value)`
- `fileResult AddNumberTofile(char* file_str, char* key, double value)`
- `fileResult AddBooleanTofile(char* file_str, char* key, int value)`
- `fileResult AddfileTofile(char* parent_file, char* key, char* child_file)`
- `fileResult AddItemToArray(char* file_array, char* item)`
- `fileResult RemoveKeyFromfile(char* file_str, char* key)`
- `fileResult RemoveItemFromArray(char* file_array, int index)`
- `fileResult Mergefile(char* file1, char* file2)`: Combina dos files

#### Utilidades
- `void FreefileResult(fileResult result)`: Libera memoria de resultados
- `void FreefileArrayResult(fileArrayResult result)`: Libera memoria de arrays

### Estructuras
```c
typedef struct {
    char* value;      // Valor obtenido
    int is_valid;     // 1 si es válido, 0 si hay error
    char* error;      // Mensaje de error (si lo hay)
} fileResult;

typedef struct {
    char** items;     // Array de elementos
    int count;        // Número de elementos
    int is_valid;     // 1 si es válido, 0 si hay error
    char* error;      // Mensaje de error (si lo hay)
} fileArrayResult;
```
