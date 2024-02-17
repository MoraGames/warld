package seed

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

type Seed struct {
	Token string
	Num1  uint64
	Num2  uint64
}

func NewSeed() *Seed {
	value := rand.Uint64()
	seed := &Seed{
		Token: strconv.FormatUint(value, 10),
		Num1:  value,
		Num2:  (value / 3) * 2,
	}
	return seed
}

func GetSeed(token string) (*Seed, error) {
	// If seed is empty, return an error
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("Invalid token: token is empty")
	}

	// Convert the seed to an unsigned integer
	value, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid token: %s", err)
	}

	// Create a new seed
	seed := &Seed{
		Token: token,
		Num1:  value,
		Num2:  (value / 3) * 2,
	}
	return seed, nil
}

func (s *Seed) ToUint64() uint64 {
	return s.Num1 + s.Num2
}

/*
func GetSeeds(token string, amount int) ([]*Seed, error) {
	// If seed is empty, return an error
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("Invalid token: token is empty")
	}

	// Convert the seed to an unsigned integer
	value, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid token: %s", err)
	}

	seeds := make([]*Seed, 0)
	for i := 0; i < amount; i++ {
		// Create a new seed
		seeds = append(seeds, &Seed{
			Token: strconv.FormatUint(value+uint64(i*64), 10),
			Num1:  value + uint64(i*64),
			Num2:  (value + uint64(i*64)/3) * 2,
		})
	}
	return seeds, nil
}
*/
