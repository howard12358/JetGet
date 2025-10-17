package test

import (
	"JetGet/backend/util"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeadRequestAndGetFileName(t *testing.T) {
	testCases := []struct {
		name             string
		handler          http.HandlerFunc
		requestPath      string
		expectedFilename string
		expectErr        bool
	}{
		{
			name: "Success with Content-Disposition",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Accept-Ranges", "bytes")
				w.Header().Set("Content-Length", "1024") // 任意一个 > 0 的值
				w.Header().Set("Content-Disposition", `attachment; filename="testfile.zip"`)
				w.WriteHeader(http.StatusOK)
			},
			requestPath:      "/download",
			expectedFilename: "testfile.zip",
			expectErr:        false,
		},
		{
			name: "Success with UTF8 Content-Disposition",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Accept-Ranges", "bytes")
				w.Header().Set("Content-Length", "2048")
				w.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''%e6%b5%8b%e8%af%95%e6%96%87%e4%bb%b6.txt`)
				w.WriteHeader(http.StatusOK)
			},
			requestPath:      "/download",
			expectedFilename: "测试文件.txt",
			expectErr:        false,
		},
		{
			name: "Fallback to URL Path",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Accept-Ranges", "bytes")
				w.Header().Set("Content-Length", "4096")
				w.WriteHeader(http.StatusOK)
			},
			requestPath:      "/files/document.pdf",
			expectedFilename: "document.pdf",
			expectErr:        false,
		},
		{
			name: "Failure with ambiguous URL and no header",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Accept-Ranges", "bytes")
				w.Header().Set("Content-Length", "512")
				w.WriteHeader(http.StatusOK)
			},
			requestPath:      "/",
			expectedFilename: "unrecognized",
			expectErr:        false,
		},
		{
			name: "HEAD request fails with server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			requestPath:      "/download",
			expectedFilename: "",
			expectErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			fullURL := server.URL + tc.requestPath
			resp, err := util.DoHeadRequest(fullURL, "")

			if (err != nil) != tc.expectErr {
				t.Fatalf("doHeadRequest() error = %v, expectErr %v", err, tc.expectErr)
			}

			if resp != nil {
				defer resp.Body.Close()
			}

			if tc.expectErr {
				return
			}

			filename := util.GetFileName(resp)

			if filename != tc.expectedFilename {
				t.Errorf("getFileName() = %v, want %v", filename, tc.expectedFilename)
			}
		})
	}
}
