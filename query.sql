SELECT COUNT(*) as total, COUNT(CASE WHEN wa_phone='' THEN 1 END) as empty_wa FROM received;
SELECT id, phone, wa_phone, LEFT(message,30) as msg FROM received WHERE wa_phone!='' ORDER BY id DESC LIMIT 5;
