CREATE TABLE users (
                         id                     serial          not null unique,
                         login                  varchar(255)    not null unique,
                         email                  varchar(255)    not null unique,
                         password_hash          varchar(255)    not null,
                         confirmation           varchar(255),   --Токен для подтверждения
                         is_confirm             boolean         default false,
                         refresh_token          varchar(255),
                         "refresh_token_expires"  int,
                         first_name             varchar(255),
                         second_name            varchar(255),
                         last_name              varchar(255),
                         image                  varchar(255),   --Путь к картинке или сама картинка в Base64
                         gender                 smallint,
                         birthday               int,
                         description            varchar(255),
                         register_date          int       not null,
                         last_date              int,
                         my_links               jsonb,
                         PRIMARY KEY (id)
);
COMMENT ON COLUMN users.confirmation    IS 'Токен для подтверждения при регистрации';
COMMENT ON COLUMN users.image           IS 'Путь к картинке или сама картинка в Base64';


CREATE TABLE boards (
                          id            serial          not null unique,
                          title         varchar(255)    not null,
                          description   varchar(255)    not null,
                          background    varchar(255),
                          owner_id      int         references users(id) default null,
                          created_date  timestamp,
                          editors       jsonb,          --Возможно не будет использоваться
                          edit_date     timestamp,
                          users         jsonb,          --Возможно не будет использоваться
                          PRIMARY KEY (id)
);
COMMENT ON COLUMN boards.editors    IS 'Возможно не будет использоваться';
COMMENT ON COLUMN boards.users      IS 'Возможно не будет использоваться';

CREATE TABLE columns (
                           id           serial          not null unique,
                           title        varchar(255)    not null,
                           description  varchar(255)    not null,
                           color        varchar(23),    --Максимум 25 символов
                           board_id     int             not null references boards(id) default null,
                           PRIMARY KEY (id)
);
COMMENT ON COLUMN columns.color      IS 'Максимум 25 символов';

CREATE TABLE tickets (
                           id           serial not null unique,
                           title        varchar(255),
                           description  varchar(255),
                           attachment   jsonb,
                           owner_id     int references users(id) default null,
                           date_created timestamp,
                           date_start   timestamp,
                           date_end     timestamp,
                           users        jsonb,  --Возможно не будет использоваться
                           tags         varchar(255),
                           level        smallint,
                           column_id    int references columns(id) on delete set null,
                           color        varchar(25),
                           PRIMARY KEY (id)
);
COMMENT ON COLUMN tickets.users    IS 'Возможно не будет использоваться';
COMMENT ON COLUMN tickets.color    IS 'Максимум 25 символов';

CREATE TABLE comments (
                          id            serial not null unique,
                          ticked_id     int,
                          description   varchar(255),
                          created_date  timestamp,
                          user_id          int references users(id) on delete set null,
                          attachment json,
                          PRIMARY KEY (id)
);
COMMENT ON COLUMN comments.attachment    IS 'Возможно не будет использоваться';

CREATE TABLE attachments (
                             id           serial          not null unique ,
                             patch        varchar(255)    not null,
                             type         varchar(255),
                             owner_id     int   default null   references users(id)      on delete set null,
                             ticked_id    int   default null   references tickets (id)   on delete set null,
                             comment_id   int   default null   references comments (id)  on delete set null,
                             PRIMARY KEY (id)
);

CREATE TABLE access (
                        user_id           int default null references users (id),
                        board_id_editor   int default null references boards (id),
                        column_id_editor  int default null references columns (id),
                        ticket_id_editor  int default null references tickets (id),
                        board_id_views    int default null references boards (id),
                        column_id_views   int default null references columns (id),
                        ticket_id_views   int default null references tickets (id)
);
