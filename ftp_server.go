package main

import (
	"net"
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
