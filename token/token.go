package token

import "fmt"

func (t *TokenMenager) AddUsage(tokens int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.tokenBalanceUsed += tokens
	return t.SaveUsage()
}

func (t *TokenMenager) Status() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.tokenBalanceLimit == 0 {
		return fmt.Sprintf("%d / ∞", t.tokenBalanceUsed)
	}

	percent := float64(t.tokenBalanceUsed) / float64(t.tokenBalanceLimit) * 100

	return fmt.Sprintf(
		"%d / %d (%.2f%%)",
		t.tokenBalanceUsed,
		t.tokenBalanceLimit,
		percent,
	)
}

func (t *TokenMenager) ResetUsage() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.tokenBalanceUsed = 0
	return t.SaveUsage()
}
