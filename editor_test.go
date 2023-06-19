package main

import (
	"fmt"
	"os"
	"testing"
)

func TestGetEditor(t *testing.T) {
	tt := []struct {
		Name      string
		EditorEnv string
		Cmd       string
		Args      []string
	}{
		{
			Name: "default",
			Cmd:  "nano",
		},
		{
			Name:      "vim",
			EditorEnv: "vim",
			Cmd:       "vim",
		},
		{
			Name:      "vim with flag",
			EditorEnv: "vim --foo",
			Cmd:       "vim",
			Args:      []string{"--foo"},
		},
		{
			Name:      "code",
			EditorEnv: "code -w",
			Cmd:       "code",
			Args:      []string{"-w"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			switch tc.EditorEnv {
			case "":
				err = os.Unsetenv("EDITOR")
			default:
				err = os.Setenv("EDITOR", tc.EditorEnv)
			}
			if err != nil {
				t.Logf("could not (un)set env: %v", err)
				t.FailNow()
			}

			cmd, args := getEditor()

			if cmd != tc.Cmd {
				t.Logf("cmd is incorrect: want %q but got %q", tc.Cmd, cmd)
				t.FailNow()
			}

			if argStr, tcArgStr := fmt.Sprint(args), fmt.Sprint(tc.Args); argStr != tcArgStr {
				t.Logf("args are incorrect: want %q but got %q", tcArgStr, argStr)
				t.FailNow()
			}
		})
	}
}
