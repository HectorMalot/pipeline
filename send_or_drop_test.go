package pipeline

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSendOrDrop(t *testing.T) {

	p := Emit(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	p = Delay(context.Background(), time.Millisecond, p)
	p = SendOrDrop(func(i int) { fmt.Printf("Dropping %d\n", i) }, p)
	p = Buffer(2, p)
	p = Delay(context.Background(), 3100*time.Microsecond, p)

	var res []int
	for v := range p {
		res = append(res, v)
	}
	require.Equal(t, 8, len(res))
}

func TestSendOrTimeout(t *testing.T) {

	p := Emit(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	p = Delay(context.Background(), time.Millisecond, p)
	p = SendOrTimeout(time.Millisecond, func(i int) { fmt.Printf("Dropping %d\n", i) }, p)
	p = Delay(context.Background(), 2*time.Millisecond, p)

	var res []int
	for v := range p {
		res = append(res, v)
	}
	require.Equal(t, 8, len(res))
}
