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

	d := NewDistributor(p)
	p1 := d.NewChan()
	p2 := d.NewChan()
	p3 := d.NewChan()

	out := Merge(p1, p2, p3)

	var result []int
	for r := range out {
		result = append(result, r)
	}
	require.Equal(t, 15, len(result))
}
