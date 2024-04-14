package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"rotaenoAPIServer/APIHandlers"
	"rotaenoAPIServer/Utils"
)

const (
	dbName    = "rotaenoAPIServer"
	tableName = "allowedObjectIDs"
)

func main() {

	var config Utils.Config

	err := Utils.InitConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	config, err = Utils.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Init Directory
	executablePath, _ := os.Executable()
	rootPath := filepath.Dir(executablePath)
	dataPath := filepath.Join(rootPath, config.FilePath)

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		_ = os.MkdirAll(dataPath, 0755)
	}

	// Init API
	apiVersion := config.APIVersion
	handleFunctions(apiVersion)

	// Init Http Server
	err = http.ListenAndServe(":"+config.ServerPort, nil)
	if err != nil {
		fmt.Println("Failed to open server: ", err)
		return
	}

}

func handleFunctions(apiVersion string) {
	http.HandleFunc("/"+apiVersion+"/users", APIHandlers.GetUsersInformationHandler)
	http.HandleFunc("/"+apiVersion+"/users/", APIHandlers.GetUsersInformationHandler)
	http.HandleFunc("/"+apiVersion+"/call/GetPurchasedItemData", APIHandlers.GetPurchasedItemDataApiHandler)
	http.HandleFunc("/"+apiVersion+"/call/CheckForUpdates", APIHandlers.CheckForUpdatesApiHandler)
	http.HandleFunc("/"+apiVersion+"/classes/CloudSave", APIHandlers.CloudSaveApiHandler)
	http.HandleFunc("/"+apiVersion+"/classes/CloudSave/", APIHandlers.CloudSaveApiHandler)
	http.HandleFunc("/"+apiVersion+"/call/CheckSafeText", APIHandlers.CheckSafeTextApiHandler)
	http.HandleFunc("/"+apiVersion+"/call/GetAllFolloweeSocialData", APIHandlers.CopyAndResendApiHandler)
	http.HandleFunc("/"+apiVersion+"/call/IncreaseFriendCap", APIHandlers.CopyAndResendApiHandler)
	http.HandleFunc("/"+apiVersion+"/call/FollowPlayer", APIHandlers.CopyAndResendApiHandler)
}
