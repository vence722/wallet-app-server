/* Create Database */
CREATE DATABASE wallet_app;

/* Use Database */
\c wallet_app;

/* Create Schema */
CREATE SCHEMA wallet_app;

/* Create Tables */
CREATE TABLE wallet_app.user (
    user_id VARCHAR(60) NOT NULL,
    user_name VARCHAR(60) UNIQUE NOT NULL,
    user_hash VARCHAR(100) NOT NULL,
    create_time TIMESTAMP NOT NULL,
    update_time TIMESTAMP,
    CONSTRAINT pk_user PRIMARY KEY(user_id)
);

CREATE TABLE wallet_app.wallet (
    wallet_id VARCHAR(60) NOT NULL,
    wallet_name VARCHAR(60) NOT NULL,
    balance NUMERIC(15, 2) NOT NULL,
    create_time TIMESTAMP NOT NULL,
    update_time TIMESTAMP,
    CONSTRAINT pk_wallet PRIMARY KEY(wallet_id)
);

CREATE TABLE wallet_app.user_wallet_bridge (
    user_id VARCHAR(60) NOT NULL,
    wallet_id VARCHAR(60) NOT NULL,
    seq INT NOT NULL,
    create_time TIMESTAMP NOT NULL,
    CONSTRAINT pk_user_wallet_bridge PRIMARY KEY(user_id, wallet_id)
);

CREATE TABLE wallet_app.txn_history (
    txn_id VARCHAR(60) NOT NULL,
    from_wallet_id VARCHAR(60) NOT NULL,
    to_wallet_id VARCHAR(60) NOT NULL,
    txn_amount NUMERIC(15, 2) NOT NULL,
    txn_time TIMESTAMP,
    CONSTRAINT pk_txn_history PRIMARY KEY(txn_id)
);

/* Create Base Data */
INSERT INTO wallet_app.user (user_id, user_name, user_hash, create_time)
VALUES
('e98f3be0-9991-471e-8bcf-d08238fa8840', 'vence.lin', 'b03ddf3ca2e714a6548e7495e2a03f5e824eaac9837cd7f159c67b90fb4b7342', '2025-06-14 12:00:00'),
('2b05751e-0607-4773-aa99-0158c00e22c2', 'mike.kwok', 'b03ddf3ca2e714a6548e7495e2a03f5e824eaac9837cd7f159c67b90fb4b7342', '2025-06-14 12:00:00'),
('250315de-dd1a-4778-bce7-edc5e9a0a036', 'angel.wong', 'b03ddf3ca2e714a6548e7495e2a03f5e824eaac9837cd7f159c67b90fb4b7342', '2025-06-14 12:00:00'),
('751bb3ea-c5b5-414d-8dae-dad6a80a1c79', 'nick.lee', 'b03ddf3ca2e714a6548e7495e2a03f5e824eaac9837cd7f159c67b90fb4b7342', '2025-06-14 12:00:00');

INSERT INTO wallet_app.wallet (wallet_id, wallet_name, balance, create_time)
VALUES
('a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d', 'default wallet', 0, '2025-06-14 12:00:00'),
('34fad474-1df7-40a1-8675-0af586d02435', 'vence wallet 1', 0, '2025-06-14 12:00:00'),
('d4598f95-4eff-421e-b6c1-186ae499b16a', 'default wallet', 0, '2025-06-14 12:00:00'),
('e5d51f9f-99d2-4768-9764-1360fe0ea55d', 'default wallet', 0, '2025-06-14 12:00:00'),
('68e95347-29ad-4324-9725-eed1feaa8594', 'default wallet', 0, '2025-06-14 12:00:00');

INSERT INTO wallet_app.user_wallet_bridge (user_id, wallet_id, seq, create_time)
VALUES
('e98f3be0-9991-471e-8bcf-d08238fa8840', 'a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d', 1, '2025-06-14 12:00:00'),
('e98f3be0-9991-471e-8bcf-d08238fa8840', '34fad474-1df7-40a1-8675-0af586d02435', 2, '2025-06-14 12:00:00'),
('2b05751e-0607-4773-aa99-0158c00e22c2', 'd4598f95-4eff-421e-b6c1-186ae499b16a', 1, '2025-06-14 12:00:00'),
('250315de-dd1a-4778-bce7-edc5e9a0a036', 'e5d51f9f-99d2-4768-9764-1360fe0ea55d', 1, '2025-06-14 12:00:00'),
('751bb3ea-c5b5-414d-8dae-dad6a80a1c79', '68e95347-29ad-4324-9725-eed1feaa8594', 1, '2025-06-14 12:00:00');