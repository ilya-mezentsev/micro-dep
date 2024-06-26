-- noinspection SqlNoDataSourceInspectionForFile

DROP TABLE IF EXISTS account;
CREATE TABLE account(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,

    registered_at BIGINT NOT NULL
);

DROP TABLE IF EXISTS author;
CREATE TABLE author(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    account_id VARCHAR(36) REFERENCES account(id),

    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(32) NOT NULL,
    registered_at BIGINT NOT NULL
);

DROP TABLE IF EXISTS auth_token;
CREATE TABLE auth_token(
    author_id VARCHAR(36) REFERENCES author(id),
    value VARCHAR(36) UNIQUE NOT NULL,

    created_at BIGINT NOT NULL,
    expired_at BIGINT NOT NULL
);

DROP TABLE IF EXISTS entity;
CREATE TABLE entity(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    account_id VARCHAR(36) REFERENCES account(id),
    author_id VARCHAR(36) REFERENCES author(id),

    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),

    UNIQUE (account_id, author_id, name)
);

DROP TABLE IF EXISTS entity_endpoint;
CREATE TABLE entity_endpoint(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,

    entity_id VARCHAR(36) REFERENCES entity(id) ON DELETE CASCADE,
    kind VARCHAR(42) NOT NULL,
    address VARCHAR(500),

    UNIQUE (entity_id, kind, address)
);

DROP TABLE IF EXISTS dependency;
CREATE TABLE dependency(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,

    from_id VARCHAR(36) REFERENCES entity(id),
    to_id VARCHAR(36) REFERENCES entity_endpoint(id)
);

-- some test data

INSERT INTO account(id, registered_at) VALUES
   ('c5d6f777-8195-4908-af7d-6e3c1cd2460e', 1699191331),
   ('32f157cc-23bc-4bc4-a40c-ad9384406809', 1699191471);

-- author for each account; passwords == md5(username)
INSERT INTO author(id, account_id, username, password, registered_at) VALUES
    ('4a7b8037-2cba-4667-bf1b-a2d227c7b80b', 'c5d6f777-8195-4908-af7d-6e3c1cd2460e', 'foo', 'acbd18db4cc2f85cedef654fccc4a4d8', 1699191331),
    ('c332c855-6b4c-4582-97f2-5ad196ff436c', '32f157cc-23bc-4bc4-a40c-ad9384406809', 'bar', '37b51d194a7513e45b56f6524f2d51f2', 1699191471);

INSERT INTO auth_token(author_id, value, created_at, expired_at) VALUES
    (
        (SELECT id FROM author LIMIT 1),
        'f93676f8f379c74cefc0d9959d843ac0', 1699191331,
        (SELECT EXTRACT(EPOCH FROM TIMESTAMP '2042-12-12'))
    ),
    (
        (SELECT id FROM author LIMIT 1 OFFSET 1),
        'b3752f1e705230fbd4ab3732357774cb', 1699191471,
        (SELECT EXTRACT(EPOCH FROM TIMESTAMP '2042-12-12'))
    );

INSERT INTO entity(account_id, author_id, name, description) VALUES
    (
        (SELECT id FROM account LIMIT 1),
        (SELECT id FROM author LIMIT 1),
        'service-1', 'some first service'
    ),
    (
        (SELECT id FROM account LIMIT 1),
        (SELECT id FROM author LIMIT 1),
        'service-2', 'some second service'
    );

INSERT INTO entity_endpoint(entity_id, kind, address) VALUES
    (
        (SELECT id FROM entity LIMIT 1),
        'get-endpoint', '/api/v1/user'
    ),
    (
        (SELECT id FROM entity LIMIT 1),
        'post-endpoint', '/api/v1/user'
    ),
    (
        (SELECT id FROM entity LIMIT 1),
        'patch-endpoint', '/api/v1/user'
    ),

    (
        (SELECT id FROM entity LIMIT 1 OFFSET 1),
        'get-endpoint', '/api/v1/comments'
    ),
    (
        (SELECT id FROM entity LIMIT 1 OFFSET 1),
        'get-endpoint', '/api/v1/comments/<id:int>'
    );
