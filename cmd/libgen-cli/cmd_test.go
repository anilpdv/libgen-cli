package libgen_cli

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func TestCommandsExist(t *testing.T) {
	cmds := []*cobra.Command{
		rootCmd,
		searchCmd,
		downloadCmd,
		downloadAllCmd,
		linkCmd,
		statusCmd,
		dbdumpsCmd,
		completionCmd,
	}

	for _, cmd := range cmds {
		if cmd == nil {
			t.Fatal("nil command found")
		}
		if cmd.Use == "" {
			t.Errorf("command has empty Use: %v", cmd)
		}
		// Verify help renders without error
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		if err := cmd.Help(); err != nil {
			t.Errorf("error running Help on %s: %v", cmd.Name(), err)
		}
	}
}

func TestMirrorFlagsExist(t *testing.T) {
	cmdsWithMirror := []*cobra.Command{
		searchCmd,
		downloadCmd,
		downloadAllCmd,
		linkCmd,
	}

	for _, cmd := range cmdsWithMirror {
		f := cmd.Flags().Lookup("mirror")
		if f == nil {
			t.Errorf("command %s missing --mirror flag", cmd.Name())
		}
		if f != nil && f.Shorthand != "m" {
			t.Errorf("command %s --mirror shorthand should be 'm', got %q", cmd.Name(), f.Shorthand)
		}
	}
}

func TestPromptTemplateExecution(t *testing.T) {
	// Verify that promptui template works on []string items without error
	promptTemplate := &promptui.SelectTemplates{
		Active:   `▸ {{ . | cyan | bold }}`,
		Selected: `{{ "✔" | green }} {{ . | cyan }}`,
	}

	funcs := template.FuncMap{
		"cyan":  func(s string) string { return s },
		"bold":  func(s string) string { return s },
		"green": func(s string) string { return s },
	}

	tmplActive, err := template.New("active").Funcs(funcs).Parse(promptTemplate.Active)
	if err != nil {
		t.Fatalf("failed to parse active template: %v", err)
	}

	buf := new(bytes.Buffer)
	sampleItem := "     643 The Turing Test by Larry Crockett | gz   | 517 KB"
	if err := tmplActive.Execute(buf, sampleItem); err != nil {
		t.Fatalf("failed to execute active template on string item: %v", err)
	}

	if !strings.Contains(buf.String(), "The Turing Test") {
		t.Errorf("rendered template missing expected content: %s", buf.String())
	}
}

func TestCommandFlagParsing(t *testing.T) {
	// Test searchCmd flags
	if err := searchCmd.Flags().Set("mirror", "https://libgen.is"); err != nil {
		t.Errorf("failed to set searchCmd mirror: %v", err)
	}
	val, _ := searchCmd.Flags().GetString("mirror")
	if val != "https://libgen.is" {
		t.Errorf("expected https://libgen.is, got %s", val)
	}
	_ = searchCmd.Flags().Set("mirror", "")

	// Test downloadCmd flags
	if err := downloadCmd.Flags().Set("mirror", "libgen.li"); err != nil {
		t.Errorf("failed to set downloadCmd mirror: %v", err)
	}
	val, _ = downloadCmd.Flags().GetString("mirror")
	if val != "libgen.li" {
		t.Errorf("expected libgen.li, got %s", val)
	}
	_ = downloadCmd.Flags().Set("mirror", "")

	// Test new pagination and type flags on searchCmd and downloadAllCmd
	for _, cmd := range []*cobra.Command{searchCmd, downloadAllCmd} {
		if f := cmd.Flags().Lookup("page"); f == nil || f.Shorthand != "P" {
			t.Errorf("command %s missing --page / -P flag", cmd.Name())
		}
		if f := cmd.Flags().Lookup("type"); f == nil || f.Shorthand != "t" {
			t.Errorf("command %s missing --type / -t flag", cmd.Name())
		}
		if f := cmd.Flags().Lookup("book-type"); f == nil {
			t.Errorf("command %s missing --book-type flag", cmd.Name())
		}
		if f := cmd.Flags().Lookup("max-pages"); f == nil {
			t.Errorf("command %s missing --max-pages flag", cmd.Name())
		}
	}
}
