package pipeline

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDistributor(t *testing.T) {
	ints := []int{}
	for i := 0; i < 5; i++ {
		ints = append(ints, i)
	}

	p := Emit(ints...)

	d := Distribute(p, 3)
	out := Merge(d...)

	var result []int
	for r := range out {
		result = append(result, r)
	}
	require.Equal(t, 15, len(result))
}
