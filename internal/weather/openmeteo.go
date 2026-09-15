package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HistoricalResponse struct {
	Hourly HourlyData `json:"hourly"`
}

type HourlyData struct {
	Time             []string  `json:"time"`
	Temperature2m    []float64 `json:"temperature_2m"`
	WindSpeed10m     []float64 `json:"wind_speed_10m"`
	WindDirection10m []float64 `json:"wind_direction_10m"`
	SolarRadiation   []float64 `json:"shortwave_radiation"`
}

type DataPoint struct {
	Time           time.Time
	TempAir        float64
	WindSpeed      float64
	WindDir        float64
	SolarRadiation float64
}

func FetchRange(lat, lon float64, from, to time.Time) ([]DataPoint, error) {
	url := fmt.Sprintf(
		"https://archive-api.open-meteo.com/v1/archive?latitude=%.4f&longitude=%.4f&start_date=%s&end_date=%s&hourly=temperature_2m,wind_speed_10m,wind_direction_10m,shortwave_radiation&timezone=America%%2FArgentina%%2FBuenos_Aires",
		lat, lon, from.Format("2006-01-02"), to.Format("2006-01-02"),
	)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("petición HTTP fallida: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var apiResp HistoricalResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode JSON fallido: %w", err)
	}

	n := len(apiResp.Hourly.Time)
	if n == 0 {
		return nil, fmt.Errorf("no hay datos horarios en la respuesta")
	}
	if len(apiResp.Hourly.Temperature2m) != n ||
		len(apiResp.Hourly.WindSpeed10m) != n ||
		len(apiResp.Hourly.WindDirection10m) != n ||
		len(apiResp.Hourly.SolarRadiation) != n {
		return nil, fmt.Errorf("respuesta meteorológica incompleta")
	}

	result := make([]DataPoint, n)
	for i := 0; i < n; i++ {
		t, err := time.ParseInLocation("2006-01-02T15:04", apiResp.Hourly.Time[i], time.FixedZone("UTC-3", -3*60*60))
		if err != nil {
			return nil, fmt.Errorf("parsear fecha meteorológica: %w", err)
		}
		result[i] = DataPoint{
			Time:           t,
			TempAir:        apiResp.Hourly.Temperature2m[i],
			WindSpeed:      apiResp.Hourly.WindSpeed10m[i] / 3.6,
			WindDir:        apiResp.Hourly.WindDirection10m[i],
			SolarRadiation: apiResp.Hourly.SolarRadiation[i],
		}
	}
	return result, nil
}
