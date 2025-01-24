package bin

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

func (s *Store) DeleteBin(binID int, userID int) error {
	_, err := s.db.Exec("DELETE FROM bins WHERE id = ? AND owner = ?;", binID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetBinById(binID int) (*types.Bin, error) {
	row := s.db.QueryRow("SELECT * FROM bins WHERE id = ?;", binID)
	bin := new(types.Bin)
	if err := row.Scan(&bin.ID, &bin.Name, &bin.OwnerID, &bin.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bin with id doesn't exist")
		}
		return nil, err
	}
	return bin, nil
}

func (s *Store) UpdateBinName(binID int, userID int, newName string) error {
	_, err := s.db.Exec("UPDATE bins SET name = ? WHERE id = ? AND owner = ?;", newName, binID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetBinsByUser(id int) ([]types.Bin, error) {
	rows, err := s.db.Query("SELECT * FROM bins WHERE owner = ?;", id)
	if err != nil {
		return nil, err
	}

	bins := make([]types.Bin, 0)
	for rows.Next() {
		bin, err := scanRowsIntoBin(rows)
		if err != nil {
			return nil, err
		}
		bins = append(bins, *bin)
	}
	return bins, nil
}

func (s *Store) CreateBin(name string, ownerID int) error {
	_, err := s.db.Exec("INSERT INTO bins(name, owner) VALUES (?,?);", name, ownerID)
	if err != nil {
		return err
	}
	return nil
}

func scanRowsIntoBin(rows *sql.Rows) (*types.Bin, error) {
	bin := new(types.Bin)
	err := rows.Scan(&bin.ID, &bin.Name, &bin.OwnerID, &bin.CreatedAt)
	if err != nil {
		return nil, err
	}
	return bin, nil
}
