package models

import "gorm.io/gorm"

// User - Пользователь системы
type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"` // Email пользователя (уникальный)
	Password string `json:"-"`                   // Хэш пароля (не выводится в JSON)
	Role     string `json:"role"`                // Роль пользователя ("admin", "user" и т.д.)
	CanalID   *uint   `json:"canal_id"` 			 // внешний ключ на канал продаж
    Canal     *Canal  `json:"canal" gorm:"foreignKey:CanalID"` // Канал продаж пользователя
}

// Canal - Канал продаж
type Canal struct {
	gorm.Model
	Name      string     `gorm:"type:varchar(100);not null"`        // Имя клиента
}

// Client - Клиент (заказчик)
type Client struct {
	gorm.Model
	Name      string     `gorm:"type:varchar(100);not null"`        // Имя клиента
	Email     string     `gorm:"type:varchar(100);unique;not null"` // Email клиента (уникальный)
	Phone     string     `gorm:"type:varchar(20)"`                  // Телефон клиента
	Address   string     `gorm:"type:text"`                         // Адрес клиента
	UserID    uint       `gorm:"not null"`                          // ID пользователя-владельца
	Contracts []Contract `gorm:"foreignKey:ClientID"`               // Договоры клиента
}

// Contract - Договор с клиентом
type Contract struct {
	gorm.Model
	ClientID  uint              `gorm:"not null"`                          // Внешний ключ к Client
	Number    string            `gorm:"type:varchar(50);not null;unique"`  // Номер договора
	Date      string            `gorm:"type:date;not null"`                // Дата договора
	Amount    float64           `gorm:"type:decimal(10,2);not null"`       // Сумма договора
	Status    string            `gorm:"type:varchar(20);not null"`         // Статус ("active", "closed" и т.д.)
	Client    Client            `gorm:"foreignKey:ClientID;references:ID"` // Клиент
	Addresses []ContractAddress `gorm:"foreignKey:ContractID"`             // Адреса, связанные с договором
}

// ContractAddress - Адрес, связанный с договором
type ContractAddress struct {
	gorm.Model
	ContractID uint     `gorm:"not null"`                            // Внешний ключ к Contract
	Address    string   `gorm:"type:text;not null"`                  // Адрес
	Type       string   `gorm:"type:varchar(50)"`                    // Тип адреса ("billing", "shipping" и т.д.)
	Contract   Contract `gorm:"foreignKey:ContractID;references:ID"` // Договор
}

// Product - Продукт/товар
type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"type:varchar(100);not null"`   // Название продукта
	Price       float64 `json:"price" gorm:"type:decimal(10,2);not null"` // Цена продукта
	Description string  `json:"description" gorm:"type:text"`             // Описание продукта
}

// Order - Заказ
type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id" gorm:"not null"`                             // ID пользователя, оформившего заказ
	ClientID   uint        `json:"client_id" gorm:"not null"`                           // ID клиента (внешний ключ)
	Client     Client      `json:"client" gorm:"foreignKey:ClientID;references:ID"`     // Клиент, оформивший заказ
	ContractID uint        `json:"contract_id" gorm:"not null"`                         // ID договора (внешний ключ)
	Contract   Contract    `json:"contract" gorm:"foreignKey:ContractID;references:ID"` // Договор, по которому оформлен заказ
	TotalPrice float64     `json:"total_price" gorm:"type:decimal(10,2);not null"`      // Общая сумма заказа
	Status     string      `json:"status" gorm:"type:varchar(20);not null"`             // Статус заказа
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`               // Позиции заказа
}

// OrderItem - Позиция заказа
type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id" gorm:"not null"`                 // ID заказа
	ProductID uint    `json:"product_id" gorm:"not null"`               // ID продукта
	Quantity  int     `json:"quantity" gorm:"not null"`                 // Количество
	Price     float64 `json:"price" gorm:"type:decimal(10,2);not null"` // Цена за единицу на момент заказа
}

