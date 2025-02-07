--liquibase formatted sql
--changeset kiper0808:KRM-001

CREATE TABLE storage (
    id binary(16) UNIQUE NOT NULL,
    hostname VARCHAR(255) UNIQUE NOT NULL COMMENT 'Хост сервера',
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--rollback drop table storage;

CREATE TABLE chunk (
    id binary(16) UNIQUE NOT NULL,
    file_id binary(16) NOT NULL COMMENT 'ID файла',
    storage_hostname varchar(64) NOT NULL COMMENT 'Hostname хранилища',
    part int NOT NULL COMMENT 'Номер куска файла',
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--rollback drop table chunk;

--changeset kiper0808:KRM-002

INSERT INTO storage (id, hostname) VALUES
    (UUID_TO_BIN('1f4a9b6a-5d8e-4b6e-b4a5-5fd7cbb24e12'), 'storage1'),
    (UUID_TO_BIN('2d3c4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f'), 'storage2'),
    (UUID_TO_BIN('3e4f5a6b-7c8d-9e0f-1a2b-3c4d5e6f7a8b'), 'storage3'),
    (UUID_TO_BIN('4f5a6b7c-8d9e-0f1a-2b3c-4d5e6f7a8b9c'), 'storage4'),
    (UUID_TO_BIN('5a6b7c8d-9e0f-1a2b-3c4d-5e6f7a8b9c0d'), 'storage5'),
    (UUID_TO_BIN('6b7c8d9e-0f1a-2b3c-4d5e-6f7a8b9c0d1e'), 'storage6');

--rollback truncate storage;
