package term

// func Width(s string) int {
// 	w := 0

// 	for _, r := range s {
// 		switch {
// 		case r == '\x1b':
// 			// skip ANSI (requires parser)

// 		case unicode.Is(unicode.Mn, r):
// 			// combining mark
// 			continue

// 		case r == 0x200D: // ZWJ
// 			continue

// 		case r >= 0xFE00 && r <= 0xFE0F:
// 			// variation selectors
// 			continue

// 		case isWide(r):
// 			w += 2

// 		default:
// 			w++
// 		}
// 	}

// 	return w
// }

// func isWide(r rune) bool {
// 	switch {
// 	case r >= 0x1100 && r <= 0x115F:
// 		return true
// 	case r >= 0x2329 && r <= 0x232A:
// 		return true
// 	case r >= 0x2E80 && r <= 0xA4CF:
// 		return true
// 	case r >= 0xAC00 && r <= 0xD7A3:
// 		return true
// 	case r >= 0xF900 && r <= 0xFAFF:
// 		return true
// 	case r >= 0xFE10 && r <= 0xFE6F:
// 		return true
// 	case r >= 0xFF00 && r <= 0xFF60:
// 		return true
// 	case r >= 0xFFE0 && r <= 0xFFE6:
// 		return true
// 	case r >= 0x1F300 && r <= 0x1FAFF:
// 		return true
// 	default:
// 		return false
// 	}
// }

// func Lines(text)
