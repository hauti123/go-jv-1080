package util

import "testing"

func TestClamp(t *testing.T) {
	if got := Clamp(5, 0, 127); got != 5 {
		t.Errorf("Clamp(5,0,127) = %d, want 5", got)
	}
	if got := Clamp(200, 0, 127); got != 127 {
		t.Errorf("Clamp(200,0,127) = %d, want 127", got)
	}
	if got := Clamp(-5, 0, 127); got != 0 {
		t.Errorf("Clamp(-5,0,127) = %d, want 0", got)
	}
	if got := Clamp(-70, -63, 63); got != -63 {
		t.Errorf("Clamp(-70,-63,63) = %d, want -63", got)
	}
	if got := Clamp(0.5, float32(0), float32(1)); got != 0.5 {
		t.Errorf("Clamp(0.5,0,1) = %v, want 0.5", got)
	}
	if got := Clamp(float32(1.5), float32(0), float32(1)); got != 1 {
		t.Errorf("Clamp(1.5,0,1) = %v, want 1", got)
	}
	if got := Clamp(byte(200), byte(32), byte(127)); got != 127 {
		t.Errorf("Clamp(byte 200,32,127) = %d, want 127", got)
	}
}

func TestEnvFlag(t *testing.T) {
	// unset -> default
	if got := EnvFlag("JV1080_TEST_UNSET", true); !got {
		t.Errorf("EnvFlag(unset, true) = false, want true")
	}
	if got := EnvFlag("JV1080_TEST_UNSET", false); got {
		t.Errorf("EnvFlag(unset, false) = true, want false")
	}
	// empty -> default
	t.Setenv("JV1080_TEST_EMPTY", "")
	if got := EnvFlag("JV1080_TEST_EMPTY", true); !got {
		t.Errorf("EnvFlag(empty, true) = false, want true")
	}
	// enabling and disabling values
	for _, v := range []string{"1", "true", "TRUE", "on", "yes", "anything"} {
		t.Setenv("JV1080_TEST_FLAG", v)
		if got := EnvFlag("JV1080_TEST_FLAG", false); !got {
			t.Errorf("EnvFlag(%q, false) = false, want true", v)
		}
	}
	for _, v := range []string{"0", "false", "FALSE", "off", "no"} {
		t.Setenv("JV1080_TEST_FLAG", v)
		if got := EnvFlag("JV1080_TEST_FLAG", true); got {
			t.Errorf("EnvFlag(%q, true) = true, want false", v)
		}
	}
}
