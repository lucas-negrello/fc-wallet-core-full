CREATE DATABASE IF NOT EXISTS wallet;
USE wallet;

CREATE TABLE IF NOT EXISTS clients
(
    id         VARCHAR(36)  NOT NULL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    created_at DATETIME     NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts
(
    id         VARCHAR(36) NOT NULL PRIMARY KEY,
    client_id  VARCHAR(36) NOT NULL,
    balance    DOUBLE      NOT NULL DEFAULT 0,
    created_at DATETIME    NOT NULL,
    CONSTRAINT fk_accounts_client FOREIGN KEY (client_id) REFERENCES clients (id)
);

CREATE TABLE IF NOT EXISTS transactions
(
    id              VARCHAR(36) NOT NULL PRIMARY KEY,
    account_id_from VARCHAR(36) NOT NULL,
    account_id_to   VARCHAR(36) NOT NULL,
    amount          DOUBLE      NOT NULL,
    created_at      DATETIME    NOT NULL,
    CONSTRAINT fk_transactions_from FOREIGN KEY (account_id_from) REFERENCES accounts (id),
    CONSTRAINT fk_transactions_to FOREIGN KEY (account_id_to) REFERENCES accounts (id)
);

-- --------------- SEEDS ------------------

INSERT INTO clients (id, name, email, created_at)
VALUES ('11111111-1111-1111-1111-111111111111', 'John Doe', 'john@wallet.com', NOW()),
       ('22222222-2222-2222-2222-222222222222', 'Jane Doe', 'jane@wallet.com', NOW());

INSERT INTO accounts (id, client_id, balance, created_at)
VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 1000, NOW()),
       ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-2222-2222-2222-222222222222', 500, NOW());