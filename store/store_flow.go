package store

import (
	"database/sql"
	"fmt"
)

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
		public TINYINT DEFAULT 0,
		downloads INT DEFAULT 0,
		ai_key_id BIGINT NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_flow_user (user_id),
		INDEX idx_flow_active (active)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`)
	if err != nil {
		return err
	}
	d.safeAddColumn("chat_flows", "ai_key_id", "BIGINT NOT NULL DEFAULT 0")
	return nil
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
	res, err := d.sql.Exec(`UPDATE chat_flows SET name=?, `+"`trigger`"+`=?, nodes=?, edges=?, updated_at=NOW() WHERE id=? AND user_id=?`,
		name, trigger, nodes, edges, id, uid)
	if err != nil { return err }
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (d *DB) DeleteFlow(id, uid int64) error {
	res, err := d.sql.Exec(`DELETE FROM chat_flows WHERE id=? AND user_id=?`, id, uid)
	if err != nil { return err }
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (d *DB) ToggleFlow(id, uid int64) (bool, error) {
	_, err := d.sql.Exec(`UPDATE chat_flows SET active=1-active WHERE id=? AND user_id=?`, id, uid)
	if err != nil {
		return false, err
	}
	var active int
	d.sql.QueryRow(`SELECT active FROM chat_flows WHERE id=? AND user_id=?`, id, uid).Scan(&active)
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
	if err != nil { return nil, err }
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

func (d *DB) migrateFlowStats() error {
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS flow_stats (flow_id BIGINT NOT NULL PRIMARY KEY, trigger_count INT DEFAULT 0, completion_count INT DEFAULT 0) ENGINE=InnoDB`)
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS flow_node_hits (flow_id BIGINT NOT NULL, node_id VARCHAR(64) NOT NULL, hit_count INT DEFAULT 0, PRIMARY KEY(flow_id, node_id)) ENGINE=InnoDB`)
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS flow_counters (counter_key VARCHAR(255) NOT NULL PRIMARY KEY, count_value INT DEFAULT 0, expires_at DATETIME NULL) ENGINE=InnoDB`)
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS flow_execution_log (id BIGINT AUTO_INCREMENT PRIMARY KEY, flow_id BIGINT NOT NULL, flow_name VARCHAR(255), phone VARCHAR(64), `+"`trigger`"+` VARCHAR(32), nodes_visited INT DEFAULT 0, replies_count INT DEFAULT 0, status VARCHAR(20) DEFAULT 'completed', error_msg TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, INDEX idx_fel_flow(flow_id)) ENGINE=InnoDB`)
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS flow_paused (phone VARCHAR(64) NOT NULL PRIMARY KEY, flow_id BIGINT NOT NULL, account_phone VARCHAR(64), variables JSON, current_node VARCHAR(64), visited_nodes JSON, contact_name VARCHAR(255), message TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, INDEX idx_fp_flow(flow_id)) ENGINE=InnoDB`)
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS flow_reviews (id BIGINT AUTO_INCREMENT PRIMARY KEY, flow_id BIGINT NOT NULL, user_id BIGINT NOT NULL, rating TINYINT NOT NULL DEFAULT 5, review TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, INDEX idx_fr_flow(flow_id)) ENGINE=InnoDB`)
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS flow_downloads (flow_id BIGINT NOT NULL, user_id BIGINT NOT NULL, PRIMARY KEY(flow_id, user_id)) ENGINE=InnoDB`)
	d.safeAddColumn("flow_counters", "expires_at", "DATETIME NULL")
	return nil
}

func (d *DB) LogFlowExecution(flowID int64, flowName, phone, trigger string, nodesVisited, repliesCount int, status, errMsg string) {
	d.sql.Exec(`INSERT INTO flow_execution_log (flow_id, flow_name, phone, `+"`trigger`"+`, nodes_visited, replies_count, status, error_msg) VALUES (?,?,?,?,?,?,?,?)`, flowID, flowName, phone, trigger, nodesVisited, repliesCount, status, errMsg)
}

func (d *DB) GetFlowExecutionLog(flowID int64, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 { limit = 50 }
	rows, err := d.sql.Query(`SELECT flow_name, phone, `+"`trigger`"+`, nodes_visited, replies_count, status, created_at FROM flow_execution_log WHERE flow_id=? ORDER BY id DESC LIMIT ?`, flowID, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var name, phone, trig, status string
		var nodes, replies int
		var created string
		rows.Scan(&name, &phone, &trig, &nodes, &replies, &status, &created)
		out = append(out, map[string]interface{}{"flow": name, "phone": phone, "trigger": trig, "nodes": nodes, "replies": replies, "status": status, "time": created})
	}
	return out, nil
}
func (d *DB) IncFlowTrigger(fid int64) { d.sql.Exec(`INSERT INTO flow_stats (flow_id,trigger_count) VALUES (?,1) ON DUPLICATE KEY UPDATE trigger_count=trigger_count+1`, fid) }
func (d *DB) IncFlowComplete(fid int64) { d.sql.Exec(`INSERT INTO flow_stats (flow_id,completion_count) VALUES (?,1) ON DUPLICATE KEY UPDATE completion_count=completion_count+1`, fid) }
func (d *DB) IncNodeHit(fid int64, nid string) { d.sql.Exec(`INSERT INTO flow_node_hits (flow_id,node_id,hit_count) VALUES (?,?,1) ON DUPLICATE KEY UPDATE hit_count=hit_count+1`, fid, nid) }
func (d *DB) GetFlowStats(fid int64) (tc, cc int, nh map[string]int) {
	nh = map[string]int{}
	d.sql.QueryRow(`SELECT IFNULL(trigger_count,0),IFNULL(completion_count,0) FROM flow_stats WHERE flow_id=?`, fid).Scan(&tc, &cc)
	rows, _ := d.sql.Query(`SELECT node_id, hit_count FROM flow_node_hits WHERE flow_id=?`, fid)
	if rows != nil { defer rows.Close(); for rows.Next() { var nid string; var h int; rows.Scan(&nid, &h); nh[nid] = h } }
	return
}
func (d *DB) IncFlowCounter(key string) int {
	d.sql.Exec(`INSERT INTO flow_counters (counter_key,count_value) VALUES (?,1) ON DUPLICATE KEY UPDATE count_value=count_value+1`, key)
	var c int; d.sql.QueryRow(`SELECT count_value FROM flow_counters WHERE counter_key=?`, key).Scan(&c); return c
}

// Flow version history
func (d *DB) migrateFlowVersions() error {
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS flow_versions (id BIGINT AUTO_INCREMENT PRIMARY KEY, flow_id BIGINT NOT NULL, name VARCHAR(255), nodes JSON, edges JSON, saved_at DATETIME DEFAULT CURRENT_TIMESTAMP, INDEX idx_fv_flow(flow_id)) ENGINE=InnoDB`)
	return nil
}
func (d *DB) SaveFlowVersion(fid int64, name, nodes, edges string) error {
	_, err := d.sql.Exec(`INSERT INTO flow_versions (flow_id, name, nodes, edges) VALUES (?,?,?,?)`, fid, name, nodes, edges)
	return err
}
func (d *DB) GetFlowVersions(fid int64) ([]map[string]interface{}, error) {
	rows, err := d.sql.Query(`SELECT id, name, saved_at FROM flow_versions WHERE flow_id=? ORDER BY id DESC LIMIT 20`, fid)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var id int64; var name, savedAt string
		rows.Scan(&id, &name, &savedAt)
		out = append(out, map[string]interface{}{"id": id, "name": name, "saved_at": savedAt})
	}
	return out, nil
}
func (d *DB) RollbackFlow(fid, versionID int64) error {
	var nodes, edges string
	err := d.sql.QueryRow(`SELECT nodes, edges FROM flow_versions WHERE id=? AND flow_id=?`, versionID, fid).Scan(&nodes, &edges)
	if err != nil { return err }
	_, err = d.sql.Exec(`UPDATE chat_flows SET nodes=?, edges=?, updated_at=NOW() WHERE id=?`, nodes, edges, fid)
	return err
}

// Marketplace
func (d *DB) migrateFlowMarket() error {
	d.safeAddColumn("chat_flows", "public", "TINYINT DEFAULT 0")
	d.safeAddColumn("chat_flows", "downloads", "INT DEFAULT 0")
	return nil
}
func (d *DB) PublishFlow(fid, uid int64) error {
	res, err := d.sql.Exec(`UPDATE chat_flows SET public=1 WHERE id=? AND user_id=?`, fid, uid)
	if err != nil { return err }
	if n, _ := res.RowsAffected(); n == 0 { return fmt.Errorf("not found") }
	return nil
}
func (d *DB) UnpublishFlow(fid, uid int64) error {
	res, err := d.sql.Exec(`UPDATE chat_flows SET public=0 WHERE id=? AND user_id=?`, fid, uid)
	if err != nil { return err }
	if n, _ := res.RowsAffected(); n == 0 { return fmt.Errorf("not found") }
	return nil
}
func (d *DB) ListPublicFlows() ([]ChatFlow, error) {
	rows, err := d.sql.Query("SELECT id, user_id, name, `trigger`, nodes, edges, active, created_at, updated_at FROM chat_flows WHERE public=1 ORDER BY downloads DESC LIMIT 50")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []ChatFlow
	for rows.Next() {
		var f ChatFlow; var trig []byte
		rows.Scan(&f.ID, &f.UserID, &f.Name, &trig, &f.NodesJSON, &f.EdgesJSON, &f.Active, &f.CreatedAt, &f.UpdatedAt)
		f.Trigger = string(trig); out = append(out, f)
	}
	return out, nil
}

// Flow AI key per flow
func (d *DB) GetFlowAIKey(fid int64) int64 {
	var aiKeyID int64
	d.sql.QueryRow(`SELECT ai_key_id FROM chat_flows WHERE id=?`, fid).Scan(&aiKeyID)
	return aiKeyID
}
func (d *DB) SetFlowAIKey(fid int64, aiKeyID int64) {
	d.sql.Exec(`UPDATE chat_flows SET ai_key_id=? WHERE id=?`, aiKeyID, fid)
}

// Paused flows — persistent across server restarts
func (d *DB) SavePausedFlow(phone string, flowID int64, accountPhone, variablesJSON, currentNode, visitedNodesJSON, contactName, message string) {
	d.sql.Exec(`REPLACE INTO flow_paused (phone, flow_id, account_phone, variables, current_node, visited_nodes, contact_name, message, created_at) VALUES (?,?,?,?,?,?,?,?,NOW())`,
		phone, flowID, accountPhone, variablesJSON, currentNode, visitedNodesJSON, contactName, message)
}
func (d *DB) LoadPausedFlow(phone string) (flowID int64, accountPhone, variablesJSON, currentNode, visitedNodesJSON, contactName, message string, err error) {
	err = d.sql.QueryRow(`SELECT flow_id, account_phone, variables, current_node, visited_nodes, contact_name, message FROM flow_paused WHERE phone=?`, phone).
		Scan(&flowID, &accountPhone, &variablesJSON, &currentNode, &visitedNodesJSON, &contactName, &message)
	return
}
func (d *DB) DeletePausedFlow(phone string) {
	d.sql.Exec(`DELETE FROM flow_paused WHERE phone=?`, phone)
}

// Flow counter with expiry (hourly windows)
func (d *DB) IncFlowCounterWithExpiry(key string, ttlHours int) int {
	d.sql.Exec(`INSERT INTO flow_counters (counter_key, count_value, expires_at) VALUES (?,1,DATE_ADD(NOW(),INTERVAL ? HOUR)) ON DUPLICATE KEY UPDATE count_value=CASE WHEN expires_at<NOW() THEN 1 ELSE count_value+1 END, expires_at=CASE WHEN expires_at<NOW() THEN DATE_ADD(NOW(),INTERVAL ? HOUR) ELSE expires_at END`, key, ttlHours, ttlHours)
	var c int
	d.sql.QueryRow(`SELECT count_value FROM flow_counters WHERE counter_key=?`, key).Scan(&c)
	return c
}

// Flow reviews
func (d *DB) AddFlowReview(flowID, uid int64, rating int, review string) error {
	_, err := d.sql.Exec(`INSERT INTO flow_reviews (flow_id, user_id, rating, review) VALUES (?,?,?,?)`, flowID, uid, rating, review)
	return err
}
func (d *DB) GetFlowReviews(flowID int64) ([]map[string]interface{}, error) {
	rows, err := d.sql.Query(`SELECT r.rating, r.review, u.name, r.created_at FROM flow_reviews r LEFT JOIN users u ON r.user_id=u.id WHERE r.flow_id=? ORDER BY r.id DESC LIMIT 50`, flowID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var rating int; var review, name, created string; var rname sql.NullString
		rows.Scan(&rating, &review, &rname, &created)
		name = rname.String
		out = append(out, map[string]interface{}{"rating": rating, "review": review, "user_name": name, "time": created})
	}
	return out, nil
}
func (d *DB) GetFlowAvgRating(flowID int64) float64 {
	var avg float64
	d.sql.QueryRow(`SELECT IFNULL(AVG(rating),0) FROM flow_reviews WHERE flow_id=?`, flowID).Scan(&avg)
	return avg
}

// Check if flow download has been counted
func (d *DB) HasDownloadedFlow(flowID int64, uid int64) bool {
	var c int
	d.sql.QueryRow(`SELECT COUNT(*) FROM flow_downloads WHERE flow_id=? AND user_id=?`, flowID, uid).Scan(&c)
	return c > 0
}
func (d *DB) MarkFlowDownloaded(flowID, uid int64) {
	d.sql.Exec(`INSERT IGNORE INTO flow_downloads (flow_id, user_id) VALUES (?,?)`, flowID, uid)
}
