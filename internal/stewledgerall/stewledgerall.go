package stewledgerall

import (
	"fmt"

	"github.com/ankitvg/stew/internal/stewledgercat"
	"github.com/ankitvg/stew/internal/stewledgers"
	"github.com/ankitvg/stew/internal/stewledgertail"
)

type Options struct {
	TargetDir string
	Limit     int
}

type Section struct {
	Name    string
	Content string
}

type Result struct {
	TargetDir string
	Sections  []Section
}

func Cat(opts Options) (Result, error) {
	ledgersResult, err := stewledgers.List(stewledgers.Options{
		TargetDir: opts.TargetDir,
	})
	if err != nil {
		return Result{}, err
	}

	sections := make([]Section, 0, len(ledgersResult.Ledgers))
	for _, ledger := range ledgersResult.Ledgers {
		catResult, err := stewledgercat.Run(stewledgercat.Options{
			TargetDir: ledgersResult.TargetDir,
			Ledger:    ledger.Name,
		})
		if err != nil {
			return Result{}, fmt.Errorf("read ledger %q: %w", ledger.Name, err)
		}

		sections = append(sections, Section{
			Name:    ledger.Name,
			Content: catResult.Content,
		})
	}

	return Result{
		TargetDir: ledgersResult.TargetDir,
		Sections:  sections,
	}, nil
}

func Tail(opts Options) (Result, error) {
	if opts.Limit <= 0 {
		return Result{}, fmt.Errorf("%w: must be greater than zero", stewledgertail.ErrInvalidLimit)
	}

	ledgersResult, err := stewledgers.List(stewledgers.Options{
		TargetDir: opts.TargetDir,
	})
	if err != nil {
		return Result{}, err
	}

	sections := make([]Section, 0, len(ledgersResult.Ledgers))
	for _, ledger := range ledgersResult.Ledgers {
		tailResult, err := stewledgertail.Run(stewledgertail.Options{
			TargetDir: ledgersResult.TargetDir,
			Ledger:    ledger.Name,
			Limit:     opts.Limit,
		})
		if err != nil {
			return Result{}, fmt.Errorf("read ledger %q: %w", ledger.Name, err)
		}

		sections = append(sections, Section{
			Name:    ledger.Name,
			Content: tailResult.Content,
		})
	}

	return Result{
		TargetDir: ledgersResult.TargetDir,
		Sections:  sections,
	}, nil
}
