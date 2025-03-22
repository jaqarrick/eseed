# eseed

A command-line tool for creating and seeding torrents using DHT and public trackers. This tool allows you to create magnet links and .torrent files from local files or directories, and immediately start seeding them.

## Features

- Create magnet links from files or directories
- Generate .torrent files
- Seeding with DHT (Distributed Hash Table) support
- Integration with public trackers for better peer discovery
- Real-time upload statistics
- Port forwarding support via UPnP
- Support for both single files and directories

## Installation

1. Make sure you have Go 1.16 or later installed
2. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/eseed.git
   cd eseed
   ```
3. Build the project:
   ```bash
   make build
   ```
   Or manually:
   ```bash
   go build -o eseed
   ```

## Development

This project includes a Makefile for common development tasks:

```bash
# Build the project
make build

# Run tests
make test

# Clean build artifacts
make clean

# Download dependencies
make deps

# Create a test file
make testfile

# Run in development mode
make dev

# Kill any running instances
make kill

# Show all available commands
make help
```

## Usage

### Basic Usage

Create a torrent and start seeding a file:
```bash
./eseed -input path/to/your/file
```

### Command Line Options

- `-input`: Path to the file or directory to create torrent from (required)
- `-magnet`: Only create magnet link (don't create .torrent file)
- `-output`: Custom path to save the .torrent file (optional, defaults to input path + .torrent)

### Examples

1. Create and seed a single file:
   ```bash
   ./eseed -input myfile.txt
   ```

2. Create only a magnet link:
   ```bash
   ./eseed -input myfile.txt -magnet
   ```

3. Create a torrent file with custom output path:
   ```bash
   ./eseed -input myfile.txt -output custom/path/file.torrent
   ```

4. Create and seed a directory:
   ```bash
   ./eseed -input my_directory
   ```

## Technical Details

- Uses the `anacrolix/torrent` library for BitTorrent functionality
- DHT is enabled by default for peer discovery
- Includes popular public trackers for better peer discovery
- Listens on port 42069 by default
- Supports UPnP for automatic port forwarding
- Piece size is set to 256KB for optimal performance

## Public Trackers

The following public trackers are included by default:
- udp://tracker.opentrackr.org:1337/announce
- udp://tracker.openbittorrent.com:6969/announce
- udp://open.stealth.si:80/announce
- udp://exodus.desync.com:6969/announce
- udp://tracker.torrent.eu.org:451/announce

## Notes

- The program will continue running until you press Ctrl+C
- Make sure port 42069 is available or change it in the code
- UPnP must be enabled on your router for automatic port forwarding
- The program requires read access to the input file/directory

## License

MIT License - feel free to use this code for any purpose.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
