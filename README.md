# cwind

CLI de primera fase para explorar ampacidad térmica dinámica (DLR) con trazas
CSV/KML y datos horarios de Open-Meteo. Los cálculos son aproximaciones de
ingeniería basadas en IEEE 738 y no sustituyen estudios eléctricos, térmicos,
meteorológicos ni decisiones operativas. Verifique entradas, supuestos y
regulación local.

```sh
go run . --version
go run . validate --trace linea.csv --conductor acsr-795-drake
go run . report --trace linea.csv --conductor acsr-795-drake --from 2024-01-01 --to 2024-01-02 --format json
```

`report` admite `text`, `json` y `csv`, `--output`, y calcula el azimut de la
traza salvo que se indique `--line-azimuth`. JSON incluye `schema_version`,
`status` y metadatos estables; CSV comienza con filas de metadatos comentadas.

El cache meteorológico es opcional y local (no se usa por defecto):

```sh
go run . report --trace linea.csv --conductor acsr-795-drake \
  --from 2024-01-01 --to 2024-01-02 --format json \
  --cache-dir .cwind-cache --cache-ttl 24h
```

Proteja el directorio de cache y elimínelo cuando ya no sea necesario. No
contiene trazas, solo respuestas del proveedor; `--cache-ttl 0` lo desactiva.

## Alcance científico

El cálculo es una estimación histórica basada en un balance térmico inspirado en
IEEE 738. No es una certificación de conformidad con la norma ni una medición
en tiempo real. La salida depende de la geometría suministrada, los parámetros
del conductor y los datos meteorológicos de Open-Meteo. No debe utilizarse como
orden operativa sin validación de ingeniería, instrumentación y regulación
aplicable.

La traza se procesa localmente. El CLI envía al proveedor meteorológico las
coordenadas representativas de la traza y las fechas solicitadas; no sube el
archivo CSV/KML.

## Desarrollo

```sh
go test ./...
go vet ./...
```

Los builds de CI cubren Windows, Linux y macOS en amd64/arm64 y ejecutan
tests con `-race` en Linux. Para reproducir una validación local:

```sh
gofmt -w .
go test ./...
go test -race ./...
go vet ./...
```

Los binarios oficiales se publican en GitHub Releases con un archivo
`checksums.txt`. Los instaladores deben verificar ese archivo antes de mover el
binario a su ubicación final.
