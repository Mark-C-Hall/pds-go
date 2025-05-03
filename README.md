# pds-go

An [ATproto](https://atproto.com) Personal Data Server implementation written in Go.

## Overview

pds-go is a lightweight and efficient Personal Data Server (PDS) implementation for the AT Protocol, written in Go. It provides the core functionality needed to participate in the decentralized social network ecosystem built on ATproto.

## Features

- Core ATproto PDS functionality
- RESTful API implementation
- Data storage and synchronization
- Authentication and identity management
- Efficient Go implementation for better performance and resource utilization

## Getting Started

### Prerequisites

- Go 1.24 or later
- git

### Installation

```bash
git clone https://github.com/yourusername/pds-go.git
cd pds-go
go mod tidy
```

### Running the Server

```bash
go run cmd/pds/main.go
```

## Development Status

This project is in early development. Contributions and feedback are welcome.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [AT Protocol](https://atproto.com) - The underlying protocol
- [Bluesky](https://bsky.app) - Original ATproto implementation