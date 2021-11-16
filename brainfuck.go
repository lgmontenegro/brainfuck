package brainfuck

import (
	"errors"
	"fmt"
	loo "lgmontenegro/brainfuck/internal/domain/loopers"
	mov "lgmontenegro/brainfuck/internal/domain/movers"
	op "lgmontenegro/brainfuck/internal/domain/operators"
	ret "lgmontenegro/brainfuck/internal/domain/retrievers"
	"strings"
)

type BrainfuckCompiler struct {
	Instructions map[string][]byte
	Operators    map[byte]func(byte) byte
	Movers       map[byte]func(int) int
	Loopers      map[byte]func(int) int
	Retrievers   map[byte]func(map[int]byte, int) string
	stackControl brainfuckStackControl
	flowControl  brainfuckLoopControl
}

type brainfuckStackControl struct {
	memAlloc    map[int]byte
	newMemAlloc bool
	pointer     int
}

type brainfuckLoopControl struct {
	loop                  int
	loopStatus            map[int]bool
	loopInstructions      []stackedInstruction
	loopInstructionsSteps map[int]int
}

type stackedInstruction struct {
	loop            int
	step            int
	instructionType string
	instruction     byte
}

func NewCompiler() (compiler BrainfuckCompiler) {
	operators := map[byte]func(byte) byte{
		byte('+'): op.Add,
		byte('-'): op.Minus,
	}

	retrievers := map[byte]func(map[int]byte, int) string{
		byte('.'): ret.Retriever,
	}

	movers := map[byte]func(int) int{
		byte('>'): mov.MoveRight,
		byte('<'): mov.MoveLeft,
	}

	loopers := map[byte]func(int) int{
		byte('['): loo.LoopInit,
		byte(']'): loo.LoopClose,
	}
	compiler.Instructions = map[string][]byte{
		"operator":  {byte('+'), byte('-')},
		"mover":     {byte('<'), byte('>')},
		"looper":    {byte(']'), byte('[')},
		"retriever": {byte('.')},
		"ignore":    {10, 32},
	}
	compiler.Operators = operators
	compiler.Retrievers = retrievers
	compiler.Movers = movers
	compiler.Loopers = loopers
	compiler.stackControl.memAlloc = map[int]byte{0: 0}
	compiler.stackControl.newMemAlloc = true
	compiler.stackControl.pointer = 0
	compiler.flowControl.loop = 0

	return compiler
}

func (b *BrainfuckCompiler) Compile(instruction byte) error {

	instructionType, err := b.getInstructionType(instruction)
	if err != nil {
		return err
	}

	instructionType = strings.ToLower(instructionType)

	if b.inLoop() && !b.flowControl.loopStatus[b.flowControl.loop] {
		step := 1
		for b.inLoop() {
			b.executeLoopInstruction(step)
			b.loopControl()
			step = b.walk(step)
		}
	}

	b.readInstruction(instructionType, instruction)

	return nil
}

func (b *BrainfuckCompiler) getInstructionType(instruction byte) (instructionType string, err error) {
	for instructionType, instructionsAvailable := range b.Instructions {
		for _, instructionAvailable := range instructionsAvailable {
			if instructionAvailable == instruction {
				return instructionType, nil
			}
		}
	}
	errMessage := fmt.Sprintf("instruction doesn't exist: %d", instruction)
	return "", errors.New(errMessage)
}

func (b *BrainfuckCompiler) readInstruction(instructionType string, instruction byte) {
	if b.inLoop() {
		if instructionType == "looper" {
			b.openCloseLoop(b.Loopers[instruction](b.flowControl.loop))
		}
		b.stackLoopInstruction(instructionType, instruction)
		b.executeInstruction(instructionType, instruction)
		b.loopControl()
		return
	}

	b.executeInstruction(instructionType, instruction)
}

func (b *BrainfuckCompiler) loopControl() {

	if b.stackControl.memAlloc[b.stackControl.pointer] == 0 && !b.stackControl.newMemAlloc {
		delete(b.flowControl.loopInstructionsSteps, b.flowControl.loop)
		b.flowControl.loopInstructions = []stackedInstruction{}
		delete(b.flowControl.loopStatus, b.flowControl.loop)

		b.flowControl.loop = b.flowControl.loop - 1
	}
}

func (b *BrainfuckCompiler) executeLoopInstruction(step int) {
	for _, instructions := range b.flowControl.loopInstructions {
		if instructions.loop == b.flowControl.loop && instructions.step == step {
			b.executeInstruction(instructions.instructionType, instructions.instruction)
			return
		}
	}
}

func (b *BrainfuckCompiler) walk(step int) (newStep int) {
	newStep = step + 1
	if newStep > b.flowControl.loopInstructionsSteps[b.flowControl.loop] {
		return 1
	}

	return newStep
}

func (b *BrainfuckCompiler) executeInstruction(instructionType string, instruction byte) {
	switch instructionType {
	case "operator":
		b.stackControl.memAlloc[b.stackControl.pointer] = b.Operators[instruction](b.stackControl.memAlloc[b.stackControl.pointer])
		b.stackControl.newMemAlloc = false
	case "retriever":
		fmt.Println(b.Retrievers[instruction](b.stackControl.memAlloc, b.stackControl.pointer))
	case "mover":
		b.stackControl.pointer = b.Movers[instruction](b.stackControl.pointer)
		b.movePointer()
	case "looper":
		loopControl := b.Loopers[instruction](b.flowControl.loop)
		b.openCloseLoop(loopControl)
	case "ignore":
		return
	default:
		return
	}
}

func (b *BrainfuckCompiler) openCloseLoop(pointerLoop int) {
	if pointerLoop > b.flowControl.loop {
		b.flowControl.loop = pointerLoop

		_, ok := b.flowControl.loopStatus[pointerLoop]
		if !ok {
			b.flowControl.loopStatus = map[int]bool{
				pointerLoop: true,
			}
			return
		}
	}

	b.flowControl.loopStatus[pointerLoop] = false
}

func (b *BrainfuckCompiler) movePointer() {
	if b.stackControl.pointer < 0 {
		b.stackControl.pointer = 0
	}

	_, ok := b.stackControl.memAlloc[b.stackControl.pointer]
	if !ok {
		b.stackControl.memAlloc[b.stackControl.pointer] = 0
		b.stackControl.newMemAlloc = true
		return
	}
}

func (b *BrainfuckCompiler) inLoop() bool {
	return b.flowControl.loop > 0
}

func (b *BrainfuckCompiler) stackLoopInstruction(instructionType string, instruction byte) {
	var lastStep int

	if b.flowControl.loopStatus[b.flowControl.loop] && instructionType != "looper" && instructionType != "ignore" {
		if len(b.flowControl.loopInstructions)-1 >= 0 {
			lastStep = b.flowControl.loopInstructions[len(b.flowControl.loopInstructions)-1].step
		}

		stackInstruction := stackedInstruction{
			loop:            b.flowControl.loop,
			step:            lastStep + 1,
			instruction:     instruction,
			instructionType: instructionType,
		}
		b.flowControl.loopInstructions = append(b.flowControl.loopInstructions, stackInstruction)

		_, ok := b.flowControl.loopInstructionsSteps[b.flowControl.loop]
		if !ok {
			b.flowControl.loopInstructionsSteps = map[int]int{
				b.flowControl.loop: stackInstruction.step,
			}
			return
		}

		b.flowControl.loopInstructionsSteps[b.flowControl.loop] = stackInstruction.step
	}
}
