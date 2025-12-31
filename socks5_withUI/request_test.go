package socks5

import (
	"bytes"
	"testing"
)

//客户端请求测试
func TestNewClientRequestMessage(t *testing.T) {
	tests := []struct{
		Version byte
		Cmd Command
		AddrType AddressType
		Address []byte
		Port []byte
		Error error
		Message ClientRequestMessage
	}{
		{
			Version: SOCKS5Version,
			Cmd: CmdConnect,
			AddrType: TypeIPv4,
			Address: []byte{192,168,1,1},
			Port: []byte{0x00, 0x50},
			Error: nil,
			Message: ClientRequestMessage{
				Cmd: CmdConnect,
				Address: "192.168.1.1",
				Port: 0x0050,
			},
		},
		{
			Version: 0x00,
			Cmd: CmdBind,
			AddrType: TypeIPv4,
			Address: []byte{192,168,1,1},
			Port: []byte{0x00, 0x50},
			Error: ErrVersionNotSupported,
			Message: ClientRequestMessage{
				Cmd: CmdBind,
				Address: "192.168.1.1",
				Port: 0x0050,
			},
		},
}

	for _,test := range tests {
	var buf bytes.Buffer
	buf.Write([]byte{SOCKS5Version, byte(test.Cmd), ReservedField, byte(test.AddrType)})
	buf.Write(test.Address)
	buf.Write(test.Port)

	message, err := NewClientRequestMessage(&buf)
	if err != test.Error {
		t.Fatalf("should get error %s , but got %s \n", test.Error, err)
		}

	if err != nil {
		return
	}
	
	if *message != test.Message {
		t.Fatalf("should get message %v ,but got %v \n", test.Message, *message)
		}
	}
}