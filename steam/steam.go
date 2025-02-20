package steam

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Merith-TK/utils/archive"
	"github.com/Merith-TK/utils/debug"
)

// SteamcmdDir is the path where SteamCMD is installed.
var (
	SteamcmdDir = ".steamcmd"
)

// StartSteamClient launches the Steam client using a command-line call.
// It opens the Steam client if installed on the system.
func StartSteamClient() error {
	cmd := exec.Command("cmd", "/C", "start", "steam://open/main")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Steam: %v", err)
	}
	return nil
}

// StopSteamClient terminates the Steam client process using taskkill.
// It forcibly closes Steam if running.
func StopSteamClient() error {
	cmd := exec.Command("taskkill", "/F", "/IM", "steam.exe")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Steam: %v", err)
	}
	return nil
}

// SetupSteamcmd ensures that SteamCMD is available on the system.
// If SteamCMD is not present, it downloads and extracts the SteamCMD installer.
func SetupSteamcmd() error {
	steamcmdPath := filepath.Join(SteamcmdDir, "steamcmd.exe")

	if _, err := os.Stat(steamcmdPath); os.IsNotExist(err) {
		fmt.Println("steamcmd not found, downloading...")

		// Download steamcmd.zip
		resp, err := http.Get("https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip")
		if err != nil {
			return fmt.Errorf("failed to download steamcmd.zip: %v", err)
		}
		defer resp.Body.Close()

		// Create steamcmd.zip file in the TEMP directory
		out, err := os.Create(os.Getenv("TEMP") + "/steamcmd.zip")
		if err != nil {
			return fmt.Errorf("failed to create steamcmd.zip: %v", err)
		}
		defer out.Close()

		// Copy the downloaded data to steamcmd.zip
		if _, err := io.Copy(out, resp.Body); err != nil {
			return fmt.Errorf("failed to write steamcmd.zip: %v", err)
		}

		// Extract steamcmd.zip to SteamcmdDir
		if err := archive.Unzip(filepath.Join(os.Getenv("TEMP"), "steamcmd.zip"), SteamcmdDir); err != nil {
			return fmt.Errorf("failed to extract steamcmd.zip: %v", err)
		}

		fmt.Println("steamcmd downloaded and extracted successfully to", SteamcmdDir)
	}

	return nil
}

// Steamcmd runs SteamCMD with the provided arguments and handles authentication.
// It returns a buffer containing the command's output.
func Steamcmd(args ...string) (bytes.Buffer, error) {
	SetupSteamcmd()
	debug.SetTitle("CMD")
	defer debug.ResetTitle()

	outputBuffer := bytes.Buffer{}

	// Check if username and password are provided
	username, err := LoadUsername()
	if err != nil {
		log.Println("Error loading username:", err)
		log.Println("Please login first using the 'login' command.")
		log.Println("Usage: resoport login <username> [password]")
		return outputBuffer, fmt.Errorf("username not found, please login first")
	}
	args = append([]string{"+login", string(username)}, args...)

	log.Println("[RESO] Running steamcmd with args:\n", args)

	// Run steamcmd with arguments
	cmd := exec.Command(SteamcmdDir+"\\steamcmd.exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// TODO: Parse output for new item uploads to extract workshop ID
	// cmd.Stdout = io.MultiWriter(os.Stdout, &outputBuffer)
	// cmd.Stderr = io.MultiWriter(os.Stderr, &outputBuffer)

	return outputBuffer, cmd.Run()
}

func SaveUsername(username string) error {
	usernameFile := filepath.Join(SteamcmdDir, "username.txt")
	if err := os.WriteFile(usernameFile, []byte(username), 0644); err != nil {
		return fmt.Errorf("failed to save username: %v", err)
	}
	return nil
}
func LoadUsername() (string, error) {
	usernameFile := filepath.Join(SteamcmdDir, "username.txt")
	content, err := os.ReadFile(usernameFile)
	if err != nil {
		return "", fmt.Errorf("failed to load username: %v", err)
	}
	return string(content), nil
}
