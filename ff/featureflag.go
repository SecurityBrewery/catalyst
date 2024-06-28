package ff

import (
	"encoding/csv"
	"os"
	"slices"
	"strings"
)

func HasDevFlag() bool {
	return has("dev")
}

func HasDemoFlag() bool {
	return has("demo")
}

func HasDummyDataFlag() bool {
	return has("dummy-data")
}

func has(flag string) bool {
	return slices.Contains(Flags(), flag)
}

func Flags() []string {
	ff := os.Getenv("FF")

	csvReader := csv.NewReader(strings.NewReader(ff))

	flags, err := csvReader.Read()
	if err != nil {
		return nil
	}

	return flags
}
