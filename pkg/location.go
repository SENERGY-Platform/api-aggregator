package pkg

import (
	"errors"
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/auth"
)

func (this *Lib) GetDevicesInLocation(token auth.Token, location string) (deviceIds []string, err error) {
	locations, err := this.GetLocations(token, []string{location})
	if err != nil {
		return nil, err
	}
	if len(locations) == 0 {
		return nil, errors.New("unknown location")
	}
	return locations[0].DeviceIds, nil
}

func (this *Lib) GetLocations(token auth.Token, locationIds []string) (locations []Location, err error) {
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: "locations",
		ListIds: &QueryListIds{
			QueryListCommons: QueryListCommons{
				Limit:    len(locationIds),
				Offset:   0,
				Rights:   "r",
				SortBy:   "name",
				SortDesc: false,
			},
			Ids: locationIds,
		},
	}, &locations)
	return
}

type Location struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Image          string   `json:"image"`
	DeviceIds      []string `json:"device_ids"`
	DeviceGroupIds []string `json:"device_group_ids"`
	RdfType        string   `json:"rdf_type"`
}
