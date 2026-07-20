package store

type ChatFlow struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	Trigger   string `json:"trigger"`   // JSON: {type, value}
	NodesJSON string `json:"-"`          // raw JSON
	EdgesJSON string `json:"-"`
	Active    int    `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (d *DB) migrateFlow() error {
	_, err := d.sql.Exec(`CREATE TABLE IF NOT EXISTS chat_flows (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		user_id BIGINT NOT NULL DEFAULT 0,
		name VARCHAR(255) NOT NULL DEFAULT 'New Flow',
		` + "`trigger`" + ` JSON NOT NULL,
		nodes JSON NOT NULL,
		edges JSON NOT NULL,
		active TINYINT NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_flow_user (user_id),
		INDEX idx_flow_active (active)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`)
	return err
}

func (d *DB) SaveFlow(uid int64, name, trigger, nodes, edges string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO chat_flows (user_id, name, `+"`trigger`"+`, nodes, edges, active) VALUES (?,?,?,?,?,1) ON DUPLICATE KEY UPDATE name=VALUES(name), `+"`trigger`"+`=VALUES(`+"`trigger`"+`), nodes=VALUES(nodes), edges=VALUES(edges)`,
		uid, name, trigger, nodes, edges)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (d *DB) UpdateFlow(id, uid int64, name, trigger, nodes, edges string) error {
	_, err := d.sql.Exec(`UPDATE chat_flows SET name=?, `+"`trigger`"+`=?, nodes=?, edges=?, updated_at=NOW() WHERE id=? AND user_id=?`,
		name, trigger, nodes, edges, id, uid)
	return err
}

func (d *DB) DeleteFlow(id, uid int64) error {
	_, err := d.sql.Exec(`DELETE FROM chat_flows WHERE id=? AND user_id=?`, id, uid)
	return err
}

func (d *DB) ToggleFlow(id, uid int64) (bool, error) {
	_, err := d.sql.Exec(`UPDATE chat_flows SET active=1-active WHERE id=? AND user_id=?`, id, uid)
	if err != nil {
		return false, err
	}
	var active int
	d.sql.QueryRow(`SELECT active FROM chat_flows WHERE id=?`, id).Scan(&active)
	return active == 1, nil
}

func (d *DB) ListFlows(uid int64) ([]ChatFlow, error) {
	rows, err := d.sql.Query("SELECT id, user_id, name, `trigger`, nodes, edges, active, created_at, updated_at FROM chat_flows WHERE user_id=? ORDER BY updated_at DESC", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ChatFlow
	for rows.Next() {
		var f ChatFlow
		var trig []byte
		rows.Scan(&f.ID, &f.UserID, &f.Name, &trig, &f.NodesJSON, &f.EdgesJSON, &f.Active, &f.CreatedAt, &f.UpdatedAt)
		f.Trigger = string(trig)
		out = append(out, f)
	}
	return out, nil
}

func (d *DB) GetFlow(id, uid int64) (*ChatFlow, error) {
	var f ChatFlow
	var trig []byte
	err := d.sql.QueryRow("SELECT id, user_id, name, `trigger`, nodes, edges, active, created_at, updated_at FROM chat_flows WHERE id=? AND user_id=?", id, uid).Scan(&f.ID, &f.UserID, &f.Name, &trig, &f.NodesJSON, &f.EdgesJSON, &f.Active, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	f.Trigger = string(trig)
	return &f, nil
}

// Raw helpers for pro package
func (d *DB) SaveFlowRaw(uid int64, name, trigger, nodes, edges string) (int64, error) {
	return d.SaveFlow(uid, name, trigger, nodes, edges)
}
func (d *DB) UpdateFlowRaw(id, uid int64, name, trigger, nodes, edges string) error {
	return d.UpdateFlow(id, uid, name, trigger, nodes, edges)
}
func (d *DB) DeleteFlowRaw(id, uid int64) error { return d.DeleteFlow(id, uid) }
func (d *DB) ToggleFlowRaw(id, uid int64) (bool, error) { return d.ToggleFlow(id, uid) }
func (d *DB) ListFlowsRaw(uid int64) ([]ChatFlow, error) { return d.ListFlows(uid) }
func (d *DB) GetFlowRaw(id, uid int64) (*ChatFlow, error) { return d.GetFlow(id, uid) }
func (d *DB) DuplicateFlowRaw(id, uid int64) (int64, error) {
	f, err := d.GetFlow(id, uid)
	if err != nil {
		return 0, err
	}
	f.Name = f.Name + " (Copy)"
	return d.SaveFlow(uid, f.Name, f.Trigger, f.NodesJSON, f.EdgesJSON)
}

// For pro flow engine: load all active flows for a user
func (d *DB) LoadActiveFlows(uid int64) ([]ChatFlow, error) {
	rows, err := d.sql.Query("SELECT id, user_id, name, `trigger`, nodes, edges, active, created_at, updated_at FROM chat_flows WHERE user_id=? AND active=1 ORDER BY updated_at DESC", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ChatFlow
	for rows.Next() {
		var f ChatFlow
		var trig []byte
		rows.Scan(&f.ID, &f.UserID, &f.Name, &trig, &f.NodesJSON, &f.EdgesJSON, &f.Active, &f.CreatedAt, &f.UpdatedAt)
		f.Trigger = string(trig)
		out = append(out, f)
	}
	return out, nil
}
