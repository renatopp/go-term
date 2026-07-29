This module provides utilities for terminal tools.
It renders in immediate mode.

Rules:
- If asked for plan, answer with the plan, dont change code
- If asked for review, analysis or general question, never change the code
- Always look for the least amount of code changes to achieve the goal
- Never commit or change remote
- Always check for consistency of names, usability, and functionality.

Internal Coding Rules:
- Avoid external dependencies
- Public methods before private methods
- Public functions before private functions
- Types before functions, but struct declaraction should be followed by its factory function
- Functions that return (value, error) must have Force* function variants that return (value) ignoring the error, if sementically possible
- Unit tests for file <name> should be in file <name>_test.go, in the same package as <name>
- For builder pattern, prefer With* methods over Set* methods, returning the builder type for chaining
  - Use As* methods for booleans (eg: AsVisible(bool))
