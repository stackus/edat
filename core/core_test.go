package core_test

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update golden file")

func getGoldenFileData(t *testing.T, fileName string) []byte {
	golden := filepath.Join("testdata", fileName+".golden")
	// if *update {
	// 	if err := ioutil.WriteFile(golden, actual, 0644); err != nil {
	// 		t.Fatalf("Error writing golden file for filename=%s: %s", fileName, err)
	// 	}
	// }
	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	return expected
}
