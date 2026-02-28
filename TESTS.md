# Latest Test Results

```text
?   	github.com/maxvanasten/gscp	[no test files]
=== RUN   TestDiagnosticsLexerUnterminatedString
--- PASS: TestDiagnosticsLexerUnterminatedString (0.00s)
=== RUN   TestDiagnosticsLexerInvalidToken
--- PASS: TestDiagnosticsLexerInvalidToken (0.00s)
=== RUN   TestDiagnosticsLexerUnterminatedBlockComment
--- PASS: TestDiagnosticsLexerUnterminatedBlockComment (0.00s)
=== RUN   TestDiagnosticsLexerSkipsLineComment
--- PASS: TestDiagnosticsLexerSkipsLineComment (0.00s)
=== RUN   TestDiagnosticsLexerSkipsBlockComment
--- PASS: TestDiagnosticsLexerSkipsBlockComment (0.00s)
=== RUN   TestDiagnosticsLexerSymbolStarts
--- PASS: TestDiagnosticsLexerSymbolStarts (0.00s)
=== RUN   TestDiagnosticsWaitFunctionCall
--- PASS: TestDiagnosticsWaitFunctionCall (0.00s)
=== RUN   TestDiagnosticsIncrementOperator
--- PASS: TestDiagnosticsIncrementOperator (0.00s)
=== RUN   TestDiagnosticsNestedBlockComment
--- PASS: TestDiagnosticsNestedBlockComment (0.00s)
=== RUN   TestDiagnosticsDoubleNegation
--- PASS: TestDiagnosticsDoubleNegation (0.00s)
=== RUN   TestDiagnosticsPercentPrefix
--- PASS: TestDiagnosticsPercentPrefix (0.00s)
=== RUN   TestDiagnosticsMissingIncludePath
--- PASS: TestDiagnosticsMissingIncludePath (0.00s)
=== RUN   TestDiagnosticsMissingWaitDuration
--- PASS: TestDiagnosticsMissingWaitDuration (0.00s)
=== RUN   TestDiagnosticsMissingUnaryOperand
--- PASS: TestDiagnosticsMissingUnaryOperand (0.00s)
=== RUN   TestDiagnosticsOperatorMissingLeftOperand
--- PASS: TestDiagnosticsOperatorMissingLeftOperand (0.00s)
=== RUN   TestDiagnosticsOperatorMissingRightOperand
--- PASS: TestDiagnosticsOperatorMissingRightOperand (0.00s)
=== RUN   TestDiagnosticsAssignmentMissingLeftSide
--- PASS: TestDiagnosticsAssignmentMissingLeftSide (0.00s)
=== RUN   TestDiagnosticsAssignmentTargetMustBeVariable
--- PASS: TestDiagnosticsAssignmentTargetMustBeVariable (0.00s)
=== RUN   TestDiagnosticsMissingClosingParen
--- PASS: TestDiagnosticsMissingClosingParen (0.00s)
=== RUN   TestDiagnosticsMissingClosingBracket
--- PASS: TestDiagnosticsMissingClosingBracket (0.00s)
=== RUN   TestDiagnosticsMissingClosingCurly
--- PASS: TestDiagnosticsMissingClosingCurly (0.00s)
=== RUN   TestDiagnosticsUnexpectedOpenCurly
--- PASS: TestDiagnosticsUnexpectedOpenCurly (0.00s)
=== RUN   TestDiagnosticsUnexpectedCloseParen
--- PASS: TestDiagnosticsUnexpectedCloseParen (0.00s)
=== RUN   TestDiagnosticsUnexpectedCloseBracket
--- PASS: TestDiagnosticsUnexpectedCloseBracket (0.00s)
=== RUN   TestDiagnosticsUnexpectedCloseCurly
--- PASS: TestDiagnosticsUnexpectedCloseCurly (0.00s)
PASS
ok  	github.com/maxvanasten/gscp/diagnostics	0.015s
=== RUN   Test_Generate_VariableReference
--- PASS: Test_Generate_VariableReference (0.00s)
=== RUN   Test_Generate_String
--- PASS: Test_Generate_String (0.00s)
=== RUN   Test_Generate_Boolean
--- PASS: Test_Generate_Boolean (0.00s)
=== RUN   Test_Generate_Number
--- PASS: Test_Generate_Number (0.00s)
=== RUN   Test_Generate_Simple_Expression
--- PASS: Test_Generate_Simple_Expression (0.00s)
=== RUN   Test_Generate_Complex_Expression
--- PASS: Test_Generate_Complex_Expression (0.00s)
=== RUN   Test_Generate_Assignment
--- PASS: Test_Generate_Assignment (0.00s)
=== RUN   Test_Generate_Unary_Expression
--- PASS: Test_Generate_Unary_Expression (0.00s)
=== RUN   Test_Generate_Lhs
--- PASS: Test_Generate_Lhs (0.00s)
=== RUN   Test_Generate_Rhs
--- PASS: Test_Generate_Rhs (0.00s)
=== RUN   Test_Generate_Array_Literal
--- PASS: Test_Generate_Array_Literal (0.00s)
=== RUN   Test_Generate_Array_Literal_Multiple
--- PASS: Test_Generate_Array_Literal_Multiple (0.00s)
=== RUN   Test_Generate_ArrayLiteral
--- PASS: Test_Generate_ArrayLiteral (0.00s)
=== RUN   Test_Generate_Array_Indexing
--- PASS: Test_Generate_Array_Indexing (0.00s)
=== RUN   Test_Generate_Array_Indexing_Assignment
--- PASS: Test_Generate_Array_Indexing_Assignment (0.00s)
=== RUN   Test_Generate_IncludeStatement
--- PASS: Test_Generate_IncludeStatement (0.00s)
=== RUN   Test_Generate_WaitStatement
--- PASS: Test_Generate_WaitStatement (0.00s)
=== RUN   Test_Generate_BreakStatement
--- PASS: Test_Generate_BreakStatement (0.00s)
=== RUN   Test_Generate_ReturnStatement
--- PASS: Test_Generate_ReturnStatement (0.00s)
=== RUN   Test_Generate_FunctionCall
--- PASS: Test_Generate_FunctionCall (0.00s)
=== RUN   Test_Generate_FunctionCall_WithArgs
--- PASS: Test_Generate_FunctionCall_WithArgs (0.00s)
=== RUN   Test_Generate_FunctionCall_Method
--- PASS: Test_Generate_FunctionCall_Method (0.00s)
=== RUN   Test_Generate_FunctionCall_Thread
--- PASS: Test_Generate_FunctionCall_Thread (0.00s)
=== RUN   Test_Generate_FunctionCall_MethodThread
--- PASS: Test_Generate_FunctionCall_MethodThread (0.00s)
=== RUN   Test_Generate_FunctionCall_Path
--- PASS: Test_Generate_FunctionCall_Path (0.00s)
=== RUN   Test_Generate_FunctionDeclaration
--- PASS: Test_Generate_FunctionDeclaration (0.00s)
=== RUN   Test_Generate_Args
--- PASS: Test_Generate_Args (0.00s)
=== RUN   Test_Generate_Scope
--- PASS: Test_Generate_Scope (0.00s)
=== RUN   Test_Generate_VectorLiteral
--- PASS: Test_Generate_VectorLiteral (0.00s)
=== RUN   Test_Generate_ForInit
--- PASS: Test_Generate_ForInit (0.00s)
=== RUN   Test_Generate_ForCondition
--- PASS: Test_Generate_ForCondition (0.00s)
=== RUN   Test_Generate_ForPost
--- PASS: Test_Generate_ForPost (0.00s)
=== RUN   Test_Generate_ForLoop
--- PASS: Test_Generate_ForLoop (0.00s)
=== RUN   Test_Generate_Condition
--- PASS: Test_Generate_Condition (0.00s)
=== RUN   Test_Generate_IfStatement
--- PASS: Test_Generate_IfStatement (0.00s)
=== RUN   Test_Generate_ElseClause
--- PASS: Test_Generate_ElseClause (0.00s)
=== RUN   Test_Generate_WhileLoop
--- PASS: Test_Generate_WhileLoop (0.00s)
=== RUN   Test_Generate_ForeachVars
--- PASS: Test_Generate_ForeachVars (0.00s)
=== RUN   Test_Generate_ForeachIter
--- PASS: Test_Generate_ForeachIter (0.00s)
=== RUN   Test_Generate_ForeachLoop
--- PASS: Test_Generate_ForeachLoop (0.00s)
=== RUN   Test_Generate_SwitchExpr
--- PASS: Test_Generate_SwitchExpr (0.00s)
=== RUN   Test_Generate_CaseClause
--- PASS: Test_Generate_CaseClause (0.00s)
=== RUN   Test_Generate_DefaultClause
--- PASS: Test_Generate_DefaultClause (0.00s)
=== RUN   Test_Generate_SwitchStatement
--- PASS: Test_Generate_SwitchStatement (0.00s)
PASS
ok  	github.com/maxvanasten/gscp/generator	0.011s
=== RUN   TestLexerSymbol
--- PASS: TestLexerSymbol (0.00s)
=== RUN   TestLexerNumber
--- PASS: TestLexerNumber (0.00s)
=== RUN   TestLexerString
--- PASS: TestLexerString (0.00s)
=== RUN   TestLexerTerminator
--- PASS: TestLexerTerminator (0.00s)
=== RUN   TestLexerComma
--- PASS: TestLexerComma (0.00s)
=== RUN   TestLexerNewline
--- PASS: TestLexerNewline (0.00s)
=== RUN   TestLexerOpenParen
--- PASS: TestLexerOpenParen (0.00s)
=== RUN   TestLexerCloseParen
--- PASS: TestLexerCloseParen (0.00s)
=== RUN   TestLexerOpenBracket
--- PASS: TestLexerOpenBracket (0.00s)
=== RUN   TestLexerCloseBracket
--- PASS: TestLexerCloseBracket (0.00s)
=== RUN   TestLexerOpenCurly
--- PASS: TestLexerOpenCurly (0.00s)
=== RUN   TestLexerCloseCurly
--- PASS: TestLexerCloseCurly (0.00s)
=== RUN   TestLexerAssignment
--- PASS: TestLexerAssignment (0.00s)
=== RUN   TestLexerCompoundAssignment
--- PASS: TestLexerCompoundAssignment (0.00s)
=== RUN   TestLexerArithmeticOperators
--- PASS: TestLexerArithmeticOperators (0.00s)
=== RUN   TestLexerComparisonOperators
--- PASS: TestLexerComparisonOperators (0.00s)
=== RUN   TestLexerLogicalOperators
--- PASS: TestLexerLogicalOperators (0.00s)
PASS
ok  	github.com/maxvanasten/gscp/lexer	0.008s
=== RUN   Test_Variable_Reference
--- PASS: Test_Variable_Reference (0.00s)
=== RUN   Test_String
--- PASS: Test_String (0.00s)
=== RUN   Test_Boolean
--- PASS: Test_Boolean (0.00s)
=== RUN   Test_Number
--- PASS: Test_Number (0.00s)
=== RUN   Test_Simple_Expression
--- PASS: Test_Simple_Expression (0.00s)
=== RUN   Test_Complex_Expression
--- PASS: Test_Complex_Expression (0.00s)
=== RUN   Test_Logical_Expression_Precedence
--- PASS: Test_Logical_Expression_Precedence (0.00s)
=== RUN   Test_Complex_Math_Expression
--- PASS: Test_Complex_Math_Expression (0.00s)
=== RUN   Test_Variable_Assignment
--- PASS: Test_Variable_Assignment (0.00s)
=== RUN   Test_Compound_Assignment
--- PASS: Test_Compound_Assignment (0.00s)
=== RUN   Test_Unary_Expression
--- PASS: Test_Unary_Expression (0.00s)
=== RUN   Test_Array_Literal_Empty
--- PASS: Test_Array_Literal_Empty (0.00s)
=== RUN   Test_Array_Literal_Multiple
--- PASS: Test_Array_Literal_Multiple (0.00s)
=== RUN   Test_Array_Indexing
--- PASS: Test_Array_Indexing (0.00s)
=== RUN   Test_Array_Index_Assignment
--- PASS: Test_Array_Index_Assignment (0.00s)
=== RUN   Test_Array_Index_Compound_Assignment
--- PASS: Test_Array_Index_Compound_Assignment (0.00s)
=== RUN   Test_Vector_Literal
--- PASS: Test_Vector_Literal (0.00s)
=== RUN   Test_Function_Call
--- PASS: Test_Function_Call (0.00s)
=== RUN   Test_Namespace_Function_Call
--- PASS: Test_Namespace_Function_Call (0.00s)
=== RUN   Test_Method_Function_Call
--- PASS: Test_Method_Function_Call (0.00s)
=== RUN   Test_Threaded_Function_Call
--- PASS: Test_Threaded_Function_Call (0.00s)
=== RUN   Test_Function_Declaration
--- PASS: Test_Function_Declaration (0.00s)
=== RUN   Test_For_Loop_Infinite
--- PASS: Test_For_Loop_Infinite (0.00s)
=== RUN   Test_For_Loop_Common
--- PASS: Test_For_Loop_Common (0.00s)
=== RUN   Test_If_Else
--- PASS: Test_If_Else (0.00s)
=== RUN   Test_While_Loop
--- PASS: Test_While_Loop (0.00s)
=== RUN   Test_Foreach_Loop
--- PASS: Test_Foreach_Loop (0.00s)
=== RUN   Test_Switch_Case_Default
--- PASS: Test_Switch_Case_Default (0.00s)
=== RUN   Test_Return_Statement
--- PASS: Test_Return_Statement (0.00s)
=== RUN   Test_IncludeStatement
--- PASS: Test_IncludeStatement (0.00s)
=== RUN   Test_WaitStatement
--- PASS: Test_WaitStatement (0.00s)
=== RUN   Test_Function_Calls
--- PASS: Test_Function_Calls (0.00s)
=== RUN   Test_Function_Call_Complex_Args
--- PASS: Test_Function_Call_Complex_Args (0.00s)
PASS
ok  	github.com/maxvanasten/gscp/parser	0.010s
```
