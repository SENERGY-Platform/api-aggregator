package lib

import (
	"errors"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

func (this *Lib) GetDevicesInLocation(jwt jwt_http_router.Jwt, location string) (deviceIds []string, err error) {
	locations, err := this.GetLocations(jwt, []string{location})
	if err != nil {
		return nil, err
	}
	if len(locations) == 0 {
		return nil, errors.New("unknown location")
	}
	return locations[0].DeviceIds, nil
}

func (this *Lib) GetLocations(jwt jwt_http_router.Jwt, locationIds []string) (locations []Location, err error) {
	err, _ = this.QueryPermissionsSearch(string(jwt.Impersonate), QueryMessage{
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
