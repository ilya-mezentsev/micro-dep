-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE IF NOT EXISTS account(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,

    registered_at BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS author(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    account_id VARCHAR(36) REFERENCES account(id),

    username VARCHAR(42) UNIQUE NOT NULL,
    password VARCHAR(32) NOT NULL,
    registered_at BIGINT NOT NULL,

    permissions_kind VARCHAR(42) NOT NULL
);

CREATE TABLE IF NOT EXISTS auth_token(
    author_id VARCHAR(36) REFERENCES author(id),
    value VARCHAR(36) UNIQUE NOT NULL,

    created_at BIGINT NOT NULL,
    expired_at BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS entity(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    author_id VARCHAR(36) REFERENCES author(id),

    name VARCHAR(100) UNIQUE NOT NULL,
    description VARCHAR(500)
);

CREATE TABLE IF NOT EXISTS entity_endpoint(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,

    entity_id VARCHAR(36) REFERENCES entity(id),
    kind VARCHAR(42) NOT NULL,
    address VARCHAR(100),

    UNIQUE (entity_id, kind, address)
);

CREATE TABLE IF NOT EXISTS dependency(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,

    from_id VARCHAR(36) REFERENCES entity(id),
    to_id VARCHAR(36) REFERENCES entity_endpoint(id)
);

CREATE OR REPLACE VIEW account_linked_entity AS
SELECT
    s.*,
    a.account_id
FROM entity s
INNER JOIN author a on s.author_id = a.id;
