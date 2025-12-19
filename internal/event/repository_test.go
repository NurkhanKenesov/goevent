package event

import "testing"

func TestNewRepo(t *testing.T) {
	_ = NewRepository(nil)
}
