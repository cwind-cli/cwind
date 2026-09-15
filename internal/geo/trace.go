package geo

import "math"

// Trace is the ordered set of vertices parsed from a line file.
type Trace []Coordinate

// ParseTraceFile parses all vertices. ParseFile remains the centroid API.
func ParseTraceFile(path string) (Trace, error) {
	return parseTraceFile(path)
}

// Azimuth returns the initial bearing (degrees clockwise from north) of a trace.
func (t Trace) Azimuth() float64 {
	if len(t) < 2 {
		return 0
	}
	a, b := t[0], t[len(t)-1]
	lat1, lat2 := a.Lat*math.Pi/180, b.Lat*math.Pi/180
	dlon := (b.Lon - a.Lon) * math.Pi / 180
	y := math.Sin(dlon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dlon)
	return math.Mod(math.Atan2(y, x)*180/math.Pi+360, 360)
}
