package entity

type User struct {
	Name    string  `json:"name" binding:"required"`
	Age     uint8   `json:"age" binding:"required,gte=1,lte=130"`
	Email   string  `json:"email" binding:"required,email"`
	Address Address `json:"address" binding:"required"`
}

type Address struct {
	State   string `json:"state" binding:"required"`
	City    string `json:"city" binding:"required"`
	Pincode int    `json:"pincode" binding:"required"`
}
