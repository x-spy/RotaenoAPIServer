package APIHandlers

import (
	"bytes"
	"io"
	"net/http"
	"rotaenoAPIServer/Utils"
	"time"
)

func CloudSaveApiHandler(w http.ResponseWriter, r *http.Request) {
	config, err := Utils.GetConfig()
	officialUrl := "https://rotaeno.leancloud.indie.xd.com/" + config.APIVersion + "/classes/CloudSave"
	var requestUrl string

	path := r.URL.Path

	base := "/" + config.APIVersion + "/classes/CloudSave"
	if len(path) > len(base) {
		extraPath := path[len(base):]
		requestUrl = officialUrl + extraPath
	} else {
		extra := r.URL.RawQuery
		requestUrl = officialUrl + "?" + extra
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read response body."+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	req, err := http.NewRequest(r.Method, requestUrl, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		http.Error(w, "Failed to create request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers
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

	// Copy response headers
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
