package term

import "testing"

func TestDecodeControlKeyShift(t *testing.T) {
	event := decodeControlKey('a')
	if event.Type != KeyRune || event.Rune != 'a' || event.Shift {
		t.Fatalf("decodeControlKey('a') = %+v, want lowercase rune, no Shift", event)
	}

	event = decodeControlKey('A')
	if event.Type != KeyRune || event.Rune != 'A' || !event.Shift {
		t.Fatalf("decodeControlKey('A') = %+v, want uppercase rune, Shift", event)
	}

	// Ctrl+letter destroys case in the control byte (0x01-0x1a), so Shift
	// must not be inferred there.
	event = decodeControlKey(1)
	if event.Type != KeyRune || event.Rune != 'a' || !event.Ctrl || event.Shift {
		t.Fatalf("decodeControlKey(1) = %+v, want Ctrl+'a', no Shift", event)
	}
}

func TestDecodeModifier(t *testing.T) {
	tests := []struct {
		mod              int
		shift, alt, ctrl bool
	}{
		{0, false, false, false},
		{1, false, false, false},
		{2, true, false, false},
		{3, false, true, false},
		{5, false, false, true},
		{8, true, true, true},
	}
	for _, tt := range tests {
		shift, alt, ctrl := decodeModifier(tt.mod)
		if shift != tt.shift || alt != tt.alt || ctrl != tt.ctrl {
			t.Errorf("decodeModifier(%d) = (%v,%v,%v), want (%v,%v,%v)",
				tt.mod, shift, alt, ctrl, tt.shift, tt.alt, tt.ctrl)
		}
	}
}

func TestSplitParams(t *testing.T) {
	tests := []struct {
		params string
		want   []int
	}{
		{"", nil},
		{"5", []int{5}},
		{"1;5", []int{1, 5}},
		{"97:65;5:1", []int{97, 5}},
		{"27;5;13", []int{27, 5, 13}},
	}
	for _, tt := range tests {
		got := splitParams(tt.params)
		if len(got) != len(tt.want) {
			t.Errorf("splitParams(%q) = %v, want %v", tt.params, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("splitParams(%q) = %v, want %v", tt.params, got, tt.want)
				break
			}
		}
	}
}

func TestDecodeArrowKey(t *testing.T) {
	event := decodeArrowKey("", 'A')
	if event.Type != KeyUp || event.Ctrl || event.Alt || event.Shift {
		t.Fatalf("decodeArrowKey(\"\", 'A') = %+v, want plain KeyUp", event)
	}

	event = decodeArrowKey("1;6", 'C')
	if event.Type != KeyRight || !event.Shift || !event.Ctrl || event.Alt {
		t.Fatalf("decodeArrowKey(\"1;6\", 'C') = %+v, want Shift+Ctrl KeyRight", event)
	}
}

func TestDecodeTildeKey(t *testing.T) {
	event, ok := decodeTildeKey("3")
	if !ok || event.Type != KeyDelete {
		t.Fatalf("decodeTildeKey(\"3\") = %+v, %v, want KeyDelete, true", event, ok)
	}

	event, ok = decodeTildeKey("5;3")
	if !ok || event.Type != KeyPgUp || !event.Alt || event.Shift || event.Ctrl {
		t.Fatalf("decodeTildeKey(\"5;3\") = %+v, %v, want Alt+KeyPgUp, true", event, ok)
	}

	event, ok = decodeTildeKey("27;5;13")
	if !ok || event.Type != KeyEnter || !event.Ctrl || event.Shift || event.Alt {
		t.Fatalf("decodeTildeKey(\"27;5;13\") = %+v, %v, want Ctrl+KeyEnter, true", event, ok)
	}

	if _, ok := decodeTildeKey("999"); ok {
		t.Fatal("decodeTildeKey(\"999\") should be unrecognized")
	}
}

func TestDecodeKittyKey(t *testing.T) {
	event, ok := decodeKittyKey("97")
	if !ok || event.Type != KeyRune || event.Rune != 'a' {
		t.Fatalf("decodeKittyKey(\"97\") = %+v, %v, want rune 'a', true", event, ok)
	}

	event, ok = decodeKittyKey("97;5")
	if !ok || event.Type != KeyRune || event.Rune != 'a' || !event.Ctrl {
		t.Fatalf("decodeKittyKey(\"97;5\") = %+v, %v, want Ctrl+'a', true", event, ok)
	}

	event, ok = decodeKittyKey("13;2")
	if !ok || event.Type != KeyEnter || !event.Shift {
		t.Fatalf("decodeKittyKey(\"13;2\") = %+v, %v, want Shift+KeyEnter, true", event, ok)
	}

	if _, ok := decodeKittyKey(""); ok {
		t.Fatal("decodeKittyKey(\"\") should be unrecognized")
	}
}

func TestIsKittyQueryReply(t *testing.T) {
	if !isKittyQueryReply([]byte("\x1b[?31u")) {
		t.Fatal("expected a valid Kitty query reply to be recognized")
	}
	if isKittyQueryReply([]byte("\x1b[A")) {
		t.Fatal("expected an arrow-key sequence not to be mistaken for a Kitty reply")
	}
	if isKittyQueryReply(nil) {
		t.Fatal("expected an empty buffer not to be mistaken for a Kitty reply")
	}
}
