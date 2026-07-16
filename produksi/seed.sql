-- ===========================================
-- ChatGo Default Seed Data
-- Run: mysql -u root -p chatgo < seed.sql
-- Safe to re-run (pakai INSERT IGNORE)
-- ===========================================

-- Default Roles
INSERT IGNORE INTO roles (name, permissions) VALUES 
('admin', 'manage_users,manage_packages,manage_waservers,manage_plugins'),
('user', '');

-- Default Packages  
INSERT IGNORE INTO packages (name, price, send_limit, device_limit, receive_limit, wa_send_limit, wa_receive_limit, wa_account_limit, contact_limit, scheduled_limit, key_limit, webhook_limit, action_limit) VALUES 
('Free', '0', 50, 0, 50, 50, 50, 0, 100, 5, 2, 2, 2),
('Pro', '29', 500, 3, 500, 500, 500, 3, 1000, 20, 10, 10, 10),
('Enterprise', '99', 9999, 10, 9999, 9999, 9999, 10, 99999, 100, 100, 100, 100);

-- Default Admin User (password: password)
INSERT IGNORE INTO users (name, email, role, country, password) VALUES
('Admin', 'admin@chatgo.test', 'admin', 'ID', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8');
