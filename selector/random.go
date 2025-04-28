package selector

import (
	"fmt"
	"math/rand"
)

type Random struct{}

func (r *Random) SelectNext(users []string, lastIndex int, counts map[string]int) (int, error) {
	if len(users) == 0 {
		return -1, fmt.Errorf("empty users list")
	}
	return rand.Intn(len(users)), nil
}
