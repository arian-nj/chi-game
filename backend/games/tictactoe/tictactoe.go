// my old code change by ai
package tictactoe

import (
	"errors"
	"fmt"
)

type Cell int32

const (
	CellEmpty Cell = 0
	CellX     Cell = 1
	CellO     Cell = 2
)

func (c Cell) IsMark() bool {
	return c == CellX || c == CellO
}

func Opponent(mark Cell) Cell {
	switch mark {
	case CellX:
		return CellO
	case CellO:
		return CellX
	default:
		return CellEmpty
	}
}

type BoardSize int

const (
	Size3 BoardSize = 3
	Size5 BoardSize = 5
)

var (
	ErrInvalidSize    = errors.New("invalid board size")
	ErrInvalidWinLen  = errors.New("invalid win length")
	ErrInvalidIndex   = errors.New("cell index out of range")
	ErrCellOccupied   = errors.New("cell already occupied")
	ErrInvalidMark    = errors.New("mark must be X or O")
	ErrInvalidCells   = errors.New("cells length does not match board size")
	ErrInvalidCellVal = errors.New("invalid cell value")
)

// WinLengthForSize returns how many marks in a row are needed to win.
func WinLengthForSize(size int) int {
	switch size {
	case 3:
		return 3
	case 5:
		return 4
	default:
		return size
	}
}

// Board is a square tic-tac-toe grid. Cells are stored row-major by index.
type Board struct {
	cells  []Cell
	size   int
	winLen int
}

// NewBoard creates an empty board for a supported size (3 or 5).
func NewBoard(size BoardSize) (*Board, error) {
	if size != Size3 && size != Size5 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidSize, size)
	}
	return NewBoardCustom(int(size), WinLengthForSize(int(size)))
}

// NewBoardCustom creates an empty board with explicit size and win length.
func NewBoardCustom(size, winLen int) (*Board, error) {
	if size < 2 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidSize, size)
	}
	if winLen < 2 || winLen > size {
		return nil, fmt.Errorf("%w: %d for size %d", ErrInvalidWinLen, winLen, size)
	}

	cells := make([]Cell, size*size)
	return &Board{cells: cells, size: size, winLen: winLen}, nil
}

// NewBoardFromCells restores a board from proto-style cell values (0=empty, 1=X, 2=O).
func NewBoardFromCells(cells []int32, size int) (*Board, error) {
	b, err := NewBoardCustom(size, WinLengthForSize(size))
	if err != nil {
		return nil, err
	}
	if len(cells) != len(b.cells) {
		return nil, ErrInvalidCells
	}
	for i, v := range cells {
		cell := Cell(v)
		if !cell.Valid() {
			return nil, fmt.Errorf("%w: %d at index %d", ErrInvalidCellVal, v, i)
		}
		b.cells[i] = cell
	}
	return b, nil
}

func (c Cell) Valid() bool {
	return c == CellEmpty || c == CellX || c == CellO
}

func (b *Board) Size() int      { return b.size }
func (b *Board) WinLength() int { return b.winLen }
func (b *Board) Len() int       { return len(b.cells) }

// Get returns the mark at index. Panics are avoided; out-of-range returns CellEmpty.
func (b *Board) Get(index int) (Cell, error) {
	if err := b.validateIndex(index); err != nil {
		return CellEmpty, err
	}
	return b.cells[index], nil
}

func (b *Board) MustGet(index int) Cell {
	cell, err := b.Get(index)
	if err != nil {
		panic(err)
	}
	return cell
}

func (b *Board) IsEmpty(index int) bool {
	cell, err := b.Get(index)
	return err == nil && cell == CellEmpty
}

// Index converts row/col to a flat cell index.
func (b *Board) Index(row, col int) (int, error) {
	if row < 0 || row >= b.size || col < 0 || col >= b.size {
		return 0, ErrInvalidIndex
	}
	return row*b.size + col, nil
}

// RowCol converts a flat cell index to row/col.
func (b *Board) RowCol(index int) (row, col int, err error) {
	if err = b.validateIndex(index); err != nil {
		return 0, 0, err
	}
	return index / b.size, index % b.size, nil
}

// Play places mark at index. On success, returns a win result when the move
// completes a line, using an incremental check from the played cell.
func (b *Board) Play(index int, mark Cell) (*WinResult, error) {
	if !mark.IsMark() {
		return nil, ErrInvalidMark
	}
	if err := b.validateIndex(index); err != nil {
		return nil, err
	}
	if b.cells[index] != CellEmpty {
		return nil, ErrCellOccupied
	}
	b.cells[index] = mark
	return b.CheckWinnerAt(index), nil
}

// HasWonAt reports whether the mark at index completes a winning line.
func (b *Board) HasWonAt(index int) bool {
	return b.CheckWinnerAt(index) != nil
}

// CheckWinnerAt checks whether the mark at index completes a winning line.
// Prefer this over CheckWinner when the last move index is known.
func (b *Board) CheckWinnerAt(index int) *WinResult {
	if err := b.validateIndex(index); err != nil {
		return nil
	}
	mark := b.cells[index]
	if !mark.IsMark() {
		return nil
	}

	row, col, err := b.RowCol(index)
	if err != nil {
		return nil
	}

	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, dir := range directions {
		count := 1 +
			b.countAlong(row, col, dir[0], dir[1], mark) +
			b.countAlong(row, col, -dir[0], -dir[1], mark)
		if count >= b.winLen {
			line := b.lineAlong(row, col, dir[0], dir[1], mark)
			return &WinResult{Winner: mark, Line: winSegment(line, index, b.winLen)}
		}
	}
	return nil
}

// AvailableMoves returns indices of all empty cells.
func (b *Board) AvailableMoves() []int {
	moves := make([]int, 0, len(b.cells))
	for i, cell := range b.cells {
		if cell == CellEmpty {
			moves = append(moves, i)
		}
	}
	return moves
}

func (b *Board) IsFull() bool {
	for _, cell := range b.cells {
		if cell == CellEmpty {
			return false
		}
	}
	return true
}

// WinResult describes a winning line on the board.
type WinResult struct {
	Winner Cell
	Line   []int
}

// CheckWinner returns the winner and winning line, or nil if there is no winner yet.
func (b *Board) CheckWinner() *WinResult {
	winLen := b.winLen

	for row := 0; row < b.size; row++ {
		for col := 0; col <= b.size-winLen; col++ {
			line := make([]int, winLen)
			for offset := range line {
				line[offset] = row*b.size + col + offset
			}
			if winner := lineWinner(b.cells, line); winner != CellEmpty {
				return &WinResult{Winner: winner, Line: line}
			}
		}
	}

	for col := 0; col < b.size; col++ {
		for row := 0; row <= b.size-winLen; row++ {
			line := make([]int, winLen)
			for offset := range line {
				line[offset] = (row+offset)*b.size + col
			}
			if winner := lineWinner(b.cells, line); winner != CellEmpty {
				return &WinResult{Winner: winner, Line: line}
			}
		}
	}

	for row := 0; row <= b.size-winLen; row++ {
		for col := 0; col <= b.size-winLen; col++ {
			mainDiag := make([]int, winLen)
			for offset := range mainDiag {
				mainDiag[offset] = (row+offset)*b.size + (col + offset)
			}
			if winner := lineWinner(b.cells, mainDiag); winner != CellEmpty {
				return &WinResult{Winner: winner, Line: mainDiag}
			}

			antiDiag := make([]int, winLen)
			for offset := range antiDiag {
				antiDiag[offset] = (row+offset)*b.size + (col + winLen - 1 - offset)
			}
			if winner := lineWinner(b.cells, antiDiag); winner != CellEmpty {
				return &WinResult{Winner: winner, Line: antiDiag}
			}
		}
	}

	return nil
}

func (b *Board) IsDraw() bool {
	return b.CheckWinner() == nil && b.IsFull()
}

func (b *Board) IsGameOver() bool {
	return b.CheckWinner() != nil || b.IsFull()
}

// Cells returns a copy of the board as proto-style int32 values.
func (b *Board) Cells() []int32 {
	out := make([]int32, len(b.cells))
	for i, cell := range b.cells {
		out[i] = int32(cell)
	}
	return out
}

// Clone returns a deep copy of the board.
func (b *Board) Clone() *Board {
	cells := make([]Cell, len(b.cells))
	copy(cells, b.cells)
	return &Board{cells: cells, size: b.size, winLen: b.winLen}
}

func (b *Board) validateIndex(index int) error {
	if index < 0 || index >= len(b.cells) {
		return ErrInvalidIndex
	}
	return nil
}

func lineWinner(cells []Cell, line []int) Cell {
	first := cells[line[0]]
	if !first.IsMark() {
		return CellEmpty
	}
	for _, idx := range line[1:] {
		if cells[idx] != first {
			return CellEmpty
		}
	}
	return first
}

func (b *Board) countAlong(row, col, dr, dc int, mark Cell) int {
	count := 0
	for i := 1; i < b.winLen; i++ {
		nr, nc := row+dr*i, col+dc*i
		if nr < 0 || nr >= b.size || nc < 0 || nc >= b.size || b.cells[nr*b.size+nc] != mark {
			break
		}
		count++
	}
	return count
}

func (b *Board) lineAlong(row, col, dr, dc int, mark Cell) []int {
	positive := make([]int, 0, b.size)
	for i := 1; ; i++ {
		nr, nc := row+dr*i, col+dc*i
		if nr < 0 || nr >= b.size || nc < 0 || nc >= b.size || b.cells[nr*b.size+nc] != mark {
			break
		}
		positive = append(positive, nr*b.size+nc)
	}

	negative := make([]int, 0, b.size)
	for i := 1; ; i++ {
		nr, nc := row-dr*i, col-dc*i
		if nr < 0 || nr >= b.size || nc < 0 || nc >= b.size || b.cells[nr*b.size+nc] != mark {
			break
		}
		negative = append(negative, nr*b.size+nc)
	}

	line := make([]int, 0, len(negative)+1+len(positive))
	for i := len(negative) - 1; i >= 0; i-- {
		line = append(line, negative[i])
	}
	line = append(line, row*b.size+col)
	line = append(line, positive...)
	return line
}

func winSegment(line []int, index, winLen int) []int {
	if len(line) <= winLen {
		return line
	}

	pos := 0
	for i, idx := range line {
		if idx == index {
			pos = i
			break
		}
	}

	start := pos
	if start > len(line)-winLen {
		start = len(line) - winLen
	}
	return line[start : start+winLen]
}
