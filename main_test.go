package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

// func convertJsonToMap(jsonStr string) (map[string][]string, error) {
// 	var response map[string][]string
// 	err := json.Unmarshal([]byte(jsonStr), &response)
// 	return response, err
// }

func TestPost(t *testing.T) {

	var jsonStr string = `{
    "https://bit.ly/2LRhr5F": [],
    "https://bit.ly/2QjDoyM": [],
    "https://bit.ly/2SMV6Mj": [],
    "https://bit.ly/2TYAZ9G": [],
    "https://bit.ly/2UpqWe1": [],
    "https://bit.ly/2YkmNhG": [],
    "https://bit.ly/31gPI2n": [],
    "https://bit.ly/3i9nFKL": [],
    "https://leetcode.com": [
        "https://leetcode.com/support",
        "https://leetcode.com/jobs",
        "https://leetcode.com/bugbounty",
        "https://leetcode.com/student",
        "https://leetcode.com/terms",
        "https://leetcode.com/privacy",
        "https://leetcode.com/region"
    ],
    "https://leetcode.com/bugbounty": [],
    "https://leetcode.com/jobs": [
        "https://www.cognitoforms.com/LeetCodeLLC/frontendarchitect",
        "https://leetcode.com/track/?data=eyJ1cmwiOiAiaHR0cHM6Ly9sZWV0Y29kZS5jb20vZGlzY3Vzcy9nZW5lcmFsLWRpc2N1c3Npb24vMTUxMTI3NC9UaGluay15b3UtZ290LXdoYXQtaXQtdGFrZXMtdG8td29yay1hdC1MZWV0Q29kZT9yZWY9am9iIiwgInV1aWQiOiAiaHR0cHM6Ly9sZWV0Y29kZS5jb20vZGlzY3Vzcy9nZW5lcmFsLWRpc2N1c3Npb24vMTUxMTI3NC9UaGluay15b3UtZ290LXdoYXQtaXQtdGFrZXMtdG8td29yay1hdC1MZWV0Q29kZT9yZWY9am9iIn0%3D",
        "https://www.cognitoforms.com/LeetCodeLLC/backendarchitect",
        "https://www.cognitoforms.com/LeetCodeLLC/seniorbackendsoftwareengineer",
        "https://www.cognitoforms.com/LeetCodeLLC/uiuxdesigner2",
        "https://www.cognitoforms.com/LeetCodeLLC/visualdesigner",
        "https://bit.ly/2QjDoyM",
        "https://bit.ly/2YkmNhG",
        "https://www.cognitoforms.com/LeetCodeLLC/FullStackEngineerIntern",
        "https://bit.ly/2SMV6Mj",
        "https://bit.ly/2UpqWe1",
        "https://www.cognitoforms.com/LeetCodeLLC/problemadder",
        "https://www.cognitoforms.com/LeetCodeLLC/internalcontesttester",
        "https://www.cognitoforms.com/LeetCodeLLC/solutionauthor",
        "https://www.cognitoforms.com/LeetCodeLLC/videocreator",
        "https://bit.ly/3i9nFKL",
        "https://bit.ly/2LRhr5F",
        "https://www.cognitoforms.com/LeetCodeLLC/ContentCreatorTechnicalWriterApplicationForm",
        "https://www.canva.com",
        "https://bit.ly/2TYAZ9G",
        "https://bit.ly/31gPI2n"
    ],
    "https://leetcode.com/privacy": [],
    "https://leetcode.com/region": [],
    "https://leetcode.com/student": [],
    "https://leetcode.com/support": [],
    "https://leetcode.com/terms": [],
    "https://leetcode.com/track/?data=eyJ1cmwiOiAiaHR0cHM6Ly9sZWV0Y29kZS5jb20vZGlzY3Vzcy9nZW5lcmFsLWRpc2N1c3Npb24vMTUxMTI3NC9UaGluay15b3UtZ290LXdoYXQtaXQtdGFrZXMtdG8td29yay1hdC1MZWV0Q29kZT9yZWY9am9iIiwgInV1aWQiOiAiaHR0cHM6Ly9sZWV0Y29kZS5jb20vZGlzY3Vzcy9nZW5lcmFsLWRpc2N1c3Npb24vMTUxMTI3NC9UaGluay15b3UtZ290LXdoYXQtaXQtdGFrZXMtdG8td29yay1hdC1MZWV0Q29kZT9yZWY9am9iIn0%3D": [],
    "https://www.canva.com": [],
    "https://www.cognitoforms.com/LeetCodeLLC/ContentCreatorTechnicalWriterApplicationForm": [],
    "https://www.cognitoforms.com/LeetCodeLLC/FullStackEngineerIntern": [],
    "https://www.cognitoforms.com/LeetCodeLLC/backendarchitect": [],
    "https://www.cognitoforms.com/LeetCodeLLC/frontendarchitect": [],
    "https://www.cognitoforms.com/LeetCodeLLC/internalcontesttester": [],
    "https://www.cognitoforms.com/LeetCodeLLC/problemadder": [],
    "https://www.cognitoforms.com/LeetCodeLLC/seniorbackendsoftwareengineer": [],
    "https://www.cognitoforms.com/LeetCodeLLC/solutionauthor": [],
    "https://www.cognitoforms.com/LeetCodeLLC/uiuxdesigner2": [],
    "https://www.cognitoforms.com/LeetCodeLLC/videocreator": [],
    "https://www.cognitoforms.com/LeetCodeLLC/visualdesigner": []
}`

	//response, _ := convertJsonToMap(jsonStr)

	f := ReqData{
		Id:  0,
		Url: "https://leetcode.com",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(f)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:9000/post/", &buf)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	result := string(bytes)
	exp := result == jsonStr
	t.Logf("Result %t", exp)

	t.Logf("POST REQUEST SUCCESSFUL")

}

func TestMockPost(t *testing.T) {

	t.Logf("Mocking POST request")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:9000/post/",
		func(req *http.Request) (*http.Response, error) {
			reqBody := make(map[string]string)
			if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}

			resp, err := httpmock.NewJsonResponse(200, reqBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	info := httpmock.GetCallCountInfo()
	t.Logf("Calls made : %d", info["POST http://localhost:9000/post/"]) //records number of calls to http://localhost:9000

}
