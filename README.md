# Fast Key Manager (fkm)

A CLI tool to manage SSH keys on Linux and Unix-based systems. Easily add SSH keys from GitHub or GitLab accounts, list existing SSH keys, and remove selected keys from your system.

## Features

- Add SSH keys manually or from GitHub/GitLab accounts
- List existing SSH keys
- Remove SSH keys
- Simple and intuitive CLI interface

## Installation

### Using Go

```bash
go install github.com/ferhatbostanci/fast-key-manager@latest
```

### Linux

#### Latest version
```bash
wget https://github.com/ferhatbostanci/fast-key-manager/releases/latest/download/fkm_linux_amd64 -O /usr/local/bin/fkm && \
    chmod +x /usr/local/bin/fkm
```

### MacOS

#### Latest version (for Apple Silicon)
```bash
wget https://github.com/ferhatbostanci/fast-key-manager/releases/latest/download/fkm_darwin_arm64 -O /usr/local/bin/fkm && \
    chmod +x /usr/local/bin/fkm
```

#### Latest version (for Intel)
```bash
wget https://github.com/ferhatbostanci/fast-key-manager/releases/latest/download/fkm_darwin_amd64 -O /usr/local/bin/fkm && \
    chmod +x /usr/local/bin/fkm
```

## Usage

### Add SSH Key
```bash
# Manual Copy-Paste (Default)
fkm add --manual

# Add from GitHub
fkm add --github <username>

# Add from GitLab
fkm add --gitlab <username>
```

### List SSH Keys
```bash
fkm list
```

### Remove SSH Key
```bash
fkm remove
```

### Check Version
```bash
fkm version
```

## Building from Source

1. Clone the repository
2. Run `go build -o fkm cmd/fkm/main.go`
3. Move the binary to your PATH
