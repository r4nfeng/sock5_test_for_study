package socks5

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var (
	ErrVersionNotSupported       = errors.New("protocol version not supported")
	ErrMethodVersionNotSupported = errors.New("sub-negotiation version not supported")
	ErrCommandNotSupported       = errors.New("command not supported")
	ErrInvalidReservedField      = errors.New("invalid reserved field")
	ErrAddressTypeNotSupported   = errors.New("address type not supported")
)

const (
	SOCKS5Version = 0x05
	ReservedField = 0x00
)

// 监听端口
type Server interface {
	Run() error
}

type SOCKS5Server struct {
	IP     string
	Port   int
	Config *Config
}

type Config struct {
	AuthMethod      Method
	PasswordChecker func(username, password string) bool
	TCPTimeout      time.Duration
}

// 初始化配置
func initConfig(config *Config) error {
	if config.AuthMethod == MethodPassword && config.PasswordChecker == nil {
		return ErrPasswordCheckNotSet
	}
	return nil
}

// 运行服务器
func (s *SOCKS5Server) Run() error {
	//Initialize server configuration
	if err := initConfig(s.Config); err != nil {
		return err
	}

	//监听在tcp的端口上，监听失败返回err，监听成功返回listen
	address := fmt.Sprintf("%s:%d", s.IP, s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	//不间断接收来自客户端的请求（这里的请求是三次握手）
	for {
		//Accept从listener到的已经完成三次握手的队列中取出一个TCP连接；返回已连接状态或者错误
		conn, err := listener.Accept()
		if err != nil {
			//发生错误时不能直接中止循环内容，否则会导致服务器停止运行，不能继续处理下一个请求
			//所以打印日志并继续下一个请求
			log.Printf("connection failure from %s: %s", conn.RemoteAddr(), err)
			continue
		}

		go func() {
			defer conn.Close() //处理完连接后关闭连接，发不发生panic都会执行，保证关闭
			if err := s.handleConnection(conn); err != nil {
				log.Printf("handle connection from %s: %s", conn.RemoteAddr(), err)
			}
		}()
	}
}

// 处理连接
func (s *SOCKS5Server) handleConnection(conn net.Conn) error {
	//协商过程
	if err := s.auth(conn); err != nil {
		return err
	}

	//请求过程
	return s.request(conn)

}

// 协商过程
func (s *SOCKS5Server) auth(conn io.ReadWriter) error {
	//读取客户端认证消息
	clientMessage, err := NewClientAuthMessage(conn)
	if err != nil {
		return err
	}

	//Only support no authentication（跳过子协商过程）
	var acceptable bool
	for _, method := range clientMessage.Methods {
		if method == s.Config.AuthMethod {
			acceptable = true
		}
	}
	if !acceptable {
		NewServerAuthMessage(conn, MethodNoAcceptable)
		return errors.New("method are not supported")
	}

	//Send no authentication required method
	if err := NewServerAuthMessage(conn, s.Config.AuthMethod); err != nil {
		return err
	}

	if s.Config.AuthMethod == MethodPassword {
		cpm, err := NewClientPasswordMessage(conn)
		if err != nil {
			return err
		}

		if !s.Config.PasswordChecker(cpm.Username, cpm.Password) {
			WriteServerPasswordMessage(conn, PasswordAuthFailure)
			return ErrPasswordAuthFailed
		}

		if err := WriteServerPasswordMessage(conn, PasswordAuthSuccess); err != nil {
			return err
		}
	}

	return nil
}

// 请求过程
func (s *SOCKS5Server) request(conn io.ReadWriter) error {
	message, err := NewClientRequestMessage(conn)
	if err != nil {
		return err
	}

	// Check if the address type is supported
	if message.AddrType == TypeIPv6 {
		WriteRequestFailureMessage(conn, ReplyAddressTypeNotSupported)
		return ErrAddressTypeNotSupported
	}

	if message.Cmd == CmdConnect {
		return s.handleTCP(conn, message)
	} else if message.Cmd == CmdUDPAssociate {
		return s.handleUDP()
	} else {
		WriteRequestFailureMessage(conn, ReplyCommandNotSupported)
		return ErrCommandNotSupported
	}
}

func (s *SOCKS5Server) handleUDP() error {
	return nil
}

func (s *SOCKS5Server) handleTCP(conn io.ReadWriter, message *ClientRequestMessage) error {
	// 请求访问目标TCP服务
	address := fmt.Sprintf("%s: %d", message.Address, message.Port)
	targetConn, err := net.DialTimeout("tcp", address, s.Config.TCPTimeout)
	if err != nil {
		WriteRequestFailureMessage(conn, ReplyHostUnreachable)
		return err
	}

	// Send success reply to client
	addrValue := targetConn.LocalAddr()
	addr := addrValue.(*net.TCPAddr)
	if err := WriteRequestSuccessMessage(conn, addr.IP, uint16(addr.Port)); err != nil {
		return err
	}

	return forward(conn, targetConn)
}

// 转发过程
func forward(conn io.ReadWriter, targetConn io.ReadWriteCloser) error {
	defer targetConn.Close()
	//双向转发
	go io.Copy(targetConn, conn)
	_, err := io.Copy(conn, targetConn)
	return err
}
