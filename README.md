# Compilación
Para generar un binario es necesario instalar Go o utilizar un contenedor de docker para construirlo:
```bash
go build evaluation.go
```
O utilizando un contenedor:
```bash
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.17 go build evaluation.go 
```
Esto generará un binario `evaluation` en la carpeta raíz.

# Uso

| Parámetro | Descripción                                                                 | Valor por defecto |
|-----------|-----------------------------------------------------------------------------|-------------------|
| -qrels    | Ruta del fichero con los juicios de relevancia                              | qrels.txt         |
| -results  | Ruta del fichero con los resultados del sistema de infomación               | results.txt       |
| -output   | Ruta del fichero donde se guardará la evaluación del sistema de información | outputs.txt       |

Se imprime en salida estandar el resultado del análisis y se genera un fichero png con las gráficas correspondientes.
