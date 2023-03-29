package main

import (
	"log"

	"github.com/netbox-community/go-netbox/netbox"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/models"
	"golang.org/x/net/context"
)

var status = "active"
var pageLimit = int64(100)

func main() {
	//Define the Netbox object
	c := netbox.NewNetboxWithAPIKey("10.10.10.6:8000", "")

	//Get List of All Active Sites
	s, err := listSites(c, status)
	if err != nil {
		log.Fatal(err)
	}

	sites := make(map[string]int64)
	if len(s.Payload.Results) < 1 {
		_, err := createSite(c, "test", "test")
		if err != nil {
			log.Printf("Could not create site test because of error %v", err)
		}
		_, err = createSite(c, "test1", "test1")
		if err != nil {
			log.Printf("Could not create site test because of error %v", err)
		}
		_, err = createSite(c, "test2", "test2")
		if err != nil {
			log.Printf("Could not create site test because of error %v", err)
		}
		_, err = createSite(c, "test3", "test3")
		if err != nil {
			log.Printf("Could not create site test because of error %v", err)
		}
	}
	for _, k := range s.Payload.Results {
		sites[*k.Name] = k.ID
	}

	//Delete all of the Sites

	for k, v := range sites {

		_, err := deleteSite(c, v)
		if err != nil {
			log.Printf("Could not delete site: %v because of error %v", k, err)
		}
	}

	//Get All the sites again
	s, err = listSites(c, status)
	if err != nil {
		log.Fatal(err)
	}

	sites = make(map[string]int64)
	for _, k := range s.Payload.Results {
		sites[*k.Name] = k.ID
	}

	//Build a new Site in Netbox
	//name := "test4"
	//site, err := createSite(c, name, name)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(site.Payload.Display)
}

func createSite(c *client.NetBoxAPI, name, slug string) (*dcim.DcimSitesCreateCreated, error) {
	siteData := models.WritableSite{Name: &name, Slug: &name}
	req := dcim.NewDcimSitesCreateParams()
	req.SetData(&siteData)
	res, err := c.Dcim.DcimSitesCreate(req, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func deleteSite(c *client.NetBoxAPI, id int64) (*dcim.DcimSitesDeleteNoContent, error) {
	ctx := context.TODO()
	siteData := dcim.DcimSitesDeleteParams{ID: id}
	siteData.Context = ctx
	res, err := c.Dcim.DcimSitesDelete(&siteData, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func listSites(c *client.NetBoxAPI, status string) (*dcim.DcimSitesListOK, error) {
	req := dcim.NewDcimSitesListParams().WithStatus(&status)
	res, err := c.Dcim.DcimSitesList(req, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}
