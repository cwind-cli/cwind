package conductors

import (
	"fmt"
	"math"
)

type Conductor struct {
	ID               string
	Name             string
	DiameterMeters   float64
	R20              float64
	Alpha            float64
	MaxTemp          float64
	StaticRating     float64
	Source           string
	RatingConditions string
}

var Database = map[string]Conductor{
	"acsr-16-2.5": {
		ID: "acsr-16-2.5", Name: "ACSR 16/2,5 (IRAM 2187-I/II)",
		DiameterMeters: 0.0054, R20: 0.00188, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 100.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-25-4": {
		ID: "acsr-25-4", Name: "ACSR 25/4 (IRAM 2187-I/II)",
		DiameterMeters: 0.0068, R20: 0.00120, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 133.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-35-6": {
		ID: "acsr-35-6", Name: "ACSR 35/6 (IRAM 2187-I/II)",
		DiameterMeters: 0.0081, R20: 0.000835, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 167.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-50-8": {
		ID: "acsr-50-8", Name: "ACSR 50/8 (IRAM 2187-I/II)",
		DiameterMeters: 0.0096, R20: 0.000595, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 206.0, Source: "CIMET, ficha IRAM 2187-I/II; confirmado también en EPE Santa Fe ETN 049",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-70-12": {
		ID: "acsr-70-12", Name: "ACSR 70/12 (IRAM 2187-I/II)",
		DiameterMeters: 0.0117, R20: 0.000413, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 260.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-95-15": {
		ID: "acsr-95-15", Name: "ACSR 95/15 (IRAM 2187-I/II)",
		DiameterMeters: 0.0136, R20: 0.000306, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 314.0, Source: "CIMET, ficha IRAM 2187-I/II; confirmado también en EPE Santa Fe ETN 049",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-120-20": {
		ID: "acsr-120-20", Name: "ACSR 120/20 (IRAM 2187-I/II)",
		DiameterMeters: 0.0155, R20: 0.000237, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 369.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-150-25": {
		ID: "acsr-150-25", Name: "ACSR 150/25 (IRAM 2187-I/II)",
		DiameterMeters: 0.0171, R20: 0.000194, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 418.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-185-30": {
		ID: "acsr-185-30", Name: "ACSR 185/30 (IRAM 2187-I/II)",
		DiameterMeters: 0.0190, R20: 0.000157, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 477.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-210-35": {
		ID: "acsr-210-35", Name: "ACSR 210/35 (IRAM 2187-I/II)",
		DiameterMeters: 0.0203, R20: 0.000138, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 518.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-240-40-cimet": {
		ID: "acsr-240-40-cimet", Name: "ACSR 240/40 (IRAM 2187-I/II, CIMET)",
		DiameterMeters: 0.0218, R20: 0.000119, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 568.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, AL SOL, V=0.6 m/s (ver alac-240-40: mismo calibre, SIN sol, IMSA)",
	},
	"acsr-300-50": {
		ID: "acsr-300-50", Name: "ACSR 300/50 (IRAM 2187-I/II)",
		DiameterMeters: 0.0244, R20: 0.0000949, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 655.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-340-30": {
		ID: "acsr-340-30", Name: "ACSR 340/30 (IRAM 2187-I/II)",
		DiameterMeters: 0.0250, R20: 0.0000851, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 696.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-380-50": {
		ID: "acsr-380-50", Name: "ACSR 380/50 (IRAM 2187-I/II)",
		DiameterMeters: 0.0270, R20: 0.0000757, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 752.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-435-55": {
		ID: "acsr-435-55", Name: "ACSR 435/55 (IRAM 2187-I/II)",
		DiameterMeters: 0.0288, R20: 0.0000660, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 816.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-550-70": {
		ID: "acsr-550-70", Name: "ACSR 550/70 (IRAM 2187-I/II)",
		DiameterMeters: 0.0324, R20: 0.0000526, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 946.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"acsr-680-85": {
		ID: "acsr-680-85", Name: "ACSR 680/85 (IRAM 2187-I/II)",
		DiameterMeters: 0.0360, R20: 0.0000426, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 1080.0, Source: "CIMET, ficha IRAM 2187-I/II",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"alac-240-40": {
		ID: "alac-240-40", Name: "ALAC 240/40 (IMSA) — familia ACSR/IRAM 2187",
		DiameterMeters: 0.02184, R20: 0.0001190, Alpha: 0.00403, MaxTemp: 80.0,
		StaticRating: 645.0, Source: "Ficha técnica IMSA, IRAM 2187",
		RatingConditions: "Tc=80°C, Tamb=40°C, SIN sol, V=0.6 m/s",
	},
	"acsr-795-drake": {
		ID: "acsr-795-drake", Name: "ACSR 795 Drake 26/7 (ASTM B-232)",
		DiameterMeters: 0.02814, R20: 0.0000732, Alpha: 0.00403, MaxTemp: 75.0,
		StaticRating: 907.0, Source: "Southwire / American Wire Group / Nexans",
		RatingConditions: "Tc=75°C, Tamb=25°C, V=0.6 m/s (2 ft/s), sol pleno, e=a=0.5, nivel del mar",
	},
	"cadla-16": {
		ID: "cadla-16", Name: "CADLA 16 mm² (IRAM 2212)",
		DiameterMeters: 0.0051, R20: 0.002070, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 97.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-25": {
		ID: "cadla-25", Name: "CADLA 25 mm² (IRAM 2212)",
		DiameterMeters: 0.0065, R20: 0.001300, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 129.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-35": {
		ID: "cadla-35", Name: "CADLA 35 mm² (IRAM 2212)",
		DiameterMeters: 0.0076, R20: 0.000944, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 158.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-50-7x": {
		ID: "cadla-50-7x", Name: "CADLA 50 mm² 7x3,02 (IRAM 2212)",
		DiameterMeters: 0.0091, R20: 0.000657, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 198.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-50-19x": {
		ID: "cadla-50-19x", Name: "CADLA 50 mm² 19x1,85 (IRAM 2212)",
		DiameterMeters: 0.0093, R20: 0.000648, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 200.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-70": {
		ID: "cadla-70", Name: "CADLA 70 mm² (IRAM 2212)",
		DiameterMeters: 0.0108, R20: 0.000480, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 242.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-95": {
		ID: "cadla-95", Name: "CADLA 95 mm² (IRAM 2212)",
		DiameterMeters: 0.0126, R20: 0.000349, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 295.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-120-19x": {
		ID: "cadla-120-19x", Name: "CADLA 120 mm² 19x2,85 (IRAM 2212)",
		DiameterMeters: 0.0143, R20: 0.000273, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 345.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-120-37x": {
		ID: "cadla-120-37x", Name: "CADLA 120 mm² 37x2,15 (IRAM 2212)",
		DiameterMeters: 0.0151, R20: 0.000247, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 367.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-150": {
		ID: "cadla-150", Name: "CADLA 150 mm² (IRAM 2212)",
		DiameterMeters: 0.0158, R20: 0.000226, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 389.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-185": {
		ID: "cadla-185", Name: "CADLA 185 mm² (IRAM 2212)",
		DiameterMeters: 0.0176, R20: 0.000180, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 448.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
	"cadla-240": {
		ID: "cadla-240", Name: "CADLA 240 mm² (IRAM 2212)",
		DiameterMeters: 0.0200, R20: 0.000141, Alpha: 0.0036, MaxTemp: 80.0,
		StaticRating: 522.0, Source: "CIMET, ficha IRAM 2212",
		RatingConditions: "Tc=80°C, Tamb=40°C, al sol, V=0.6 m/s",
	},
}

func Get(id string) (Conductor, bool) {
	c, ok := Database[id]
	return c, ok
}

func List() []Conductor {
	result := make([]Conductor, 0, len(Database))
	for _, c := range Database {
		result = append(result, c)
	}
	return result
}

func (c Conductor) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("conductor sin identificador")
	}
	if c.DiameterMeters <= 0 || math.IsNaN(c.DiameterMeters) || math.IsInf(c.DiameterMeters, 0) {
		return fmt.Errorf("conductor %q: diámetro inválido", c.ID)
	}
	if c.R20 <= 0 || math.IsNaN(c.R20) || math.IsInf(c.R20, 0) {
		return fmt.Errorf("conductor %q: resistencia R20 inválida", c.ID)
	}
	if c.Alpha < 0 || math.IsNaN(c.Alpha) || math.IsInf(c.Alpha, 0) {
		return fmt.Errorf("conductor %q: coeficiente térmico inválido", c.ID)
	}
	if c.MaxTemp <= -273.15 || math.IsNaN(c.MaxTemp) || math.IsInf(c.MaxTemp, 0) {
		return fmt.Errorf("conductor %q: temperatura máxima inválida", c.ID)
	}
	if c.StaticRating <= 0 || math.IsNaN(c.StaticRating) || math.IsInf(c.StaticRating, 0) {
		return fmt.Errorf("conductor %q: ampacidad estática inválida", c.ID)
	}
	return nil
}
