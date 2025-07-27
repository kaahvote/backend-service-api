ALTER TABLE sessions
ALTER COLUMN public_id TYPE UUID USING public_id::uuid;