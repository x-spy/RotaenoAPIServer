package APIHandlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rotaenoAPIServer/Utils"
)

type keyResponse struct {
	Result struct {
		Keys map[string]string `json:"keys"`
	} `json:"result"`
}

func GetPurchasedItemDataApiHandler(w http.ResponseWriter, r *http.Request) {
	var responseJson keyResponse
	config, err := Utils.GetConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//The object ID of the owners of shared keys
	baseObjectID := config.PurchaseInformation.BaseObjectID
	baseKey := Utils.GetTransmissionKey(baseObjectID)

	responseJson.Result.Keys = config.PurchaseInformation.Keys

	xLcSession := r.Header.Get("X-LC-Session")
	if xLcSession == "" {
		http.Error(w, "Invalid session id.", http.StatusBadRequest)
		return
	}
	fmt.Println("User X-LC Session: " + xLcSession)
	executablePath, _ := os.Executable()
	rootPath := filepath.Dir(executablePath)
	dataPath := filepath.Join(rootPath, config.FilePath)
	filePath := filepath.Join(dataPath, xLcSession)
	userObjectIDBytes, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Couldn't find object id for this session token."+err.Error(), http.StatusBadRequest)
		return
	}
	userObjectID := string(userObjectIDBytes)
	fmt.Println("ObjectID: " + userObjectID)
	// Check current object id available
	objectIDMap := make(map[string]bool)
	for _, objectID := range config.AllowedObjectID {
		objectIDMap[objectID] = true
	}
	if _, exists := objectIDMap[userObjectID]; exists {
		key := Utils.GetTransmissionKey(userObjectID)

		// Process keys
		for keyName, base64Value := range responseJson.Result.Keys {
			// Decrypt
			decryptedKey, err := Utils.RotaenoDecryptFromBase64(base64Value, baseKey)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Re-encrypt
			newEncryptedKey, err := Utils.RotaenoEncryptToBase64(decryptedKey, key)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Update the corresponding fields of the structure
			responseJson.Result.Keys[keyName] = newEncryptedKey
		}

		response, err := json.Marshal(responseJson)
		if err != nil {
			http.Error(w, "Failed to marshal response."+err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			http.Error(w, "Failed to write response."+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		fmt.Println("Refused to provide keys. Object ID: " + userObjectID)
		w.WriteHeader(http.StatusForbidden)
		errString := "You do not have permission to get keys."
		_, err := w.Write([]byte(errString))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
