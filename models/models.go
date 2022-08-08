package models

type Domain struct {
	Name       string `json:"name" binding:"required"`
	NameServer string `json:"name_server" binding:"required"`
	NSIp       string `json:"nsip" binding:"required"`
	Ttl        string `json:"ttl"`
}

type Subdomain struct {
	Name string `json:"name" binding:"required"`
	Ip   string `json:"ip" binding:"required"`
}

type Record struct {
	Domain string `json:"domain" binding:"required"`
	Ttl    string `json:"ttl"`
	Class  string `json:"class"`
	Type   string `json:"type"`
}

type SOARecord struct {
	Record
	NameServer string `json:"name_server" binding:"required"`
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
	NameServer string `json:"name_server" binding:"required"`
}

type MXRecord struct {
	Record
	Priority    int    `json:"priority" binding:"required"`
	EmailServer string `json:"email_server" binding:"required"`
}

// TODO: Add more types of records
