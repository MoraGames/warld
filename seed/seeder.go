package seed

import "math/rand/v2"

type Seeder struct {
	seed   *Seed
	Random *rand.Rand
}

func NewSeeder(seed *Seed) *Seeder {
	return &Seeder{
		seed:   seed,
		Random: rand.New(rand.NewPCG(seed.Num1, seed.Num2)),
	}
}

func (s *Seeder) GetSeed() *Seed {
	return s.seed
}

func (s *Seeder) SetSeed(seed *Seed) {
	s.seed = seed
	s.Random = rand.New(rand.NewPCG(seed.Num1, seed.Num2))
}
