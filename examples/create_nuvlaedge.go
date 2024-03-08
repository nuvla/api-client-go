package main

import (
	apiclientgo "api-client-go"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

type NuvlaBoxMinResourceData struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	RefreshInterval   int      `json:"refresh-interval"`
	HeartbeatInterval int      `json:"heartbeat-interval"`
	VpnServerId       string   `json:"vpn-server-id"`
	Tags              []string `json:"tags"`
}

func NewNuvlaBoxMinResourceData() map[string]interface{} {
	nb := &NuvlaBoxMinResourceData{
		Name:              "[Test] NuvlaEdge Go library",
		Description:       "Test NuvlaBox for Go API Client",
		RefreshInterval:   60,
		HeartbeatInterval: 20,
		VpnServerId:       "infrastructure-service/eb8e09c2-8387-4f6d-86a4-ff5ddf3d07d7",
		Tags:              []string{"test", "go"},
	}
	data := make(map[string]interface{})
	dataB, _ := json.Marshal(nb)
	_ = json.Unmarshal(dataB, &data)
	return data
}

func main() {
	c := apiclientgo.NewUserClient("https://nuvla.io", false, false)
	err := c.Client.LoginApiKeys("credential/<UUID>", "<SECRET>")
	if err != nil {
		fmt.Println(err)
	}

	resId, err := c.AddNuvlaEdge(NewNuvlaBoxMinResourceData())
	if err != nil {
		log.Errorf("Error creating NuvlaBox resource: %s", err)
		os.Exit(1)
	}
	fmt.Printf("NuvlaBox resource ID: %s\n", resId)

	res, err := c.GetNuvlaEdge(resId.String(), nil)
	if err != nil {
		log.Errorf("Error getting NuvlaBox %s resource: %s", resId, err)
		os.Exit(1)
	}
	log.Infof("NuvlaBox resource: %v", res.Data)
}
