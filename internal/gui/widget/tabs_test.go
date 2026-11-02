package widget

import (
	"testing"
)

func TestTabsSelect(t *testing.T) {
	tb := NewTabs([]string{"COMMON", "TONE 1", "TONE 2", "TONE 3", "TONE 4"})
	if tb.Selected != 0 {
		t.Errorf("initial Selected = %d, want 0", tb.Selected)
	}
	var got []int
	tb.OnSelect = func(i int) { got = append(got, i) }
	tb.selectTab(2)
	if tb.Selected != 2 {
		t.Errorf("Selected after selectTab(2) = %d, want 2", tb.Selected)
	}
	tb.selectTab(2) // no change, no callback
	tb.selectTab(-1)
	tb.selectTab(99)
	if len(got) != 1 || got[0] != 2 {
		t.Errorf("OnSelect calls = %v, want [2]", got)
	}
	// programmatic Select must not call OnSelect
	tb.Select(4)
	if tb.Selected != 4 || len(got) != 1 {
		t.Errorf("Select(4): Selected=%d calls=%v, want 4/[2]", tb.Selected, got)
	}
}
