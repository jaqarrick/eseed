package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

// Common public trackers
var defaultTrackers = []string{
	"udp://tracker.opentrackr.org:1337/announce",
	"udp://tracker.openbittorrent.com:6969/announce",
	"udp://open.stealth.si:80/announce",
	"udp://exodus.desync.com:6969/announce",
	"udp://tracker.torrent.eu.org:451/announce",
}

func main() {
	var (
		inputPath  = flag.String("input", "", "Path to file or directory to create torrent from")
		magnetOnly = flag.Bool("magnet", false, "Only create magnet link (don't create .torrent file)")
		outputPath = flag.String("output", "", "Path to save .torrent file (optional, defaults to input path + .torrent)")
	)

	flag.Parse()

	if *inputPath == "" {
		log.Fatal("Please provide an input path using -input flag")
	}

	// Verify input path exists
	_, err := os.Stat(*inputPath)
	if err != nil {
		log.Fatalf("Error accessing input path: %v", err)
	}

	// Create torrent metainfo
	mi := &metainfo.MetaInfo{}

	// Create the info dictionary
	info := metainfo.Info{
		PieceLength: 256 * 1024, // 256 KB pieces
	}

	// Add the file/directory
	err = info.BuildFromFilePath(*inputPath)
	if err != nil {
		log.Fatalf("Error building torrent info: %v", err)
	}

	// Set the info in metainfo
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		log.Fatalf("Error encoding torrent info: %v", err)
	}

	// Add trackers to metainfo
	mi.AnnounceList = [][]string{}
	for _, tracker := range defaultTrackers {
		mi.AnnounceList = append(mi.AnnounceList, []string{tracker})
	}
	mi.Announce = defaultTrackers[0]

	// Generate magnet link
	magnet, err := mi.MagnetV2()
	if err != nil {
		log.Fatalf("Error generating magnet link: %v", err)
	}
	fmt.Printf("Magnet link: %v\n", magnet.String())

	if !*magnetOnly {
		// If no output path specified, use input path + .torrent
		if *outputPath == "" {
			*outputPath = *inputPath + ".torrent"
		}

		// Save .torrent file
		file, err := os.Create(*outputPath)
		if err != nil {
			log.Fatalf("Error creating torrent file: %v", err)
		}
		defer file.Close()

		// Write the torrent file
		err = mi.Write(file)
		if err != nil {
			log.Fatalf("Error writing torrent file: %v", err)
		}
		fmt.Printf("Created torrent file: %s\n", *outputPath)
	}

	// Configure and create torrent client for seeding
	cfg := torrent.NewDefaultClientConfig()
	cfg.Seed = true
	cfg.DisableTrackers = false // Enable trackers
	cfg.NoDHT = false           // Enable DHT
	cfg.DataDir = filepath.Dir(*inputPath)

	// Enable port forwarding and set listening port
	cfg.UpnpID = "eseed"
	cfg.ListenPort = 42069 // Default port for seeding

	// Create client with DHT and trackers enabled
	client, err := torrent.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating torrent client: %v", err)
	}
	defer client.Close()

	// Add the torrent for seeding
	t, err := client.AddMagnet(magnet.String())
	if err != nil {
		log.Fatalf("Error adding magnet for seeding: %v", err)
	}

	fmt.Println("Starting to seed... (Press Ctrl+C to stop)")
	fmt.Printf("DHT and trackers enabled, listening on port %d\n", cfg.ListenPort)
	fmt.Printf("Info hash: %s\n", t.InfoHash().String())
	fmt.Println("Using trackers:")
	for _, tracker := range defaultTrackers {
		fmt.Printf("  - %s\n", tracker)
	}

	// Wait indefinitely while seeding
	<-t.GotInfo()
	t.DownloadAll()

	// Print initial status
	fmt.Println("Torrent ready, waiting for peers...")

	for {
		time.Sleep(time.Second)
		stats := t.Stats()
		activePeers := len(t.PeerConns())
		fmt.Printf("\rSeeding... Upload: %.2f KB/s, Total: %.2f MB, Active Peers: %d    ",
			float64(stats.BytesWrittenData.Int64())/1024,
			float64(stats.BytesWrittenData.Int64())/(1024*1024),
			activePeers)
	}
}
