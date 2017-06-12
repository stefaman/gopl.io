package word_test

import(
	// "fmt"
	"testing"
	"math/rand"
	"time"
	"unicode"
	word "gopl.io/ch11/palindrome"
)

func palindrome(rng *rand.Rand) string {
	const n = 100
	str := make([]rune, n)
	for i, j :=0, n-1; i < j; i++ {
		r := rune(rng.Int31n(unicode.MaxRune))//rune() can be omitted
		str[i] = r
		if unicode.IsLetter(r) {
			str[j] = r
			j--
		}
	}
	return string(str)
}

func noPalindrome(rng *rand.Rand) string {
	const n = 100
	str := []rune(palindrome(rng))
	pos := uint32(rng.Int31n(n))
	str[pos] = 'a'
	str[n-1-pos]= 'b'

	return string(str)
}

func TestRandPalindrome(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 1000; i++ {
		str := palindrome(rng)
		if !word.IsPalindrome(str) {
			t.Errorf("for seed %d, IsPalindrome(%q) is false", seed, str)
		}

	}
}

func TestRandNoPalindrome(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 1000; i++ {
		str := noPalindrome(rng)
		if word.IsPalindrome(str) {
			t.Errorf("for seed %d, IsPalindrome(%q) is true", seed, str)
		}
	}
}

func BenchmarkPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		word.IsPalindrome("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	}
}

func ExampleIsPalindrome() {
	fmt.Println(IsPalindrome("A man, a plan, a canal: Panama"))
	fmt.Println(IsPalindrome("palindrome"))
	// Output:
	// true
	// false
}
