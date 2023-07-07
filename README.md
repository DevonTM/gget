# Gget

A command-line interface (CLI) downloader written in Golang that supports resuming downloads.

## Features

- [x] Resuming downloads
- [x] Custom HTTP Headers
- [x] HTTP and SOCKS5 Proxies
- [x] Netscape Cookies
- [x] Concurrent Chunk Downloading

## Build Instructions

To build the CLI downloader, follow the instructions below:

1. Clone the repository
    ```bash
    git clone github.com/DevonTM/gget
    ```
2. Navigate to the cloned repository
    ```bash
    cd gget
    ```
3. Build the project
    ```bash
    make build
    ```

## Usage

To use Gget, run the following command:

```bash
gget [OPTIONS] URL
```

Available options are:

| Option | Description |
| --- | --- |
| `-url` | URL to download, can be set from last argument |
| `-o` | Output path, default current directory |
| `-O` | Output filename, default from server |
| `-f` | Force download, overwrites existing file |
| `-c` | Chunk size in bytes, default 1M |
| `-j` | Maximum number of concurrent chunks download, default 4 |
| `-t` | Maximum retry count, default 3 |
| `-r` | Set HTTP referer |
| `-ua` | Set HTTP user agent |
| `-p` | Proxy URL |
| `-C` | Netscape cookie file |
| `-H` | Custom HTTP header |
| `-no-h2` | Disable HTTP/2 |
| `-help` | Show help |
| `-version` | Show version |

## Examples
    
```bash
gget -o file.zip https://example.com/file.zip
```
```bash
gget -o foo -O bar file.zip https://example.com/file.zip
```
```bash
gget -p socks5://127.0.0.1:1080 https://example.com/file.zip
```
```bash
gget -H "Accept:text/plain" -H "User-Agent:gget/1.0" https://example.com/file.zip
```
```bash
gget -C cookies.txt https://example.com/file.zip
```

## TODO

- [ ] Add tests
- [ ] Support user/password for proxy
- [ ] Better error messages

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
