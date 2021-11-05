package stringx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasItem(t *testing.T) {
	tests := []struct {
		caseDesc  string
		giveArray []string
		giveItem  string
		wantRet   bool
	}{
		{
			caseDesc:  "existed",
			giveArray: []string{"fine", "bad"},
			giveItem:  "bad",
			wantRet:   true,
		},
		{
			caseDesc:  "not existed",
			giveArray: []string{"fine", "bad"},
			giveItem:  "not",
			wantRet:   false,
		},
		{
			caseDesc:  "nil",
			giveArray: nil,
			giveItem:  "not",
			wantRet:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := HasItem(tc.giveArray, tc.giveItem)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestIsSuperset(t *testing.T) {
	tests := []struct {
		caseDesc   string
		giveSrc    []string
		giveSubset []string
		wantRet    bool
	}{
		{
			caseDesc:   "is",
			giveSrc:    []string{"a", "b", "c"},
			giveSubset: []string{"a", "c"},
			wantRet:    true,
		},
		{
			caseDesc:   "not is",
			giveSrc:    []string{"a", "b", "c"},
			giveSubset: []string{"d", "c"},
			wantRet:    false,
		},
		{
			caseDesc:   "empty",
			giveSrc:    nil,
			giveSubset: nil,
			wantRet:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := IsSuperset(tc.giveSrc, tc.giveSubset)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveA    []string
		giveB    []string
		wantRet  []string
	}{
		{
			caseDesc: "sanity",
			giveA:    []string{"a", "b", "c"},
			giveB:    []string{"a", "c", "e"},
			wantRet:  []string{"a", "c"},
		},
		{
			caseDesc: "a-empty",
			giveA:    nil,
			giveB:    []string{"a", "c"},
			wantRet:  nil,
		},
		{
			caseDesc: "b-empty",
			giveA:    []string{"a", "b", "c"},
			giveB:    nil,
			wantRet:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := Intersect(tc.giveA, tc.giveB)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}
