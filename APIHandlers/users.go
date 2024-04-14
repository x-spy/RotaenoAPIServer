package APIHandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"rotaenoAPIServer/Utils"
	"strconv"
	"time"
)

type userInformation struct {
	UpdatedAt               time.Time   `json:"updatedAt"`
	MobilePhoneVerified     bool        `json:"mobilePhoneVerified"`
	SessionToken            string      `json:"sessionToken"`
	MirrorFreeGemAmount     int         `json:"mirrorFreeGemAmount"`
	MirrorGemAmount         int         `json:"mirrorGemAmount"`
	ShortId                 string      `json:"shortId"`
	AuthData                authData    `json:"authData"`
	EmailVerified           bool        `json:"emailVerified"`
	MirrorInventoryItemSids []string    `json:"mirrorInventoryItemSids"`
	LoginRecord             loginRecord `json:"loginRecord"`
	Username                string      `json:"username"`
	ObjectId                string      `json:"objectId"`
	MirrorPaidGemAmount     int         `json:"mirrorPaidGemAmount"`
	CreatedAt               time.Time   `json:"createdAt"`
}

type authData struct {
	Xdg xdg `json:"xdg"`
}

type xdg struct {
	Username    string `json:"username"`
	Detail      detail `json:"detail"`
	Uid         string `json:"uid"`
	DeviceId    string `json:"device_id"`
	AccessToken string `json:"access_token"`
}

type detail struct {
	Source     int       `json:"source"`
	Avatar     string    `json:"avatar"`
	UserId     string    `json:"userId"`
	IsGuest    bool      `json:"isGuest"`
	UserCode   string    `json:"userCode"`
	AppId      int       `json:"appId"`
	RegistIp   string    `json:"registIp"`
	RegistTime time.Time `json:"registTime"`
	NickName   string    `json:"nickName"`
	LoginType  int       `json:"loginType"`
}

type loginRecord struct {
	Devices map[string]device `json:"Devices"`
}

type device struct {
	LastLoginTime date `json:"LastLoginTime"`
}

type date struct {
	Type string    `json:"__type"`
	Iso  time.Time `json:"iso"`
}

func GetUsersInformationHandler(w http.ResponseWriter, r *http.Request) {

	config, err := Utils.GetConfig()

	officialUrl := "https://rotaeno.leancloud.indie.xd.com/" + config.APIVersion + "/users"
	var requestUrl string

	path := r.URL.Path
	base := "/" + config.APIVersion + "/users"
	if len(path) > len(base) {
		extraPath := path[len(base):]
		requestUrl = officialUrl + extraPath
	} else {
		requestUrl = officialUrl
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
	// Prevent encoding errors under linux (Windows does not require this modification).
	req.Header.Set("Accept-Encoding", "")

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

	var userInfo userInformation
	err = json.Unmarshal(respData, &userInfo)
	if err != nil {
		http.Error(w, "Failed to process JSON data."+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get inventory item information through the provided keys
	inventoryItems := make([]string, 0)
	for keyName := range config.PurchaseInformation.Keys {
		if keyName == "main" {
			continue
		}
		inventoryItems = append(inventoryItems, keyName)
	}
	userInfo.MirrorInventoryItemSids = inventoryItems

	newUserInfo, err := json.Marshal(userInfo)
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
	w.Header().Set("Content-Length", strconv.Itoa(len(newUserInfo)))
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(newUserInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	//User data saving: Read the x-lc session token, and save the object id locally with x-lc session token as the file name.
	sessionToken := userInfo.SessionToken
	userObjectID := userInfo.ObjectId
	var file *os.File

	if sessionToken == "" {
		return
	}
	executablePath, _ := os.Executable()
	rootPath := filepath.Dir(executablePath)
	dataPath := filepath.Join(rootPath, config.FilePath)
	filePath := filepath.Join(dataPath, sessionToken)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err = os.Create(filePath)
		if err != nil {
			http.Error(w, "Couldn't create file on server."+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_, err = file.Write([]byte(userObjectID))
	fmt.Println("Object ID Saved for object id: " + userObjectID + " as session token: " + sessionToken)
	if err != nil {
		http.Error(w, "Failed to write file."+err.Error(), http.StatusInternalServerError)
		return
	}

}
