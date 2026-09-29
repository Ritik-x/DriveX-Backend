CREATE TABLE files(

    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

folder_id UUID REFERENCES folders(id) ON DELETE SET NULL,
  name VARCHAR(255) NOT NULL,

 original_name VARCHAR(255) NOT NULL,
 storage_key TEXT NOT NULL UNIQUE,

    mime_type VARCHAR(100) NOT NULL,

    deleted_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()



);

CREATE INDEX idx_files_owner_id
ON files(owner_id);

CREATE INDEX idx_files_folder_id
ON files(folder_id);

CREATE INDEX idx_files_deleted_at
ON files(deleted_at);