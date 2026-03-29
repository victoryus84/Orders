package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UUIDModel provides a UUID field and a shared BeforeCreate hook.
type UUIDModel struct {
	UUID string `gorm:"type:uuid;uniqueIndex;default null"`
}

// Directories - Directoare

// ********** User - Utilizatorul sistemului **********
type User struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Email     string    `gorm:"unique;not null"`           // Email-ul utilizatorului (unic)
	Password  string    `gorm:"not null"`                  // Hash-ul parolei (nu se afișează în JSON)
	Role      string    `gorm:"type:varchar(20);not null"` // Rolul utilizatorului ("admin", "user" etc.)
	Channels  []Channel `gorm:"many2many:user_channels;"`  // Canalele de vânzări la care are acces utilizatorul
}

// ****************************************************

// ********** Channel - Canal de vânzări **********
type Channel struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string `gorm:"type:varchar(100);not null"` // Numele canalului
	Description string `gorm:"type:text"`                  // Descrierea canalului
	Users       []User `gorm:"many2many:user_channels;"`   // Utilizatorii care aparțin acestui canal
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
	ClientTypeID  uint       `gorm:"not null"`                          // Foreign key to ClientType
	ClientType    ClientType `gorm:"foreignKey:ClientTypeID;not null"`  // Tipul clientului ("individual", "company", etc.)
	Name          string     `gorm:"type:varchar(100);not null"`        // Numele clientului
	FiscalID      string     `gorm:"type:varchar(15);unique;not null"`  // Codul fiscal al clientului (unic)
	Email         string     `gorm:"type:varchar(100);unique;not null"` // Email-ul clientului (unic)
	Phone         string     `gorm:"type:varchar(50)"`                  // Telefonul clientului
	Address 	  string     `gorm:"type:text"`                         // Adresa clientului
	PostalAddress string     `gorm:"type:text"`                   		// Adresa postala a clientului
	Contracts     []Contract `gorm:"foreignKey:ClientID"`               // Contractele clientului
	Addresses     []ClientAddress `gorm:"foreignKey:ClientID"`          // Adresele asociate clientului
}

// ****************************************************

// ********** ClientAddress - Adresă asociată clientului **********
type ClientAddress struct {
	gorm.Model
	UUIDModel  `gorm:"embedded"`
	Address    string   `gorm:"type:text;not null"`                  // Adresa
	Type       string   `gorm:"type:varchar(50)"`                    // Tipul adresei ("billing", "shipping" etc.)
	ClientID uint     `gorm:"not null"`                              // Cheie externă către Client
	Client   Client `gorm:"foreignKey:ClientID;references:ID"`       // Clientul
	OwnerID    uint     `gorm:"not null"`                            // ID-ul ownerului (utilizatorului)
	Owner      User     `gorm:"foreignKey:OwnerID;references:ID"`    // Ownerul adresei
}

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
	Address    string   `gorm:"type:text;not null"`                  // Adresa
	Type       string   `gorm:"type:varchar(50)"`                    // Tipul adresei ("billing", "shipping" etc.)
	ContractID uint     `gorm:"not null"`                            // Cheie externă către Contract
	Contract   Contract `gorm:"foreignKey:ContractID;references:ID"` // Contractul
	OwnerID    uint     `gorm:"not null"`                            // ID-ul ownerului (utilizatorului)
	Owner      User     `gorm:"foreignKey:OwnerID;references:ID"`    // Ownerul adresei
}

// ****************************************************

// ********** ProductGroup - Grupa de Produse **********
type ProductGroup struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string    `gorm:"type:varchar(100);not null;unique"` // Numele grupei (ex: "Băuturi", "Electronice")
	Description string    `gorm:"type:text"`                         // Descrierea grupei
	Products    []Product `gorm:"foreignKey:ProductGroupID"`         // O grupă are mai multe produse
}

// ****************************************************

// ********** Product - Produs **********
type Product struct {
	gorm.Model
	UUIDModel      `gorm:"embedded"`
	Name           string       `gorm:"type:varchar(100);not null"`              // Numele produsului
	Price          float64      `gorm:"type:decimal(10,2);default:0.0"`          // Prețul produsului
	Description    string       `gorm:"type:text"`                               // Descrierea produsului
	ProductGroupID uint         `gorm:"not null"`                                // ID-ul grupei de produse
	ProductGroup   ProductGroup `gorm:"foreignKey:ProductGroupID;references:ID"` // Grupa de produse din care face parte
	UnitID         uint         `gorm:"not null"`                                // ID-ul unității de măsură
	Unit           Unit         `gorm:"foreignKey:UnitID;references:ID"`         // Unitatea de măsură a produsului
	VatTaxID       uint         `gorm:"not null"`                                // ID-ul taxei VAT
	VatTax         VatTax       `gorm:"foreignKey:VatTaxID;references:ID"`       // Taxa VAT a produsului
}

// ****************************************************

// ********** VatRate - TVA **********
type VatTax struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string  `gorm:"type:varchar(100);not null"`  // Numele taxei
	Rate        float64 `gorm:"type:decimal(10,2);not null"` // Rata taxei
	Description string  `gorm:"type:text"`                   // Descrierea taxei
}

// ****************************************************

// ********** IncomeTax - Taxă pe venit **********
type IncomeTax struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string  `gorm:"type:varchar(100);not null"`  // Numele taxei
	Rate        float64 `gorm:"type:decimal(10,2);not null"` // Rata taxei
	Description string  `gorm:"type:text"`                   // Descrierea taxei
}

// ****************************************************

// ********** Units of Measurement - Unitati de măsură **********
type Unit struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string  `gorm:"type:varchar(50);not null"`               // Numele unității de măsură (ex: "buc", "kg")
	Description string  `gorm:"type:text"`                               // Descrierea unității de măsură
	Coefficient float64 `gorm:"type:decimal(10,4);not null;default:1.0"` // Coeficient de conversie față de unitatea de bază

}

// ****************************************************

// ********** Price type of products - Tipuri de pret **********
type PriceType struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string `gorm:"type:varchar(50);not null"` // Numele tipului de preț (ex: "Cu amamuntul", "En-gros")
	Description string `gorm:"type:text"`                 // Descrierea tipiului de preț
}

// ****************************************************

// ********** Price of products - Preturi producte **********
type PriceProduct struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	ProductID   uint      `gorm:"not null"`                             // Cheie externă către Product
	Product     Product   `gorm:"foreignKey:ProductID;references:ID"`   // Produsul
	PriceTypeID uint      `gorm:"not null"`                             // Cheie externă către PriceType
	PriceType   PriceType `gorm:"foreignKey:PriceTypeID;references:ID"` // Tipul de preț
	Price       float64   `gorm:"type:decimal(10,2);not null"`          // Prețul pentru acest tip
}

// ****************************************************

// Documents - Documente
// ********** Order - Comandă **********
type Order struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	OwnerID     uint        `gorm:"not null"`                            // ID-ul ownerului (utilizatorului)
	Owner       User        `gorm:"foreignKey:OwnerID;references:ID"`    // Ownerul comenzii
	ClientID    uint        `gorm:"not null"`                            // ID-ul clientului (cheie externă)
	Client      Client      `gorm:"foreignKey:ClientID;references:ID"`   // Clientul care a plasat comanda
	PriceTypeID uint        `gorm:"not null"`                            // ID-ul tipului de preț (cheie externă)
	PriceType   PriceType   `gorm:"foreignKey:PriceTypeID"`              // Tipul de preț al comenzii
	ContractID  uint        `gorm:"not null"`                            // ID-ul contractului (cheie externă)
	Contract    Contract    `gorm:"foreignKey:ContractID;references:ID"` // Contractul asociat comenzii
	TotalPrice  float64     `gorm:"type:decimal(10,2);not null"`         // Suma totală a comenzii
	Status      string      `gorm:"type:varchar(20);not null"`           // Statusul comenzii
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`                  // Pozițiile comenzii
}

// ****************************************************

// ********** OrderItem - Poziție comandă **********
type OrderItem struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	OrderID     uint    `gorm:"not null"`                           // ID-ul comenzii
	ProductID   uint    `gorm:"not null"`                           // ID-ul produsului
	Product     Product `gorm:"foreignKey:ProductID;references:ID"` // Produsul asociat poziției
	Quantity    float64 `gorm:"type:decimal(10,3);not null"`        // Cantitatea
	Price       float64 `gorm:"type:decimal(10,2);not null"`        // Prețul unitar la momentul comenzii
	UnitID      uint    `gorm:"not null"`                           // ID-ul unității de măsură
	Unit        Unit    `gorm:"foreignKey:UnitID;references:ID"`    // Unitatea de măsură asociată poziției
	UnitName    string  `gorm:"type:varchar(20)"`                   // Stocăm "KG" sau "BUC"
	VatTaxID    uint    `gorm:"not null"`                           // ID-ul taxei VAT
	VatTax      VatTax  `gorm:"foreignKey:VatTaxID;references:ID"`  // Taxa VAT asociată poziției
	VatRate     float64 `gorm:"type:decimal(10,2);not null"`        // Rata TVA-ului (preluată din VatTax)
	Summ        float64 `gorm:"type:decimal(10,2);not null"`        // Suma totală pentru poziție (Price * Quantity)
	VatSumm     float64 `gorm:"type:decimal(10,2);not null"`        // Valoarea TVA-ului în bani
	SummWithVat float64 `gorm:"type:decimal(10,2);not null"`        // Suma totală pentru poziție cu TVA (Summ + VAT)
}

// ****************************************************

// Hooks - Hook-uri GORM
// BeforeCreate hook pentru UUIDModel - generează un UUID dacă nu este deja setat

func (m *UUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.UUID == "" {
		m.UUID = uuid.New().String()
	}
	return
}

// AfterCreate hook pentru User - dacă este primul utilizator creat, îi setează rolul de admin
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		tx.Model(u).Update("role", "admin")
	}
	return
}
