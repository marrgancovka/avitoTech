CREATE TABLE IF NOT EXISTS "user" (
    id UUID NOT NULL PRIMARY KEY,
    username TEXT CONSTRAINT  username_length CHECK ( char_length(username) <= 30 ) NOT NULL,
    password_hash TEXT CONSTRAINT passwordHash_length CHECK ( char_length(password_hash) <= 64) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS tag (
    id SERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS feature (
    id SERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS banner (
    id SERIAL,
    content JSONB NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    id_feature INTEGER NOT NULL,
    FOREIGN KEY (id_feature) REFERENCES feature(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT PK_BannerFeature PRIMARY KEY (id, id_feature)
);

CREATE TABLE IF NOT EXISTS banner_feature_tag (
    id_banner INTEGER NOT NULL ,
    id_tag INTEGER NOT NULL,
    id_feature INTEGER NOT NULL ,
    FOREIGN KEY (id_banner, id_feature) REFERENCES banner(id, id_feature) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT PK_FeatureTag PRIMARY KEY (id_feature, id_tag)
);

INSERT INTO "user" (id, username, password_hash, is_admin) VALUES ('c6a88f63-866a-4d0c-9254-d0c0d5ef8f50', 'admin','5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5', true);
INSERT INTO "user" (id, username, password_hash, is_admin) VALUES ('c6a88f63-866a-4d0c-9254-d0c0d5ef8f51', 'user','5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5', false);
-- 12345