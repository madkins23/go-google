package sheets

import (
	"strconv"
)

// ColToAlpha converts a column number into a column alpha code.
// Column numbers are zero based.
func ColToAlpha(column int) string {
	var result string

	column++

	for column > 0 {
		digit := (column - 1) % 26
		column = column / 26
		result = string('A'+int32(digit)) + result
	}

	return result
}

// RowCol converts a row and column number into a cell address.
// Row and column numbers are zero based.
// The cell address is for the "current" sheet.
func RowCol(row, column int) string {
	return ColToAlpha(column) + strconv.Itoa(row)
}

// SheetRowCol converts a sheet name and row and column numbers into a cell address.
// Row and column numbers are zero based.
func SheetRowCol(sheet string, row, column int) string {
	return sheet + "!" + ColToAlpha(column) + strconv.Itoa(row)
}
