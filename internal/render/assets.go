package render

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func avatarDataURI(pathOrURL string) (string, error) {
	pathOrURL = strings.TrimSpace(pathOrURL)
	if pathOrURL == "" {
		return "", nil
	}
	if strings.HasPrefix(pathOrURL, "data:") {
		return pathOrURL, nil
	}
	var reader io.ReadCloser
	var contentType string
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		resp, err := http.Get(pathOrURL)
		if err != nil {
			return "", err
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			return "", fmt.Errorf("failed to fetch icon: %s", resp.Status)
		}
		reader = resp.Body
		contentType = resp.Header.Get("Content-Type")
	} else {
		file, err := os.Open(pathOrURL)
		if err != nil {
			return "", err
		}
		reader = file
		contentType = contentTypeFromExt(pathOrURL)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", contentType, encoded), nil
}

func contentTypeFromExt(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	default:
		return ""
	}
}
