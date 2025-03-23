package main

import (
	"errors"
	"fmt"
)

type sudoku struct {
	data [9][9]int
	row  [9]int
	col  [9]int
	grid [9]int
}

func new_sudoku(data [][]int) (*sudoku, error) {
	if len(data) != 9 || len(data[0]) != 9 {
		return nil, errors.New("invalid original data length")
	}
	s := &sudoku{}
	for i := range data {
		for j := range data[i] {
			if data[i][j] > 9 || data[i][j] < 0 {
				return nil, fmt.Errorf("invalid data content, row: %d, col: %d", i, j)
			}
			if data[i][j] == 0 {
				continue
			}
			if !s.try_put(i, j, data[i][j]) {
				return nil, fmt.Errorf("conflict data content, row: %d, col: %d", i, j)
			}
		}
	}
	return s, nil
}

func (s *sudoku) get_next_arr() [][2]int {
	res := make([][2]int, 0, 81)
	for i := range s.data {
		for j := range s.data[i] {
			if s.data[i][j] == 0 {
				res = append(res, [2]int{i, j})
			}
		}
	}
	return res

}

func (s *sudoku) solve(next [][2]int, now int, res *[][9][9]int) {
	if now < 0 || now >= len(next) {
		return
	}
	candidates := s.get_available(next[now][0], next[now][1])
	if len(candidates) == 0 {
		return
	}
	for _, v := range candidates {
		if !s.try_put(next[now][0], next[now][1], v) {
			continue
		}
		if now == len(next)-1 {
			*res = append(*res, s.data)
			s.unset(next[now][0], next[now][1])
			continue
		}
		s.solve(next, now+1, res)
		s.unset(next[now][0], next[now][1])
	}
}

func (s *sudoku) get_available(row_idx, col_idx int) []int {
	if s == nil ||
		row_idx < 0 || row_idx > 8 ||
		col_idx < 0 || col_idx > 8 {
		return nil
	}
	res := make([]int, 0, 9)
	r, c, g := s.row[row_idx], s.col[col_idx], s.grid[3*(row_idx/3)+col_idx/3]
	for i := range make([]any, 9) {
		if r>>(i+1)&1 == 0 &&
			c>>(i+1)&1 == 0 &&
			g>>(i+1)&1 == 0 {
			res = append(res, i+1)
		}
	}
	if len(res) == 0 {
		return nil
	}
	return res
}

func (s *sudoku) try_put(row_idx, col_idx, val int) bool {
	if s == nil || val < 1 || val > 9 ||
		row_idx < 0 || row_idx > 8 ||
		col_idx < 0 || col_idx > 8 {
		return false
	}
	if !s.check_row(row_idx, val) ||
		!s.check_col(col_idx, val) ||
		!s.check_grid(row_idx, col_idx, val) {
		return false
	}
	s.data[row_idx][col_idx] = val
	s.set_row(row_idx, val)
	s.set_col(col_idx, val)
	s.set_grid(row_idx, col_idx, val)
	return true
}

func (s *sudoku) unset(row_idx, col_idx int) {
	if s == nil ||
		row_idx < 0 || row_idx > 8 ||
		col_idx < 0 || col_idx > 8 {
		return
	}
	if s.data[row_idx][col_idx] == 0 {
		return
	}
	val := s.data[row_idx][col_idx]
	s.data[row_idx][col_idx] = 0
	s.unset_row(row_idx, val)
	s.unset_col(col_idx, val)
	s.unset_grid(row_idx, col_idx, val)
}

func (s *sudoku) unset_row(row_idx, val int) {
	s.row[row_idx] &^= 1 << val
}

func (s *sudoku) unset_col(col_idx, val int) {
	s.col[col_idx] &^= 1 << val
}

func (s *sudoku) unset_grid(row_idx, col_idx, val int) {
	s.grid[3*(row_idx/3)+col_idx/3] &^= 1 << val
}

func (s *sudoku) set_row(row_idx, val int) {
	s.row[row_idx] |= 1 << val
}

func (s *sudoku) set_col(col_idx, val int) {
	s.col[col_idx] |= 1 << val
}

func (s *sudoku) set_grid(row_idx, col_idx, val int) {
	s.grid[3*(row_idx/3)+col_idx/3] |= 1 << val
}

func (s *sudoku) check_row(row_idx, val int) bool {
	return (s.row[row_idx]>>val)&1 == 0
}

func (s *sudoku) check_col(col_idx, val int) bool {
	return (s.col[col_idx]>>val)&1 == 0
}

func (s *sudoku) check_grid(row_idx, col_idx, val int) bool {
	return (s.grid[3*(row_idx/3)+col_idx/3]>>val)&1 == 0
}

func main() {
	s, err := new_sudoku([][]int{
		{0, 0, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	})
	if err != nil {
		panic(err)
	}
	next := s.get_next_arr()
	if len(next) == 0 {
		fmt.Println("already solved sudoku")
	}
	fmt.Println(next)
	res := make([][9][9]int, 0, 100)
	s.solve(next, 0, &res)
	if len(res) > 0 {
		fmt.Println("solved with:")
		fmt.Println(res)
	} else {
		fmt.Println("can't solve:")
	}
}
