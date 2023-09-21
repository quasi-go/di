package di

import "testing"

func TestTypeId(t *testing.T) {
	if ObjectTypeId(&Thing1{}) != TypeId[Thing1]() {
		t.Error("Type IDs should match")
	}
}
