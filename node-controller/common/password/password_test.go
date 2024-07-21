package password

import "testing"

func TestPasswd(t *testing.T) {
	s, _ := Encode("abcd1234")
	t.Log(s)
}
