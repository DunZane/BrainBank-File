CREATE TABLE file
(
    id               VARCHAR(36)   NOT NULL,
    name             VARCHAR(255)  NOT NULL DEFAULT 'unknown',
    type        VARCHAR(10)   DEFAULT 'pdf',
    size             BIGINT        NOT NULL DEFAULT 0,
    path             VARCHAR(1024) NOT NULL DEFAULT '',
    storage_provider VARCHAR(50)   DEFAULT 'local',
    created_at       TIMESTAMP                              DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP                              DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    owner_id         BIGINT        NOT NULL DEFAULT 0,
    status           ENUM ('active', 'pending', 'deleted') DEFAULT 'pending',
    tags             JSON,
    checksum         CHAR(32),               -- MD5 校验和
    PRIMARY KEY (`id`)
);
