package parser

import (
	l "github.com/maxvanasten/gscp/lexer"
)

type nodeSpan struct {
	line        int
	col         int
	endLine     int
	endCol      int
	startOffset int
	endOffset   int
	valid       bool
}

func spanFromToken(token l.Token) nodeSpan {
	return nodeSpan{
		line:        token.Line,
		col:         token.Col,
		endLine:     token.EndLine,
		endCol:      token.EndCol,
		startOffset: token.StartOffset,
		endOffset:   token.EndOffset,
		valid:       true,
	}
}

func spanFromTokens(tokens []l.Token) nodeSpan {
	if len(tokens) == 0 {
		return nodeSpan{}
	}
	start := spanFromToken(tokens[0])
	end := spanFromToken(tokens[len(tokens)-1])
	return mergeSpan(start, end)
}

func spanFromNode(node Node) nodeSpan {
	if !node.spanValid {
		return nodeSpan{}
	}
	return nodeSpan{
		line:        node.Line,
		col:         node.Col,
		endLine:     node.endLine,
		endCol:      node.endCol,
		startOffset: node.startOffset,
		endOffset:   node.endOffset,
		valid:       true,
	}
}

func spanFromNodes(nodes []Node) nodeSpan {
	if len(nodes) == 0 {
		return nodeSpan{}
	}
	start := spanFromNode(nodes[0])
	end := spanFromNode(nodes[len(nodes)-1])
	return mergeSpan(start, end)
}

func mergeSpan(start nodeSpan, end nodeSpan) nodeSpan {
	if !start.valid {
		return end
	}
	if !end.valid {
		return start
	}
	return nodeSpan{
		line:        start.line,
		col:         start.col,
		endLine:     end.endLine,
		endCol:      end.endCol,
		startOffset: start.startOffset,
		endOffset:   end.endOffset,
		valid:       true,
	}
}

func nodeWithSpan(nodeType string, data NodeData, children []Node, span nodeSpan) Node {
	length := 0
	if span.valid {
		length = span.endOffset - span.startOffset + 1
	}
	return Node{
		Type:        nodeType,
		Data:        data,
		Children:    children,
		Line:        span.line,
		Col:         span.col,
		Length:      length,
		startOffset: span.startOffset,
		endOffset:   span.endOffset,
		endLine:     span.endLine,
		endCol:      span.endCol,
		spanValid:   span.valid,
	}
}
