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

gscp outputs a JSON object on STDOUT containing both the AST and diagnostics. You can format the JSON nicely with `jq` and `bat`, for example: `./gscp -p input.gsc | jq | bat -l json` or `./gscp -p input.gsc | jq .ast`

## Documentation

For a detailed walk-through of how gscp works internally, see [`docs/gscp-internals.md`](./docs/gscp-internals.md).

## Application

gscp can be the backbone for future projects like a gsc language server or doing complex code analysis on the original codebase.

## Latest test results

You can find the latest test results [here](./TESTS.md)
