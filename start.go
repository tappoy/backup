package main

import (
	"github.com/tappoy/env"

	"os"
	"os/signal"
	"syscall"
	"time"
	"strconv"
	"net"
)

func readFromSocket(sock net.Listener) (string, error) {
	conn, err := sock.Accept()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// read line
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func start() int {
	VaultPrefix = env.Getenv("VAULT_PREFIX", "")
	if VaultPrefix == "" {
		env.Errf("VAULT_PREFIX is not set\n")
		RunHelpMessage()
		env.Exit(1)
	}

	// create pid file
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)
	err := os.WriteFile(PidFile, []byte(pidStr), 0644)
	if err != nil {
		env.Errf("failed to create pid file: %v\n", err)
		return 1
	}
	defer os.Remove(PidFile)

	// create socket
	sock, err := net.ListenUnix("unix", &net.UnixAddr{Name: SocketFile, Net: "unix"})
	if err != nil {
		env.Errf("failed to create socket: %v\n", err)
		return 1
	}
	defer sock.Close()

	// make channel for commands
	cmd := make(chan string)

	// read line from socket
	go func() {
		for {
			// pass through to channel
			line, err := readFromSocket(sock)
			if err != nil {
				env.Notice("failed to read from socket: %v\n", err)
				continue
			}
			env.Errf("received: %s\n", line)
			cmd <- line
		}
	}()

	// make channel for termination signals
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// last executed time
	var lastExecuted time.Time

	// Output to stderr until the server starts.
	env.Errf("service started\n")

	// After the server starts, output to a log file.
	env.Info("service started\n")

	// loop
	for {
		// sleep 1 minute
		select {
		case <-time.After(1 * time.Minute):
			env.Debug("wake up\n")
			// check if it is time to execute (2:00 AM)
			now := time.Now()
			if now.Hour() == 2 && lastExecuted.Day() != now.Day() {
				// execute
				env.Info("doBackup started by schedule\n")
				err := doBackup()
				if err != nil {
					env.Error("failed to doBackup: %v\n", err)
				} else {
					lastExecuted = now
				}
			}
		case line := <-cmd:
			// execute command
			env.Errf("command received: %s\n", line)
			if line == "stop" {
				env.Info("service stopped by command\n")
				return 0
			}
			doCommand(line)
			continue
		case <-sigterm:
			env.Errf("received termination signal\n")
			env.Info("service stopped by signal\n")
			return 0
		}
	}
}

func doBackup() error {
	// doBackup
	return nil
}

func doCommand(line string) {
	// doCommand
}
