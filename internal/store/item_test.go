package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKey(t *testing.T) {
	item1 := Item{Date: "2021-10-30", AdvertId: "123456"}
	key1 := item1.getKey()

	item2 := Item{Date: "2021-10-30", AdvertId: "123456"}
	key2 := item2.getKey()

	item3 := Item{Date: "2021-10-30", AdvertId: "123456"}
	key3 := item3.getKey()

	assert.True(t, key1 == key2 && key2 == key3)

	item4 := Item{Date: "2021-10-30", AdvertId: "654321"}
	key4 := item4.getKey()

	item5 := Item{Date: "2021-10-29", AdvertId: "123456"}
	key5 := item5.getKey()

	assert.True(t, key3 != key4)
	assert.True(t, key3 != key5)
}
