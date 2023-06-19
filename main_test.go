package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCLI(t *testing.T) {
	tmp := tmpHome(t)

	t.Run("stdin", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Logf("could not open pipe: %v", err)
			t.FailNow()
		}
		os.Stdin = r

		w.WriteString("foo bar baz")
		w.Close()
		runCLI([]string{"foo/bar.baz"})

		cfg := readConfig()
		snippets := readSnippets(cfg)

		if len(snippets) != 1 {
			t.Logf("snippet count is incorrect: got %d but want 1", len(snippets))
			t.FailNow()
		}

		fn := filepath.Join(tmp, "foo/bar.baz")
		fi, err := os.Open(fn)
		if err != nil {
			t.Logf("could not open test file: %v", err)
			t.FailNow()
		}
		defer fi.Close()

		content, err := io.ReadAll(fi)
		if err != nil {
			t.Logf("could not read test file: %v", err)
			t.FailNow()
		}

		if string(content) != "foo bar baz" {
			t.Logf(`snippet is incorrect: got %q but want "foo bar baz"`, string(content))
			t.FailNow()
		}
	})

	t.Run("stdout", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Logf("could not open pipe: %v", err)
			t.FailNow()
		}
		os.Stdout = w
		runCLI([]string{"foo/bar.baz"})
		w.Close()
		out, err := io.ReadAll(r)
		if err != nil {
			t.Log("could not read stdout")
			t.FailNow()
		}

		if string(out) != "foo bar baz" {
			t.Logf(`snippet is incorrect: got %q but want "foo bar baz"`, string(out))
			t.FailNow()
		}
	})

	t.Run("list", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Logf("could not open pipe: %v", err)
			t.FailNow()
		}
		os.Stdout = w
		runCLI([]string{"list"})
		w.Close()
		out, err := io.ReadAll(r)
		if err != nil {
			t.Log("could not read stdout")
			t.FailNow()
		}

		if string(out) != "foo/bar.baz\n" {
			t.Logf(`snippet is incorrect: got %q but want "foo/bar.baz\n"`, string(out))
			t.FailNow()
		}
	})
}

func TestScan(t *testing.T) {
	tmp := tmpHome(t)

	cfg := readConfig()
	snippets := readSnippets(cfg)
	snippets = scanSnippets(cfg, snippets)
	initNum := len(snippets)

	tmpSnippetFolder := filepath.Join(tmp, "foo")
	tmpSnippet := filepath.Join(tmpSnippetFolder, "bar.baz")
	if err := os.MkdirAll(tmpSnippetFolder, os.ModePerm); err != nil {
		t.Logf("could not create snippet folder: %v", err)
		t.FailNow()
	}
	if err := os.WriteFile(tmpSnippet, []byte("foo bar baz"), os.ModePerm); err != nil {
		t.Logf("could not create snippet: %v", err)
		t.FailNow()
	}

	snippets = scanSnippets(cfg, snippets)
	if len(snippets) != initNum+1 {
		t.Logf("incorrect number of snippets after initial scanning: want %d but got %d", initNum+1, len(snippets))
		t.FailNow()
	}

	if err := os.Remove(tmpSnippet); err != nil {
		t.Logf("could not remove snippet: %v", err)
		t.FailNow()
	}

	snippets = scanSnippets(cfg, snippets)
	if len(snippets) != initNum {
		t.Logf("incorrect number of snippets after follow-up scanning: want %d but got %d", initNum, len(snippets))
		t.FailNow()
	}
}

func tmpHome(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()
	if err := os.Setenv("NAP_HOME", tmp); err != nil {
		t.Log("could not set NAP_HOME")
		t.FailNow()
	}
	return tmp
}
