-- Migration: 002_add_used_reference_image.sql
-- Purpose: Add used_reference_image column to track whether an image was generated using a reference image
-- Default FALSE ensures existing records are backfilled without a data migration

ALTER TABLE images ADD COLUMN used_reference_image BOOLEAN DEFAULT FALSE;
