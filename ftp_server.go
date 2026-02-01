package main

import (
	"fmt"
	"log"
	"net"
	"os"
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
