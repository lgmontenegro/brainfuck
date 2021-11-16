package brainfuck

import (
	"testing"
)

func TestNewCompiler(t *testing.T) {

	operators := map[byte]func(byte) byte{
		byte('+'): func(b byte) byte { return b + 2 },
		byte('-'): func(b byte) byte { return b - 3 },
	}

	defaultCompiler := NewCompiler()
	modifiedCompiler := defaultCompiler
	modifiedCompiler.Operators = operators

	addedSignalsCompiler := defaultCompiler
	operators['*'] = func(b byte) byte { return b * 2 }
	operators['/'] = func(b byte) byte { return b / 2 }
	addedSignalsCompiler.Operators = operators

	tests := []struct {
		name               string
		wantCompiler       BrainfuckCompiler
		wantOperatorResult byte
		operatorSignal     byte
		operatorInput      byte
	}{
		{
			name:               "common_use",
			wantCompiler:       defaultCompiler,
			wantOperatorResult: 2,
			operatorInput:      1,
			operatorSignal:     byte('+'),
		},
		{
			name:               "modified_operator",
			wantCompiler:       modifiedCompiler,
			wantOperatorResult: 3,
			operatorInput:      1,
			operatorSignal:     byte('+'),
		},
		{
			name:               "added_operator",
			wantCompiler:       addedSignalsCompiler,
			wantOperatorResult: 6,
			operatorInput:      3,
			operatorSignal:     byte('*'),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addTest := tt.wantCompiler.Operators[tt.operatorSignal](tt.operatorInput)
			if addTest != tt.wantOperatorResult {
				t.Log("result not expected", tt.wantOperatorResult, addTest)
				t.Fail()
			}
		})
	}
}

func TestBrainfuckCompiler_Compile(t *testing.T) {
	type fields struct {
		Instructions map[string][]byte
		Operators    map[byte]func(byte) byte
		Movers       map[byte]func(int) int
		Loopers      map[byte]func(int) int
		Retrievers   map[byte]func(map[int]byte, int) string
		stackControl brainfuckStackControl
		flowControl  brainfuckLoopControl
	}
	type args struct {
		instruction byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BrainfuckCompiler{
				Instructions: tt.fields.Instructions,
				Operators:    tt.fields.Operators,
				Movers:       tt.fields.Movers,
				Loopers:      tt.fields.Loopers,
				Retrievers:   tt.fields.Retrievers,
				stackControl: tt.fields.stackControl,
				flowControl:  tt.fields.flowControl,
			}
			if err := b.Compile(tt.args.instruction); (err != nil) != tt.wantErr {
				t.Errorf("BrainfuckCompiler.Compile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
