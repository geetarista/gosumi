package gosumi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// BaseURL defines the default API endpoint
	BaseURL = "fmipmobile.icloud.com"
)

var (
	defaultBody = []byte(`{
		"clientContext":{
			"appName": "FindMyiPhone",
			"appVersion": "1.4",
			"buildVersion": "145",
			"deviceUDID": "0000000000000000000000000000000000000000",
			"inactiveTime": 2147483647,
			"osVersion": "4.2.1",
			"personID": 0,
			"productType": "iPad1,1"
		}
	}`)
)

// Device represents an iOS device
type Device struct {
	ID               string          `json:"id"`
	BatteryLevel     float32         `json:"batteryLevel"`
	BatteryStatus    string          `json:"batteryStatus"`
	CanWipeAfterLock bool            `json:"canWipeAfterLock"`
	Class            string          `json:"deviceClass"`
	Color            string          `json:"deviceColor"`
	DarkWake         bool            `json:"darkWake"`
	DisplayName      string          `json:"deviceDisplayName"`
	FamilyShare      bool            `json:"fmlyShare"`
	Features         map[string]bool `json:"features"`
	LocationCapable  bool            `json:"locationCapable"`
	LocationEnabled  bool            `json:"locationEnabled"`
	Locating         bool            `json:"isLocating"`
	Location         `json:"location"`
	LockedTimestamp  float64 `json:"lockedTimestamp"`
	LostDevice       string  `json:"lostDevice"`
	LostModeCapable  bool    `json:"lostModeCapable"`
	LostModeEnabled  bool    `json:"lostModeEnabled"`
	LostTimestamp    string  `json:"lostTimestamp"`
	Mac              bool    `json:"isMac"`
	MaxMsgChar       int     `json:"maxMsgChar"`
	// Mesg             string  `json:"mesg"`
	// Msg              string  `json:"msg"`
	Model            string `json:"deviceModel"`
	ModelDisplayName string `json:"modelDisplayName"`
	Name             string `json:"name"`
	PasscodeLength   int    `json:"passcodeLength"`
	PRSID            string `json:"prsId"`
	RawDeviceModel   string `json:"rawDeviceModel"`
	RemoteWipe       string `json:"remoteWipe"`
	// SND              string  `json:"snd"`
	Status         string  `json:"deviceStatus"`
	ThisDevice     bool    `json:"thisDevice"`
	TrackingInfo   string  `json:"trackingInfo"`
	WipeInProgress bool    `json:"wipeInProgress"`
	WipedTimestamp float64 `json:"wipedTimestamp"`
}

// ICloud is an iCloud account
type ICloud struct {
	Email     string
	Password  string
	Devices   []*Device `json:"content"`
	Partition string
}

// A Location represents a specific location for a device
type Location struct {
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Timestamp          int64   `json:"timestamp"`
	LocationType       string  `json:"locationType"`
	HorizontalAccuracy float32 `json:"horizontalAccuracy"`
	PositionType       string  `json:"positionType"`
	Finished           bool    `json:"locationFinished"`
	Inaccurate         bool    `json:"isInaccurate"`
	Old                bool    `json:"isOld"`
}

// New returns a device with the attached username/password
func New(email, password string) (icloud *ICloud, err error) {
	icloud = &ICloud{
		Email:     email,
		Partition: BaseURL,
		Password:  password,
	}

	if err = icloud.fetchPartition(); err != nil {
		return
	}

	err = icloud.updateDevices()

	return
}

func (i *ICloud) fetchPartition() error {
	println("fetchPartition")
	_, err := i.post(fmt.Sprintf("/fmipservice/device/%s/initClient", i.Email), defaultBody)
	if err != nil {
		return err
	}

	return nil
}

func (i *ICloud) post(uri string, data []byte) (body []byte, err error) {
	url := fmt.Sprintf("https://%s%s", i.Partition, uri)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-type", "application/json; charset=utf-8")
	req.Header.Set("X-Apple-Find-Api-Ver", "2.0")
	req.Header.Set("X-Apple-Authscheme", "UserIdGuest")
	req.Header.Set("X-Apple-Realm-Support", "1.2")
	req.Header.Set("User-agent", "Find iPhone/1.2 MeKit (iPad: iPhone OS/4.2.1)")
	req.Header.Set("X-Client-Name", "iPad")
	req.Header.Set("X-Client-UUID", "0cf3dc501ff812adb0b202baed4f37274b210853")
	req.Header.Set("Accept-Language", "en-us")
	req.SetBasicAuth(i.Email, i.Password)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if partition := resp.Header.Get("X-Apple-MMe-Host"); partition != "" {
		i.Partition = partition
	}

	if resp.StatusCode > 200 && resp.StatusCode != 330 {
		log.Println(resp.StatusCode)
		return nil, errors.New("default text from ResponseContentFilter")
	}

	body, err = ioutil.ReadAll(resp.Body)

	println(string(body))

	return
}

func (i *ICloud) updateDevices() (err error) {
	println("updateDevices")
	resp, err := i.post(fmt.Sprintf("/fmipservice/device/%s/initClient", i.Email), defaultBody)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &i)

	return
}
