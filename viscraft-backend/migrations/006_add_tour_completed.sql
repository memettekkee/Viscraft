-- Migration: 006_add_tour_completed.sql
-- Purpose: Add tour_completed flag to users table for per-user onboarding state

ALTER TABLE users ADD COLUMN IF NOT EXISTS tour_completed BOOLEAN NOT NULL DEFAULT FALSE;
