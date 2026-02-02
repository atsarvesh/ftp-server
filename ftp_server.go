package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// constants for default FTP credentials

const (
	UserDefault = "user"
	PassDefault = "pass"
)

// transfer modes

const (
	BIN = iota // binary mode
	ASC        // ascii mode
)

// FTP response codes

const (
	FTPDataOpen         = 150
	FTPOK               = 200
	FTPInfo             = 211
	FTPSyst             = 215
	FTPNewConn          = 220
	FTPGoodbye          = 221
	FTPTransferComplete = 226
	FTPLoginSuccess     = 230
	FTPUserOKPassReq    = 331
	FTPDataConnFail     = 425
	FTPUnknownCmd       = 500
	FTPSyntaxErr        = 501
	FTPAuthErr          = 530
	FTPFileErr          = 550
)

// ClientContext holds information about the connected FTP client

type ClientContext struct {
	username       string
	password       string
	authenticated  bool
	dataPort       int
	transferMode   int
	clientDataIP   string
	dataSocketConn net.Conn
}

// global listener for the FTP server

var serverListener net.Listener

// formatFileMetadata formats the file metadata similar to Unix 'ls -l' output

func formatFileMetadata(filename string) string {

	info, err := os.Stat(filename)
	if err != nil {
		log.Printf("failed to get file info for %s: %v", filename, err)
		return fmt.Sprintf("error retrieving file info for %s", filename)
	}

	// file permissions

	mode := info.Mode()
	permissions := mode.String()

	nlink := 1 // number of hard links, default to 1

	// owner and group placeholders

	owner := "unknown"
	group := "unknown"

	// file size

	size := info.Size()

	// last modified time

	modTime := info.ModTime().Format(time.UnixDate)

	return fmt.Sprintf("%s %d %s  %s %d %s %s\r\n", permissions, nlink, owner, group, size, modTime, filename)
}

// setupActiveDataConnection establishes an active mode data connection to the client

func setupActiveDataConnection(ctx *ClientContext) (net.Conn, error) {

	address := fmt.Sprintf("%s:%d", ctx.clientDataIP, ctx.dataPort)

	conn, err := net.Dial("tcp", address, 10*time.Second)

	if err != nil {
		log.Printf("failed to establish data connection: %w", err)
	}
	return conn, nil
}

// sendFTPResponse sends an FTP response to the client

func sendFTPResponse(conn net.Conn, statusCode int, message string) {

	response := fmt.Sprintf("%d %s\r\n", statusCode, message)

	fmt.Printf("sent: %s", response)

	_, err := conn.Write([]byte(response))

	if err != nil {
		log.Printf("failed to send response: %v", err)
	}
}

// getDataIPPort parses the PORT command to extract the client's data IP and port

func getDataIPPort(buffer string, ctx *ClientContext) error {

	if !strings.HasPrefix(buffer, "PORT ") {
		return fmt.Errorf("invalid PORT command")
	}

	params := strings.TrimSpace(buffer[5:])

	parts := strings.Split(params, ",")

	if len(parts) != 6 {
		return fmt.Errorf("invalid PORT command parameters")
	}

	var nums [6]int

	for i, part := range parts {

		num, err := strconv.Atoi(strings.TrimSpace(part))

		if err != nil {
			return fmt.Errorf("invalid PORT command parameter: %s", part)
		}

		nums[i] = num
	}

	ctx.clientDataIP = fmt.Sprintf("%d.%d.%d.%d", nums[0], nums[1], nums[2], nums[3])

	ctx.dataPort = nums[4]*256 + nums[5]

	return nil
}

// handleUser processes the USER command from the client

func handleUser(conn net.Conn, buffer string, ctx *ClientContext) {

	username := strings.TrimSpace(buffer[5:])

	username = strings.TrimRight(username, "\r\n")

	ctx.username = username

	fmt.Printf("Username: %s\n", ctx.username)

	sendFTPResponse(conn, FTPUserOKPassReq, "User name okay, need password.\r\n")
}

// handlePass processes the PASS command from the client

func handlePass(conn net.Conn, buffer string, ctx *ClientContext) {

	password := strings.TrimSpace(buffer[5:])

	password = strings.TrimRight(password, "\r\n")

	ctx.password = password

	fmt.Printf("Password provided.\n")

	if ctx.username == "anon" {

		ctx.authenticated = true

		sendFTPResponse(conn, FTPLoginSuccess, "Login successful for anonymous user.\r\n")

	} else if ctx.username == UserDefault && ctx.password == PassDefault {

		ctx.authenticated = true

		sendFTPResponse(conn, FTPLoginSuccess, "Login successful.\r\n")
	} else {

		ctx.authenticated = false

		sendFTPResponse(conn, FTPAuthErr, "Authentication failed.\r\n")
	}
}
