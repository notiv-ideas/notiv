package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var tmpTestDir string = "/tmp/data"
var tmpLocalDir string = filepath.Join(".", ".notiv")

func TestMain(m *testing.M) {

	// Reset test home
	os.RemoveAll(tmpLocalDir)
	os.RemoveAll(tmpTestDir)

	code := m.Run()
	os.Exit(code)
}

func TestMainOutput(t *testing.T) {
	type preliminaryTestCase struct {
		name   string
		cmd    []string
		expect []string
	}
	// Define preliminary scenarios
	preliminaries := []preliminaryTestCase{
		{name: "notiv opening",
			cmd: []string{"notiv"},
			expect: []string{
				`notiv ` + appVersion + ` — Cryptographically versioned key-value CLI (powered by graviton)`,
				"Usage: notiv [global flags] <command> [subcommand] [arguments...]",
				"Run 'notiv help' or 'notiv --help' for detailed usage.",
			}},
		// INVALID CMD
		{name: "does not accept invalid commands",
			cmd:    []string{"notiv", "invalidCommand"},
			expect: []string{"notiv:", "unknown command:"}},

		// VERSION CMD
		{name: "does version check",
			cmd:    []string{"notiv", "version"},
			expect: []string{appVersion}},

		// HELP CMD
		{name: "does diplay help message",
			cmd:    []string{"notiv", "help"},
			expect: []string{"notiv", "Built on graviton", "See 'notiv help' or 'notiv --help' for more details."}},
		{name: "does display a help message help flag",
			cmd:    []string{"notiv", "--help"},
			expect: []string{"notiv", "Built on graviton", "See 'notiv help' or 'notiv --help' for more details."}},

		// INIT CMD
		{name: "init does accept password flag",
			cmd:    []string{"notiv", "init", "-password=password"},
			expect: []string{"config.json created"}},
		{name: "init does accept share flag with a password flag",
			cmd:    []string{"notiv", "--share=" + tmpTestDir, "init", "-password=password"},
			expect: []string{"config.json created"}},
		{name: "re-init notiv configs does not overwrite exisiting",
			cmd:    []string{"notiv", "init"},
			expect: []string{"notiv configs already exist", notivDir}},

		// SHARE CMD
		// 	there is a whole world of in-memory shares that we haven't explored
		// 	and much of that is going to need to be explored at some point.
		{name: "create new share does have a default data directory",
			cmd:    []string{"notiv", "share", "new"},
			expect: []string{"Share created at", shareDir + "/data"}},
		{name: "create new share does accept named arguement",
			cmd:    []string{"notiv", "share", "new", "foo"},
			expect: []string{"Share created at", shareDir + "/foo"}},
		{name: "create new share does accept path arguement ",
			cmd:    []string{"notiv", "share", "new", tmpTestDir},
			expect: []string{"Share created at", shareDir}},
		{name: "create new share with --share flag does not overwrite existing directory",
			cmd:    []string{"notiv", "--share=" + tmpTestDir, "share", "new"},
			expect: []string{"Share already created at", tmpTestDir}},
		{name: "create new share with --share flag does accept name argument ",
			cmd:    []string{"notiv", "--share=" + tmpTestDir, "share", "new", "foo"},
			expect: []string{"Share created at", tmpTestDir}},

		// CONFIG CMD
		{name: "config usage does show usage and examples",
			cmd:    []string{"notiv", "config"},
			expect: []string{"Usage:", "Examples:"}},
		{name: "config list does match in-memory cfg",
			cmd: []string{"notiv", "config", "list"},
			expect: func() []string {
				var exp []string
				for k, v := range cfg {
					exp = append(exp, strings.Join([]string{k, v}, " "))
				}
				return exp
			}()},
		{name: "config edits are persistent to disk",
			cmd:    []string{"notiv", "config", "edit", "defaultEncryptPolicy", "false"},
			expect: []string{"notiv config:", "defaultEncryptPolicy", "false"}},
		{name: "config edits are persistent to disk",
			cmd:    []string{"notiv", "config", "edit", "defaultEncryptPolicy", "true"},
			expect: []string{"notiv config:", "defaultEncryptPolicy", "true"}},
		{name: "config edits are persistent to disk",
			cmd:    []string{"notiv", "config", "edit", "defaultEncryptPolicy", "foo"},
			expect: []string{"notiv config:", "defaultEncryptPolicy", "must be either true or false"}},
		{name: "config edit does not change defaultEncryptedCHECKSUM",
			cmd:    []string{"notiv", "config", "edit", "defaultEncryptedCHECKSUM", "\"something else\""},
			expect: []string{"Cannot modify defaultEncryptedCHECKSUM", "use 'config password new'"}},
	}

	for _, tc := range preliminaries {
		t.Run(tc.name, func(t *testing.T) {
			output := captureMainOutput(t, tc.cmd)
			// Perform assertions on the captured output
			for _, expected := range tc.expect {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output:\n'%s'\ngot: %s", expected, output)
				}
			}
		})
	}

	type primaryTestCase struct {
		name string
		test func(t *testing.T)
	}
	primaries := []primaryTestCase{
		{name: "put does return encrypted value with ok status and current version 1",
			test: func(t *testing.T) {
				cmd := []string{"notiv",
					"put", "alice", "email@example.com", "foo",
					// a note here about the password
					"--password=password", // obviously we would want to use the password entry
					// however, there currently isn't a convenient way to handle this
				}
				expect := []string{"ok", "alice", "foo", "current tree version: 1"}
				output := captureMainOutput(t, cmd)
				t.Run("returns value encrypted", func(t *testing.T) {
					values := strings.Split(output, " ")
					initialValue := cmd[3]
					encryptedValue := values[2]
					if initialValue == encryptedValue {
						t.Errorf("initial value was %s, got %s", initialValue, encryptedValue)
					}

					decryptedValue := decrypt(encryptedValue, readPassword())
					if initialValue != decryptedValue {
						t.Errorf("initial value was %s, got %s", initialValue, decryptedValue)
					}
				})
				t.Run("returns ok and current version 1", func(t *testing.T) {
					rangeExpectations(t, output, expect)
				})
			}},
		{name: "get does return encrypted value status with current version 1",
			test: func(t *testing.T) {
				cmd := []string{"notiv", "-n", "get", "alice", "foo"}
				expect := []string{"email@example.com"}
				output := captureMainOutput(t, cmd)
				var decryptedValue string
				t.Run("returns value encrypted", func(t *testing.T) {
					values := strings.Split(output, " ")
					initialValue := "email@example.com"
					encryptedValue := values[0]
					if initialValue == encryptedValue {
						t.Errorf("initial value was %s, got %s", initialValue, encryptedValue)
					}
					*password_flag = "password"
					decryptedValue = decrypt(encryptedValue, readPassword())
					if initialValue != decryptedValue {
						t.Errorf("initial value was %s, got %s", initialValue, decryptedValue)
					}
				})
				t.Run("returns values and current version 1", func(t *testing.T) {
					rangeExpectations(t, decryptedValue, expect)
				})
			}},
		{name: "put does return encrypted value with ok status and current tree version",
			test: func(t *testing.T) {
				iterations := 10
				// var version string
				var count int
				// var output string
				for range iterations { // should be the same value
					count++
					cmd := []string{"notiv", "put", "bob", "email@example.com", "bar",
						// a note here about the password
						"--password=password", // obviously we would want to use the password entry
						// however, there currently isn't a convenient way to handle this
					}
					expect := []string{"ok", "bob", "bar", "current tree version:", fmt.Sprintf("%d", count)}
					output := captureMainOutput(t, cmd)
					t.Run("returns value encrypted", func(t *testing.T) {
						values := strings.Split(output, " ")
						initialValue := cmd[3]
						encryptedValue := values[2]
						if initialValue == encryptedValue {
							t.Errorf("initial value was %s, got %s", initialValue, encryptedValue)
						}

						decryptedValue := decrypt(encryptedValue, readPassword())
						if initialValue != decryptedValue {
							t.Errorf("initial value was %s, got %s", initialValue, decryptedValue)
						}
					})
					t.Run("returns ok and current version 1", func(t *testing.T) {
						rangeExpectations(t, output, expect)
					})
				}
			}},
	}

	for _, tc := range primaries {
		t.Run(tc.name, tc.test)
	}
}

func rangeExpectations(t *testing.T, output string, expect []string) {
	for _, e := range expect {
		if !strings.Contains(output, e) {
			t.Errorf("captured output:\n%swanted: %s", output, e)
		}
	}
}

func captureMainOutput(t *testing.T, args []string) string {
	// Backup original os.Stdout
	originalStdout := os.Stdout

	// Create a pipe to capture the output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("Failed to create pipe:", err)
	}
	os.Stdout = w // Redirect os.Stdout to the pipe

	// Simulate the command-line arguments
	os.Args = args

	// Use a WaitGroup to wait for all output to finish
	main()

	// Close the write end of the pipe
	if err := w.Close(); err != nil {
		t.Fatal("Failed to close the write end of the pipe:", err)
	}

	// Read the captured output from the pipe
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatal("Failed to read from pipe:", err)
	}

	// Restore original os.Stdout
	os.Stdout = originalStdout

	// Get the captured output
	output := buf.String()

	// Log the captured output for debugging purposes
	// t.Log("Captured Output:", output)

	return output
}
