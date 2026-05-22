package cache

import (
	"encoding/json"
	"testing"
)

func TestJSONRoundtrip_InMemory(t *testing.T) {
	type testStruct struct {
		A int
		B string
	}

	val := testStruct{A: 42, B: "hello"}
	data := struct{ A int; B string }{}
	// Simulate what Cache.Set/Get does
	enc, err := json.Marshal(val)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := json.Unmarshal(enc, &data); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if data.A != 42 || data.B != "hello" {
		t.Errorf("roundtrip failed: got %+v", data)
	}
}

func TestJSONRoundtrip_Slice(t *testing.T) {
	val := []int{1, 2, 3}
	data, err := json.Marshal(val)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded []int
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(decoded) != 3 {
		t.Errorf("expected 3 items, got %d", len(decoded))
	}
}
