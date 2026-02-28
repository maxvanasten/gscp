package diagnostics

type Diagnostic struct {
	Message  string
	Line     int
	Col      int
	EndLine  int
	EndCol   int
	Severity string
}

func New(message string, line int, col int, endLine int, endCol int, severity string) Diagnostic {
	return Diagnostic{message, line, col, endLine, endCol, severity}
}
