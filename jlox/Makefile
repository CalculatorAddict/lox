.PHONY: build generate jlox

build: generate
	@javac -d bin src/com/craftinginterpreters/lox/*

generate:
	@javac -d bin src/com/craftinginterpreters/tool/*
	@java -cp bin com.craftinginterpreters.tool.GenerateAst src/com/craftinginterpreters/lox

jlox:
	@java -cp bin com.craftinginterpreters.lox.Lox $(file)
