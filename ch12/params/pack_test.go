package params

import "testing"

func TestPack(t *testing.T) {
	var data = struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}{
		[]string{"cai", "guo", "zhun"},
		19870515,
		true,
	}

	printData(t, &data)
	printData(t, 1)
}

func printData(t *testing.T, data interface{}) {
	str, err := Pack(data)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(str)

	}
}
