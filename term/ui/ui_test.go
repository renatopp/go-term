package ui

import "testing"

func TestEventsWrapsIntoMultiEvent(t *testing.T) {
	got := Events("a", "b")

	multi, ok := got.(MultiEvent)
	if !ok {
		t.Fatalf("Events returned %T, want MultiEvent", got)
	}
	if len(multi) != 2 || multi[0] != "a" || multi[1] != "b" {
		t.Fatalf("unexpected MultiEvent contents: %#v", multi)
	}
}

func TestEventsWithNoArgsWrapsIntoEmptyMultiEvent(t *testing.T) {
	got := Events()

	multi, ok := got.(MultiEvent)
	if !ok {
		t.Fatalf("Events returned %T, want MultiEvent", got)
	}
	if len(multi) != 0 {
		t.Fatalf("expected empty MultiEvent, got %#v", multi)
	}
}
