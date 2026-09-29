CREATE TABLE folders (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

 owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   parent_id UUID REFERENCES folders(id) ON DELETE CASCADE,
     name VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_folders_owner_id
ON folders(owner_id);


CREATE INDEX idx_folders_parent_id
ON folders(parent_id);