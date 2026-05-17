-- +goose Up
INSERT INTO users
values ('37d44725-7659-4e1f-b225-0757e9c7223c', 'TEST_USER', 'test_user@gnom.com', '$2a$10$Z8sp7T0HE.VOIEAFRHhUXuGCDZgil7lT1fwI.ygX8dYvzRqyBrHoO', '[{"ProviderName":"telegram","Target":"8212543139"}]'::jsonb);

-- +goose Down
DELETE FROM users
WHERE user_id = '37d44725-7659-4e1f-b225-0757e9c7223c'