package charcount_test

import(
	"fmt"
	"strings"
	"testing"
	"gopl.io/ch11/charcount"
)

type result struct {
	counts map[rune]int
	invalid int
}

func compare(get, want result) error{
	if get.invalid != want.invalid {
		return fmt.Errorf(" get %d invalid runes, want %d", get.invalid, want.invalid)
	}
	if len(get.counts) != len(want.counts) {
		return fmt.Errorf("counts length difference")
	}
	for k, v := range want.counts {
		if get.counts[k] != v {
			return fmt.Errorf("%q rune get %d, want %d", k, get.counts[k], want.counts[k])
		}
	}
	return nil
}

var tests = []struct {
	text string
	res result
}{
	{ "abc\n蔡蔡蔡cba",
		result{map[rune]int{'a': 2, 'b': 2, 'c': 2, '\n': 1, '蔡': 3}, 0},
	},
	{ "\x9f\n蔡吃🍌",
		result{map[rune]int{'\n': 1, '蔡': 1, '吃': 1, '🍌': 1}, 1},
	},
	{ "\xf0\x9f\x8d\xf0\x9f\x8d\x8c",
		result{map[rune]int{'🍌': 1}, 3},
	},
}
func TestCount(t *testing.T) {
	for _, test := range tests {
		counts, invalid, err := charcount.Count(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("test %q: %v", test.text, err)
		}
		if err := compare(result{counts, invalid}, test.res); err != nil {
			t.Errorf("test %q: %v", test.text, err)
		}
	}
}
