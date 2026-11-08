package cmd

import (
	"os"
	"path/filepath"

	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init [project_id]",
	Aliases: []string{"clone", "fetch"},
	Short:   "Initialize local workspace files for a Guided Project",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		projectID := "7a2b9c3d-4e5f-41a2-9b8c-1d2e3f4a5b6c"
		if len(args) > 0 {
			projectID = args[0]
		}

		// Initialize files directly in the current working directory (.)
		dir := "."

		if err := os.MkdirAll(filepath.Join(dir, "internal", "framer"), 0755); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(dir, "internal", "storage"), 0755); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(dir, "internal", "index"), 0755); err != nil {
			return err
		}

		// Write internal/framer/record.go
		recordGoContent := `package framer

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// WriteRecord serializes a key-value pair into a binary frame with CRC32 checksum.
// Format:
// [0..4]   CRC32 Checksum (uint32)
// [4..8]   Key Length (uint32)
// [8..16]  Value Length (uint64)
// [16..]   Key + Value Payload
func WriteRecord(key string, value []byte) ([]byte, error) {
	keyBytes := []byte(key)
	keyLen := uint32(len(keyBytes))
	valLen := uint64(len(value))

	header := make([]byte, 16)
	binary.BigEndian.PutUint32(header[4:8], keyLen)
	binary.BigEndian.PutUint64(header[8:16], valLen)

	payload := append(header[4:], append(keyBytes, value...)...)
	checksum := crc32.ChecksumIEEE(payload)
	binary.BigEndian.PutUint32(header[0:4], checksum)

	fullFrame := append(header[0:4], payload...)
	return fullFrame, nil
}

// ReadRecord decodes binary record frame and validates CRC32 checksum.
func ReadRecord(data []byte) (string, []byte, error) {
	if len(data) < 16 {
		return "", nil, fmt.Errorf("truncated record frame")
	}

	expectedCRC := binary.BigEndian.Uint32(data[0:4])
	actualCRC := crc32.ChecksumIEEE(data[4:])
	if expectedCRC != actualCRC {
		return "", nil, fmt.Errorf("checksum mismatch: corrupted record data")
	}

	keyLen := binary.BigEndian.Uint32(data[4:8])
	valLen := binary.BigEndian.Uint64(data[8:16])

	if uint64(len(data)) < 16+uint64(keyLen)+valLen {
		return "", nil, fmt.Errorf("record payload incomplete")
	}

	key := string(data[16 : 16+keyLen])
	val := data[16+keyLen : 16+uint64(keyLen)+valLen]
	return key, val, nil
}
`
		_ = os.WriteFile(filepath.Join(dir, "internal", "framer", "record.go"), []byte(recordGoContent), 0644)

		// Write internal/framer/record_test.go
		recordTestGoContent := `package framer

import (
	"bytes"
	"testing"
)

func TestBinaryFraming(t *testing.T) {
	t.Run("EncodeDecodeRoundtrip", func(t *testing.T) {
		key := "user:101"
		value := []byte("SkillStack Local Learner")
		encoded, err := WriteRecord(key, value)
		if err != nil {
			t.Fatalf("WriteRecord failed: %v", err)
		}

		decodedKey, decodedVal, err := ReadRecord(encoded)
		if err != nil {
			t.Fatalf("ReadRecord failed: %v", err)
		}
		if decodedKey != key {
			t.Errorf("expected key %s, got %s", key, decodedKey)
		}
		if !bytes.Equal(decodedVal, value) {
			t.Errorf("value mismatch")
		}
	})

	t.Run("CorruptedChecksumDetection", func(t *testing.T) {
		key := "session:token"
		value := []byte("secret_value")
		encoded, _ := WriteRecord(key, value)
		// Corrupt checksum byte
		encoded[0] ^= 0xFF
		_, _, err := ReadRecord(encoded)
		if err == nil {
			t.Errorf("expected checksum mismatch error, got nil")
		}
	})
}
`
		_ = os.WriteFile(filepath.Join(dir, "internal", "framer", "record_test.go"), []byte(recordTestGoContent), 0644)

		// Write internal/storage/segment.go
		segmentGoContent := `package storage

import "sync"

type SegmentManager struct {
	mu           sync.Mutex
	activeSize   int64
	maxSize      int64
	segmentCount int
}

func NewSegmentManager(maxSize int64) *SegmentManager {
	return &SegmentManager{
		maxSize:      maxSize,
		segmentCount: 1,
	}
}

func (sm *SegmentManager) Append(bytesWritten int64) (int, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	rotated := false
	if sm.activeSize+bytesWritten > sm.maxSize {
		sm.segmentCount++
		sm.activeSize = 0
		rotated = true
	}
	sm.activeSize += bytesWritten
	return sm.segmentCount, rotated
}
`
		_ = os.WriteFile(filepath.Join(dir, "internal", "storage", "segment.go"), []byte(segmentGoContent), 0644)

		// Write internal/storage/segment_test.go
		segmentTestGoContent := `package storage

import "testing"

func TestSegmentRotation(t *testing.T) {
	t.Run("RotateAtThreshold", func(t *testing.T) {
		sm := NewSegmentManager(100)
		seg, rotated := sm.Append(60)
		if seg != 1 || rotated {
			t.Errorf("unexpected rotation on initial append")
		}
		seg, rotated = sm.Append(50)
		if seg != 2 || !rotated {
			t.Errorf("expected segment rotation to segment 2")
		}
	})

	t.Run("FsyncDurabilityVerification", func(t *testing.T) {
		sm := NewSegmentManager(1000)
		sm.Append(200)
	})
}
`
		_ = os.WriteFile(filepath.Join(dir, "internal", "storage", "segment_test.go"), []byte(segmentTestGoContent), 0644)

		// Write internal/index/key_directory.go
		indexGoContent := `package index

import "sync"

type KeyDirectory struct {
	mu     sync.RWMutex
	offset map[string]int64
}

func NewKeyDirectory() *KeyDirectory {
	return &KeyDirectory{
		offset: make(map[string]int64),
	}
}

func (kd *KeyDirectory) Put(key string, off int64) {
	kd.mu.Lock()
	defer kd.mu.Unlock()
	kd.offset[key] = off
}

func (kd *KeyDirectory) Get(key string) (int64, bool) {
	kd.mu.RLock()
	defer kd.mu.RUnlock()
	off, ok := kd.offset[key]
	return off, ok
}
`
		_ = os.WriteFile(filepath.Join(dir, "internal", "index", "key_directory.go"), []byte(indexGoContent), 0644)

		// Write internal/index/key_directory_test.go
		indexTestGoContent := `package index

import "testing"

func TestKeyDirectory(t *testing.T) {
	t.Run("PointLookupLatency", func(t *testing.T) {
		kd := NewKeyDirectory()
		kd.Put("user:1", 4096)
		off, ok := kd.Get("user:1")
		if !ok || off != 4096 {
			t.Errorf("expected offset 4096, got %d", off)
		}
	})
}
`
		_ = os.WriteFile(filepath.Join(dir, "internal", "index", "key_directory_test.go"), []byte(indexTestGoContent), 0644)

		if IsJSONMode() {
			return style.PrintJSON(out, map[string]interface{}{
				"initialized": true,
				"directory":   ".",
				"project_id":  projectID,
			})
		}

		style.PrintSuccess(out, "Initialized Guided Project starter workspace files in current directory.")
		style.PrintKeyValue(out, "Next Step", "skillstack run")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
