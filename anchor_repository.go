package junkboy

import "database/sql"

type AnchorSQLiteRepository struct {
	db *sql.DB
}

func NewAnchorSQLiteRepository(db *sql.DB) *AnchorSQLiteRepository {
	return &AnchorSQLiteRepository{
		db: db,
	}
}

func (r *AnchorSQLiteRepository) AddAnchor(a Anchor) (int, error) {
	stmt, err := r.db.Prepare("INSERT INTO anchors (url) VALUES (?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(a.URL)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *AnchorSQLiteRepository) UpdateAnchor(a Anchor) error {
	stmt, err := r.db.Prepare("UPDATE anchors SET url=? WHERE id=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(a.URL, a.ID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *AnchorSQLiteRepository) GetAnchor(id int) (Anchor, error) {
	row := r.db.QueryRow("SELECT id, url FROM anchors WHERE id=?", id)

	anchor := Anchor{}
	err := row.Scan(&anchor.ID, &anchor.URL)

	if err != nil {
		return anchor, err
	}

	return anchor, err
}

func (r *AnchorSQLiteRepository) GetAnchors() ([]Anchor, error) {
	rows, err := r.db.Query("SELECT * FROM anchors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	anchors := []Anchor{}

	for rows.Next() {
		anchor := Anchor{}
		err = rows.Scan(&anchor.ID, &anchor.URL)

		if err != nil {
			return nil, err
		}

		anchors = append(anchors, anchor)
	}

	return anchors, nil
}

func (r *AnchorSQLiteRepository) DeleteAnchor(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM anchors WHERE id=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
