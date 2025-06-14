/* Create Database */
CREATE DATABASE wallet_app;

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

CREATE TABLE wallet_app.user_activity (
    user_act_id VARCHAR(60) NOT NULL,
    user_id VARCHAR(60) NOT NULL,
    user_act_type VARCHAR(10) NOT NULL,
    user_act_detail VARCHAR(255) NOT NULL,
    user_wallet_id VARCHAR(60),
    user_act_time TIMESTAMP NOT NULL,
    CONSTRAINT pk_user_activity PRIMARY KEY(user_act_id)
);

CREATE TABLE wallet_app.txn_history (
    txn_id VARCHAR(60) NOT NULL,
    from_wallet_id VARCHAR(60) NOT NULL,
    to_wallet_id VARCHAR(60) NOT NULL,
    txn_type VARCHAR(10) NOT NULL,
    txn_amount NUMERIC(15, 2) NOT NULL,
    txn_time TIMESTAMP,
    CONSTRAINT pk_txn_history PRIMARY KEY(txn_id)
);