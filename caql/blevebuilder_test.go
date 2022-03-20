package caql_test

import (
	"testing"

	"github.com/SecurityBrewery/catalyst/caql"
)

func TestBleveBuilder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		saql           string
		wantBleve      string
		wantParseErr   bool
		wantRebuildErr bool
	}{
		{name: "Search 1", saql: `"Bob"`, wantBleve: `"Bob"`},
		{name: "Search 2", saql: `"Bob" AND title == 'Name'`, wantBleve: `"Bob" title:"Name"`},
		{name: "Search 3", saql: `"Bob" OR title == 'Name'`, wantRebuildErr: true},
		{name: "Search 4", saql: `title == 'malware' AND 'wannacry'`, wantBleve: `title:"malware" "wannacry"`},
	}
	for _, tt := range tests {
		tt := tt

		parser := &caql.Parser{}

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expr, err := parser.Parse(tt.saql)
			if (err != nil) != tt.wantParseErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantParseErr)
				if expr != nil {
					t.Error(expr.String())
				}

				return
			}
			if err != nil {
				return
			}

			got, err := expr.BleveString()
			if (err != nil) != tt.wantRebuildErr {
				t.Error(expr.String())
				t.Errorf("String() error = %v, wantErr %v", err, tt.wantParseErr)

				return
			}
			if err != nil {
				return
			}
			if got != tt.wantBleve {
				t.Errorf("String() got = %v, want %v", got, tt.wantBleve)
			}
		})
	}
}
