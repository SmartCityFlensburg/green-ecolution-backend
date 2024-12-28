package entities

type Watermark struct {
	Resistance int `json:"resistance"`
	Centibar   int `json:"centibar"`
	Depth      int `json:"depth"`
}

type MqttPayload struct {
	Device      string      `json:"device"`
	Battery     float64     `json:"battery"`
	Humidity    float64     `json:"humidity"`
	Temperature float64     `json:"temperature"`
	Latitude    float64     `json:"latitude"`
	Longitude   float64     `json:"longitude"`
	Watermarks  []Watermark `json:"watermarks"`
}
