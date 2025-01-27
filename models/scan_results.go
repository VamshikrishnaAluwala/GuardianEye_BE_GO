package models

import "time"

//type ScanResult struct {
//	ID        uint      `gorm:"primarykey" json:"id"`
//	CompanyID uint      `gorm:"index;not null" json:"company_id"`
//	Company   Company   `gorm:"foreignkey:CompanyID" json:"-"`
//	Domain    string    `gorm:"index" json:"domain"`
//	Results   []string  `gorm:"type:json" json:"results"`
//	Params    string    `gorm:"type:json" json:"params"`
//	CreatedAt time.Time `json:"created_at"`
//	Error     string    `json:"error,omitempty"`
//}

type ScanResult struct {
	ID         uint        `gorm:"primarykey" json:"id"`
	CompanyID  uint        `gorm:"index;not null" json:"company_id"`
	Company    Company     `gorm:"foreignkey:CompanyID" json:"-"`
	Domain     string      `gorm:"index" json:"domain"`
	CreatedAt  time.Time   `json:"created_at"`
	Error      string      `json:"error,omitempty"`
	Subdomains []Subdomain `gorm:"foreignKey:ScanResultID" json:"subdomains,omitempty"`
}

type Subdomain struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	ScanResultID uint       `gorm:"index;not null" json:"scan_result_id"`
	ScanResult   ScanResult `gorm:"foreignKey:ScanResultID" json:"-"`
	Subdomain    string     `gorm:"index" json:"subdomain"`
	CreatedAt    time.Time  `json:"created_at"`
}

type ScanParams struct {
	CommandArgs []string `json:"command_args"`
}

type BatchScanRequest struct {
	Domains []string   `json:"domains"`
	Params  ScanParams `json:"params"`
}

type ScanResultFilter struct {
	CompanyID uint   `form:"-"` // Set internally, not from request
	Domain    string `form:"domain"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Limit     int    `form:"limit,default=10"`
	Offset    int    `form:"offset,default=0"`
}
