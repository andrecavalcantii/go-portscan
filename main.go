// goscan: port scanner TCP concorrente escrito em Go.
//
// Use somente em máquinas suas ou com autorização explícita.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"goscan/scanner"
)

func main() {
	host := flag.String("host", "127.0.0.1", "IP ou domínio do alvo")
	ports := flag.String("ports", "1-1024", `portas: "80", "22,80,443" ou "1-1024"`)
	workers := flag.Int("workers", 200, "quantidade de goroutines testando ao mesmo tempo")
	timeout := flag.Duration("timeout", 500*time.Millisecond, "tempo máximo por porta (ex.: 300ms, 1s)")
	asJSON := flag.Bool("json", false, "imprime o resultado em JSON")
	flag.Parse()

	list, err := scanner.ParsePorts(*ports)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	start := time.Now()
	open := scanner.OpenOnly(scanner.Scan(*host, list, *workers, *timeout))
	elapsed := time.Since(start).Round(time.Millisecond)

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(open)
		return
	}

	fmt.Printf("Escaneando %s (%d portas, %d workers)...\n\n", *host, len(list), *workers)
	if len(open) == 0 {
		fmt.Println("Nenhuma porta aberta encontrada.")
	}
	for _, r := range open {
		service := r.Service
		if service == "" {
			service = "desconhecido"
		}
		fmt.Printf("  %-6d aberta   %s\n", r.Port, service)
	}
	fmt.Printf("\n%d porta(s) aberta(s) em %s\n", len(open), elapsed)
}
