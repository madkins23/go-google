package sheets

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Headers collects result of parsing headers from first row of sheet.
type Headers struct {
	colAlpha  map[string]string
	colNumber map[string]int
	maxAlpha  string
}

// ParseHeaders takes a row of cells containing headers and an optional list of required headers.
// Returns a Headers object containing the cell column alpha and number for each header.
// Also tracks the alpha of the furthest right non-empty column.
func ParseHeaders(headers []interface{}, required ...string) (*Headers, error) {
	errs := make([]string, 0)
	hdrs := &Headers{
		colAlpha:  make(map[string]string),
		colNumber: make(map[string]int),
	}

	for column, value := range headers {
		header, ok := value.(string)
		if !ok {
			errs = append(errs, fmt.Sprintf("Unable to get string from cell value: %v", value))
			continue
		}

		header = strings.ToLower(header)
		hdrs.colAlpha[header] = ColToAlpha(column)
		hdrs.colNumber[header] = column
		hdrs.maxAlpha = hdrs.colAlpha[header]
	}

	for _, reqd := range required {
		_, okAlpha := hdrs.colAlpha[reqd]
		_, okNumber := hdrs.colNumber[reqd]
		if !okAlpha || !okNumber {
			errs = append(errs, fmt.Sprintf("Missing required header '%s'", reqd))
		}
	}

	if len(errs) < 1 {
		return hdrs, nil
	}

	return hdrs, errors.Errorf("Parsing header row of sheet:\n  " + strings.Join(errs, "\n  "))
}

// Alpha returns the alpha column address for the specified headers.
func (hdrs *Headers) Alpha(header string) string {
	return hdrs.colAlpha[header]
}

// MaxAlpha returns the furthest right alpha column address in the header row.
func (hdrs *Headers) MaxAlpha() string {
	return hdrs.maxAlpha
}

// Number returns the zero-based column index for the header.
func (hdrs *Headers) Number(header string) int {
	return hdrs.colNumber[header]
}
