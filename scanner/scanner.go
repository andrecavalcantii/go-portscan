// Package scanner testa portas TCP de forma concorrente.
package scanner

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Result guarda o resultado de uma porta testada.
type Result struct {
	Port    int    `json:"port"`
	Open    bool   `json:"open"`
	Service string `json:"service,omitempty"`
}

// services mapeia portas comuns para o nome do serviço.
var services = map[int]string{
	21: "FTP", 22: "SSH", 23: "Telnet", 25: "SMTP", 53: "DNS",
	80: "HTTP", 110: "POP3", 135: "MSRPC", 139: "NetBIOS", 143: "IMAP",
	443: "HTTPS", 445: "SMB", 3306: "MySQL", 3389: "RDP",
	5432: "PostgreSQL", 6379: "Redis", 8080: "HTTP-Alt", 27017: "MongoDB",
}

// ServiceName devolve o nome do serviço conhecido para a porta, ou "".
func ServiceName(port int) string {
	return services[port]
}

// ParsePorts converte "22,80,443" ou "1-1024" (ou os dois misturados) em uma lista de portas.
func ParsePorts(spec string) ([]int, error) {
	seen := map[int]bool{}
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		start, end, err := parseRange(part)
		if err != nil {
			return nil, err
		}
		for p := start; p <= end; p++ {
			seen[p] = true
		}
	}
	if len(seen) == 0 {
		return nil, fmt.Errorf("nenhuma porta informada")
	}
	ports := make([]int, 0, len(seen))
	for p := range seen {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	return ports, nil
}

func parseRange(part string) (int, int, error) {
	bounds := strings.SplitN(part, "-", 2)
	start, err := parsePort(bounds[0])
	if err != nil {
		return 0, 0, err
	}
	end := start
	if len(bounds) == 2 {
		end, err = parsePort(bounds[1])
		if err != nil {
			return 0, 0, err
		}
	}
	if start > end {
		return 0, 0, fmt.Errorf("intervalo inválido: %q", part)
	}
	return start, end, nil
}

func parsePort(s string) (int, error) {
	p, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || p < 1 || p > 65535 {
		return 0, fmt.Errorf("porta inválida: %q (use 1-65535)", s)
	}
	return p, nil
}

// ScanPort tenta abrir uma conexão TCP; se conseguir, a porta está aberta.
func ScanPort(host string, port int, timeout time.Duration) Result {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return Result{Port: port}
	}
	conn.Close()
	return Result{Port: port, Open: true, Service: ServiceName(port)}
}

// Scan testa todas as portas usando um pool de workers (goroutines)
// e devolve os resultados ordenados por porta.
func Scan(host string, ports []int, workers int, timeout time.Duration) []Result {
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan int)
	results := make(chan Result)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				results <- ScanPort(host, port, timeout)
			}
		}()
	}

	go func() {
		for _, p := range ports {
			jobs <- p
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	all := make([]Result, 0, len(ports))
	for r := range results {
		all = append(all, r)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Port < all[j].Port })
	return all
}

// OpenOnly filtra só as portas abertas.
func OpenOnly(results []Result) []Result {
	open := []Result{}
	for _, r := range results {
		if r.Open {
			open = append(open, r)
		}
	}
	return open
}
