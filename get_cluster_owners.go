package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/IBM-Cloud/bluemix-go"

	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"
	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/session"
)

func main() {
	c := new(bluemix.Config)
	flag.StringVar(&c.Region, "region", "", " au-syd, eu-de, eu-gb, us-east, us-south")

	var allRegions bool
	flag.BoolVar(&allRegions, "all", false, "All Regions")

	var org string
	flag.StringVar(&org, "org", "", "Bluemix Organization")

	var space string
	flag.StringVar(&space, "space", "", "Bluemix Space")

	flag.Parse()

	if (org == "" || space == "") || (c.Region == "" && allRegions != true) {
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

	var kubernetesRegions []string
	if allRegions == true {
		kubernetesRegions = []string{"au-syd", "eu-de", "eu-gb", "us-east", "us-south"}
	} else {
		kubernetesRegions = []string{c.Region}
	}
	clusterCount := 0
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
		sort.Slice(out, func(i, j int) bool {
			return out[i].OwnerEmail < out[j].OwnerEmail
		})
		for _, c := range out {
			clusterCount++
			fmt.Println(clusterCount, "[", c.DataCenter, "]", "[", c.OwnerEmail, "]", "[", c.Name, "]")
		}
	}
}