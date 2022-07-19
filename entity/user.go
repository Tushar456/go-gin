package entity

type User struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Address Address `json:"address"`
}

type Address struct {
	State   string `json:"state"`
	City    string `json:"city"`
	Pincode int    `json:"pincode"`
}
