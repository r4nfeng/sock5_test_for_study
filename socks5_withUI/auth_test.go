package socks5

import (
	"bytes"
	"log"
	"reflect"
	"testing"
)

// 单元测试
func TestNewClientAuthMessage(t *testing.T) {
	//成功
	t.Run("should generate a message", func(t *testing.T) {
		//用依赖注入假设成功的io.Reader
		b := []byte{SOCKS5Version, 2, MethodNoAuth, MethodPassword}
		r := bytes.NewReader(b)

		message, err := NewClientAuthMessage(r)
		if err != nil {
			t.Fatalf("want error = nil but got %s", err)
		}

		//验证message的每个字段是否有效
		if message.Version != SOCKS5Version {
			t.Fatalf("want version = %d but got %d", SOCKS5Version, message.Version)
		}

		if message.NMethods != 2 {
			t.Fatalf("want nmethods = 2 but got %d", message.NMethods)
		}

		//判断底层数据是否相等
		if !reflect.DeepEqual(message.Methods, []Method{MethodNoAuth, MethodPassword}) {
			t.Fatalf("want methods = %v but got %v", []Method{MethodNoAuth, MethodPassword}, message.Methods)
		}
	})

	//失败案例：比如长度不够
	t.Run("methods length is shorter than nmethods", func(t *testing.T) {
		b := []byte{SOCKS5Version, 2, MethodNoAuth, MethodPassword}
		r := bytes.NewReader(b)

		_, err := NewClientAuthMessage(r)
		if err == nil {
			t.Fatalf("want error but got nil")
		}
	})
}

// 测试服务器认证消息的生成
func TestNewServerAuthMessage(t *testing.T) {
	t.Run("should send noauth", func(t *testing.T) {
		var buf bytes.Buffer
		err := NewServerAuthMessage(&buf, MethodNoAuth)
		if err != nil {
			t.Fatalf("want  get nil error but got %s", err)
		}

		got := buf.Bytes()
		if !reflect.DeepEqual(got, []byte{SOCKS5Version, MethodNoAuth}) {
			t.Fatalf("want send %v, but send %v", []byte{SOCKS5Version, MethodNoAuth}, got)
		}
	})

	t.Run("should send no acceptable", func(t *testing.T) {
		var buf bytes.Buffer
		err := NewServerAuthMessage(&buf, MethodNoAcceptable)
		if err != nil {
			t.Fatalf("want error = nil but got %s", err)
		}

		got := buf.Bytes()
		if !reflect.DeepEqual(got, []byte{SOCKS5Version, MethodNoAcceptable}) {
			t.Fatalf("want send %v, but send %v", []byte{SOCKS5Version, MethodNoAcceptable}, got)
		}
	})
}

func TestNewClientPasswordMessage(t *testing.T) {
	t.Run("valid password auth message", func(t *testing.T) {
		//构造一个包含用户名和密码的报文
		username, password := "testuser", "testpass"
		var buf bytes.Buffer
		buf.Write([]byte{PasswordMethodVersion, 5}) //版本和用户名长度
		buf.WriteString(username)
		buf.WriteByte(6) //密码长度
		buf.WriteString(password)

		message, err := NewClientPasswordMessage(&buf)
		if err != nil {
			log.Fatalf("want error = nil but got %s", err)
		}

		want := ClientPasswordMessage{
			Username: username,
			Password: password,
		}
		if *message != want {
			log.Fatalf("want message = %v but got %v", want, message)
		}
	})
}
