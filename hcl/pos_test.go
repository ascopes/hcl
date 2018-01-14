package hcl

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestRangeOver(t *testing.T) {
	tests := []struct {
		A    Range
		B    Range
		Want Range
	}{
		{
			Range{ //   ##
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
		},
		{
			Range{ // ####
				Start: Pos{Byte: 0, Line: 1, Column: 1},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ // #####
				Start: Pos{Byte: 0, Line: 1, Column: 1},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
		},
		{
			Range{ //   ####
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 6, Line: 1, Column: 7},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  #####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 6, Line: 1, Column: 7},
			},
		},
		{
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //   ##
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
		},
		{
			Range{ //  ###
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
		},
		{
			Range{ //   ###
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
		},
		{
			Range{ //  ####
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ####
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ####
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
		},
		{
			Range{ // ##
				Start: Pos{Byte: 0, Line: 1, Column: 1},
				End:   Pos{Byte: 2, Line: 1, Column: 3},
			},
			Range{ //     ##
				Start: Pos{Byte: 4, Line: 1, Column: 5},
				End:   Pos{Byte: 6, Line: 1, Column: 7},
			},
			Range{ // ######
				Start: Pos{Byte: 0, Line: 1, Column: 1},
				End:   Pos{Byte: 6, Line: 1, Column: 7},
			},
		},
		{
			Range{ //     ##
				Start: Pos{Byte: 4, Line: 1, Column: 5},
				End:   Pos{Byte: 6, Line: 1, Column: 7},
			},
			Range{ // ##
				Start: Pos{Byte: 0, Line: 1, Column: 1},
				End:   Pos{Byte: 2, Line: 1, Column: 3},
			},
			Range{ // ######
				Start: Pos{Byte: 0, Line: 1, Column: 1},
				End:   Pos{Byte: 6, Line: 1, Column: 7},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s<=>%s", test.A, test.B), func(t *testing.T) {
			got := RangeOver(test.A, test.B)
			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf(
					"wrong result\nA   : %-10s %s\nB   : %-10s %s\ngot : %-10s %s\nwant: %-10s %s",
					visRangeOffsets(test.A), test.A,
					visRangeOffsets(test.B), test.B,
					visRangeOffsets(got), got,
					visRangeOffsets(test.Want), test.Want,
				)
			}
		})
	}
}

func TestPosOverlap(t *testing.T) {
	tests := []struct {
		A    Range
		B    Range
		Want Range
	}{
		{
			Range{ //   ##
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //   ##
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
		},
		{
			Range{ // ####
				Start: Pos{Byte: 0, Line: 1, Column: 1},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ###
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
		},
		{
			Range{ //   ####
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 6, Line: 1, Column: 7},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //   ###
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
		},
		{
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //   ##
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
			Range{ //   ##
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
		},
		{
			Range{ //  ###
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ###
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
		},
		{
			Range{ //   ###
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ####
				Start: Pos{Byte: 1, Line: 1, Column: 2},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //   ###
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
		},
		{
			Range{ //  ####
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ####
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
			Range{ //  ####
				Start: Pos{Byte: 2, Line: 1, Column: 3},
				End:   Pos{Byte: 5, Line: 1, Column: 6},
			},
		},
		{
			Range{ // ##
				Start: Pos{Byte: 0, Line: 1, Column: 1},
				End:   Pos{Byte: 2, Line: 1, Column: 3},
			},
			Range{ //     ##
				Start: Pos{Byte: 4, Line: 1, Column: 5},
				End:   Pos{Byte: 6, Line: 1, Column: 7},
			},
			Range{ // (no overlap)
				Start: Pos{Byte: 0, Line: 1, Column: 1},
				End:   Pos{Byte: 0, Line: 1, Column: 1},
			},
		},
		{
			Range{ //     ##
				Start: Pos{Byte: 4, Line: 1, Column: 5},
				End:   Pos{Byte: 6, Line: 1, Column: 7},
			},
			Range{ // ##
				Start: Pos{Byte: 0, Line: 1, Column: 1},
				End:   Pos{Byte: 2, Line: 1, Column: 3},
			},
			Range{ // (no overlap)
				Start: Pos{Byte: 4, Line: 1, Column: 5},
				End:   Pos{Byte: 4, Line: 1, Column: 5},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s<=>%s", test.A, test.B), func(t *testing.T) {
			got := test.A.Overlap(test.B)
			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf(
					"wrong result\nA   : %-10s %s\nB   : %-10s %s\ngot : %-10s %s\nwant: %-10s %s",
					visRangeOffsets(test.A), test.A,
					visRangeOffsets(test.B), test.B,
					visRangeOffsets(got), got,
					visRangeOffsets(test.Want), test.Want,
				)
			}
		})
	}
}

// visRangeOffsets is a helper that produces a visual representation of the
// start and end byte offsets of the given range, which can then be stacked
// with the same for other ranges to more easily see how the ranges relate
// to one another.
func visRangeOffsets(rng Range) string {
	var buf bytes.Buffer
	if rng.End.Byte < rng.Start.Byte {
		// Should never happen, but we'll visualize it anyway so we can
		// more easily debug failing tests.
		for i := 0; i < rng.End.Byte; i++ {
			buf.WriteByte(' ')
		}
		for i := rng.End.Byte; i < rng.Start.Byte; i++ {
			buf.WriteByte('!')
		}
		return buf.String()
	}

	for i := 0; i < rng.Start.Byte; i++ {
		buf.WriteByte(' ')
	}
	for i := rng.Start.Byte; i < rng.End.Byte; i++ {
		buf.WriteByte('#')
	}
	return buf.String()
}
