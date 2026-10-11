package json

import (
	"strings"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestUnmarshalTupleLength(t *testing.T) {
	ty := cty.Tuple([]cty.Type{cty.Number, cty.Number})

	t.Run("short", func(t *testing.T) {
		_, err := Unmarshal([]byte(`[1]`), ty)
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "not enough tuple elements") {
			t.Fatalf("got %v", err)
		}
		pe, ok := err.(cty.PathError)
		if !ok {
			t.Fatalf("expected PathError, got %T", err)
		}
		if len(pe.Path) != 0 {
			t.Fatalf("top-level path %#v", pe.Path)
		}
	})

	t.Run("empty", func(t *testing.T) {
		_, err := Unmarshal([]byte(`[]`), ty)
		if err == nil || !strings.Contains(err.Error(), "not enough tuple elements") {
			t.Fatalf("got %v", err)
		}
	})

	t.Run("long", func(t *testing.T) {
		_, err := Unmarshal([]byte(`[1,2,3]`), ty)
		if err == nil || !strings.Contains(err.Error(), "too many tuple elements") {
			t.Fatalf("got %v", err)
		}
	})

	t.Run("exact", func(t *testing.T) {
		got, err := Unmarshal([]byte(`[1,2]`), ty)
		if err != nil {
			t.Fatal(err)
		}
		want := cty.TupleVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2)})
		if !got.RawEquals(want) {
			t.Fatalf("got %#v", got)
		}
	})

	t.Run("nested short", func(t *testing.T) {
		nested := cty.Object(map[string]cty.Type{
			"a": cty.Object(map[string]cty.Type{
				"b": ty,
			}),
		})
		_, err := Unmarshal([]byte(`{"a":{"b":[1]}}`), nested)
		if err == nil {
			t.Fatal("expected an error")
		}
		pe, ok := err.(cty.PathError)
		if !ok {
			t.Fatalf("expected PathError, got %T: %v", err, err)
		}
		want := cty.Path{
			cty.GetAttrStep{Name: "a"},
			cty.GetAttrStep{Name: "b"},
		}
		if len(pe.Path) != len(want) {
			t.Fatalf("path %#v", pe.Path)
		}
		for i := range want {
			step, ok := pe.Path[i].(cty.GetAttrStep)
			if !ok || step.Name != want[i].(cty.GetAttrStep).Name {
				t.Fatalf("path[%d] = %#v", i, pe.Path[i])
			}
		}
	})
}
