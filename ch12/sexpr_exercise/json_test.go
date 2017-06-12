//e12.4 e12.6

package sexpr

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"
)

func testJSON(t *testing.T, v interface{}) {
	// Encode it
	data, err := JSONIndent(v)
	if err != nil {
		t.Logf("JSONIndent failed: %v", err)
	} else {
		t.Logf("JSONIdent() =\n%s\n", data)
	}
	// cmp with json and fmt
	str, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("json.MarshalIndent() =\n%s\n", str)
	}
	t.Logf("#v =\n%#v\n", v)
}

func TestJSONCMP(t *testing.T) {
	// Encode it
	v := vJ
	var cmp J
	data, err := JSONIndent(v)
	if err != nil {
		t.Fatalf("JSONIndent failed: %v", err)
	}

	// Decode it
	if err := json.Unmarshal(data, &cmp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Check equality.
	if !reflect.DeepEqual(cmp, v) {
		fmt.Printf("JSONIndent() =\n%s\n", data)
		fmt.Printf("Unmarshal() =\n%#v\n", cmp)
		t.Fatal("not equal")
	}

}

type J struct {
	Float float64
	I     int
	B     [3]bool
	M     [2]Movie
}

var vJ = J{
	math.Pi,
	315,
	[3]bool{true},
	[2]Movie{strangelove, strangelove},
}

func TestJSONPrint(t *testing.T) {
	// testJSON(t, strangelove)
	// testJSON(t, testValue)

	testJSON(t, os.Stdout)
	// testJSON(t, struct{ F []float32 }{[]float32{2.3}})
	// testJSON(t, []float32{2.3})
	// testJSON(t, Movie{Actor: map[string]string{}})
	// testJSON(t, struct{}{})
	// testJSON(t, []Movie{strangelove, strangelove})

}
