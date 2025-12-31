package socks5

import (
	"errors"
	"io"
)

type ClientAuthMessage struct {
	Version  byte     //版本
	NMethods byte     //客户端支持的认证方法数量
	Methods  []Method //可变长的method
} //接收端结构体

type ClientPasswordMessage struct {
	Username string
	Password string
} //用户名密码认证报文体

type Method = byte

const (
	MethodNoAuth       Method = 0x00
	MethodGSSAPI       Method = 0x01
	MethodPassword     Method = 0x02
	MethodNoAcceptable Method = 0xFF
)

const (
	PasswordMethodVersion byte = 0x01
	PasswordAuthSuccess   byte = 0x00
	PasswordAuthFailure   byte = 0x01
)

var (
	ErrPasswordCheckNotSet = errors.New("Password auth method selected, but no password checker set")
	ErrPasswordAuthFailed  = errors.New("Username/password authentication failed")
)

// 从流中读取报文并生成相应的报文体
func NewClientAuthMessage(conn io.Reader) (*ClientAuthMessage, error) {
	//Read version and nmethods
	buf := make([]byte, 2)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}

	//Validate version
	if buf[0] != SOCKS5Version { //版本不是socks5
		return nil, ErrVersionNotSupported
	}
	//Read methods
	nMethods := buf[1]
	buf = make([]byte, nMethods)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}

	return &ClientAuthMessage{
		Version:  SOCKS5Version,
		NMethods: nMethods,
		Methods:  buf,
	}, nil
}

func NewServerAuthMessage(conn io.Writer, method Method) error {
	buf := []byte{SOCKS5Version, byte(method)}
	_, err := conn.Write(buf) //io.writer的Write方法一定会写入len(buf)字符，不同于io.Reader.Read()
	return err
}

func NewClientPasswordMessage(conn io.Reader) (*ClientPasswordMessage, error) {
	//Read version
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	version, usernameLen := buf[0], buf[1]
	if version != PasswordMethodVersion {
		return nil, ErrPasswordAuthFailed
	}

	// Read username,password length
	buf = make([]byte, usernameLen+1)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	username, passwordLen := string(buf[:len(buf)-1]), buf[len(buf)-1]

	// Read password
	if len(buf) < int(passwordLen) {
		buf = make([]byte, passwordLen)
	}
	if _, err := io.ReadFull(conn, buf[:passwordLen]); err != nil {
		return nil, err
	}

	return &ClientPasswordMessage{
		Username: username,
		Password: string(buf[:passwordLen]),
	}, nil
}

func WriteServerPasswordMessage(conn io.Writer, status byte) error {
	_, err := conn.Write([]byte{PasswordMethodVersion, status})
	return err
}
