package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UUIDModel provides a UUID field and a shared BeforeCreate hook.
type UUIDModel struct {
	UUID string `gorm:"type:uuid;uniqueIndex;default null"`
}

func (m *UUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.UUID == "" {
		m.UUID = uuid.New().String()
	}
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		tx.Model(u).Update("role", "admin")
	}
	return
}

// Directories - Directoare
// ********** User - Utilizatorul sistemului **********
type User struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Email     string `gorm:"unique;not null"` // Email-ul utilizatorului (unic)
	Password  string `gorm:"not null"`        // Hash-ul parolei (nu se afișează în JSON)
	Role      string // Rolul utilizatorului ("admin", "user" etc.)
	CanalID   *uint  // Cheie externă către canalul de vânzări
	Canal     *Canal `gorm:"foreignKey:CanalID"` // Canalul de vânzări al utilizatorului
}

// ****************************************************

// ********** Canal - Canal de vânzări **********
type Canal struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Name      string `gorm:"type:varchar(100);not null"` // Numele canalului
}

// ****************************************************

// ********** Client - Client (beneficiar) **********
type ClientType struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Name      string `gorm:"type:varchar(20);not null"` // Tipul clientului ("individual", "company", etc.)
}

// ****************************************************

// ********** Client - Client (beneficiar) **********
type Client struct {
	gorm.Model
	UUIDModel    `gorm:"embedded"`
	ClientTypeID uint       `gorm:"not null"`                          // Foreign key to ClientType
	ClientType   ClientType `gorm:"foreignKey:ClientTypeID;not null"`  // Tipul clientului ("individual", "company", etc.)
	Name         string     `gorm:"type:varchar(100);not null"`        // Numele clientului
	FiscalID     string     `gorm:"type:varchar(15);unique;not null"`  // Codul fiscal al clientului (unic)
	Email        string     `gorm:"type:varchar(100);unique;not null"` // Email-ul clientului (unic)
	Phone        string     `gorm:"type:varchar(50)"`                  // Telefonul clientului
	Address      string     `gorm:"type:text"`                         // Adresa clientului
	Contracts    []Contract `gorm:"foreignKey:ClientID"`               // Contractele clientului
}

// ****************************************************

// ********** Contract - Contract cu clientul **********
type Contract struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Number    string            `gorm:"type:varchar(50);not null;unique"`  // Numărul contractului
	Name      string            `gorm:"type:varchar(100);not null"`        // Numele contractului
	Date      string            `gorm:"type:date;not null"`                // Data contractului
	Amount    float64           `gorm:"type:decimal(10,2);not null"`       // Suma contractului
	Status    string            `gorm:"type:varchar(20);not null"`         // Statutul ("active", "closed" etc.)
	ClientID  uint              `gorm:"not null"`                          // Cheie externă către Client
	Client    Client            `gorm:"foreignKey:ClientID;references:ID"` // Clientul
	OwnerID   uint              `gorm:"not null"`                          // ID-ul ownerului (utilizatorului)
	Owner     User              `gorm:"foreignKey:OwnerID;references:ID"`  // Ownerul contractului
	Addresses []ContractAddress `gorm:"foreignKey:ContractID"`             // Adresele asociate contractului
}

// ****************************************************

// ********** ContractAddress - Adresă asociată contractului **********
type ContractAddress struct {
	gorm.Model
	UUIDModel  `gorm:"embedded"`
	ContractID uint     `gorm:"not null"`                            // Cheie externă către Contract
	Address    string   `gorm:"type:text;not null"`                  // Adresa
	Type       string   `gorm:"type:varchar(50)"`                    // Tipul adresei ("billing", "shipping" etc.)
	Contract   Contract `gorm:"foreignKey:ContractID;references:ID"` // Contractul
	OwnerID    uint     `gorm:"not null"`                            // ID-ul ownerului (utilizatorului)
	Owner      User     `gorm:"foreignKey:OwnerID;references:ID"`    // Ownerul adresei
}

// ****************************************************

// ********** Product - Produs **********
type Product struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string  `gorm:"type:varchar(100);not null"`        // Numele produsului
	Price       float64 `gorm:"type:decimal(10,2);not null"`       // Prețul produsului
	Description string  `gorm:"type:text"`                         // Descrierea produsului
	UnitID      uint    `gorm:"not null"`                          // ID-ul unității de măsură
	Unit        Unit    `gorm:"foreignKey:UnitID;references:ID"`   // Unitatea de măsură a produsului
	VatTaxID    uint    `gorm:"not null"`                          // ID-ul taxei VAT
	VatTax      VatTax  `gorm:"foreignKey:VatTaxID;references:ID"` // Taxa VAT a produsului
}

// ****************************************************

// ********** VatRate - TVA **********
type VatTax struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string  `gorm:"type:varchar(100);not null"`  // Numele taxei
	Rate        float32 `gorm:"type:decimal(10,2);not null"` // Rata taxei
	Description string  `gorm:"type:text"`                   // Descrierea taxei
}

// ****************************************************

// ********** Product - Produs **********
type IncomeTax struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string  `gorm:"type:varchar(100);not null"`  // Numele taxei
	Rate        float32 `gorm:"type:decimal(10,2);not null"` // Rata taxei
	Description string  `gorm:"type:text"`                   // Descrierea taxei
}

// ****************************************************

// ********** Units of Measurement **********
type Unit struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Name      string `gorm:"type:varchar(50);not null"` // Numele unității de măsură (ex: "buc", "kg")
}

// Documents - Documente
// ********** Order - Comandă **********
type Order struct {
	gorm.Model
	UUIDModel  `gorm:"embedded"`
	OwnerID    uint        `gorm:"not null"`                            // ID-ul ownerului (utilizatorului)
	Owner      User        `gorm:"foreignKey:OwnerID;references:ID"`    // Ownerul comenzii
	ClientID   uint        `gorm:"not null"`                            // ID-ul clientului (cheie externă)
	Client     Client      `gorm:"foreignKey:ClientID;references:ID"`   // Clientul care a plasat comanda
	ContractID uint        `gorm:"not null"`                            // ID-ul contractului (cheie externă)
	Contract   Contract    `gorm:"foreignKey:ContractID;references:ID"` // Contractul asociat comenzii
	TotalPrice float64     `gorm:"type:decimal(10,2);not null"`         // Suma totală a comenzii
	Status     string      `gorm:"type:varchar(20);not null"`           // Statusul comenzii
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`                  // Pozițiile comenzii
}

// ****************************************************

// ********** OrderItem - Poziție comandă **********
type OrderItem struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	OrderID     uint    `gorm:"not null"`                           // ID-ul comenzii
	ProductID   uint    `gorm:"not null"`                           // ID-ul produsului
	Product     Product `gorm:"foreignKey:ProductID;references:ID"` // Produsul asociat poziției
	Quantity    int     `gorm:"not null"`                           // Cantitatea
	Price       float64 `gorm:"type:decimal(10,2);not null"`        // Prețul unitar la momentul comenzii
	UnitID      uint    `gorm:"not null"`                           // ID-ul unității de măsură
	Unit        Unit    `gorm:"foreignKey:UnitID;references:ID"`    // Unitatea de măsură asociată poziției
	VatTaxID    uint    `gorm:"not null"`                           // ID-ul taxei VAT
	VatTax      VatTax  `gorm:"foreignKey:VatTaxID;references:ID"`  // Taxa VAT asociată poziției
	Summ        float64 `gorm:"type:decimal(10,2);not null"`        // Suma totală pentru poziție (Price * Quantity)
	SummWithVat float64 `gorm:"type:decimal(10,2);not null"`        // Suma totală pentru poziție cu TVA (Summ + VAT)
}

// ****************************************************
