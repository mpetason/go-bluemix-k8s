package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

func main() {
	sess := session.New(os.Getenv("SL_USERNAME"), os.Getenv("SL_APIKEY"))

	//sess.Debug = true

	doListBlockVolumes(sess)
}

func doListBlockVolumes(sess *session.Session) {
	// Get the Account service for Block Storage
	service := services.GetAccountService(sess)

	// List Block Storage
	fileStorage, err := service.Limit(500).GetNetworkStorage()

	if err != nil {
		fmt.Printf("Error retrieving File Storage from account: %s\n", err)
	} else {
		counter := 0
		var notes map[string]interface{}

		for _, fileStorage := range fileStorage {

			notes = make(map[string]interface{})
			if fileStorage.Notes != nil {
				json.Unmarshal([]byte(*fileStorage.Notes), &notes)
			}

			if _, ok := notes["cluster"]; ok {
				counter++
				fmt.Println(counter, "ID:", *fileStorage.Id, "Name:", *fileStorage.Username, "Cluster", notes["cluster"], "PV:", notes["pv"], "PVC:", notes["pvc"])
			}
		}
	}
}
