package model

type ShortDevice struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DeviceTypeId string `json:"device_type_id"`
}

type ShortDeviceType struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	DeviceClassId string `json:"device_class_id"`
}

type DeviceClass struct {
	Id    string `json:"id"`
	Image string `json:"image"`
	Name  string `json:"name"`
}
