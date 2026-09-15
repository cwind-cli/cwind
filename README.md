# cwind

`cwind` es una herramienta CLI para evaluar ampacidad térmica dinámica (DLR) a
partir de una traza georreferenciada, un conductor y un rango de fechas. El
objetivo es apoyar análisis históricos y comparación de escenarios meteorológicos,
no reemplazar mediciones de campo, estudios normativos ni decisiones operativas
sin criterio profesional.

## Alcance y rigor científico

El cálculo es un modelo aproximado basado en un balance térmico inspirado en
IEEE 738, con parámetros de ingeniería fijos y datos históricos de Open-Meteo.
Esto significa que:

- la salida depende de la geometría de la traza,
- del conductor seleccionado,
- de la zona horaria y de los datos meteorológicos del período consultado,
- y de los supuestos del modelo interno.

No debe interpretarse como una medición en tiempo real, validación de
conformidad ni orden operativa. La herramienta es útil para análisis,
comparación de escenarios y revisión técnica preliminar.

La traza se procesa localmente. El CLI solo envía coordenadas representativas y
fechas al proveedor meteorológico; no sube el archivo CSV/KML completo.

## Instalación

Se distribuye como binario compilado y se publica en GitHub Releases junto con
`checksums.txt`.

```sh
# ver la versión instalada
cwind --version

# o compilar localmente
# go build -o cwind .
```

## Comandos disponibles

```sh
cwind --help
cwind validate --help
cwind report --help
cwind conductors list
cwind doctor
```

### 1) `doctor`

Comprueba que la base de conductores y la instalación local estén disponibles.

```sh
cwind doctor
```

### 2) `conductors list`

Lista los conductores registrados en la base local.

```sh
cwind conductors list
```

Ejemplo de salida:

```text
ID                 CONDUCTOR                                    AMPACIDAD ESTÁTICA
acsr-120-20        ACSR 120/20 (IRAM 2187-I/II)                 369 A
acsr-150-25        ACSR 150/25 (IRAM 2187-I/II)                 418 A
...
```

### 3) `validate`

Valida que la traza sea legible y que el conductor exista y cumpla con los
parámetros mínimos esperados, sin consultar datos meteorológicos.

```sh
cwind validate --trace linea.csv --conductor acsr-795-drake
```

Esto es útil para confirmar rápidamente que la entrada está bien formada antes
del cálculo completo.

### 4) `report`

Genera un reporte horario de ampacidad dinámica para un rango de fechas.

```sh
cwind report \
  --trace linea.csv \
  --conductor acsr-795-drake \
  --from 2024-01-01 \
  --to 2024-01-02 \
  --format json
```

Se puede especificar `text`, `json` o `csv`.

## Ejemplos de uso

### Reporte mínimo

```sh
cwind report \
  --trace linea.csv \
  --conductor acsr-795-drake \
  --from 2024-01-01 \
  --to 2024-01-02
```

### Guardar salida en archivo

```sh
cwind report \
  --trace linea.csv \
  --conductor acsr-795-drake \
  --from 2024-01-01 \
  --to 2024-01-02 \
  --format json \
  --output reporte.json
```

### Salida CSV para Excel o análisis externo

```sh
cwind report \
  --trace linea.csv \
  --conductor acsr-795-drake \
  --from 2024-01-01 \
  --to 2024-01-02 \
  --format csv > reporte.csv
```

## Flags importantes del `report`

### `--line-azimuth`

Permite fijar manualmente el azimut de la línea en grados. Si no se indica,
`cwind` calcula el azimut a partir de la traza.

```sh
cwind report --trace linea.csv --conductor acsr-795-drake \
  --from 2024-01-01 --to 2024-01-02 \
  --line-azimuth 60
```

### `--timezone`

Define la zona horaria usada para interpretar las fechas y para la consulta de
weather. Debe usar un identificador IANA válido, por ejemplo:

```sh
--timezone America/Argentina/Buenos_Aires
--timezone UTC
```

### `--solar-elevation`

Ajusta la elevación solar asumida por el modelo, en grados, entre 0 y 90.

```sh
cwind report --trace linea.csv --conductor acsr-795-drake \
  --from 2024-01-01 --to 2024-01-02 \
  --solar-elevation 45
```

### `--cache-dir` y `--cache-ttl`

El cache meteorológico es local y opcional. Si se activa, el CLI guarda las
respuestas del proveedor para reutilizarlas en ejecuciones posteriores.

```sh
cwind report \
  --trace linea.csv \
  --conductor acsr-795-drake \
  --from 2024-01-01 \
  --to 2024-01-02 \
  --cache-dir .cwind-cache \
  --cache-ttl 24h
```

- `--cache-dir ""` desactiva el cache.
- `--cache-ttl 0` lo desactiva en forma explícita.
- El cache no contiene la traza ni un archivo de trabajo completo; solo guarda
  respuestas meteorológicas recuperadas del proveedor.

## Salida y estructura del reporte

### JSON

El formato JSON incluye metadatos estables como:

- `schema_version`
- `status`
- `cwind_version`
- `weather_source`
- `model`
- `conductor_id`
- `latitude`
- `longitude`
- `line_azimuth_deg`
- `timezone`
- `results`

### CSV

El csv incluye:

- cabeceras de metadata comentadas con `#`
- columnas estables para hora, temperatura, viento, radiación solar y ampacidad

## Flujo recomendado

1. Ejecutar `cwind validate` para verificar la traza y el conductor.
2. Ejecutar `cwind report` con el período de interés.
3. Revisar el resumen y, si corresponde, exportar JSON/CSV.
4. Usar la salida solo como soporte técnico y comparar con estudios locales,
   mediciones de campo y criterios operativos.

## Desarrollo

```sh
go test ./...
go vet ./...
```

Para facilitar la validación local:

```sh
gofmt -w .
go test ./...
go test -race ./...
go vet ./...
```

Los builds de CI cubren Windows, Linux y macOS en amd64/arm64. Los binarios
oficiales se publican en GitHub Releases con un archivo `checksums.txt`, que los
instaladores deben verificar antes de instalar el binario final.
