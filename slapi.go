package main

import (
	"fmt"
	"os"

	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

func main() {
	sess := session.New(os.Getenv("SL_USERNAME"), os.Getenv("SL_APIKEY"))

	sess.Debug = true

	doListBlockVolumes(sess)
}

func doListBlockVolumes(sess *session.Session) {
	service := services.GetAccountService(sess)

	fileStorage, err := service.Mask("id;capacity;createDate").Limit(100).GetPortableStorageVolumes()
	if err != nil {
		fmt.Printf("Error retrieving File Storage from account: %s\n", err)
		return
	} else {
		fmt.Println("File Storage:")
	}

	for _, fileStorage := range fileStorage {
		fmt.Printf("\tID: [%d] -- Capacity: %d GB -- Creation Date: %s\n", *fileStorage.Id, *fileStorage.Capacity, *fileStorage.CreateDate)
	}

}
