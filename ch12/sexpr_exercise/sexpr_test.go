//e12.4 e12.6

package sexpr

import (
	"io"
	"reflect"
	"testing"
)

type Movie struct {
	Title, Subtitle string
	Year            int
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
	// Hot             float64
	IsOk bool
}

var strangelove = Movie{
	Title:    "Dr. Strangelove",
	Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
	Year:     1964,
	Actor: map[string]string{
		"Dr. Strangelove":            "Peter Sellers",
		"Grp. Capt. Lionel Mandrake": "Peter Sellers",
		"Pres. Merkin Muffley":       "Peter Sellers",
		"Gen. Buck Turgidson":        "George C. Scott",
		"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
		`Maj. T.J. "King" Kong`:      "Slim Pickens",
	},
	Oscars: []string{
		"Best Actor (Nomin.)",
		"Best Adapted Screenplay (Nomin.)",
		"Best Director (Nomin.)",
		"Best Picture (Nomin.)",
	},
}

type T struct {
	Writer io.Writer
	F      *Movie
	FLoat  float64
	Pt     **int
	I      int
	B      bool
	Ba     [3]bool
	// M     []Movie
}

var i = 33
var ip = &i

var testValue = T{
	I(44),
	&strangelove,
	3.14,
	&ip,
	33,
	true,
	[3]bool{true},
	// []Movie{strangelove, strangelove},
}

//因为map 输出的无序，无法固定的比较Marshal输出
func testSexpr(t *testing.T, v interface{}) {
	// Encode it
	data, err := Marshal(v)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	t.Logf("Marshal() = %s\n", data)

	// Pretty-print it:
	data, err = MarshalIndent(v)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("MarshalIdent() =\n%s\n", data)

	var value = reflect.New(reflect.TypeOf(v)).Elem() //接口类型field = nil, 丢失信息了；Unmarshal将失败如果v包含接口类型filed
	// t.Log(reflect.TypeOf(value.Interface()))//* <type of v>

	if err := Unmarshal(data, value.Addr().Interface()); err != nil {
		t.Errorf("Unmarshal failed: %v", err)
	} else {
		// t.Logf("Unmarshal() = %#v\n", movie)
		// Check equality.
		if !reflect.DeepEqual(value.Interface(), v) {
			t.Error("not equal")
			str, _ := MarshalIndent(value.Interface())
			t.Logf("Unmarshal() =\n%s\n", str)
		}
	}
}

//实现了再testSexpr 中比较
func TestSexprCMP(t *testing.T) {

	// Encode it
	v := testValue
	data, err := MarshalIndent(v)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	// Decode it
	var cmp = T{Writer: I(0)} // 包含了接口类型field明确给出具体类型信息，
	if err := Unmarshal(data, &cmp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	// t.Logf("Unmarshal() = %+v\n", movie)
	// t.Logf("Unmarshal() = %#v\n", movie)
	// Check equality.
	if !reflect.DeepEqual(cmp, v) {
		t.Fatal("not equal")
	}
}

func TestSexprPrint(t *testing.T) {
	testSexpr(t, strangelove)
	// testSexpr(t, testValue)
	// testSexpr(t, &ip)
	// testSexpr(t, 33)
	testSexpr(t, T{})
	testSexpr(t, struct{}{})
}
