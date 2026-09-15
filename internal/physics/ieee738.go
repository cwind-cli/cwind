package physics

import (
	"math"
)

const (
	solarAbsorptivity = 0.9
	emissivity        = 0.9
	stefanBoltzmann   = 5.67e-8
	airDensity        = 1.165
	airViscosity      = 1.85e-5
	airThermalCond    = 0.0263
	airSpecificHeat   = 1005.0
)

type WeatherData struct {
	TempAir        float64
	WindSpeed      float64
	WindDir        float64
	LineAzimuth    float64
	SolarRadiation float64
	SolarElevation float64
}

func Calculate(w WeatherData, c Conductor) float64 {
	if w.WindSpeed < 0 {
		w.WindSpeed = 0
	}
	windAngle := math.Abs(w.WindDir - w.LineAzimuth)
	if windAngle > 180 {
		windAngle = 360 - windAngle
	}
	windAngleRad := windAngle * math.Pi / 180.0

	diameter := c.DiameterMeters
	reynolds := airDensity * w.WindSpeed * diameter / airViscosity

	var nusselt float64
	if reynolds <= 0 {
		nusselt = 0.618 * math.Pow(grashof(diameter, c.MaxTemp-w.TempAir), 0.25)
	} else {
		angleFactor := 0.5 + 0.5*math.Cos(windAngleRad)
		nusselt = angleFactor * (0.35 + 0.56*math.Pow(reynolds, 0.52)) * math.Pow(prandtl(), 0.33)
	}

	hc := nusselt * airThermalCond / diameter
	qc := hc * math.Pi * diameter * (c.MaxTemp - w.TempAir)

	hr := emissivity * stefanBoltzmann * math.Pi * diameter *
		((c.MaxTemp+273.15)*(c.MaxTemp+273.15)*(c.MaxTemp+273.15)*(c.MaxTemp+273.15) -
			(w.TempAir+273.15)*(w.TempAir+273.15)*(w.TempAir+273.15)*(w.TempAir+273.15))

	solar := w.SolarRadiation
	if solar < 0 {
		solar = 0
	}
	solarElevation := w.SolarElevation
	if solarElevation <= 0 {
		solarElevation = 60
	}
	if solarElevation > 90 {
		solarElevation = 90
	}
	qs := solarAbsorptivity * solar * diameter * math.Sin(solarElevation*math.Pi/180)

	r := c.R20 * (1 + c.Alpha*(c.MaxTemp-20.0))

	heatBalance := (qc + hr - qs) / r
	if heatBalance <= 0 {
		return 0
	}
	ampacity := math.Sqrt(heatBalance)
	return ampacity
}

func grashof(diameter, deltaT float64) float64 {
	if deltaT < 0 {
		deltaT = 0
	}
	beta := 1.0 / (273.15 + 20.0)
	return 9.81 * beta * deltaT * math.Pow(diameter, 3) * math.Pow(airDensity, 2) / math.Pow(airViscosity, 2)
}

func prandtl() float64 {
	return airViscosity * airSpecificHeat / airThermalCond
}

type Conductor struct {
	DiameterMeters float64
	R20            float64
	Alpha          float64
	MaxTemp        float64
}
