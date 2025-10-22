package models

type Program struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	Name           string `gorm:"column:name" json:"nama"`
	Deskripsi      string `json:"deskripsi"`
	BentukManfaat  string `json:"bentuk_manfaat"`
	ManfaatLengkap string `json:"manfaat_lengkap"`
}
