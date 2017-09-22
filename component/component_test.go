package component

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestComponentWriteMetaJSON(t *testing.T) {
	dir, err := ioutil.TempDir("", "meta_test")
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(dir)

	c, err := New(
		OptionName("test_deployment/test_component"),
		OptionWorkingDirectory(dir))
	if err != nil {
		t.Error("expected no errors, got", err.Error())
	}

	if err != c.WriteMetaJSON() {
		t.Error("expected no errors writing file, got", err.Error())
	}
}
