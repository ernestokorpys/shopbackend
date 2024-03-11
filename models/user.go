package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Usuario representa la estructura de un usuario en la base de datos
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserName   string             `bson:"userName"`
	Email      string             `bson:"email"`
	Password   []byte             `bson:"password"`
	ShopCar    ShopCar            `bson:"shopCar"`
	CreditCard CreditCard         `bson:"creditCard"`
}

// ShopCar representa la estructura de un carrito de compras
type ShopCar struct {
	Products []Product `bson:"products"`
}

// Producto representa la estructura de un producto
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ProductName string             `bson:"productName"`
	Cost        string             `bson:"cost"`
	Picture     string             `bson:"picture"`
}

// CreditCard representa la estructura de una tarjeta de crédito
type CreditCard struct {
	CardNumber     string `bson:"cardNumber"`
	ExpirationDate string `bson:"expirationDate"`
	CVV            string `bson:"cvv"`
}

// func (user *User) SetPassword(password string) {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	if err != nil {
// 		log.Println("Error hashing password:", err)
// 		return
// 	}
// 	user.Password = hashedPassword
// }
