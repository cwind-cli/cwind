package cmd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cwind-cli/cwind/internal/conductors"
	"github.com/cwind-cli/cwind/internal/geo"
	"github.com/cwind-cli/cwind/internal/physics"
	"github.com/cwind-cli/cwind/internal/version"
	"github.com/cwind-cli/cwind/internal/weather"
	"github.com/spf13/cobra"
)

type reportResult struct {
	weather.DataPoint
	Ampacity float64 `json:"ampacity_a"`
}
type reportDocument struct {
	SchemaVersion  string         `json:"schema_version"`
	Status         string         `json:"status"`
	Errors         []string       `json:"errors,omitempty"`
	CwindVersion   string         `json:"cwind_version"`
	WeatherSource  string         `json:"weather_source"`
	Model          string         `json:"model"`
	Conductor      string         `json:"conductor"`
	Latitude       float64        `json:"latitude"`
	Longitude      float64        `json:"longitude"`
	ConductorID    string         `json:"conductor_id"`
	LineAzimuth    float64        `json:"line_azimuth_deg"`
	SolarElevation float64        `json:"solar_elevation_deg"`
	Timezone       string         `json:"timezone"`
	Results        []reportResult `json:"results"`
}

var fetchWeather = weather.FetchRangeContextWithCache

var reportCmd = &cobra.Command{
	Use: "report", Short: "Genera un reporte horario de ampacidad dinámica (DLR)",
	Args: cobra.NoArgs,
	RunE: runReport,
}

func runReport(cmd *cobra.Command, args []string) error {
	tracePath, _ := cmd.Flags().GetString("trace")
	id, _ := cmd.Flags().GetString("conductor")
	fromText, _ := cmd.Flags().GetString("from")
	toText, _ := cmd.Flags().GetString("to")
	azimuth, _ := cmd.Flags().GetFloat64("line-azimuth")
	azimuthSet := cmd.Flags().Changed("line-azimuth")
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")
	timezone, _ := cmd.Flags().GetString("timezone")
	solarElevation, _ := cmd.Flags().GetFloat64("solar-elevation")
	cacheDir, _ := cmd.Flags().GetString("cache-dir")
	cacheTTL, _ := cmd.Flags().GetDuration("cache-ttl")
	if strings.TrimSpace(tracePath) == "" {
		return fmt.Errorf("falta --trace")
	}
	if format != "text" && format != "json" && format != "csv" {
		return fmt.Errorf("--format debe ser text, json o csv")
	}
	if solarElevation <= 0 || solarElevation > 90 {
		return fmt.Errorf("--solar-elevation debe estar entre 0 y 90 grados (mayor que cero)")
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("--timezone inválida: %w", err)
	}
	from, err := time.ParseInLocation("2006-01-02", fromText, location)
	if err != nil {
		return fmt.Errorf("--from debe tener formato YYYY-MM-DD")
	}
	to, err := time.ParseInLocation("2006-01-02", toText, from.Location())
	if err != nil {
		return fmt.Errorf("--to debe tener formato YYYY-MM-DD")
	}
	if to.Before(from) || to.Sub(from) > 24*time.Hour {
		return fmt.Errorf("el rango debe ser de una o dos fechas consecutivas")
	}
	trace, err := geo.ParseTraceFile(tracePath)
	if err != nil {
		return fmt.Errorf("traza: %w", err)
	}
	coord := geo.Coordinate{Lat: 0, Lon: 0}
	for _, p := range trace {
		coord.Lat += p.Lat
		coord.Lon += p.Lon
	}
	coord.Lat /= float64(len(trace))
	coord.Lon /= float64(len(trace))
	if !azimuthSet {
		azimuth = trace.Azimuth()
	}
	if azimuth < 0 || azimuth >= 360 {
		return fmt.Errorf("--line-azimuth debe estar entre 0 y 360 grados")
	}
	cond, ok := conductors.Get(id)
	if id == "" {
		cond, id, err = chooseConductor()
		ok = err == nil
	}
	if !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("conductor %q no encontrado", id)
	}
	if err := cond.Validate(); err != nil {
		return err
	}
	data, err := fetchWeather(cmd.Context(), coord.Lat, coord.Lon, from, to, timezone, cacheDir, cacheTTL)
	if err != nil {
		return fmt.Errorf("weather: %w", err)
	}
	if len(data) == 0 {
		return fmt.Errorf("la respuesta meteorológica no contiene datos")
	}
	results := make([]reportResult, 0, len(data))
	for _, point := range data {
		a := physics.Calculate(physics.WeatherData{TempAir: point.TempAir, WindSpeed: point.WindSpeed, WindDir: point.WindDir, LineAzimuth: azimuth, SolarRadiation: point.SolarRadiation, SolarElevation: solarElevation}, physics.Conductor{DiameterMeters: cond.DiameterMeters, R20: cond.R20, Alpha: cond.Alpha, MaxTemp: cond.MaxTemp})
		results = append(results, reportResult{DataPoint: point, Ampacity: a})
	}
	doc := reportDocument{
		SchemaVersion: "1",
		Status:        "ok",
		CwindVersion:  version.Value,
		WeatherSource: "Open-Meteo archive API",
		Model:         "Modelo térmico aproximado basado en IEEE 738",
		Conductor:     cond.Name,
		ConductorID:   id,
		Latitude:      coord.Lat,
		Longitude:     coord.Lon,
		LineAzimuth:   azimuth,
		Timezone:      timezone,
		Results:       results,
	}
	var w io.Writer = cmd.OutOrStdout()
	var file *os.File
	if output != "" {
		file, err = os.Create(output)
		if err != nil {
			return fmt.Errorf("abrir --output: %w", err)
		}
		defer file.Close()
		w = file
	}
	return writeReport(w, format, doc, cond)
}

func chooseConductor() (conductors.Conductor, string, error) {
	list := conductors.List()
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	fmt.Fprintln(os.Stderr, "Seleccione un conductor:")
	for i, c := range list {
		fmt.Fprintf(os.Stderr, "  %d) %s — %s (%.0f A)\n", i+1, c.ID, c.Name, c.StaticRating)
	}
	fmt.Fprint(os.Stderr, "Opción: ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return conductors.Conductor{}, "", fmt.Errorf("no se pudo leer la selección")
	}
	n, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || n < 1 || n > len(list) {
		return conductors.Conductor{}, "", fmt.Errorf("selección de conductor inválida")
	}
	return list[n-1], list[n-1].ID, nil
}

func writeReport(w io.Writer, format string, doc reportDocument, cond conductors.Conductor) error {
	if len(doc.Results) == 0 {
		return fmt.Errorf("no hay resultados para escribir")
	}
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(doc)
	case "csv":
		c := csv.NewWriter(w)
		// Comment records keep metadata stable while preserving spreadsheet compatibility.
		metadata := [][]string{{"# cwind_schema_version", doc.SchemaVersion}, {"# status", doc.Status},
			{"# cwind_version", doc.CwindVersion}, {"# weather_source", doc.WeatherSource},
			{"# model", doc.Model}, {"# conductor_id", doc.ConductorID}, {"# timezone", doc.Timezone},
			{"# latitude", fmt.Sprintf("%.6f", doc.Latitude)}, {"# longitude", fmt.Sprintf("%.6f", doc.Longitude)},
			{"# line_azimuth_deg", fmt.Sprintf("%.3f", doc.LineAzimuth)}}
		for _, row := range metadata {
			if err := c.Write(row); err != nil {
				return err
			}
		}
		if err := c.Write([]string{"time", "temperature_c", "wind_speed_mps", "wind_direction_deg", "solar_wm2", "ampacity_a"}); err != nil {
			return err
		}
		for _, r := range doc.Results {
			if err := c.Write([]string{r.Time.Format(time.RFC3339), fmt.Sprintf("%.3f", r.TempAir), fmt.Sprintf("%.3f", r.WindSpeed), fmt.Sprintf("%.3f", r.WindDir), fmt.Sprintf("%.3f", r.SolarRadiation), fmt.Sprintf("%.3f", r.Ampacity)}); err != nil {
				return err
			}
		}
		c.Flush()
		return c.Error()
	default:
		fmt.Fprintf(w, "cwind · línea · %s\n  %.6f, %.6f · azimut %.1f°\n\n", doc.Conductor, doc.Latitude, doc.Longitude, doc.LineAzimuth)
		fmt.Fprintf(w, "  hora   t°C   viento      sol       dlr     vs estática\n")
		var sum, min, max float64
		min, max = doc.Results[0].Ampacity, doc.Results[0].Ampacity
		riskHours := 0
		for _, r := range doc.Results {
			sum += r.Ampacity
			if r.Ampacity < min {
				min = r.Ampacity
			}
			if r.Ampacity > max {
				max = r.Ampacity
			}
			delta := (r.Ampacity/cond.StaticRating - 1) * 100
			risk := ""
			if r.Ampacity < cond.StaticRating {
				risk = "  riesgo"
				riskHours++
			}
			fmt.Fprintf(w, "  %s  %4.1f  %5.1f m/s  %8.0f  %6.0f A  %+4.0f%%%s\n", r.Time.Format("15:04"), r.TempAir, r.WindSpeed, r.SolarRadiation, r.Ampacity, delta, risk)
		}
		fmt.Fprintf(w, "\n  resumen (%d h)\n  media    %.0f A\n  mínimo   %.0f A\n  máximo   %.0f A\n  riesgo   %d h (%.1f%%)\n",
			len(doc.Results), sum/float64(len(doc.Results)), min, max, riskHours,
			float64(riskHours)/float64(len(doc.Results))*100)
		return nil
	}
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().String("trace", "", "Archivo CSV o KML con coordenadas")
	reportCmd.Flags().String("conductor", "", "ID del conductor")
	reportCmd.Flags().String("from", "", "Fecha inicial YYYY-MM-DD")
	reportCmd.Flags().String("to", "", "Fecha final YYYY-MM-DD")
	reportCmd.Flags().Float64("line-azimuth", 0, "Azimut de la línea en grados (por defecto se calcula de la traza)")
	reportCmd.Flags().String("format", "text", "Formato de salida: text, json o csv")
	reportCmd.Flags().StringP("output", "o", "", "Archivo de salida (por defecto stdout)")
	reportCmd.Flags().String("timezone", "America/Argentina/Buenos_Aires", "Zona horaria IANA para fechas y weather")
	reportCmd.Flags().String("cache-dir", "", "Directorio de cache local (opt-in; vacío desactiva)")
	reportCmd.Flags().Duration("cache-ttl", 24*time.Hour, "Vigencia del cache local (0 desactiva)")
	reportCmd.Flags().Float64("solar-elevation", 60, "Elevación solar usada por el modelo, en grados (0-90)")
	_ = reportCmd.MarkFlagRequired("trace")
	_ = reportCmd.MarkFlagRequired("from")
	_ = reportCmd.MarkFlagRequired("to")
}
