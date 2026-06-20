-- Migration: 003_scenes.sql
-- Purpose: Replace images table with scenes table; add genre/art_style to projects

-- Add genre and art_style columns to projects
ALTER TABLE projects ADD COLUMN genre VARCHAR(100) NOT NULL DEFAULT 'fantasy';
ALTER TABLE projects ADD COLUMN art_style VARCHAR(255) NOT NULL DEFAULT '';

-- Create scenes table
CREATE TABLE scenes (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id              UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id                 UUID REFERENCES users(id),
    order_index             INT NOT NULL,
    prompt                  TEXT NOT NULL,
    reference_scene_id      UUID REFERENCES scenes(id),
    used_uploaded_reference  BOOLEAN DEFAULT FALSE,
    status                  VARCHAR(50) DEFAULT 'processing',
    file_path               VARCHAR(500),
    file_url                VARCHAR(500),
    error_code              VARCHAR(20),
    created_at              TIMESTAMP DEFAULT NOW()
);

-- Indexes for query performance
CREATE INDEX idx_scenes_project_id ON scenes(project_id);
CREATE INDEX idx_scenes_order ON scenes(project_id, order_index);

-- Drop old images table and its indexes
DROP TABLE IF EXISTS images CASCADE;
