INSERT INTO users (password, email, role) VALUES
  ('$2a$10$fntioAJ76duQ6EmibwvKEe/F7AQrYdBpAgLKTGyf465d/hKQVrNte', 'misha1@example.com', 'admin')
ON CONFLICT (email) DO NOTHING;

INSERT INTO users (password, email, role) VALUES
  ('$2a$10$fntioAJ76duQ6EmibwvKEe/F7AQrYdBpAgLKTGyf465d/hKQVrNte', 'misha2@example.com', 'user')
ON CONFLICT (email) DO NOTHING;