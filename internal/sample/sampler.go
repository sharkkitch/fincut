// Package sample provides line-level sampling for structured log streams.
// It supports rate-based (1-in-N) and random probability sampling.
package sample

import (
	"fmt"
	"math/rand"
)

// Options configures the Sampler behaviour.
type Options struct {
	// Rate keeps every Nth line (e.g. Rate=3 keeps lines 1,4,7,...). 0 disables rate sampling.
	Rate int
	// Probability is a value in (0,1] for random sampling. 0 disables probability sampling.
	Probability float64
	// Seed is used to initialise the random source. 0 means non-deterministic.
	Seed int64
}

// Sampler filters a slice of lines according to the configured strategy.
type Sampler struct {
	opts Options
	rng  *rand.Rand
}

// New constructs a Sampler from opts, returning an error if the options are invalid.
func New(opts Options) (*Sampler, error) {
	if opts.Rate < 0 {
		return nil, fmt.Errorf("sample: rate must be >= 0, got %d", opts.Rate)
	}
	if opts.Probability < 0 || opts.Probability > 1 {
		return nil, fmt.Errorf("sample: probability must be in [0,1], got %f", opts.Probability)
	}
	if opts.Rate == 0 && opts.Probability == 0 {
		return nil, fmt.Errorf("sample: at least one of Rate or Probability must be set")
	}
	if opts.Rate > 0 && opts.Probability > 0 {
		return nil, fmt.Errorf("sample: Rate and Probability are mutually exclusive")
	}

	var rng *rand.Rand
	if opts.Probability > 0 {
		//nolint:gosec // deterministic seed is intentional for reproducibility
		rng = rand.New(rand.NewSource(opts.Seed))
	}

	return &Sampler{opts: opts, rng: rng}, nil
}

// Apply returns the sampled subset of lines.
func (s *Sampler) Apply(lines []string) []string {
	if len(lines) == 0 {
		return []string{}
	}

	if s.opts.Rate > 0 {
		return s.applyRate(lines)
	}
	return s.applyProbability(lines)
}

func (s *Sampler) applyRate(lines []string) []string {
	out := make([]string, 0, len(lines)/s.opts.Rate+1)
	for i, l := range lines {
		if i%s.opts.Rate == 0 {
			out = append(out, l)
		}
	}
	return out
}

func (s *Sampler) applyProbability(lines []string) []string {
	out := make([]string, 0, int(float64(len(lines))*s.opts.Probability)+1)
	for _, l := range lines {
		if s.rng.Float64() < s.opts.Probability {
			out = append(out, l)
		}
	}
	return out
}
