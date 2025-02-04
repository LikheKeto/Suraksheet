package document

import (
	"database/sql"
	"fmt"

	"github.com/LikheKeto/Suraksheet/types"
	"github.com/lib/pq"
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
	row := s.db.QueryRow("SELECT b.owner FROM documents d JOIN bins b ON d.bin = b.id WHERE d.id = $1;", id)
	var ownerID int
	if err := row.Scan(&ownerID); err != nil {
		return 0, err
	}
	return ownerID, nil
}

func (s *Store) GetDocumentsInBin(binID int) ([]types.Document, error) {
	rows, err := s.db.Query("SELECT * FROM documents WHERE bin = $1;", binID)
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
	_, err := s.db.Exec("DELETE FROM documents WHERE id = $1;", id)
	return err
}

func (s *Store) ReferenceNameExistsInBin(name string, binID int) error {
	row := s.db.QueryRow("SELECT id FROM documents WHERE bin = $1 AND referenceName = $2;", binID, name)
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
	row := s.db.QueryRow("SELECT * FROM documents WHERE id = $1;", id)
	doc := new(types.Document)
	if err := row.Scan(&doc.ID, &doc.Name, &doc.ReferenceName,
		&doc.BinID, &doc.Url, &doc.Extract, &doc.CreatedAt, &doc.Language); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Store) InsertDocument(doc types.Document) (*types.Document, error) {
	var newDoc types.Document

	query := `
		INSERT INTO documents (name, referenceName, bin, url, language)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, referenceName, bin, url, extract, createdAt, language;
	`
	err := s.db.QueryRow(query, doc.Name, doc.ReferenceName, doc.BinID, doc.Url, doc.Language).Scan(
		&newDoc.ID, &newDoc.Name, &newDoc.ReferenceName,
		&newDoc.BinID, &newDoc.Url, &newDoc.Extract, &newDoc.CreatedAt, &newDoc.Language,
	)
	if err != nil {
		return nil, err
	}

	return &newDoc, nil
}

func (s *Store) UpdateDocumentName(id int, name string) error {
	_, err := s.db.Exec("UPDATE documents SET referenceName = $1 WHERE id = $2;", name, id)
	if err != nil {
		return fmt.Errorf("unable to update document: %v", err)
	}
	return nil
}

func (s *Store) FetchDocumentsFromDB(docIDs []int) ([]*types.Document, error) {
	var documents []*types.Document

	query := "SELECT * FROM documents WHERE id = ANY($1)"
	rows, err := s.db.Query(query, pq.Array(docIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		doc, err := scanRowsIntoDocument(rows)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func scanRowsIntoDocument(rows *sql.Rows) (*types.Document, error) {
	doc := new(types.Document)
	err := rows.Scan(&doc.ID, &doc.Name, &doc.ReferenceName,
		&doc.BinID, &doc.Url, &doc.Extract, &doc.CreatedAt, &doc.Language)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
