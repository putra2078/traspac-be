INSERT INTO positions (name, department_id)
VALUES 
('Backend Developer', 1),
('Frontend Developer', 1),
('HR Specialist', 1),
('Finance Analyst', 1)
ON CONFLICT (name) DO NOTHING;
