# 🛡️ VPN ADS Router for Linux

A Go-based proxy that emulates the functionality of the TwinCAT ADS Router on Linux. Seamlessly forwards Beckhoff ADS traffic across VPN connections like [NetBird](https://www.netbird.io), enabling engineering tools to interact with PLCs without manual route configuration.

## ✨ Features

- 🔁 **ADS Proxy Router**: Fully emulates an AMS Router for TwinCAT Engineering Tools
- 🌐 **Network Scanner**: Automatically discovers Beckhoff PLCs on the local subnet
- 🧠 **Dynamic Routing**: Rewrites AMS source/destination NetIDs transparently
- 🔒 **VPN Aware**: Integrates with NetBird, supports peer isolation and discovery
- ⚙️ **Zero Config on Clients**: No need to define AMS routes in TwinCAT or on PLCs

---

## 🗃️ Project Structure

```
vpn-ads-router/
├── cmd/                   # Main application entrypoint
├── internal/
│   ├── ads/              # Core ADS routing logic
│   ├── network/          # Discovery & scanning
│   ├── proxy/            # TCP proxy and session management
│   ├── vpn/              # VPN (NetBird) integration
│   └── config/           # Config loading & defaults
├── pkg/                  # Reusable packages (AMS, logging, utils)
├── scripts/              # Helper scripts for build/deploy
├── configs/              # YAML configuration files
├── deployments/          # NixOS & Docker support
├── docs/                 # Protocol, VPN setup, and architecture notes
├── tests/                # Unit tests
└── README.md
```

---

## 🚀 Getting Started

### 📦 Requirements

- Go 1.21+
- Linux system (tested on NixOS)
- NetBird VPN (or similar peer-to-peer mesh VPN)
- Beckhoff PLC (for testing)
- TwinCAT Engineering Tools (on remote client)

### 🛠️ Installation

```bash
git clone https://github.com/rosco/vpn-ads-router.git
cd vpn-ads-router
go build -o bin/vpn-ads-router ./cmd/vpn-ads-router
```

### 🧪 Running

Update `configs/config.yaml` as needed:

```yaml
ads:
  listen_port: 48898
  local_netid: "5.100.200.1.1.1"
  scan_subnet: "192.168.1.0/24"
  plc_port: 8016

vpn:
  mode: netbird
  isolate_peers: true
```

Then run:

```bash
./bin/vpn-ads-router
```

---

## 📚 Documentation

- [`docs/ads_protocol.md`](docs/ads_protocol.md): TwinCAT AMS/ADS protocol internals
- [`docs/vpn_setup.md`](docs/vpn_setup.md): NetBird VPN isolation setup
- [`docs/architecture.md`](docs/architecture.md): Router design overview

---

## 🧰 Development

```bash
go test ./...
```

Contributions welcome via PRs or issues!

---

