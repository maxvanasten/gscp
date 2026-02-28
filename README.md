# gscp

gscp is a small lexer and parser for the .gsc scripting language used in older Call of Duty games. gscp takes raw .gsc code and turns it into an abstract syntax tree. You can find a small demo input and output file in `./demo`.

## Installation

### Download binary
1. Go to the releases tab and download one of the gscp binaries

### Build from source
```bash
#Clone the repo
git clone https://github.com/maxvanasten/gscp
#Build the parser
cd ./gscp/ && go build
#Run the parser
./gscp input.gsc
```

## Usage
```bash

# Parse GSC file into AST and output the result on STDOUT
./gscp -p input_file.gsc
# Generate GSC file from AST JSON
./gscp -g input_ast.json

```

gscp outputs the abstract syntax tree in JSON format on STDOUT. You can format the JSON nicely with `jq` and `bat`, for example: `./gscp -p input.gsc | jq | bat -l json`

## Application

gscp can be the backbone for future projects like a gsc language server or doing complex code analysis on the original codebase.

## Latest test results

```
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
ok      github.com/maxvanasten/gscp/lexer       0.005s
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
ok      github.com/maxvanasten/gscp/parser      0.006s
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
ok      github.com/maxvanasten/gscp/generator   0.006s
```
