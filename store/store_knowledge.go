package store

import "encoding/json"

// KnowledgeEntry is a FAQ/document entry used as AI training context.
type KnowledgeEntry struct {
	ID      int64
	Title   string
	Content string // JSON: {"rows": [{"question":"...","answer":"...","category":"..."}]}
	Active  bool
	Created string
}

// KnowledgeRow is a single Q&A pair extracted from the content JSON.
type KnowledgeRow struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Category string `json:"category"`
}

func (d *DB) migrateKnowledge() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS knowledge (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL DEFAULT '',
			content TEXT NOT NULL,
			is_active TINYINT NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) AddKnowledge(title, content string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO knowledge (title, content) VALUES (?, ?)`, title, content)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
func (d *DB) UpdateKnowledge(id int64, title, content string) error {
	_, err := d.sql.Exec(`UPDATE knowledge SET title=?, content=? WHERE id=?`, title, content, id)
	return err
}
func (d *DB) DeleteKnowledge(id int64) error { return d.del("knowledge", id) }
func (d *DB) ToggleKnowledge(id int64) error {
	_, err := d.sql.Exec(`UPDATE knowledge SET is_active = 1 - is_active WHERE id=?`, id)
	return err
}
func (d *DB) ListKnowledge() ([]KnowledgeEntry, error) {
	rows, err := d.sql.Query(`SELECT id, title, content, is_active, created_at FROM knowledge ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []KnowledgeEntry
	for rows.Next() {
		var k KnowledgeEntry
		var active int
		if err := rows.Scan(&k.ID, &k.Title, &k.Content, &active, &k.Created); err != nil {
			return nil, err
		}
		k.Active = active == 1
		out = append(out, k)
	}
	return out, nil
}

// ActiveKnowledgeRows returns all flattened Q&A rows from active knowledge entries.
func (d *DB) ActiveKnowledgeRows() []KnowledgeRow {
	entries, err := d.sql.Query(`SELECT content FROM knowledge WHERE is_active=1 LIMIT 50`)
	if err != nil {
		return nil
	}
	defer entries.Close()
	var out []KnowledgeRow
	for entries.Next() {
		var raw string
		if err := entries.Scan(&raw); err != nil {
			continue
		}
		var wrapper struct {
			Rows []KnowledgeRow `json:"rows"`
		}
		if err := json.Unmarshal([]byte(raw), &wrapper); err != nil {
			continue
		}
		out = append(out, wrapper.Rows...)
	}
	return out
}
