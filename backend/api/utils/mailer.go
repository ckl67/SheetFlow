package utils

import (
	"backend/api/config"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
)

func SendMail(to, subject, body string) error {
	cfg := config.Config().Smtp

	if cfg.Enabled != "true" {
		return nil
	}

	addr := net.JoinHostPort(cfg.HostServerAddr, strconv.Itoa(cfg.HostServerPort))

	auth := smtp.PlainAuth(
		"",
		cfg.Username,
		cfg.Password,
		cfg.HostServerAddr,
	)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n"+
			"%s\r\n",
		cfg.From,
		to,
		subject,
		body,
	))

	// Connexion TLS explicite
	tlsConfig := &tls.Config{
		ServerName: cfg.HostServerAddr,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, cfg.HostServerAddr)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(cfg.From); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	return w.Close()
}
