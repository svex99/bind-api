package models

type Record struct {
	Domain string `json:"domain" binding:"required"`
	Ttl    string `json:"ttl"`
	Class  string `json:"class"`
	Type   string `json:"type"`
}

type SOARecord struct {
	Record
	NameServer string `json:"nameServer" binding:"required"`
	Admin      string `json:"admin"`
	Serial     uint   `json:"serial"`
	Refresh    string `json:"refresh"`
	Retry      string `json:"retry"`
	Expire     string `json:"expire"`
	Minimum    string `json:"minimum"`
}

type ARecord struct {
	Record
	Ip string `json:"ip" binding:"required"`
}

type AAAARecord struct {
	Record
	Ip string `json:"ip" binding:"required"`
}

type NSRecord struct {
	Record
	NameServer string `json:"nameServer" binding:"required"`
}

type MXRecord struct {
	Record
	Priority    int    `json:"priority" binding:"required"`
	EmailServer string `json:"emailServer" binding:"required"`
}

// TODO: Add more types of records
