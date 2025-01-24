package document

import (
	"database/sql"
	"fmt"

	"github.com/LikheKeto/Suraksheet/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetDocumentOwner(id int) (int, error) {
	row := s.db.QueryRow("SELECT b.owner FROM documents d JOIN bins b ON d.bin = b.id WHERE d.id =?;", id)
	var ownerID int
	if err := row.Scan(&ownerID); err != nil {
		return 0, err
	}
	return ownerID, nil

}

func (s *Store) GetDocumentsInBin(binID int) ([]types.Document, error) {
	rows, err := s.db.Query("SELECT * FROM documents WHERE bin = ?;", binID)
	if err != nil {
		return nil, err
	}
	docs := make([]types.Document, 0)
	for rows.Next() {
		doc, err := scanRowsIntoDocument(rows)
		if err != nil {
			return nil, err
		}
		docs = append(docs, *doc)
	}
	return docs, nil
}

func (s *Store) DeleteDocumentByID(id int) error {
	_, err := s.db.Exec("DELETE FROM documents WHERE id = ?;", id)
	return err
}

func (s *Store) ReferenceNameExistsInBin(name string, binID int) error {
	row := s.db.QueryRow("SELECT id FROM documents WHERE bin = ? AND referenceName = ?;", binID, name)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return fmt.Errorf("document with reference name already exists")
}

func (s *Store) GetDocumentByID(id int) (*types.Document, error) {
	row := s.db.QueryRow("SELECT * FROM documents WHERE id = ?;", id)
	doc := new(types.Document)
	if err := row.Scan(&doc.ID, &doc.Name, &doc.ReferenceName,
		&doc.BinID, &doc.Url, &doc.Extract, &doc.CreatedAt); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Store) InsertDocument(doc types.Document) (int64, error) {
	query := "INSERT INTO documents(name, referenceName, bin, url) VALUES (?,?,?,?)"
	res, err := s.db.Exec(query, doc.Name, doc.ReferenceName, doc.BinID, doc.Url)
	if err != nil {
		return 0, err
	}
	docId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return docId, nil
}

func (s *Store) UpdateDocumentName(id int, name string) error {
	_, err := s.db.Exec("UPDATE documents SET referenceName = ? WHERE id = ?;", name, id)
	if err != nil {
		return fmt.Errorf("unable to update document: %v", err)
	}
	return nil
}

func scanRowsIntoDocument(rows *sql.Rows) (*types.Document, error) {
	doc := new(types.Document)
	err := rows.Scan(&doc.ID, &doc.Name, &doc.ReferenceName,
		&doc.BinID, &doc.Url, &doc.Extract, &doc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
