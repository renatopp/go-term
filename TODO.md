
[x] Review
  [x] Add OnKeyboard
  [x] Refactor OnResize
  [x] Add OnEvent
  [x] Add signal listening (SIGWINCH, SIGINT, etc.)
[ ] UI
  [x] Decide layouting (~~retained~~ vs immediate mode)
  [x] Program/layout event loop
  [x] Buffer structure
  [x] Mouse and Keyboard events
    [x] Global raw stdin listeners and fanout
  [x] Layouting
  [ ] Widgets
    [ ] Content
      [x] Boxes
      [x] Lists
      [ ] Tables
      [ ] Trees
      [ ] Graphs
    [ ] Dynamic
      [x] Spinners
      [x] Progress bars
      [x] Counters
      [ ] Gauges
    [ ] Inputs
      [x] Confirm prompt
      [x] Input prompt
      [x] Select prompt
      [ ] Multi-select prompt
      [x] Password prompt
      [x] Number prompt
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

[ ] Review dos widges:
  [ ] Adicionar width and height
  [ ] Default styling
  [ ] Make widgets standalone

[ ] Background and Forecolor Detection (dark/light mode)
[ ] Clipboard read/write (OSC 52)
[ ] Scroll region
[ ] Draw only visible part
[ ] Hyperlinks (OSC 8 clickable links)
[ ] Mouse to component feedback
[ ] Paste detection
