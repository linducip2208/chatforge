-- ===========================================
-- ChatGo Default Seed Data
-- Run: mysql -u root -p chatgo < seed.sql
-- Safe to re-run (INSERT IGNORE)
-- ===========================================

-- Admin role: all 40 permissions
INSERT IGNORE INTO roles (name, permissions) VALUES 
('admin', 'manage_users,manage_roles,manage_packages,manage_vouchers,manage_subscriptions,manage_transactions,manage_payouts,manage_pages,manage_marketing,manage_languages,manage_waservers,manage_gateways,manage_shorteners,manage_plugins,manage_meta,manage_metatemplates,wa_send,wa_broadcast,wa_scheduled,wa_sent,wa_received,wa_inbox,wa_status,wa_autoreply,wa_ai_keys,wa_ai_plugins,wa_knowledge,wa_contacts,wa_groups,wa_unsub,wa_templates,wa_apikeys,wa_webhooks,wa_logger,wa_settings,wa_docs,wa_hosts,wa_ussd,wa_impersonate'),
('user', 'wa_send,wa_inbox,wa_templates,wa_apikeys,wa_settings,wa_docs');

-- Packages with meta_limit
INSERT IGNORE INTO packages (name, price, send_limit, device_limit, receive_limit, wa_send_limit, wa_receive_limit, wa_account_limit, contact_limit, scheduled_limit, key_limit, webhook_limit, action_limit, meta_limit) VALUES 
('Free', '0', 50, 0, 50, 50, 50, 0, 100, 5, 2, 2, 2, 0),
('Pro', '29', 500, 3, 500, 500, 500, 3, 1000, 20, 10, 10, 10, 2),
('Enterprise', '99', 9999, 10, 9999, 9999, 9999, 10, 99999, 100, 100, 100, 100, 10);
