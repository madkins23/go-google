package sheets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColToA1(t *testing.T) {
	assert.Equal(t, "A", ColToAlpha(0))
	assert.Equal(t, "C", ColToAlpha(2))
	assert.Equal(t, "BC", ColToAlpha(54))
}

func TestRowCol(t *testing.T) {
	assert.Equal(t, "A1", RowCol(1, 0))
	assert.Equal(t, "C9", RowCol(9, 2))
	assert.Equal(t, "BC13", RowCol(13, 54))
}

func TestSheetRowCol(t *testing.T) {
	assert.Equal(t, "title!A1", SheetRowCol("title", 1, 0))
	assert.Equal(t, "title!C9", SheetRowCol("title", 9, 2))
	assert.Equal(t, "title!BC13", SheetRowCol("title", 13, 54))
}
