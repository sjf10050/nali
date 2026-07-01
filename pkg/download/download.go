package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/zu1k/nali/pkg/common"
)

func Download(ctx context.Context, filePath string, urls ...string) (data []byte, err error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("未指定下载 url")
	}

	// Stream download: write to temp file, then rename
	tmpPath := filePath + ".tmp"

	// Try each URL until one succeeds
	var lastErr error
	for _, url := range urls {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("User-Agent", common.UserAgent)

		resp, err := common.GetHttpClient().Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
			continue
		}

		f, err := os.Create(tmpPath)
		if err != nil {
			resp.Body.Close()
			lastErr = err
			continue
		}

		// Cap the download size so a broken/hostile server can't fill the disk.
		const limit = int64(common.MaxResponseSize)
		n, err := io.Copy(f, io.LimitReader(resp.Body, limit+1))
		resp.Body.Close()
		f.Close()

		if err == nil && n > limit {
			err = fmt.Errorf("download exceeds max size %d bytes", limit)
		}
		if err != nil {
			os.Remove(tmpPath)
			lastErr = err
			continue
		}

		// Success — rename temp to target
		os.Remove(filePath)
		if err := os.Rename(tmpPath, filePath); err != nil {
			lastErr = err
			continue
		}

		// Read back for validation (existing callers expect []byte return)
		os.Remove(tmpPath) // cleanup
		data, err = os.ReadFile(filePath)
		return data, err
	}

	return nil, lastErr
}
