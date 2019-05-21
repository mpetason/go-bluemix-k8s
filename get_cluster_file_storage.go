package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"
	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/softlayer/softlayer-go/services"
	slSession "github.com/softlayer/softlayer-go/session"
)

func main() {
	c := new(bluemix.Config)

	var org string
	flag.StringVar(&org, "org", "", "Bluemix Organization")

	var space string
	flag.StringVar(&space, "space", "", "Bluemix Space")

	flag.Parse()

	kubernetesRegions := []string{"au-syd", "eu-de", "eu-gb", "us-east", "us-south"}

	if org == "" || space == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	client, err := mccpv2.New(sess)

	if err != nil {
		log.Fatal(err)
	}

	region := sess.Config.Region
	orgAPI := client.Organizations()
	myorg, err := orgAPI.FindByName(org, region)

	if err != nil {
		log.Fatal(err)
	}

	spaceAPI := client.Spaces()
	myspace, err := spaceAPI.FindByNameInOrg(myorg.GUID, space, region)

	if err != nil {
		log.Fatal(err)
	}

	accClient, err := accountv2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	accountAPI := accClient.Accounts()
	myAccount, err := accountAPI.FindByOrg(myorg.GUID, region)
	if err != nil {
		log.Fatal(err)
	}

	validClusters := make(map[string]bool)

	for _, r := range kubernetesRegions {
		target := v1.ClusterTargetHeader{
			OrgID:     myorg.GUID,
			SpaceID:   myspace.GUID,
			AccountID: myAccount.GUID,
			Region:    r,
		}

		clusterClient, err := v1.New(sess)
		if err != nil {
			log.Fatal(err)
		}
		clustersAPI := clusterClient.Clusters()

		out, err := clustersAPI.List(target)
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range out {
			validClusters[c.ID] = true
		}
	}

	softlayerSession := slSession.New(os.Getenv("SL_USERNAME"), os.Getenv("SL_APIKEY"))
	clusterID := doListBlockVolumes(softlayerSession)

	for cluster, storageID := range clusterID {
		if !validClusters[cluster] {
			fmt.Println("[I]", cluster, storageID)
		}
	}
}

func doListBlockVolumes(sess *slSession.Session) map[string][]string {
	// Get the Account service for Block Storage
	service := services.GetAccountService(sess)

	// List Block Storage
	fileStorage, err := service.Limit(500).GetNetworkStorage()

	// Create slice to return
	storageList := make(map[string][]string)

	if err != nil {
		fmt.Printf("Error retrieving File Storage from account: %s\n", err)
	} else {
		var notes map[string]interface{}

		for _, fileStorage := range fileStorage {

			notes = make(map[string]interface{})
			if fileStorage.Notes != nil {
				json.Unmarshal([]byte(*fileStorage.Notes), &notes)
			}

			if _, ok := notes["cluster"]; ok {
				clusterID := notes["cluster"].(string)
				storageList[clusterID] = append(storageList[clusterID], strconv.Itoa(*fileStorage.Id))
			}
		}
	}
	return storageList
}
