package generator

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerator_GenString(t *testing.T) {
	tests := []struct {
		name    string
		genData GenData
	}{
		{
			name: "default alphabet, len 10",
			genData: GenData{
				Alphabet: DefaultAlphabet,
				Length:   DefaultLength,
			},
		},
		{

			name: "default alphabet, len 7",
			genData: GenData{
				Alphabet: DefaultAlphabet,
				Length:   7,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			generator := NewGeneratorWithData(test.genData)
			output := generator.GenString()
			require.Equal(t, test.genData.Length, len(output))

		})
	}
}
