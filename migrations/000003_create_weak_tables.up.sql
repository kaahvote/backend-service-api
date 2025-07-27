CREATE TABLE sessions(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    public_id TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    voting_policy_id INT NOT NULL,
    voters_policy_id INT NOT NULL,
    candidate_policy_id INT NOT NULL,
    created_by INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE candidates(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    name TEXT,
    image_url TEXT,
    user_id INT,
    session_id INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE voters(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    user_id INT, 
    session_id INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL

);

CREATE TABLE votes(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    voter_id INT,
    candidate_id INT NOT NULL,
    session_id INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE flows(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    session_id INT NOT NULL,
    state_id INT NOT NULL,
    comment TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);