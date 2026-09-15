package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cwind-cli/cwind/internal/conductors"
	"github.com/cwind-cli/cwind/internal/geo"
	"github.com/cwind-cli/cwind/internal/physics"
	"github.com/cwind-cli/cwind/internal/weather"
	"github.com/spf13/cobra"
)

type reportResult struct {
	weather.DataPoint
	ampacity float64
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Genera un reporte horario de ampacidad dinámica (DLR)",
	Long: `Genera un reporte de ampacidad a partir de una traza y un período de hasta 48 horas.

La traza debe ser un CSV con cabecera lat,lon o Latitud_DD,Longitud_DD y coordenadas decimales,
o un archivo KML con coordenadas en el formato lon,lat,alt.
Las fechas deben usar el formato YYYY-MM-DD.`,
	Run: func(cmd *cobra.Command, args []string) {
		trace, _ := cmd.Flags().GetString("trace")
		conductorID, _ := cmd.Flags().GetString("conductor")
		fromText, _ := cmd.Flags().GetString("from")
		toText, _ := cmd.Flags().GetString("to")

		if trace == "" {
			fail("[ERROR] Falta --trace (archivo CSV o KML con coordenadas)")
		}
		from, err := time.ParseInLocation("2006-01-02", fromText, time.FixedZone("UTC-3", -3*60*60))
		if err != nil {
			fail("[ERROR] --from debe tener formato YYYY-MM-DD")
		}
		to, err := time.ParseInLocation("2006-01-02", toText, from.Location())
		if err != nil {
			fail("[ERROR] --to debe tener formato YYYY-MM-DD")
		}
		if to.Before(from) || to.Sub(from) > 24*time.Hour {
			fail("[ERROR] El rango debe ser inclusivo y cubrir como máximo 48 horas (dos fechas)")
		}

		coord, err := geo.ParseFile(trace)
		if err != nil {
			fail(fmt.Sprintf("[ERROR] Traza: %v", err))
		}
		cond, ok := conductors.Get(conductorID)
		if conductorID == "" {
			cond, conductorID, err = chooseConductor()
			ok = err == nil
		}
		if !ok {
			if err != nil {
				fail(fmt.Sprintf("[ERROR] %v", err))
			}
			fail(fmt.Sprintf("[ERROR] Conductor '%s' no encontrado", conductorID))
		}
		data, err := weather.FetchRange(coord.Lat, coord.Lon, from, to)
		if err != nil {
			fail(fmt.Sprintf("[ERROR] Weather: %v", err))
		}

		results := make([]reportResult, 0, len(data))
		for _, point := range data {
			ampacity := physics.Calculate(physics.WeatherData{
				TempAir: point.TempAir, WindSpeed: point.WindSpeed,
				WindDir: point.WindDir, LineAzimuth: 60,
				SolarRadiation: point.SolarRadiation,
			}, physics.Conductor{
				DiameterMeters: cond.DiameterMeters, R20: cond.R20,
				Alpha: cond.Alpha, MaxTemp: cond.MaxTemp,
			})
			results = append(results, reportResult{DataPoint: point, ampacity: ampacity})
		}
		printReport(coord, cond, results, from, to)
	},
}

func chooseConductor() (conductors.Conductor, string, error) {
	list := conductors.List()
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })

	fmt.Println("Seleccione un conductor:")
	for i, conductor := range list {
		fmt.Printf("  %d) %s — %s (%.0f A)\n", i+1, conductor.ID, conductor.Name, conductor.StaticRating)
	}
	fmt.Print("Opción: ")

	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return conductors.Conductor{}, "", fmt.Errorf("no se pudo leer la selección")
	}
	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || choice < 1 || choice > len(list) {
		return conductors.Conductor{}, "", fmt.Errorf("selección de conductor inválida")
	}
	return list[choice-1], list[choice-1].ID, nil
}

func fail(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func printReport(coord geo.Coordinate, cond conductors.Conductor, results []reportResult, from, to time.Time) {
	values := make([]float64, len(results))
	riskHours := 0
	sum, min, max := 0.0, results[0].ampacity, results[0].ampacity
	minResult, maxResult := results[0], results[0]
	for i, item := range results {
		values[i] = item.ampacity
		sum += item.ampacity
		if item.ampacity < min {
			min, minResult = item.ampacity, item
		}
		if item.ampacity > max {
			max, maxResult = item.ampacity, item
		}
		if item.ampacity < cond.StaticRating {
			riskHours++
		}
	}
	sort.Float64s(values)
	p10 := values[int(float64(len(values)-1)*0.10)]
	p50 := values[int(float64(len(values)-1)*0.50)]
	p90 := values[int(float64(len(values)-1)*0.90)]
	location := time.FixedZone("UTC-3", -3*60*60)
	fmt.Printf("cwind · línea · %s\n", cond.Name)
	fmt.Printf("  %.6f, %.6f · estática %.0f A\n\n", coord.Lat, coord.Lon, cond.StaticRating)
	fmt.Printf("  %s–%s, 00:00–23:00 (%d h)                       UTC-3\n\n",
		from.Format("02 Jan"), to.Format("02 Jan"), len(results))
	fmt.Println("  hora   t°C   viento      sol       dlr     vs estática")
	for _, item := range results {
		delta := (item.ampacity/cond.StaticRating - 1) * 100
		risk := ""
		if item.ampacity < cond.StaticRating {
			risk = "  riesgo"
		}
		solar := "—"
		if item.SolarRadiation > 0 {
			solar = fmt.Sprintf("%.0f W/m²", item.SolarRadiation)
		}
		fmt.Printf("  %s  %4.1f  %5.1f m/s  %8s  %6.0f A  %+4.0f%%%s\n",
			item.Time.In(location).Format("15:04"), item.TempAir, item.WindSpeed,
			solar, item.ampacity, delta, risk)
	}
	fmt.Printf("\n  resumen (%d h)\n", len(results))
	fmt.Printf("  media    %.0f A  %+0.0f%% vs estática\n", sum/float64(len(results)), (sum/float64(len(results))/cond.StaticRating-1)*100)
	fmt.Printf("  mínimo   %.0f A  %+0.0f%%   %s\n", min, (min/cond.StaticRating-1)*100, minResult.Time.In(location).Format("02 Jan 15:04"))
	fmt.Printf("  máximo   %.0f A  %+0.0f%%   %s\n", max, (max/cond.StaticRating-1)*100, maxResult.Time.In(location).Format("02 Jan 15:04"))
	fmt.Printf("\n  p10 %.0f A · p50 %.0f A · p90 %.0f A\n", p10, p50, p90)
	fmt.Printf("  %d h de riesgo (%.1f%%) sobre %d h\n", riskHours, float64(riskHours)/float64(len(results))*100, len(results))
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().String("trace", "", "Archivo CSV o KML con coordenadas de la línea (requerido)")
	reportCmd.Flags().String("conductor", "", "ID del conductor; si se omite, se muestra un menú")
	reportCmd.Flags().String("from", "", "Fecha inicial YYYY-MM-DD (requerida)")
	reportCmd.Flags().String("to", "", "Fecha final YYYY-MM-DD (requerida)")
	_ = reportCmd.MarkFlagRequired("trace")
	_ = reportCmd.MarkFlagRequired("from")
	_ = reportCmd.MarkFlagRequired("to")
}
