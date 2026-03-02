# gscp internals

This document explains how the gscp codebase works under the hood. It is meant to be a precise, code-accurate reference for contributors and readers who want to understand the lexer, parser, AST, generator, diagnostics, and test flow.

## Repository map

- `main.go` - CLI entrypoint. Routes between parse and generate modes.
- `lexer/lexer.go` - Lexer that turns raw GSC into tokens with position metadata.
- `parser/parser.go` - Parser that turns tokens into an AST with diagnostics.
- `generator/generator.go` - Generator that converts AST back to GSC source.
- `diagnostics/diagnostics.go` - Diagnostic data structure.
- `demo/` - Sample GSC input, AST JSON output, and regenerated GSC.
- `test.sh`, `TESTS.md` - Test runner and latest test output.

## CLI behavior and data flow

gscp has two modes, both defined in `main.go`.

### Parse mode (`-p`)

Command:

```bash
./gscp -p input_file.gsc
```

Flow:

1. Read the file into memory.
2. Run the lexer (`lexer.NewLexer`) to produce tokens and lexer diagnostics.
3. Marshal then unmarshal tokens into `[]lexer.Token`. This ensures the data is JSON-compatible and not tied to internal buffers.
4. Run the parser (`parser.Parse`) to produce the AST and parser diagnostics.
5. Output JSON to STDOUT with shape:

```json
{
  "ast": [ ...nodes... ],
  "diagnostics": [ ...diagnostics... ]
}
```

### Generate mode (`-g`)

Command:

```bash
./gscp -g input_ast.json
```

Flow:

1. Read AST JSON into `[]parser.Node`.
2. Generate GSC source by calling `generator.Generate` for each top-level node.
3. Join with newlines and print to STDOUT.

## Tokens and diagnostics

### Token structure

Token fields are defined in `lexer/lexer.go`:

```go
type Token struct {
  Type    TokenType
  Content string
  Line    int
  Col     int
  EndLine int
  EndCol  int
  StartOffset int
  EndOffset   int
}
```

Positions are 1-based and represent the start and end of the token in the input. Offsets are 0-based byte indices into the input and include the full token span.

### Token types

`TokenType` includes structural and semantic categories, used throughout parsing:

- `SYMBOL`, `NUMBER`, `STRING`
- `TERMINATOR` (`;`), `COMMA` (`,`), `COLON` (`:`), `NEWLINE` (`\n`)
- `OPEN_PAREN` `(`, `CLOSE_PAREN` `)`
- `OPEN_BRACKET` `[`, `CLOSE_BRACKET` `]`
- `OPEN_CURLY` `{`, `CLOSE_CURLY` `}`
- `ASSIGNMENT` (`=`, `+=`, `-=`, `*=`, `/=`, `&=`, `|=`)
- `OPERATOR` (`+`, `-`, `*`, `/`, `%`, `^`, `~`, `?`, `==`, `!=`, `<=`, `>=`, `&&`, `||`, `++`, `--`, `<`, `>`)

### Diagnostic structure

`diagnostics/diagnostics.go` defines:

```go
type Diagnostic struct {
  Message  string
  Line     int
  Col      int
  EndLine  int
  EndCol   int
  Severity string
}
```

Diagnostics are produced by both lexer and parser and merged in `main.go` before output. Severity values are plain strings (e.g., `"error"`).

## Lexer internals

The lexer consumes raw bytes and emits tokens while tracking positions.

### Core loop

- `NewLexer` initializes with `line=1`, `col=1` and iterates `Next()` until `EOF()`.
- `Next()` calls `HandleCharacter` for the current byte and advances the index and position using `advancePosition`.
- `EOF()` flushes any buffered symbol/number via `HandleBuffer`.

### Buffering rules

The lexer builds a buffer for symbols and numbers. The buffer is flushed when a delimiter or operator is encountered.

- If the buffer is all digits (optionally one decimal `.`), it becomes a `NUMBER` token.
- Otherwise, the buffer must pass `isSymbolStart` or it produces an `invalid token` diagnostic.

### Symbol rules (`isSymbolStart`)

Valid symbol starts include:

- Letters (`a-z`, `A-Z`), `_`, `#`.
- `.` when followed by a letter, `_`, or `#` (method-like forms).
- `::` when followed by a letter or `_` (namespace-like forms).

### Strings

- Strings are delimited by `"` and support escaping via backslash.
- If the closing quote is missing, an `unterminated string literal` diagnostic is emitted.

### Comments

The lexer recognizes three styles:

- Line comments: `// ...` to end of line.
- Line comments: `#/ ...` to end of line.
- Block comments: `/# ... #/` with nesting support.

If a block comment is not closed, the lexer emits an `unterminated block comment` diagnostic.

### Operators and assignments

Operator recognition is centralized in `handleOperatorToken`:

- `++` and `--` are emitted as `OPERATOR`.
- Compound assignments (`+=`, `-=`, `*=`, `/=`, `&=`, `|=`, `%=`) are emitted as `ASSIGNMENT`.
- Comparisons (`==`, `!=`, `<=`, `>=`) are emitted as `OPERATOR`.
- Logical (`&&`, `||`) and bitwise/arithmetics map to `OPERATOR`.
- Single `=` becomes `ASSIGNMENT`.

### Structural tokens and separators

- `(`, `)`, `[`, `]`, `{`, `}` map directly to open/close tokens.
- `;`, `,`, `:` and newline create `TERMINATOR`, `COMMA`, `COLON`, `NEWLINE` tokens.
- `:` is special-cased: `::` is treated as part of a symbol buffer, not a `COLON` token.

## Parser internals

The parser is a single-pass, token-driven parser that builds a tree of `Node` values.

### AST node structure

From `parser/parser.go`:

```go
type NodeData struct {
  VarName      string `json:"variable_name,omitempty"`
  FunctionName string `json:"function_name,omitempty"`
  Path         string `json:"path,omitempty"`
  Operator     string `json:"operator,omitempty"`
  Delay        string `json:"delay,omitempty"`
  Thread       bool   `json:"thread,omitempty"`
  Method       string `json:"method,omitempty"`
  Index        string `json:"index,omitempty"`
  Content      string `json:"content,omitempty"`
}

type Node struct {
  Type     string   `json:"type"`
  Data     NodeData `json:"data"`
  Children []Node   `json:"children,omitempty"`
  Line     int      `json:"line"`
  Col      int      `json:"col"`
  Length   int      `json:"length"`
}
```

### Parsing model

`Parse(tokens []Token)` loops over tokens and appends nodes to an output slice. It also accumulates diagnostics. The parser is permissive: it prefers to emit diagnostics and continue rather than halt.

Helper utilities:

- `tokensUntilMatchingClose` - grabs tokens until the matching closing delimiter.
- `splitTopLevel` - splits tokens on a delimiter at depth 0 (parentheses/brackets).
- `topLevelIndex` - finds a token at depth 0 (used for ternary `:`).
- `parseScope` - parses `{ ... }` into a `scope` node and emits missing-brace diagnostics.
- `parseOperatorToken` - handles unary, ternary, increment/decrement, and binary operators.

### Node types emitted

Common node types include:

- Literals: `string`, `number`, `boolean`
- References: `variable_reference` (with optional `Index`)
- Operators: `unary_expression`, `expression`, `ternary_expression`
- Statements: `assignment`, `return_statement`, `break_statement`, `wait_statement`, `include_statement`
- Callables: `function_call`, `function_declaration`, `args`
- Containers: `array_literal`, `vector_literal`, `scope`
- Control flow: `for_header`, `for_init`, `for_condition`, `for_post`, `for_loop`, `if_header`, `if_statement`, `else_header`, `else_clause`, `while_header`, `while_loop`, `foreach_header`, `foreach_loop`, `switch_header`, `switch_statement`, `case_clause`, `default_clause`, `do_header`, `do_while_loop`
- Error placeholders: `operator`, `open_curly` (used when unexpected tokens appear)

Some header nodes (`*_header`) are intermediate and are replaced by full statements when their following scope is parsed. If a header is not followed by `{ ... }`, the header node can remain in the output.

### Keywords and symbols

Keywords are not a separate token type. They are parsed from `SYMBOL` content:

- `#include` emits `include_statement` with `Path` from the following symbol.
- `wait` emits `wait_statement` when followed by a number/symbol; otherwise it can be a function call.
- `thread` emits `thread_keyword` used by function-call parsing.
- `true`/`false` emit `boolean`.
- `break` emits `break_statement`.
- `return` consumes tokens until newline/terminator and emits `return_statement` with children.
- `case` / `default` emit `case_clause` / `default_clause`.
- `else` / `do` emit header nodes for later scope parsing.

### Parentheses handling

Parentheses are interpreted based on context:

- If the previous node is `variable_reference` and its name is `for`, the parenthesized content is split into `for_init`, `for_condition`, and `for_post` segments based on top-level `;`.
- If the previous node is `variable_reference` with name `if`, `while`, `foreach`, or `switch`, the content becomes a header node.
- If the previous node is a normal `variable_reference`, the parentheses become function-call arguments.
- Otherwise, if the parenthesized content contains a top-level comma, it becomes a `vector_literal`.
- Otherwise, it is treated as grouping and parsed directly into the output.

### Brackets handling

- If the previous node is `variable_reference`, `[...]` becomes an index appended to the variable (stored in `NodeData.Index`).
- Otherwise, `[...]` becomes an `array_literal` with elements split by top-level commas.

### Curly braces and scopes

`{ ... }` always parse to a `scope` node. When a header precedes `{`, the parser emits the appropriate statement node:

- `function_call` + scope -> `function_declaration`
- `for_header` + scope -> `for_loop`
- `if_header` + scope -> `if_statement`
- `while_header` + scope -> `while_loop`
- `foreach_header` + scope -> `foreach_loop`
- `switch_header` + scope -> `switch_statement`
- `else_header` + scope -> `else_clause`
- `do_header` + scope -> `do_while_loop` (with optional `while (cond);` that follows)

If `{` is encountered unexpectedly, the parser emits an `unexpected {` diagnostic and inserts an `open_curly` placeholder node.

### Expressions and operators

Operator handling is centralized in `parseOperatorToken`:

- Unary operators: `!`, `!!`, `-`, `&`, `~`, `%`.
- Increment/decrement: `++`, `--` are converted into assignments (`x = x + 1` or `x = x - 1`).
- Ternary: `? :` splits tokens at top-level `:` into `condition`, `true_expr`, and `false_expr`.
- Binary expressions: `lhs OP rhs` where `lhs` is a prior node and `rhs` is parsed from tokens until newline/terminator.

Precedence note: `&&` has a special-case guard that prevents it from consuming a `||` at the same nesting depth, so `&&` binds more tightly than `||` in this parser.

### Assignments

Assignments require the previous node to be `variable_reference`. If not, a diagnostic is emitted and a placeholder `assignment` node is added.

For compound assignments (`+=`, `-=`, `*=`, `/=`), the parser builds an `expression` node and stores it as the assignment child.

## Generator internals

The generator walks the AST and produces formatted GSC source. It is intentionally simple and deterministic.

### Formatting rules

- `Indent` is two spaces.
- Most statement nodes emit trailing semicolons.
- `scope` nodes indent each child line and ensure function calls inside a scope end with `;`.
- Multi-line blocks (`if`, `for`, `while`, `switch`, `function_declaration`) print braces on their own lines.

### Special cases

- `joinInlineChildren` strips trailing semicolons from children and joins them with separators.
- A `variable_reference` of `#` followed by a string becomes `#"string"` when generating.
- `switch_statement` indents case/default labels one level and statements inside cases two levels.

### Node-to-source mapping highlights

- `function_call` -> `method thread path::name(args);`
- `function_declaration` -> `name(args) { ... }`
- `array_literal` -> `[a, b, c]`
- `vector_literal` -> `(a, b, c)`
- `assignment` -> `var = expr;` (indexing included when `Index` is present)

## Diagnostics: when they appear

The parser and lexer emit diagnostics for missing operands, missing delimiters, and malformed input. Examples include:

- `missing include path` when `#include` is not followed by a symbol.
- `missing wait duration` when `wait` is not followed by a duration or symbol.
- `missing unary operand`, `operator missing left-hand operand`, `operator missing right-hand operand`.
- `missing closing )`, `]`, or `}`.
- `unexpected )`, `]`, or `}`.
- `unexpected {`.

Diagnostics are never fatal to parsing; they are collected and returned alongside the AST.

## Test and demo workflow

- `demo/gobblegums.gsc` is the sample input.
- `demo/ast.json` is a captured AST output.
- `demo/generated.gsc` is generated output from the AST.
- `test.sh` runs unit tests across all packages.
- `TESTS.md` contains the latest test output snapshot.

To run tests:

```bash
./test.sh
```

## Extending the codebase safely

- Update both parser and generator when you introduce a new node type.
- Add a parser test and a generator test to lock in behavior.
- Keep tokenization rules in sync with any new syntax.
- Prefer emitting diagnostics and continuing rather than failing hard.
