package environment

type Location struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Image          string   `json:"image"`
	DeviceIds      []string `json:"device_ids"`
	DeviceGroupIds []string `json:"device_group_ids"`
	RdfType        string   `json:"rdf_type"`
}

type Device struct {
	Id           string `json:"id"`
	LocalId      string `json:"local_id"`
	Name         string `json:"name"`
	DeviceTypeId string `json:"device_type_id"`
}

type LocationCommand struct {
	Command  string   `json:"command"`
	Id       string   `json:"id"`
	Owner    string   `json:"owner"`
	Location Location `json:"location"`
}

type DeviceCommand struct {
	Command string `json:"command"`
	Id      string `json:"id"`
	Owner   string `json:"owner"`
	Device  Device `json:"device"`
}
