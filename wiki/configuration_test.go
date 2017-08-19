package wiki

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const filecontents = `
{
	"foo": "bar"
}
`

func Test_readConfigurationImpl(t *testing.T) {

	dir, err := ioutil.TempDir("", "configuration")
	if err != nil {
		t.Fatalf("Couldn't make tempdir %v", err)
	}
	defer os.RemoveAll(dir)

	testfile := filepath.Join(dir, ".wikinewrc")
	fd, err := os.Create(testfile)
	if err != nil {
		t.Fatalf("Couldn't open the tmpfile %v", err)
	}
	if _, err := io.WriteString(fd, filecontents); err != nil {
		fd.Close()
		t.Fatalf("Couldn't write the tmpfile %v", err)
	}
	fd.Close()

	config, err := readConfigurationImpl(dir)
	if err != nil {
		t.Fatalf("couldn't read and parase config %v", err)
	}

	if b, ok := config["foo"]; !ok || b != "bar" {
		t.Errorf("was expecting a serialized json, didn't get one")
	}
}
