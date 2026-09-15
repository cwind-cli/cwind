package physics

import (
	"math"
	"testing"
)

func TestCalculateReturnsPositiveAmpacity(t *testing.T) {
	got := Calculate(WeatherData{
		TempAir: 25, WindSpeed: 2, WindDir: 60, LineAzimuth: 60,
	}, Conductor{
		DiameterMeters: 0.0218, R20: 0.000119, Alpha: 0.004, MaxTemp: 80,
	})
	if got <= 0 || math.IsNaN(got) || math.IsInf(got, 0) {
		t.Fatalf("Calculate() = %v, want finite positive value", got)
	}
}

func TestCalculateSolarRadiationReducesAmpacity(t *testing.T) {
	conductor := Conductor{
		DiameterMeters: 0.0218, R20: 0.000119, Alpha: 0.004, MaxTemp: 80,
	}
	base := WeatherData{TempAir: 30, WindSpeed: 1, WindDir: 60, LineAzimuth: 60}
	withoutSolar := Calculate(base, conductor)
	withSolar := Calculate(WeatherData{
		TempAir: 30, WindSpeed: 1, WindDir: 60, LineAzimuth: 60,
		SolarRadiation: 900,
	}, conductor)
	if withSolar >= withoutSolar {
		t.Fatalf("solar radiation increased ampacity: without=%v with=%v", withoutSolar, withSolar)
	}
}

func TestCalculateHandlesHotAirWithoutNaN(t *testing.T) {
	got := Calculate(WeatherData{
		TempAir: 100, WindSpeed: 0, LineAzimuth: 0,
	}, Conductor{
		DiameterMeters: 0.0218, R20: 0.000119, Alpha: 0.004, MaxTemp: 80,
	})
	if math.IsNaN(got) || math.IsInf(got, 0) || got < 0 {
		t.Fatalf("Calculate() = %v, want finite non-negative value", got)
	}
}
