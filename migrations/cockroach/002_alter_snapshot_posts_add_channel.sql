ALTER TABLE ledger
ADD CONSTRAINT ledger_id_unique UNIQUE (id);
ALTER TABLE snapshot_users
ADD CONSTRAINT snapshot_users_id_unique UNIQUE (id);
ALTER TABLE snapshot_posts
ADD CONSTRAINT snapshot_posts_id_unique UNIQUE (id);
ALTER TABLE snapshot_posts
ADD CONSTRAINT fk_idowner_ref_snapshot_users FOREIGN KEY (idowner) REFERENCES snapshot_users(id);
CREATE TABLE snapshot_channels (
    id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    tag STRING NOT NULL
);
INSERT INTO snapshot_channels (id, tag)
VALUES(
        '00000000-0000-0000-0000-000000000000',
        'default'
    );
ALTER TABLE snapshot_posts
ADD COLUMN idchannel UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' REFERENCES snapshot_channels(id);