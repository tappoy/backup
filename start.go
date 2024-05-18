package main

import (
	"github.com/tappoy/env"

	"os"
	"os/signal"
	"syscall"
	"time"
	"strconv"
	"net"
	"errors"
	"strings"
)

var ErrInterrupted = errors.New("Interrupted")

func accept(sock net.Listener) (string, error) {
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

func readLineFromSocket(sock net.Listener, cmd chan string, sigterm chan os.Signal) {
	errCount := 0
	for {
		select {
		case <-sigterm:
			return
		default:
			line, err := accept(sock)
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					env.ENotice("socket closed. readLineFromSocket stopped")
					return
				}
				if errCount == 1 {
					env.EError("failed to read from socket: %v", err)
				}
				errCount++
				continue
			}
			errCount = 0
			env.EDebug("received: %s\n", line)

			// trim line
			line = strings.TrimSpace(line)

			// send to channel
			cmd <- line
		}
	}
}

func start() int {
	// check if the service is already running
	if _, err := os.Stat(PidFile); err == nil {
		env.Errf("pid file exists. the service is already running\n")
		env.Errf("pid file: %v\n", PidFile)
		return 1
	}

	// check if the socket file exists
	if _, err := os.Stat(SocketFile); err == nil {
		env.Errf("socket file exists. the service is already running\n")
		env.Errf("socket file: %v\n", SocketFile)
		return 1
	}

	// get VAULT_PREFIX
	VaultPrefix = env.Getenv("VAULT_PREFIX", "")
	if VaultPrefix == "" {
		env.Errf("VAULT_PREFIX is not set\n")
		RunHelpMessage()
		return 1
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
	defer func () {
		sock.Close()
		os.Remove(SocketFile)
	}()

	// make channel for commands
	cmd := make(chan string)
	defer close(cmd)

	// make channel for termination signals
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigterm)
		close(sigterm)
	}()

	// read line from socket
	go readLineFromSocket(sock, cmd, sigterm)

	// last executed time
	var lastExecuted time.Time

	env.EInfo("service started")
	defer env.EInfo("service stopped")

	// loop
	for {
		select {
		case <-time.After(1 * time.Minute):
			env.EDebug("wake up")
			// check if it is time to execute (2:00 AM)
			now := time.Now()
			if checkTime(now, lastExecuted) {
				// execute
				env.EInfo("backup started by schedule")
				c := make(chan error)
				go backup(c, cmd, sigterm)
				err := <-c
				close(c)

				// check error
				if err == ErrInterrupted {
					env.EInfo("backup interrupted")
					return 0
				}
				if err != nil {
					env.EError("failed to backup: %v", err)
					continue
				}
				lastExecuted = now
			}
		case line := <-cmd:
			// execute command
			env.EDebug("command received: %s", line)
			if line == "stop" {
				env.EDebug("service stopped by command")
				return 0
			}
			doCommand(line)
			continue
		case <-sigterm:
			// interrupted
			env.EDebug("service stopped by signal")
			return 0
		}
	}
}

func checkTime(now time.Time, lastExecuted time.Time) bool {
	// check if it is time to execute (2:00 AM)
	if now.Hour() == 2 && lastExecuted.Day() != now.Day() {
		return true
	}
	return false
}

func backup(c chan error, cmd chan string, sigterm chan os.Signal) {
	for {
		select {
		case <-sigterm:
			env.EInfo("backup interrupted by signal")
			c <- ErrInterrupted
			return
		case line := <-cmd:
			// execute command
			env.EInfo("command received while backup: %s", line)
			if line == "stop" {
				env.EInfo("backup interrupted by stop command")
				c <- ErrInterrupted
				return
			} else {
				env.EInfo("command ignored while backup: %s", line)
			}
		default:
			// do something
			time.Sleep(1 * time.Second)
			env.Errf(".")
		}
	}
}

func doCommand(line string) {
	// doCommand
}
