# Enigma Protocol Go

Go implementation of the Enigma Protocol.

This is a enigma protocol server, that allows client to connect and send messages to each other and exchange public keys. The specification of the protocol can be found [here](https://github.com/PradyumnaKrishna/enigma-protocol/blob/main/SPECIFICATION.md).

Use the enigma protocol client application to connect to this server. The client is made for secure and end to end encrypted communication. The client application can be found [here](https://github.com/PradyumnaKrishna/enigma-protocol).

There is also a Python implementation of Enigma Protocol available.

## Getting Started

Download the dependencies using the following command:

```bash
go mod download
```

Build the project using the following command:

```bash
go build -o app cmd/main.go
```

Run the project using the following command:

```bash
./app
```

### Docker Container

You can also run the server using the provided Dockerfile. Build the image using the following command:

```bash
docker build -t enigma-protocol-go .
```

Use the provided environment variables to configure the server and run the container using the following command:

```bash
docker run -p 5000:5000 enigma-protocol-go
```

Also available on github container registry:

```bash
docker pull ghcr.io/pradyumnakrishna/enigma-protocol-go:latest
```

## Usage

The server will start on `localhost:5000`. You can provide configuration using environment variables. Here are the available environment variables:

- `PORT`: The port on which the server will run. Default is `5000`.
- `DATABASE_PATH`: The path to the database file, uses sqlite3 database. Default is `./sqlite3.db`.
- `ALLOWED_ORIGINS`: Comma separated list of allowed origins for CORS.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.
