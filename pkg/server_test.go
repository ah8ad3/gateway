package pkg

import "testing"

func TestRUN(t *testing.T) {
	RUN("127.0.0.1", "3000", "v2")
	RUN("127.0.0.1", "3000", "v1")
	//RUN("", "", "V1")
}
