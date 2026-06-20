-- Migration: 004_product_ad_pivot.sql
-- Purpose: Pivot from concept art (scenes) to product ad photography.
-- Renames genre → product_category, art_style → visual_style on projects.
-- Adds user_prompt/generated_prompt split on scenes.
-- Creates prompt_options table for dynamic form options.

ALTER TABLE projects RENAME COLUMN genre TO product_category;
ALTER TABLE projects RENAME COLUMN art_style TO visual_style;

ALTER TABLE projects ALTER COLUMN product_category SET DEFAULT 'general';

ALTER TABLE scenes RENAME COLUMN prompt TO user_prompt;
ALTER TABLE scenes ADD COLUMN generated_prompt TEXT;

CREATE TABLE prompt_options (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category     VARCHAR(50) NOT NULL,
    label        VARCHAR(100) NOT NULL,
    prompt_value TEXT NOT NULL,
    sort_order   INT DEFAULT 0
);

CREATE INDEX idx_prompt_options_category ON prompt_options(category, sort_order);

-- Seed initial prompt options
INSERT INTO prompt_options (category, label, prompt_value, sort_order) VALUES
('background', 'White marble',     'white marble surface with subtle veining', 1),
('background', 'Wooden table',     'warm wooden table with natural texture', 2),
('background', 'Studio white',     'clean white studio background', 3),
('background', 'Tropical',         'lush tropical leaves background', 4),
('background', 'Dark moody',       'dark moody background with smoke effect', 5),
('lighting',   'Natural window',   'soft natural window light from left', 1),
('lighting',   'Studio softbox',   'even studio softbox lighting, no harsh shadows', 2),
('lighting',   'Dramatic side',    'dramatic side lighting with deep shadows', 3),
('lighting',   'Golden hour',      'golden hour warm sunlight', 4),
('mood',       'Luxury',           'luxury premium feel, elegant composition', 1),
('mood',       'Fresh & clean',    'fresh clean aesthetic, bright and airy', 2),
('mood',       'Minimal',          'minimal and modern', 3),
('mood',       'Playful',          'vibrant playful energy, bold colors', 4),
('angle',      'Eye-level',        'eye-level product shot', 1),
('angle',      '45 degrees',       '45-degree angle shot', 2),
('angle',      'Flat lay',         'flat lay top-down view', 3),
('angle',      'Close-up',         'close-up macro detail shot', 4),
('props',      'Flowers',          'with scattered fresh flowers nearby', 1),
('props',      'Fresh fruits',     'with fresh tropical fruits nearby', 2),
('props',      'Linen fabric',     'with linen fabric draped beside', 3),
('props',      'Water drops',      'with water droplets on surface', 4);
