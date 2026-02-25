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
./gscp input.gsc
```

gscp outputs the abstract syntax tree on STDOUT accesibility. I like to read the output of gscp by using `jq` and `bat` to format the json nicely, by running gscp like so: `./gscp input.gsc | jq | bat -l json`

## Application

gscp can be the backbone for future projects like a gsc language server or doing complex code analysis on the original codebase.
