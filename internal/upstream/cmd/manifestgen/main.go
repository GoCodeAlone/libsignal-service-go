package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GoCodeAlone/libsignal-service-go/internal/upstream"
)

type artifactSpec struct {
	upstreamPath string
	localPath    string
	mode         string
}

var artifacts = []artifactSpec{
	{"rust/net/grpc/proto/TextSecure.proto", "proto/signal/net/grpc/TextSecure.proto", "generated"},
	{"rust/net/grpc/proto/google/rpc/error_details.proto", "proto/signal/net/grpc/google/rpc/error_details.proto", "generated"},
	{"rust/net/grpc/proto/google/rpc/status.proto", "proto/signal/net/grpc/google/rpc/status.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/account.proto", "proto/signal/net/grpc/org/signal/chat/account.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/backups.proto", "proto/signal/net/grpc/org/signal/chat/backups.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/challenge.proto", "proto/signal/net/grpc/org/signal/chat/challenge.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/common.proto", "proto/signal/net/grpc/org/signal/chat/common.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/credentials.proto", "proto/signal/net/grpc/org/signal/chat/credentials.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/device.proto", "proto/signal/net/grpc/org/signal/chat/device.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/errors.proto", "proto/signal/net/grpc/org/signal/chat/errors.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/keys.proto", "proto/signal/net/grpc/org/signal/chat/keys.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/messages.proto", "proto/signal/net/grpc/org/signal/chat/messages.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/profile.proto", "proto/signal/net/grpc/org/signal/chat/profile.proto", "generated"},
	{"rust/net/grpc/proto/org/signal/chat/require.proto", "proto/signal/net/grpc/org/signal/chat/require.proto", "generated-support"},
	{"rust/net/grpc/proto/org/signal/chat/tag.proto", "proto/signal/net/grpc/org/signal/chat/tag.proto", "generated-support"},
	{"rust/net/src/proto/chat_websocket.proto", "proto/signal/net/src/chat_websocket.proto", "generated"},
}

func main() {
	tag := flag.String("tag", "", "upstream signalapp/libsignal tag")
	out := flag.String("out", "", "manifest output path")
	printTag := flag.Bool("print-tag", false, "print current manifest tag")
	flag.Parse()

	if *printTag {
		m, err := readManifest("internal/upstream/manifest.json")
		must(err)
		fmt.Println(m.UpstreamTag)
		return
	}
	if *tag == "" || *out == "" {
		flag.Usage()
		os.Exit(2)
	}

	blobSHAs, err := upstreamBlobSHAs(*tag)
	must(err)

	m := upstream.Manifest{
		UpstreamRepo:             "signalapp/libsignal",
		UpstreamTag:              *tag,
		DescriptorChecksumSHA256: descriptorChecksum(),
	}

	for _, spec := range artifacts {
		content, err := os.ReadFile(spec.localPath)
		must(err)
		sum := sha256.Sum256(content)
		m.Artifacts = append(m.Artifacts, upstream.Artifact{
			UpstreamPath: spec.upstreamPath,
			LocalPath:    spec.localPath,
			BlobSHA:      blobSHAs[spec.upstreamPath],
			SHA256:       hex.EncodeToString(sum[:]),
			Mode:         spec.mode,
			Header:       headerSnippet(content),
		})
	}

	data, err := json.MarshalIndent(m, "", "  ")
	must(err)
	data = append(data, '\n')
	must(os.MkdirAll(filepath.Dir(*out), 0o755))
	must(os.WriteFile(*out, data, 0o644))
}

func readManifest(path string) (upstream.Manifest, error) {
	var m upstream.Manifest
	data, err := os.ReadFile(path)
	if err != nil {
		return m, err
	}
	return m, json.Unmarshal(data, &m)
}

func upstreamBlobSHAs(tag string) (map[string]string, error) {
	cmd := exec.Command("gh", "api", "repos/signalapp/libsignal/git/trees/"+tag+"?recursive=1", "--jq", ".tree[] | [.path,.sha] | @tsv")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("query upstream tree: %w", err)
	}
	result := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		fields := strings.Split(line, "\t")
		if len(fields) == 2 {
			result[fields[0]] = fields[1]
		}
	}
	for _, spec := range artifacts {
		if result[spec.upstreamPath] == "" {
			return nil, fmt.Errorf("missing upstream blob sha for %s", spec.upstreamPath)
		}
	}
	return result, nil
}

func descriptorChecksum() string {
	cmd := exec.Command("buf", "build", "-o", "-")
	out, err := cmd.Output()
	must(err)
	sum := sha256.Sum256(out)
	return hex.EncodeToString(sum[:])
}

func headerSnippet(content []byte) string {
	content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
	lines := strings.Split(string(content), "\n")
	var picked []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" && len(picked) == 0 {
			continue
		}
		if strings.HasPrefix(trimmed, "syntax =") {
			break
		}
		picked = append(picked, line)
		if len(picked) >= 8 {
			break
		}
	}
	return strings.TrimSpace(strings.Join(picked, "\n"))
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

