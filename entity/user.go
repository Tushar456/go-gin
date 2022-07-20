package entity

type User struct {
	UserName string  `json:"username" binding:"required,alphanum"`
	Password string  `json:"password" binding:"required,alphanum"`
	Name     string  `json:"name" binding:"required,alphanum"`
	Age      uint8   `json:"age" binding:"required,gte=1,lte=130"`
	Email    string  `json:"email" binding:"required,email"`
	Phone    string  `json:"phone" binding:"required,e164"`
	DOB      string  `json:"dob" binding:"required,datetime=2006-01-02"`
	Address  Address `json:"address" binding:"required"`
}

type Address struct {
	State   string `json:"state" binding:"required,alphanum"`
	City    string `json:"city" binding:"required,alphanum"`
	Pincode int    `json:"pincode" binding:"required,number"`
}
