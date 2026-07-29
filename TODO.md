
[ ] Background and Forecolor Detection (dark/light mode)
[ ] UI
  [x] Decide layouting (~~retained~~ vs immediate mode)
  [x] Program/layout event loop
  [x] Buffer structure
  [ ] Mouse and Keyboard events
    [ ] Global raw stdin listeners and fanout
  [ ] Layouting
  [ ] Widgets
    [ ] Content
      [ ] Boxes
      [ ] Lists
      [ ] Tables
      [ ] Trees
      [ ] Graphs
    [ ] Dynamic
      [ ] Spinners
      [ ] Progress bars
      [ ] Counters
      [ ] Timers
      [ ] Badges
      [ ] Gauges
    [ ] Inputs
      [ ] Confirm prompt
      [ ] Input prompt
      [ ] Select prompt
      [ ] Multi-select prompt
      [ ] Password prompt
      [ ] Number prompt
      [ ] Fuzzy search prompt
[x] Styles
  [x] Colors Types and Initialization (ansi, 256, 24-bit, hex)
  [x] Color Constants
  [x] Style Type
  [x] Style Constants
[x] Cursor
  [x] Hide cursor
  [x] Show cursor
  [x] Save cursor position
  [x] Restore cursor position
  [x] Move cursor (up, down, left, right, to)
  [x] Get cursor position
[x] Screen
  [x] Get size (width, height)
  [x] Alternate screen (enter, exit)
  [x] Clear screen
  [x] Clear line
  [x] Clear from cursor to end of line
  [x] Clear from cursor to end of screen
  [x] Resize events

  


- Scroll region set/reset, scroll up/down N lines
- Set title (window/tab title)
- Query terminal capabilities (color support, size via ioctl/ANSI query)
- Save/restore entire screen buffer
- Set background/foreground default colors, reset colors
- Bell/beep
- Synchronized output (begin/end sync for flicker-free redraw)
- Detect terminal type (xterm, tmux, Windows Terminal, etc.)

- Set cursor shape/style (block, underline, bar, blinking/steady)
- Cursor visibility query (is it currently hidden?)
- Move to column only (CHA)
- Move relative by N (not just single step)
- Cursor color

- Input handling: raw mode toggle, read keypresses/escape sequences, mouse events (click, drag, scroll), paste detection (bracketed paste)
- Styling: color palettes (16/256/truecolor), text attributes (bold, italic, underline, strikethrough), style composition/reset helpers
- Terminal detection: is TTY, is piped, supports ANSI (esp. Windows legacy consoles)
- Signal handling: SIGWINCH (resize), SIGINT cleanup (restore terminal state on exit)
- Buffered/double-buffered rendering to reduce flicker
- Unicode/wide-character width handling (for CJK, emoji) for accurate layout
- Hyperlinks (OSC 8 clickable links)
- Clipboard read/write (OSC 52)
- Cross-platform abstraction (Windows Console API vs ANSI/VT for Unix)
