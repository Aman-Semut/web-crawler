package dbController

import (
	"testing"
)

func TestdbController(t *Testing.T) {

	testCases := []struct {
		name    string
		input   map[string][]string
		baseUrl string
		output  map[string][]string
	}{
		{
			name: "TC#1",
			input: map[string][]string{
				"A": {
					"1",
					"2",
					"3",
				},
				"B": {},
			},
			baseUrl: "unit_test/TC#1",
			output: map[string][]string{
				"A": {
					"1",
					"2",
					"3",
				},
				"B": {},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			done := AddData(tc.input, tc.baseUrl)

			if done != true {
				t.Errorf("Testcase %s FAILED", tc.name)
				return
			} else {
				res := GetData(tc.baseUrl)
				if res == nil {
					t.Errorf("Testcase %s FAILED", tc.name)
					return
				} else {
					for k, v := range tc.output {
						for _,s := range res[k]{
							if s!= v {
                                t.Errorf("Testcase %s FAILED", tc.name)
                                return
                            }
                        }
						}
					}
					t.Logf("Testcase %s PASSED", tc.name)
				}

			}

		})
    }

}
