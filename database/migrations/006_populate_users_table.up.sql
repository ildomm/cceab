INSERT INTO users (id, email, balance, last_game_result_at, games_result_validated, created_at)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'user1@example.com', 0, NOW() - INTERVAL '5 days', TRUE, NOW() - INTERVAL '1 month'),
    ('22222222-2222-2222-2222-222222222222', 'user2@example.com', 0, NOW() - INTERVAL '7 days', FALSE, NOW() - INTERVAL '2 months'),
    ('33333333-3333-3333-3333-333333333333', 'user3@example.com', 0, NOW() - INTERVAL '3 days', TRUE, NOW() - INTERVAL '3 months'),
    ('44444444-4444-4444-4444-444444444444', 'user4@example.com', 0, NOW() - INTERVAL '9 days', FALSE, NOW() - INTERVAL '4 months'),
    ('55555555-5555-5555-5555-555555555555', 'user5@example.com', 0, NOW() - INTERVAL '2 days', TRUE, NOW() - INTERVAL '5 months'),
    ('66666666-6666-6666-6666-666666666666', 'user6@example.com', 0, NOW() - INTERVAL '6 days', FALSE, NOW() - INTERVAL '6 months'),
    ('77777777-7777-7777-7777-777777777777', 'user7@example.com', 0, NOW() - INTERVAL '4 days', TRUE, NOW() - INTERVAL '7 months'),
    ('88888888-8888-8888-8888-888888888888', 'user8@example.com', 0, NOW() - INTERVAL '8 days', FALSE, NOW() - INTERVAL '8 months'),
    ('99999999-9999-9999-9999-999999999999', 'user9@example.com', 0, NOW() - INTERVAL '1 days', TRUE, NOW() - INTERVAL '9 months');
