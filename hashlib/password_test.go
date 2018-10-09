package hashlib

import (
	"testing"
)

func TestHash512AndEncodeBase64(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"angryMonkey", "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="},
		{"", "z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg=="},
		{"12345678", "+lhdichR3TOKcNz1Naoqkv7ng23Wr/EiZYPojgmWKT8WvACcZSgm4PxccGaVoDzdzjcvE57/TROVnabx9dPqvg=="},
	}

	for _, test := range tests {
		if got := Hash512AndEncodeBase64(test.input); got != test.want {
			t.Errorf("Hash512AndEncodeBase64(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
