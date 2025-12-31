package socks5

import(
	"bytes"
	"net"
	"reflect"
	"testing"
)

func TestAuth(t *testing.T) {
	server := SOCKS5Server{
		Config: &Config{
			AuthMethod: MethodNoAuth,
		},
	}
	
	t.Run("a valid client auth message", func(t *testing.T) {
		var buf bytes.Buffer
		buf.Write([]byte{SOCKS5Version, 2, MethodNoAuth, MethodPassword})
		if err := server.auth(&buf); err != nil {
			t.Fatalf("want error = nil but got %s", err)
		}

		want := []byte{SOCKS5Version, MethodNoAuth}
		got := buf.Bytes()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("want response = %v but got %v", want, got)
		}	
	})

	t.Run("an invalid client auth message", func(t *testing.T) {
		var buf bytes.Buffer
		buf.Write([]byte{SOCKS5Version, 2, MethodNoAuth})
		if err := server.auth(&buf); err == nil {
			t.Fatalf("want error but got nil")
		}
	})
}

func TestWriteRequestSuccessMessage(t *testing.T) {
	var buf bytes.Buffer
	ip := net.IP([]byte{192,168,1,1})

	err := WriteRequestSuccessMessage(&buf, ip, 0x0439)
	if err != nil {
		t.Fatalf("error while writing: %s", err)
	}

	want := []byte{SOCKS5Version, ReplySucceeded, ReservedField, TypeIPv4, 192,168,1,1, 0x04, 0x39}
	got := buf.Bytes()
	if !reflect.DeepEqual(want,got) {
		t.Fatalf("message mismatch: want %v but got %v", want, got)
	}
}