package models

// DataFormatVersion
// This denotes the data format structure of the models.
// We must need to update the version if there is any change in the structures in this file
const DataFormatVersion = "1.1.0"
type Status string
const (
	TripStart    Status = "New"
	TripRunning  Status = "In progress"
	TripFinished Status = "Completed"
)

type Trip struct {
	Id            string
	Status        Status
	LastTime      interface{}
	LastPosition  *Position
	StartTime     interface{}
	StartPosition *Position
	EndTime       interface{}
	EndPosition   *Position
	AvgSpeed      float64
	DataCount     int
	Distance      float64
	PrevPosition *Position
}
type CellInfo struct {
	BaseStationId  interface{} `json:"base_station_id,omitempty"`
	CId            interface{} `json:"cid,omitempty"`
	LAC            interface{} `json:"lac,omitempty"`
	MCC            interface{} `json:"mcc,omitempty"`
	MNC            interface{} `json:"mnc,omitempty"`
	SignalStrength interface{} `json:"signal_strength,omitempty"`
	UseLCellId     interface{} `json:"use_long_cid,omitempty"`
}

type GPS struct {
	SatellitesCount interface{} `json:"no_of_satellites,omitempty"`
	HDOP            interface{} `json:"hdop,omitempty"`
	Position        *Position   `json:"position,omitempty"`
}

type Position struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Altitude  float64 `json:"alt"`
	IsValid   bool    `json:"is_valid"`
}

type Velocity struct {
	Speed     float64 `json:"speed"`
	Direction float64 `json:"direction"`
}

type Voltage struct {
	External interface{} `json:"external,omitempty"`
	Internal interface{} `json:"internal,omitempty"`
}

type Fuel struct {
	Consumption interface{} `json:"consumption,omitempty"`
	Percentage  interface{} `json:"percentage,omitempty"`
}

type TrackingData struct {
	FormatVersion     string                  `default:"1.1.0" json:"format_version"`
	SourceId          string                  `json:"source_id"`
	PacketType        interface{}             `json:"packet_type"`
	ServerTime        interface{}             `json:"server_time"`
	SourceTime        interface{}             `json:"source_time,omitempty"`
	GSMSignalStrength interface{}             `json:"gsm_signal_strength,omitempty"`
	RFID              interface{}             `json:"rfid_id,omitempty"`
	PacketID          interface{}             `json:"packet_id"`
	Ignition          interface{}             `json:"ignition,omitempty"`
	Mileage           interface{}             `json:"mileage,omitempty"`
	GPS               *GPS                    `json:"gps,omitempty"`
	Velocity          *Velocity               `json:"velocity,omitempty"`
	CellInfos         *[]CellInfo             `json:"cell_info,omitempty"`
	Voltage           *Voltage                `json:"voltage,omitempty"`
	Fuel              *Fuel                   `json:"fuel,omitempty"`
	OtherData         *map[string]interface{} `json:"other_data,omitempty"`
	OBD               *map[string]interface{} `json:"obd,omitempty"`
	Events            *map[string]interface{} `json:"events,omitempty"`
	//TODO: Remove below elements later
	Imei         string      `json:"imei"`
	PositionTime interface{} `json:"position_time,omitempty"`
}

type DBType int

const (
	MongoDB DBType = iota
	Postgres
)

const (
	LatestFormattedDataTableName string = "latest_formatted_data"
	TrackingDataTableName        string = "tracking_data"
	RawDataTableName             string = "raw_data"
)

type RawDataResult struct {
	SourceId   string `json:"source_id,omitempty"`
	ServerTime int64  `json:"server_time"`
	Data       string `json:"data"`
	PacketId   string `json:"packet_id"`
}

const (
	AdapterQueuesIdentifier       string = "ADAPTER"
	StoreQueuesIdentifier         string = "STORE"
	BackendRouterQueuesIdentifier string = "BACKEND_ROUTER"
	PublisherQueuesIdentifier     string = "PUBLISHER"
)

type Event struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	SourceId   string                 `json:"source_id"`
	ServerTime interface{}            `json:"server_time"`
	SourceTime interface{}            `json:"source_time"`
	Data       map[string]interface{} `json:"data"`
}
