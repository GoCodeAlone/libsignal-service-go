package upstream

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestManifestCoversCopiedArtifacts(t *testing.T) {
	m := readTestManifest(t)
	if m.UpstreamRepo != "signalapp/libsignal" {
		t.Fatalf("upstream repo = %q", m.UpstreamRepo)
	}
	if m.UpstreamTag != "v0.96.4" {
		t.Fatalf("upstream tag = %q, want v0.96.4", m.UpstreamTag)
	}
	if len(m.Artifacts) == 0 {
		t.Fatal("manifest has no artifacts")
	}
	for _, a := range m.Artifacts {
		if a.UpstreamPath == "" || a.LocalPath == "" || a.BlobSHA == "" || a.SHA256 == "" || a.Mode == "" {
			t.Fatalf("incomplete artifact entry: %+v", a)
		}
		data, err := os.ReadFile("../../" + a.LocalPath)
		if err != nil {
			t.Fatalf("%s missing: %v", a.LocalPath, err)
		}
		sum := sha256.Sum256(data)
		if got := hex.EncodeToString(sum[:]); got != a.SHA256 {
			t.Fatalf("%s sha256 = %s, want %s", a.LocalPath, got, a.SHA256)
		}
		if strings.Contains(a.Header, "Signal Messenger") && !strings.Contains(a.Header, "SPDX-License-Identifier") {
			t.Fatalf("%s header records copyright without SPDX: %q", a.LocalPath, a.Header)
		}
	}
}

func TestDescriptorChecksum(t *testing.T) {
	m := readTestManifest(t)
	if m.DescriptorChecksumSHA256 == "" {
		t.Fatal("descriptor checksum is empty")
	}
	cmd := exec.Command("buf", "build", "-o", "-")
	cmd.Dir = "../.."
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("buf build: %v", err)
	}
	sum := sha256.Sum256(out)
	if got := hex.EncodeToString(sum[:]); got != m.DescriptorChecksumSHA256 {
		t.Fatalf("descriptor checksum = %s, want %s", got, m.DescriptorChecksumSHA256)
	}
}

func readTestManifest(t *testing.T) Manifest {
	t.Helper()
	data, err := os.ReadFile("manifest.json")
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	return m
}

