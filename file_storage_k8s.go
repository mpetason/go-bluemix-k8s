package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
		return
	}
	counter := 0
	for _, fileStorage := range fileStorage {

		notes, _ := json.Marshal(fileStorage.Notes)

		if strings.Contains(string(notes), "cluster") {
			counter += 1
			clusters := string(notes)
			clustersSplit := strings.Split(clusters, ",")
			var clustersSplitClean []string
			clustersSplitClean = append(clustersSplitClean, strings.Replace(clustersSplit[2], "\\", "", -1),
				strings.Replace(clustersSplit[5], "\\", "", -1),
				strings.Replace(clustersSplit[6], "\\", "", -1))
			fmt.Println(counter, "ID:", *fileStorage.Id, clustersSplitClean[0], clustersSplitClean[1], clustersSplitClean[2])
		}
	}

}
