//go:build integration

package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/thesouldev/goboxd/internal/runs"
	"github.com/thesouldev/goboxd/server"
)

func TestRunPython3(t *testing.T) {
	body := `{
		"language": "py3",
		"source": "name = input()\nprint(\"Hello, \" + name + \"!\")\n",
		"tests": [{"stdin": "Ada\n", "expected_stdout": "Hello, Ada!\n"}]
	}`
	assertAcceptedRun(t, body)
}

func TestRunCPP(t *testing.T) {
	body := `{
		"language": "cpp",
		"source": "#include <iostream>\n#include <string>\n\nint main() {\n    std::string name;\n    std::getline(std::cin, name);\n    std::cout << \"Hello, \" << name << \"!\\n\";\n    return 0;\n}\n",
		"build": {"flags": ["-std=c++17", "-Wall"]},
		"tests": [{"stdin": "Ada\n", "expected_stdout": "Hello, Ada!\n"}]
	}`
	assertAcceptedRun(t, body)
}

func TestRunJava(t *testing.T) {
	body := `{
		"language": "java",
		"source_filename": "Main.java",
		"artifact_filename": "Main",
		"source": "import java.io.BufferedReader;\nimport java.io.InputStreamReader;\n\npublic class Main {\n    public static void main(String[] args) throws Exception {\n        BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));\n        String name = reader.readLine();\n        System.out.println(\"Hello, \" + name + \"!\");\n    }\n}\n",
		"tests": [{"stdin": "Ada\n", "expected_stdout": "Hello, Ada!\n"}]
	}`
	assertAcceptedRun(t, body)
}

func assertAcceptedRun(t *testing.T, body string) {
	t.Helper()

	testServer := httptest.NewServer(server.NewMux())
	defer testServer.Close()

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(testServer.URL+"/run", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("POST /run: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var response runs.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != runs.StatusAccepted {
		t.Fatalf("status = %q, want %q: %#v", response.Status, runs.StatusAccepted, response)
	}
	if len(response.Tests) != 1 {
		t.Fatalf("test count = %d, want 1", len(response.Tests))
	}
	if response.Tests[0].Status != runs.StatusAccepted {
		t.Fatalf("test status = %q, want %q: %#v", response.Tests[0].Status, runs.StatusAccepted, response.Tests[0])
	}
}
