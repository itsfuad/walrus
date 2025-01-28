package parser

import (
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
)

func parseImplStmt(p *Parser) ast.Node {

	start := p.eat().Start

	typeName := p.expect(lexer.IDENTIFIER_TOKEN)

	p.expect(lexer.OPEN_CURLY)

	methods := make([]ast.MethodToImplement, 0)

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		IsPrivate := false

		if p.currentTokenKind() == lexer.PRIVATE_TOKEN {
			IsPrivate = true
			p.eat()
		}

		p.expect(lexer.FUNCTION_TOKEN)

		fnName := p.expect(lexer.IDENTIFIER_TOKEN)

		params, ret := parseFunctionSignature(p)

		body := parseBlock(p)

		method := ast.MethodToImplement{
			IsPrivate: IsPrivate,
			FunctionDeclStmt: ast.FunctionDeclStmt{
				Identifier: ast.IdentifierExpr{
					Name: fnName.Value,
					Location: ast.Location{
						Start: fnName.Start,
						End:   fnName.End,
					},
				},
				FunctionLiteral: ast.FunctionLiteral{
					Params:     params,
					ReturnType: ret,
					Body:       body,
					Location: ast.Location{
						Start: fnName.Start,
						End:   body.End,
					},
				},
			},
		}

		methods = append(methods, method)
	}

	end := p.expect(lexer.CLOSE_CURLY).End

	return ast.ImplStmt{
		ImplFor: ast.IdentifierExpr{
			Name: typeName.Value,
			Location: ast.Location{
				Start: typeName.Start,
				End:   typeName.End,
			},
		},
		Methods: methods,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}
