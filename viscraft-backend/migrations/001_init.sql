-- Migration: 001_init.sql
-- Purpose: Initialize the Viscraft database schema
-- Creates the core tables: users, projects, images
-- Includes foreign key relationships with cascade deletes and performance indexes

CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,
    name        VARCHAR(255),
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE projects (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE images (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id     UUID REFERENCES users(id),
    prompt      TEXT NOT NULL,
    prompt_hash VARCHAR(64),
    genre       VARCHAR(100),
    asset_type  VARCHAR(100),
    mood        VARCHAR(100),
    status      VARCHAR(50) DEFAULT 'processing',
    file_path   VARCHAR(500),
    error_code  VARCHAR(20),
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Indexes for query performance
CREATE INDEX idx_images_project_id ON images(project_id);
CREATE INDEX idx_images_prompt_hash ON images(prompt_hash);
CREATE INDEX idx_images_user_id ON images(user_id);
