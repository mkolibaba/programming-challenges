package main

import (
	"bytes"
	"fmt"
)

type TUIBuilder struct {
	buffer *bytes.Buffer
}

func NewTUIBuilder() *TUIBuilder {
	return &TUIBuilder{&bytes.Buffer{}}
}

func (b *TUIBuilder) Line(text string) *TUIBuilder {
	b.buffer.WriteString(text)
	b.buffer.WriteString("\n")
	return b
}

func (b *TUIBuilder) Formatted(format string, a ...any) *TUIBuilder {
	if _, err := fmt.Fprintf(b.buffer, format, a...); err != nil {
		panic(err) // хз что тут может случиться
	}
	return b
}

func (b *TUIBuilder) EmptyLine() *TUIBuilder {
	b.buffer.WriteString("\n")
	return b
}

func (b *TUIBuilder) With(callback func(b *TUIBuilder)) *TUIBuilder {
	callback(b)
	return b
}

func (b *TUIBuilder) WithBuffer(callback func(buffer *bytes.Buffer)) *TUIBuilder {
	callback(b.buffer)
	return b
}

func (b *TUIBuilder) Build() string {
	return b.buffer.String()
}
