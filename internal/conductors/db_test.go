package conductors

import "testing"

func TestDatabaseConductorsAreValid(t *testing.T) {
	for id, conductor := range Database {
		if err := conductor.Validate(); err != nil {
			t.Errorf("%s: %v", id, err)
		}
	}
}
