package diagnostics

type Diagnostic struct {
	Message  string
	Line     int
	Col      int
	Severity string
}

func New(message string, line int, col int, severity string) Diagnostic {
	return Diagnostic{message, line, col, severity}
}
