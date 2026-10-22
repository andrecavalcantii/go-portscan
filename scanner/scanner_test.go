package scanner

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestParsePorts(t *testing.T) {
	cases := []struct {
		spec string
		want []int
	}{
		{"80", []int{80}},
		{"22,80,443", []int{22, 80, 443}},
		{"1-3", []int{1, 2, 3}},
		{"443, 20-22, 22", []int{20, 21, 22, 443}},
	}
	for _, c := range cases {
		got, err := ParsePorts(c.spec)
		if err != nil {
			t.Fatalf("ParsePorts(%q) erro inesperado: %v", c.spec, err)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ParsePorts(%q) = %v, esperado %v", c.spec, got, c.want)
		}
	}
}

func TestParsePortsInvalid(t *testing.T) {
	for _, spec := range []string{"", "0", "70000", "abc", "100-50"} {
		if _, err := ParsePorts(spec); err == nil {
			t.Errorf("ParsePorts(%q) deveria dar erro", spec)
		}
	}
}

// Abre um servidor local numa porta aleatória e confere se o scanner a encontra.
func TestScanFindsOpenPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	open := OpenOnly(Scan("127.0.0.1", []int{port}, 4, time.Second))
	if len(open) != 1 || open[0].Port != port {
		t.Fatalf("esperava a porta %d aberta, veio %v", port, open)
	}
}

func TestScanClosedPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close() // fecha antes de escanear

	if r := ScanPort("127.0.0.1", port, 300*time.Millisecond); r.Open {
		t.Fatalf("porta %d deveria estar fechada", port)
	}
}
