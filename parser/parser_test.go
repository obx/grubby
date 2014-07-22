package parser_test

import (
	"github.com/grubby/grubby/ast"
	"github.com/grubby/grubby/parser"

	. "github.com/grubby/grubby/parser/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("goyacc parser", func() {
	var (
		lexer parser.RubyLexer
	)

	JustBeforeEach(func() {
		parser.DebugStatements = []string{}
		parser.Statements = make([]ast.Node, 0)
		Expect(parser.RubyParse(lexer)).To(BeSuccessful())
	})

	Describe("parsing an integer", func() {
		BeforeEach(func() {
			lexer = parser.NewLexer("5")
		})

		It("works, mostly", func() {
			Expect(parser.Statements).To(Equal([]ast.Node{
				ast.ConstantInt{Value: 5},
			}))
		})
	})

	Describe("parsing a float", func() {
		BeforeEach(func() {
			lexer = parser.NewLexer("123.4567")
		})

		It("works, mostly", func() {
			Expect(parser.Statements).To(Equal([]ast.Node{
				ast.ConstantFloat{Value: 123.4567},
			}))
		})
	})

	Describe("strings", func() {
		BeforeEach(func() {
			lexer = parser.NewLexer("'hello world'")
		})

		It("returns a SimpleString struct", func() {
			Expect(parser.Statements).To(Equal([]ast.Node{
				ast.SimpleString{Value: "'hello world'"},
			}))
		})
	})

	Describe("symbols", func() {
		BeforeEach(func() {
			lexer = parser.NewLexer(":foo")
		})

		It("returns an ast.Symbol", func() {
			Expect(parser.Statements).To(Equal([]ast.Node{
				ast.Symbol{Name: "foo"},
			}))
		})
	})

	Describe("parsing multiple lines", func() {
		BeforeEach(func() {
			lexer = parser.NewLexer(":foo\n:bar")
		})

		It("returns multiple nodes", func() {
			Expect(parser.Statements).To(Equal([]ast.Node{
				ast.Symbol{Name: "foo"},
				ast.Symbol{Name: "bar"},
			}))
		})
	})

	Describe("variables", func() {
		BeforeEach(func() {
			lexer = parser.NewLexer("foo")
		})

		It("returns a bare reference", func() {
			Expect(parser.Statements).To(Equal([]ast.Node{
				ast.BareReference{Name: "foo"},
			}))
		})
	})

	Describe("call expressions", func() {
		Context("without parens", func() {
			BeforeEach(func() {
				lexer = parser.NewLexer("puts 'foo'")
			})

			It("returns a call expression with one arg", func() {
				Expect(parser.Statements).To(Equal([]ast.Node{
					ast.CallExpression{
						Func: ast.BareReference{Name: "puts"},
						Args: []ast.Node{ast.SimpleString{Value: "'foo'"}},
					},
				}))
			})
		})

		Context("with parens", func() {
			BeforeEach(func() {
				lexer = parser.NewLexer("puts('foo', 'bar', 'baz')")
			})

			It("returns a call expression with args", func() {
				Expect(parser.Statements).To(Equal([]ast.Node{
					ast.CallExpression{
						Func: ast.BareReference{Name: "puts"},
						Args: []ast.Node{
							ast.SimpleString{Value: "'foo'"},
							ast.SimpleString{Value: "'bar'"},
							ast.SimpleString{Value: "'baz'"},
						},
					},
				}))
			})
		})

		Context("without args", func() {
			BeforeEach(func() {
				lexer = parser.NewLexer("puts()")
			})

			It("returns a call expression without args", func() {
				Expect(parser.Statements).To(Equal([]ast.Node{
					ast.CallExpression{
						Func: ast.BareReference{Name: "puts"},
						Args: []ast.Node{},
					},
				}))
			})
		})
	})

	Describe("whitespace", func() {
		BeforeEach(func() {
			lexer = parser.NewLexer(`
puts()
`)
		})

		It("parses just fine", func() {
			Expect(parser.Statements).To(Equal([]ast.Node{
				ast.CallExpression{
					Func: ast.BareReference{Name: "puts"},
					Args: []ast.Node{},
				},
			}))
		})
	})
})
