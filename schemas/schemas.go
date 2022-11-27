package schemas

type ZoneData struct {
	Origin     string `json:"origin"`
	Ttl        string `json:"ttl"`
	NameServer string `json:"nameServer"`
	Admin      string `json:"admin"`
	Refresh    uint   `json:"refresh" binding:"gt=0"`
	Retry      uint   `json:"retry" binding:"gt=0"`
	Expire     uint   `json:"expire" binding:"gt=0"`
	Minimum    uint   `json:"minimum" binding:"gt=0"`
}
