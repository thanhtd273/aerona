CREATE TABLE resource (
      resource_id BIGSERIAL PRIMARY KEY,
      name VARCHAR(75) NOT NULL,
      description VARCHAR(255),
      path VARCHAR(255) NOT NULL UNIQUE,
      create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      modified_date TIMESTAMP,
      status INTEGER NOT NULL
);

CREATE TABLE _action (
     action_id BIGSERIAL PRIMARY KEY,
     name VARCHAR(75) NOT NULL,
     code INTEGER NOT NULL UNIQUE,
     description VARCHAR(255),
     create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
     modified_date TIMESTAMP,
     status INTEGER NOT NULL
);

CREATE TABLE _permission (
     permission_id BIGSERIAL PRIMARY KEY,
     name VARCHAR(75) NOT NULL,
     description VARCHAR(255),
     create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
     modified_date TIMESTAMP,
     status INTEGER NOT NULL
);

CREATE TABLE permission_resource_action (
    id BIGSERIAL PRIMARY KEY,
    permission_id BIGINT NOT NULL,
    resource_id BIGINT NOT NULL,
    action_id BIGINT NOT NULL,
    deleted BOOLEAN,

    CONSTRAINT permission_fk FOREIGN KEY (permission_id) REFERENCES _permission(permission_id),
    CONSTRAINT resource_fk FOREIGN KEY (resource_id) REFERENCES resource(resource_id),
    CONSTRAINT action_fk FOREIGN KEY (action_id) REFERENCES _action(action_id)
);

CREATE TABLE _user (
   user_id BIGSERIAL PRIMARY KEY,
   email VARCHAR(75) NOT NULL UNIQUE,
   name VARCHAR(75),
   password TEXT NOT NULL,
   avatar TEXT,
    password_reset_code INTEGER,
   password_reset_expiration_time TIMESTAMP,

   facebook_access_token VARCHAR(255),
   google_access_token VARCHAR(255),

   create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   modified_date TIMESTAMP,
   status INTEGER NOT NULL
);

CREATE TABLE _role (
   role_id BIGSERIAL PRIMARY KEY,
   "name" VARCHAR(75) NOT NULL,
   description VARCHAR(255),

   create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   modified_date TIMESTAMP,
   status INTEGER NOT NULL
);

CREATE TABLE user_role (
   id BIGSERIAL PRIMARY KEY,
   user_id BIGINT NOT NULL,
   role_id BIGINT NOT null,

   CONSTRAINT role_fk FOREIGN KEY (role_id) REFERENCES _role(role_id),
   CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES _user(user_id)
);

CREATE TABLE role_permission (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,

    CONSTRAINT role_fk FOREIGN KEY (role_id) REFERENCES _role(role_id),
    CONSTRAINT permission_fk FOREIGN KEY (permission_id) REFERENCES _permission(permission_id)
);