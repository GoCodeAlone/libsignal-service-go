package upstream

type Manifest struct {
	UpstreamRepo             string     `json:"upstream_repo"`
	UpstreamTag              string     `json:"upstream_tag"`
	DescriptorChecksumSHA256 string     `json:"descriptor_checksum_sha256"`
	Artifacts                []Artifact `json:"artifacts"`
}

type Artifact struct {
	UpstreamPath string `json:"upstream_path"`
	LocalPath    string `json:"local_path"`
	BlobSHA      string `json:"blob_sha"`
	SHA256       string `json:"sha256"`
	Mode         string `json:"mode"`
	Header       string `json:"header"`
}

