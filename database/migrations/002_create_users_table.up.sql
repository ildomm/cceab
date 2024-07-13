CREATE TABLE IF NOT EXISTS users (
    id                      UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    email                   VARCHAR NOT NULL,
    balance                 DECIMAL(10,2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
    last_game_result_at     TIMESTAMP(6) WITHOUT TIME ZONE NULL,
    games_result_validated  BOOLEAN NULL,
    created_at              TIMESTAMP(6) WITHOUT TIME ZONE NOT NULL
);