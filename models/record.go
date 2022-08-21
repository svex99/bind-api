package models

type Record struct {
	Id         uint64 `json:"id" gorm:"primaryKey"`
	Ttl        string `json:"ttl" binding:"min=2"`
	Class      string `json:"class" binding:"min=1"`
	DomainName string `json:"domainName" binding:"min=1"`
	Domain     Domain `gorm:"references:Name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type SOARecord struct {
	Record
	NameServer string `json:"nameServer" binding:"min=1"`
	Admin      string `json:"admin" binding:"min=1"`
	Serial     uint   `json:"serial"`
	Refresh    uint   `json:"refresh" binding:"gt=0"`
	Retry      uint   `json:"retry" binding:"gt=0"`
	Expire     uint   `json:"expire" binding:"gt=0"`
	Minimum    uint   `json:"minimum" binding:"gt=0"`
}

type NSRecord struct {
	Record
	NameServer string `json:"nameServer" binding:"min=1"`
}

type ARecord struct {
	Record
	Name string `json:"name" binding:"min=1"`
	Ip   string `json:"ip" binding:"ipv4"`
}

// type AAAARecord struct {
// 	Record
// 	Name string `json:"name" binding:"min=1"`
// 	Ip   string `json:"ip" binding:"ipv6"`
// }

type MXRecord struct {
	Record
	Priority    uint   `json:"priority" binding:"gt=0"`
	EmailServer string `json:"emailServer" binding:"min=1"`
}

// TODO: Add more types of records
