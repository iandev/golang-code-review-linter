package main

import (
	"fmt"
	"strings"
	"unicode"

	goparser "./parser"
	"github.com/antlr/antlr4/runtime/Go/antlr" // The antlr library
)

type styleLinterListener struct {
	*goparser.BaseGolangListener
}

func (l *styleLinterListener) EnterConstSpec(ctx *goparser.ConstSpecContext) {
	decl := ctx.GetChild(0).(antlr.ParseTree).GetText()
	if strings.Contains(decl, "_") || strings.ToUpper(decl) == decl {
		line := ctx.GetStart().GetLine()
		fmt.Printf("Invalid constant declaration %s on line %d\n", decl, line)
	}
}

func (s *styleLinterListener) EnterReceiver(ctx *goparser.ReceiverContext) {
	decl := ctx.GetChild(0).GetChild(1).GetChild(0).GetChild(0).(antlr.ParseTree).GetText()
	typ := ctx.GetChild(0).GetChild(1).GetChild(0).GetChild(1).(antlr.ParseTree).GetText()

	declLen := len(decl)
	typeLen := len(strings.FieldsFunc(strings.Replace(typ, "*", "", 1), func(c rune) bool {
		return unicode.IsUpper(c)
	}))

	if declLen > typeLen {
		line := ctx.GetStart().GetLine()
		fmt.Printf("Receiver declaration %s is too long on line %d\n", decl, line)
	}
}

var src = `package foo

	const Invalid_Constant = "123"

	const (
		Invalid___Constant = 1
		ValidConstant = 2
	)

	const ACONSTANT = "123"

	type fooStruct1 struct {}

	func (f *fooStruct1) FooBar() {}
	func (fs *fooStruct1) FooBar2() {}
	func (fsb *fooStruct1) FooBar3() {}
	func (fsba fooStruct1) FooBar3() {}
`

func main() {
	is := antlr.NewInputStream(src)

	lexer := goparser.NewGolangLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := goparser.NewGolangParser(stream)
	p.BaseRecognizer.RemoveErrorListeners()

	tree := p.SourceFile()
	antlr.ParseTreeWalkerDefault.Walk(&styleLinterListener{}, tree)
}
