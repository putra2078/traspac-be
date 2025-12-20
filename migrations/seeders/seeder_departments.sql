INSERT INTO departments (name, slug)
VALUES
('Human Resource', 'hr'),
('Engineering', 'engineer'),
('Finance', 'finance')
ON CONFLICT (name) DO NOTHING;
