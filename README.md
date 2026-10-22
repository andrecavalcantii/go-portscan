# goscan

![CI](https://github.com/SEU-USUARIO/goscan/actions/workflows/ci.yml/badge.svg)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

Port scanner TCP **concorrente** escrito em Go. Usa um pool de goroutines para testar centenas de portas ao mesmo tempo e identifica serviços comuns (SSH, HTTP, MySQL, RDP...).

> 🇺🇸 *A fast, concurrent TCP port scanner written in Go, using a goroutine worker pool.*

```
$ goscan -host 127.0.0.1 -ports 1-9000

Escaneando 127.0.0.1 (9000 portas, 200 workers)...

  135    aberta   MSRPC
  445    aberta   SMB
  8080   aberta   HTTP-Alt

3 porta(s) aberta(s) em 210ms
```

## ⚠️ Aviso legal

Use **somente** em máquinas suas ou com autorização por escrito. Escanear sistemas de terceiros sem permissão pode ser crime. Para treinar, use `127.0.0.1` ou `scanme.nmap.org` (disponibilizado pelo projeto Nmap para testes).

## Instalação

Requer [Go 1.22+](https://go.dev/dl/). No Windows: `winget install GoLang.Go`

```bash
git clone https://github.com/SEU-USUARIO/goscan.git
cd goscan
go build -o goscan.exe .   # no Linux/Mac: go build -o goscan .
```

## Uso

```bash
goscan                                   # 127.0.0.1, portas 1-1024
goscan -host scanme.nmap.org -ports 22,80,443
goscan -ports 1-65535 -workers 500 -timeout 300ms
goscan -ports 80,443 -json               # saída em JSON
```

| Flag | Padrão | Descrição |
|---|---|---|
| `-host` | `127.0.0.1` | IP ou domínio do alvo |
| `-ports` | `1-1024` | `80`, `22,80,443`, `1-1024` ou misturado (`22,8000-8100`) |
| `-workers` | `200` | goroutines testando ao mesmo tempo |
| `-timeout` | `500ms` | tempo máximo de espera por porta |
| `-json` | `false` | imprime o resultado em JSON |

## Como funciona

1. `ParsePorts` transforma a especificação (`22,80,1-100`) em uma lista sem duplicatas
2. `Scan` cria N workers (goroutines) que leem portas de um canal `jobs`
3. Cada worker chama `net.DialTimeout("tcp", host:porta)`: se a conexão abre, a porta está **aberta**
4. Os resultados voltam por um canal `results`, são ordenados e filtrados

```
 portas ──► [jobs] ──► worker 1 ─┐
                   ├─► worker 2 ─┼──► [results] ──► ordena ──► imprime
                   └─► worker N ─┘
```

Sem concorrência, 1024 portas fechadas com timeout de 500ms levariam ~8 minutos. Com 200 workers, poucos segundos.

## Testes

```bash
go test ./...
```

Os testes sobem um servidor TCP local numa porta aleatória e conferem se o scanner o encontra.

## Roadmap

- [ ] Banner grabbing (ler a resposta do serviço)
- [ ] Escanear faixas de IP (`192.168.0.0/24`)
- [ ] Saída em CSV
- [ ] Barra de progresso

## Licença

[MIT](LICENSE) © 2026 André Silveira Cavalcanti
