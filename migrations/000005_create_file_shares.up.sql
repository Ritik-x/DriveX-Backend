--  CREATE TABLE file_shares(
--         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

--     file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,

--     owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     shared_with_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--   permission VARCHAR(20) NOT NULL
--         CHECK (permission IN ('viewer', 'editor')),


--          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

--     UNIQUE(file_id, shared_with_user_id)
--  );

--  CREATE INDEX idx_file_shares_file_id
-- ON file_shares(file_id);

-- CREATE INDEX idx_file_shares_shared_user_id
-- ON file_shares(shared_with_user_id);s


CREATE TABLE file_shares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,

    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    shared_with_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    permission VARCHAR(20) NOT NULL
        CHECK (permission IN ('viewer', 'editor')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(file_id, shared_with_user_id)
);

CREATE INDEX idx_file_shares_file_id
ON file_shares(file_id);

CREATE INDEX idx_file_shares_shared_user_id
ON file_shares(shared_with_user_id);