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

CREATE TABLE IF NOT EXISTS service(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    owner_id VARCHAR(36) REFERENCES author(id),

    name VARCHAR(100) NOT NULL,
    description VARCHAR(500)
);

CREATE TABLE IF NOT EXISTS service_endpoint_info(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,

    kind VARCHAR(42) NOT NULL,
    address VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS service_endpoint(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,

    service_id VARCHAR(36) REFERENCES service(id),
    type_id VARCHAR(36) REFERENCES service_endpoint_info(id)
);

CREATE TABLE IF NOT EXISTS dependency(
    id VARCHAR(36) DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,

    from_id VARCHAR(36) REFERENCES service(id),
    to_id VARCHAR(36) REFERENCES service_endpoint(id)
);
