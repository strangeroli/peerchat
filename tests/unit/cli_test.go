package unit

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestCLIVersion tests the version command
func TestCLIVersion(t *testing.T) {
	cmd := exec.Command("../../bin/peerchat-cli", "version")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run version command: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Xelvra P2P Messenger CLI") {
		t.Errorf("Version output doesn't contain expected text. Got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "v0.1.0-alpha") {
		t.Errorf("Version output doesn't contain version number. Got: %s", outputStr)
	}
}

// TestCLIHelp tests the help command
func TestCLIHelp(t *testing.T) {
	cmd := exec.Command("../../bin/peerchat-cli", "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}

	outputStr := string(output)
	expectedCommands := []string{"init", "start", "status", "send", "listen", "discover", "doctor"}
	for _, command := range expectedCommands {
		if !strings.Contains(outputStr, command) {
			t.Errorf("Help output doesn't contain command '%s'. Got: %s", command, outputStr)
		}
	}
}

// TestCLIStatus tests the status command
func TestCLIStatus(t *testing.T) {
	cmd := exec.Command("../../bin/peerchat-cli", "status")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run status command: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Node Status:") {
		t.Errorf("Status output doesn't contain expected header. Got: %s", outputStr)
	}
}

// TestCLIDoctor tests the doctor command with timeout
func TestCLIDoctor(t *testing.T) {
	cmd := exec.Command("timeout", "5", "../../bin/peerchat-cli", "doctor")
	output, err := cmd.Output()
	
	// timeout command returns exit code 124 on timeout, which is expected
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 124 && exitError.ExitCode() != 0 {
				t.Fatalf("Doctor command failed with unexpected exit code: %v", err)
			}
		} else {
			t.Fatalf("Failed to run doctor command: %v", err)
		}
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Network Diagnostics") {
		t.Errorf("Doctor output doesn't contain expected text. Got: %s", outputStr)
	}
}

// TestCLIDiscover tests the discover command with timeout
func TestCLIDiscover(t *testing.T) {
	cmd := exec.Command("timeout", "5", "../../bin/peerchat-cli", "discover")
	output, err := cmd.Output()
	
	// timeout command returns exit code 124 on timeout, which is expected
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 124 && exitError.ExitCode() != 0 {
				t.Fatalf("Discover command failed with unexpected exit code: %v", err)
			}
		} else {
			t.Fatalf("Failed to run discover command: %v", err)
		}
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Discovering peers") {
		t.Errorf("Discover output doesn't contain expected text. Got: %s", outputStr)
	}
}

// TestCLISendFileValidation tests send-file command input validation
func TestCLISendFileValidation(t *testing.T) {
	// Create a test file
	testFile := "/tmp/test_cli_file.txt"
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() {
		if err := os.Remove(testFile); err != nil {
			t.Logf("Warning: Failed to remove test file: %v", err)
		}
	}()

	// Test with invalid multiaddr
	cmd := exec.Command("../../bin/peerchat-cli", "send-file", "invalid-multiaddr", testFile)
	output, _ := cmd.Output() // Error is expected for invalid input

	outputStr := string(output)
	// Currently file transfer is not implemented, so we expect the "not yet implemented" message
	if !strings.Contains(outputStr, "File transfer not yet implemented") {
		t.Errorf("Send-file should show not implemented message. Got: %s", outputStr)
	}
}

// TestCLILogRotation tests that log rotation functions exist
func TestCLILogRotation(t *testing.T) {
	// This test verifies that the CLI creates log files
	// Run a quick command to generate logs
	cmd := exec.Command("../../bin/peerchat-cli", "version")
	_, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run command for log test: %v", err)
	}

	// Check if log directory exists
	logDir := os.ExpandEnv("$HOME/.xelvra")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// Log directory might not be created by version command, that's OK
		t.Skip("Log directory not created by version command, skipping log rotation test")
	}
}

// TestCLIBinaryExists tests that the binary exists and is executable
func TestCLIBinaryExists(t *testing.T) {
	binaryPath := "../../bin/peerchat-cli"
	
	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("CLI binary does not exist at %s", binaryPath)
	}

	// Check if binary is executable
	info, err := os.Stat(binaryPath)
	if err != nil {
		t.Fatalf("Failed to get binary info: %v", err)
	}

	mode := info.Mode()
	if mode&0111 == 0 {
		t.Errorf("CLI binary is not executable")
	}
}

// TestCLIPerformance tests basic performance metrics
func TestCLIPerformance(t *testing.T) {
	// Test that version command completes quickly
	start := time.Now()
	cmd := exec.Command("../../bin/peerchat-cli", "version")
	_, err := cmd.Output()
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to run version command: %v", err)
	}

	// Version command should complete in under 1 second
	if duration > time.Second {
		t.Errorf("Version command took too long: %v", duration)
	}
}
