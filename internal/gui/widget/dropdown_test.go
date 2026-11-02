package widget

import (
	"testing"
)

func TestDropdownInitialAndSet(t *testing.T) {
	d := NewDropdown("WAVE", "patch/temp/tone/1/number", []string{"001 A", "002 B", "003 C"})
	if d.Value() != 0 {
		t.Errorf("initial value = %d, want 0", d.Value())
	}
	d.Set(2)
	if d.Value() != 2 {
		t.Errorf("Set(2) = %d, want 2", d.Value())
	}
	d.Set(99)
	if d.Value() != 2 {
		t.Errorf("Set(99) = %d, want 2 (clamped)", d.Value())
	}
	d.Set(-1)
	if d.Value() != 0 {
		t.Errorf("Set(-1) = %d, want 0 (clamped)", d.Value())
	}
}

func TestDropdownChoose(t *testing.T) {
	d := NewDropdown("WAVE", "patch/temp/tone/1/number", []string{"001 A", "002 B", "003 C"})
	var got []int
	d.OnChange = func(i int) { got = append(got, i) }
	d.open = true
	d.choose(1)
	if d.Value() != 1 {
		t.Errorf("choose(1) value = %d, want 1", d.Value())
	}
	if d.open {
		t.Errorf("choose did not close the list")
	}
	if len(got) != 1 || got[0] != 1 {
		t.Errorf("OnChange calls = %v, want [1]", got)
	}
	// programmatic Set must not call OnChange
	d.Set(2)
	if len(got) != 1 {
		t.Errorf("Set triggered OnChange: %v", got)
	}
}

func TestDropdownSetOptionsClamps(t *testing.T) {
	d := NewDropdown("WAVE", "patch/temp/tone/1/number", []string{"a", "b", "c", "d", "e"})
	d.Set(4)
	if len(d.tags) != 6 { // button + 5 options
		t.Fatalf("tags = %d, want 6", len(d.tags))
	}
	d.SetOptions([]string{"x", "y"})
	if d.Value() != 1 {
		t.Errorf("value after shrinking options = %d, want 1 (clamped)", d.Value())
	}
	if len(d.Options) != 2 || len(d.tags) != 3 {
		t.Errorf("SetOptions did not replace options/tags: %d/%d", len(d.Options), len(d.tags))
	}
}
