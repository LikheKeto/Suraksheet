package types

import (
	"time"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type BinStore interface {
	GetBinsByUser(id int) ([]Bin, error)
	CreateBin(name string, ownerID int) (*Bin, error)
	GetBinById(id int) (*Bin, error)
	UpdateBinName(id int, userID int, name string) error
	DeleteBin(id int, userID int) error
}

type DocumentStore interface {
	InsertDocument(doc Document) (*Document, error)
	GetDocumentByID(id int) (*Document, error)
	UpdateDocumentName(id int, name string) error
	ReferenceNameExistsInBin(name string, binID int) error
	DeleteDocumentByID(id int) error
	GetDocumentsInBin(binID int) ([]Document, error)
	GetDocumentOwner(id int) (int, error)
	FetchDocumentsFromDB(docIDs []int) ([]*Document, error)
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type Bin struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	OwnerID   int       `json:"owner"`
	CreatedAt time.Time `json:"createdAt"`
}

type Document struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	ReferenceName string    `json:"referenceName"`
	BinID         int       `json:"bin"`
	Url           string    `json:"url"`
	Extract       string    `json:"extract"`
	CreatedAt     time.Time `json:"createdAt"`
	Language      string    `json:"language"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=120"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateBinPayload struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type EditBinPayload struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type DeleteBinDocPayload struct {
	Id int `json:"id" validate:"required"`
}

type EditDocumentPayload struct {
	Id            int    `json:"id" validate:"required"`
	ReferenceName string `json:"referenceName" validate:"required"`
}
