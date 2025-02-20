package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Merith-TK/Resonite.Portable/steam"
	"github.com/Merith-TK/utils/debug"
)

func main() {
	flag.Parse()

	var appid = "2519830" // Intentionally set so that this project can be forked with ease

	args := flag.Args()
	debug.Print("Args: ", args)

	if len(args) == 0 {
		curPWD, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
			return
		}
		curPWD = filepath.ToSlash(curPWD)
		curPWD, err = filepath.Abs(curPWD)
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
			return
		}

		if _, err := os.Stat(filepath.Join(curPWD, ".steamcmd/steamapps/common/Resonite/Resonite.exe")); os.IsNotExist(err) {
			if len(args) == 0 {
				steam.Steamcmd("+app_license_request", appid, "+app_update", appid, "validate", "+quit")
				return
			}
		}
		// Run the desired game with the provided arguments
		cmd := exec.Command(".steamcmd/steamapps/common/Resonite/Resonite.exe", "-DataPath", filepath.Join(curPWD, "Data"), "-CachePath", filepath.Join(curPWD, "Cache"), "--UserPath", filepath.Join(curPWD, "User"), "--LogPath", filepath.Join(curPWD, "Logs"), "--AppID", appid)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to run Resonite: %v", err)
		}
		return
	}

	switch args[0] {
	case "update":
		steam.Steamcmd("+app_update", appid, "validate", "+quit")
	case "login":
		if len(args) > 2 {
			steam.SaveUsername(args[1])
			steam.Steamcmd("+login", args[1], args[2], "+quit")
		} else if len(args) == 2 {
			steam.SaveUsername(args[1])
			steam.Steamcmd("+login", args[1], "+quit")
		} else {
			// No username or password provided, use anonymous login
			log.Println("Error: No username or password provided. Please provide at least a username.")
			log.Println("Usage: resoport login <username> [password]")
		}
		return
	}

}
