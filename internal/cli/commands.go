package cli

import (
	"github.com/spf13/cobra"
)

// CreateRootCommand creates the root command with all subcommands
func CreateRootCommand(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "peerchat-cli",
		Short: "Xelvra P2P Messenger CLI - Decentralized messaging without servers",
		Long: `Xelvra P2P Messenger CLI - Decentralized messaging without servers

QUICK START:
  1. peerchat-cli init     # Create your identity
  2. peerchat-cli doctor   # Test network connectivity
  3. peerchat-cli start    # Start interactive chat

STANDALONE COMMANDS (no running node required):
  init, doctor, version, manual, help

INTERACTIVE COMMANDS (available in chat mode):
  /help, /peers, /discover, /connect, /status, /quit

NODE-DEPENDENT COMMANDS (require running node):
  send, discover, status, listen

Performance targets:
- Latency: <50ms for direct connections
- Memory: <20MB idle usage
- CPU: <1% idle usage`,
		Version: version,
	}

	// Add subcommands
	rootCmd.AddCommand(createInitCommand())
	rootCmd.AddCommand(createStartCommand())
	rootCmd.AddCommand(createStatusCommand())
	rootCmd.AddCommand(createVersionCommand(version))
	rootCmd.AddCommand(createSendCommand())
	rootCmd.AddCommand(createConnectCommand())
	rootCmd.AddCommand(createListenCommand())
	rootCmd.AddCommand(createDiscoverCommand())
	rootCmd.AddCommand(createIdCommand())
	rootCmd.AddCommand(createProfileCommand())
	rootCmd.AddCommand(createSendFileCommand())
	rootCmd.AddCommand(createStopCommand())
	rootCmd.AddCommand(createSetupCommand())
	rootCmd.AddCommand(createDoctorCommand())
	rootCmd.AddCommand(createManualCommand(version))

	return rootCmd
}

// createInitCommand creates the init command
func createInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Xelvra identity and configuration",
		Run:   RunInit,
	}
}

// createStartCommand creates the start command
func createStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the P2P node and begin interactive chat",
		Run:   RunStart,
	}
	cmd.Flags().Bool("daemon", false, "Run as background daemon")
	return cmd
}

// createStatusCommand creates the status command
func createStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Display current node status and statistics",
		Run:   RunStatus,
	}
}

// createVersionCommand creates the version command
func createVersionCommand(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			RunVersion(version)
		},
	}
}

// createSendCommand creates the send command
func createSendCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "send [peer_id] [message]",
		Short: "Send a message to a peer",
		Args:  cobra.ExactArgs(2),
		Run:   RunSend,
	}
}

// createConnectCommand creates the connect command
func createConnectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "connect [peer_id]",
		Short: "Connect to a peer",
		Args:  cobra.ExactArgs(1),
		Run:   RunConnect,
	}
}

// createListenCommand creates the listen command
func createListenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "listen",
		Short: "Listen for incoming messages (passive mode)",
		Run:   RunListen,
	}
}

// createDiscoverCommand creates the discover command
func createDiscoverCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "discover",
		Short: "Discover peers in the network",
		Run:   RunDiscover,
	}
}

// createIdCommand creates the id command
func createIdCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "id",
		Short: "Show your identity information",
		Run:   RunShowID,
	}
}

// createProfileCommand creates the profile command
func createProfileCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "profile [peer_id]",
		Short: "Show profile information for a peer",
		Args:  cobra.ExactArgs(1),
		Run:   RunProfile,
	}
}

// createSendFileCommand creates the send-file command
func createSendFileCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "send-file [peer_id] [file_path]",
		Short: "Send a file to a peer",
		Args:  cobra.ExactArgs(2),
		Run:   RunSendFile,
	}
}

// createStopCommand creates the stop command
func createStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the running P2P node",
		Run:   RunStop,
	}
}

// createSetupCommand creates the setup command
func createSetupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Run the interactive setup wizard",
		Run:   RunSetup,
	}
}

// createDoctorCommand creates the doctor command
func createDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose and fix network issues",
		Run:   RunDoctor,
	}
}

// createManualCommand creates the manual command
func createManualCommand(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "manual",
		Short: "Show comprehensive manual",
		Run: func(cmd *cobra.Command, args []string) {
			RunManual(version)
		},
	}
}
