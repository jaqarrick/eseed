package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

func TestCreateTorrentFromFile(t *testing.T) {
	// Create a temporary test file
	content := []byte("test content for torrent")
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Create metainfo
	mi := &metainfo.MetaInfo{}
	info := metainfo.Info{
		PieceLength: 256 * 1024, // 256 KB pieces
	}

	// Add the file
	err = info.BuildFromFilePath(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to build torrent info: %v", err)
	}

	// Set the info in metainfo
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		t.Fatalf("Error encoding torrent info: %v", err)
	}

	// Test magnet link generation
	magnet, err := mi.MagnetV2()
	if err != nil {
		t.Fatalf("Failed to generate magnet link: %v", err)
	}
	magnetStr := magnet.String()
	if !strings.Contains(magnetStr, "magnet:?") {
		t.Error("Invalid magnet link format")
	}
}

func TestCreateTorrentFromDirectory(t *testing.T) {
	// Create a temporary directory with some files
	tmpdir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	// Create a few test files in the directory
	files := []string{"file1.txt", "file2.txt", "subdir/file3.txt"}
	for _, fname := range files {
		fpath := filepath.Join(tmpdir, fname)
		if strings.Contains(fname, "/") {
			if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
				t.Fatalf("Failed to create subdirectory: %v", err)
			}
		}
		if err := os.WriteFile(fpath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", fname, err)
		}
	}

	// Create metainfo
	mi := &metainfo.MetaInfo{}
	info := metainfo.Info{
		PieceLength: 256 * 1024, // 256 KB pieces
	}

	// Add the directory
	err = info.BuildFromFilePath(tmpdir)
	if err != nil {
		t.Fatalf("Failed to build torrent info: %v", err)
	}

	// Set the info in metainfo
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		t.Fatalf("Error encoding torrent info: %v", err)
	}

	// Test magnet link generation
	magnet, err := mi.MagnetV2()
	if err != nil {
		t.Fatalf("Failed to generate magnet link: %v", err)
	}
	magnetStr := magnet.String()
	if !strings.Contains(magnetStr, "magnet:?") {
		t.Error("Invalid magnet link format")
	}
}

func TestInvalidPath(t *testing.T) {
	info := metainfo.Info{
		PieceLength: 256 * 1024,
	}

	// Try to add a non-existent path
	err := info.BuildFromFilePath("/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for non-existent path, but got nil")
	}
}

func TestMagnetLinkFormat(t *testing.T) {
	// Create a temporary test file
	content := []byte("test content for magnet link")
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Create metainfo
	mi := &metainfo.MetaInfo{}
	info := metainfo.Info{
		PieceLength: 256 * 1024,
	}

	// Add the file
	err = info.BuildFromFilePath(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to build torrent info: %v", err)
	}

	// Set the info in metainfo
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		t.Fatalf("Error encoding torrent info: %v", err)
	}

	// Test magnet link format
	magnet, err := mi.MagnetV2()
	if err != nil {
		t.Fatalf("Failed to generate magnet link: %v", err)
	}
	magnetStr := magnet.String()

	// Check magnet link format
	if !strings.HasPrefix(magnetStr, "magnet:?") {
		t.Error("Magnet link should start with 'magnet:?'")
	}
	if !strings.Contains(magnetStr, "xt=urn:btih:") {
		t.Error("Magnet link should contain hash (xt=urn:btih:)")
	}
}

func TestDHTConfiguration(t *testing.T) {
	// Create a temporary test file
	content := []byte("test content for DHT")
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Create metainfo
	mi := &metainfo.MetaInfo{}
	info := metainfo.Info{
		PieceLength: 256 * 1024,
	}

	err = info.BuildFromFilePath(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to build torrent info: %v", err)
	}

	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		t.Fatalf("Error encoding torrent info: %v", err)
	}

	// Generate magnet link
	magnet, err := mi.MagnetV2()
	if err != nil {
		t.Fatalf("Failed to generate magnet link: %v", err)
	}

	// Test client configuration
	cfg := torrent.NewDefaultClientConfig()
	cfg.Seed = true
	cfg.DisableTrackers = true
	cfg.NoDHT = false
	cfg.DataDir = filepath.Dir(tmpfile.Name())
	cfg.UpnpID = "eseed"
	cfg.ListenPort = 42069

	// Verify DHT is enabled
	if cfg.NoDHT {
		t.Error("DHT should be enabled")
	}

	// Verify trackers are disabled
	if !cfg.DisableTrackers {
		t.Error("Trackers should be disabled")
	}

	// Verify seeding is enabled
	if !cfg.Seed {
		t.Error("Seeding should be enabled")
	}

	// Verify port configuration
	if cfg.ListenPort != 42069 {
		t.Errorf("Expected port 42069, got %d", cfg.ListenPort)
	}

	// Create client with a timeout
	client, err := torrent.NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create torrent client: %v", err)
	}
	defer client.Close()

	// Add torrent for seeding with a timeout
	torrent, err := client.AddMagnet(magnet.String())
	if err != nil {
		t.Fatalf("Failed to add magnet for seeding: %v", err)
	}

	// Verify torrent was added successfully
	if torrent == nil {
		t.Error("Torrent should not be nil")
	}

	// Verify magnet link was parsed correctly
	expectedHash := strings.TrimPrefix(magnet.InfoHash.String(), "Some(")
	expectedHash = strings.TrimSuffix(expectedHash, ")")
	actualHash := torrent.InfoHash().String()
	if expectedHash != actualHash {
		t.Errorf("Torrent info hash mismatch. Expected: %s, Got: %s", expectedHash, actualHash)
	}
}
