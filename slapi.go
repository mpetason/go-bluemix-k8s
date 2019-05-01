package main

import (
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
	fileStorage, err := service.Mask("id;notes").Limit(100).GetNetworkStorage()
	if err != nil {
		fmt.Printf("Error retrieving File Storage from account: %s\n", err)
		return
	} else {
		fmt.Println("File Storage:")
	}

	for _, fileStorage := range fileStorage {
		fmt.Printf("\tID: %d -- Notes: %s\n", *fileStorage.Id, *fileStorage.Notes)
	}

}
