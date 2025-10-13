package main


import (
	"fmt"
	"io/ioutil"
)

func WriteStringsToFiles(values, names []string) error {
	if len(values) != len(names) {
		return fmt.Errorf("arrays must have the same length")
	}
	for i := 0; i < len(values); i++ {
		if err := ioutil.WriteFile(names[i], []byte(values[i]), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", names[i], err)
		}
	}
	return nil
}

func main() {
	values := []string{Test_header_1 ,
Test_header_1_dash ,
Test_header_2 ,
Test_header_3 ,
Test_header_4 ,
Test_header_5 ,
Test_header_6 ,
Test_header_6_dash ,
Test_header_7 ,
Test_header_8 ,
Test_header_9 ,
Test_header_9_dash ,
Test_header_10_dash ,
Test_header_10 ,
}

	names := []string{"Test_header_1",
"Test_header_1_dash",
"Test_header_2",
"Test_header_3",
"Test_header_4",
"Test_header_5",
"Test_header_6",
"Test_header_6_dash",
"Test_header_7",
"Test_header_8",
"Test_header_9",
"Test_header_9_dash",
"Test_header_10_dash",
"Test_header_10",
}
	if err := WriteStringsToFiles(values, names); err != nil {
		fmt.Println("Error:", err)
	}
}
// go build ./testhelpers