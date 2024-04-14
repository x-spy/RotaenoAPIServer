package APIHandlers

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

func CheckForUpdatesApiHandler(w http.ResponseWriter, r *http.Request) {
	officialUrl := "https://rotaeno.leancloud.indie.xd.com/1.1/call/CheckForUpdates"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read response body."+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	req, err := http.NewRequest(r.Method, officialUrl, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		http.Error(w, "Failed to create request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//复制Header
	for k, vv := range r.Header {
		for _, v := range vv {
			if k == "Host" {
				req.Header.Set("Host", "rotaeno.leancloud.indie.xd.com")
			} else {
				req.Header.Add(k, v)
			}
		}
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send data to rotaeno server."+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 复制响应头部
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(respData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
