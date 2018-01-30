package litecoinpool

import (
	"io/ioutil"
	"log"
	"testing"
)

func loadApi() *[]byte {
	dat, err := ioutil.ReadFile("api_test.json")
	if err != nil {
		log.Println(err)
	}
	return &dat
}

/*
	Just enough to ensure parsing works
*/
func TestParseJson(t *testing.T) {
	data := loadApi()

	poolData, err := parseJson(data)
	if err != nil {
		t.Fatal(err)
	}

	length := len(poolData.Workers)
	if length != 2 {
		t.Fatal("Length != 2")
	}
}
