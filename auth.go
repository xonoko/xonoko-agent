package main

import (
	"io"
	"net"
	"bufio"
	"crypto/rand"
	"encoding/hex"
    "github.com/Sirupsen/logrus"
	"crypto/hmac"
	"crypto/sha256"
	"github.com/go-ini/ini"
	"strconv"
)

type Auth struct {
	Type  string
	Id    int
	Nonce string
	Hash  string
}

func proxyAuth(buf *bufio.Reader, conn net.Conn, key string) bool {
	result := false
	//Wait for sNonce
	sNonce := waitForString(buf)
	logger.WithFields(logrus.Fields{"sNonce": sNonce}).Debug("Received server Nonce")

	//Send cNonce
	cNonce := nonce()
	logger.WithFields(logrus.Fields{"cNonce": cNonce}).Debug("Send client Nonce")
	io.WriteString(conn, cNonce + "\n")

	//Send client hash
	cHash := cryptoHash(sNonce, cNonce, key)
	logger.WithFields(logrus.Fields{"cHash": cHash}).Debug("Send client Hash")
	io.WriteString(conn, cHash + "\n")

	answer := waitForString(buf)
	if answer == "ok" {
		result = true
		logrus.Info("Authenticated with server")
	} else {
		result = false
		logrus.Warn("Auth failure")
	}
	return result
}

/*
cryptoHash returns a cryptographic hash created from two nonces and a secret string.
 */
func cryptoHash(sNonce string, cNonce string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(sNonce))
	h.Write([]byte(cNonce))
	return hex.EncodeToString(h.Sum(nil))
}

/*
nonce creates a new nonce in the form of a hex string.
 */
func nonce() string {
	nonce := make([]byte, 32)
	_, err := rand.Read(nonce)
	strNonce := hex.EncodeToString(nonce)
	if err != nil {
		logger.WithFields(logrus.Fields{"Error": err}).Fatal("Error generating nonce")
	}
	return strNonce
}

/*
Auth is a rpc handler for authentification request
 */
func (t *Agent) Auth(args *Auth, reply *Auth) error {
	logger.Info("Auth request from server")

	cfg, err := ini.Load("config/agent.cfg")
	if err != nil {
		logger.WithFields(logrus.Fields{"path" : "config/agent.cfg", "err": err}).Error("Config file not found")
		return err
	}
	var id int
	var key string
	if cfg.Section("").HasKey("AGENTID") && cfg.Section("").HasKey("AGENTKEY"){
		//Install has been completed, auth as agent
		id = cfg.Section("").Key("AGENTID").MustInt()
		key = cfg.Section("").Key("AGENTKEY").MustString("key")
		reply.Type = "agent"
	} else {
		//Need to install to get agent credentials
		id = cfg.Section("").Key("INSTALLID").MustInt()
		key = cfg.Section("").Key("INSTALLKEY").MustString("key")
		reply.Type = "installer"
	}
	cNonce := nonce()
	reply.Id = id
	reply.Nonce = cNonce
	reply.Hash = cryptoHash(args.Nonce, cNonce, key)
	err = cfg.SaveTo("config/agent.cfg")
	if err != nil {
		logger.WithFields(logrus.Fields{"path" : "config/agent.cfg", "err": err}).Error("Can't write config file")
		return err
	}
	return nil
}

/*
Auth is a rpc handler for authentification request
 */
func (t *Agent) Create(args *Auth, reply *Bool) error {
	logger.Info("Getting agent credentials")
	cfg, err := ini.Load("config/agent.cfg")
	if err != nil {
		logger.WithFields(logrus.Fields{"path" : "config/agent.cfg", "err": err}).Error("Config file not found")
		return err
	}
	cfg.Section("").NewKey("AGENTID", strconv.Itoa(args.Id))
	cfg.Section("").NewKey("AGENTKEY", args.Nonce)
	cfg.SaveTo("config/agent.cfg")
	logger.WithFields(logrus.Fields{"id" : args.Id, "key" : args.Nonce}).Info("Agent conf edited")
	reply.Value = true
	return nil
}

