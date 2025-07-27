ALTER TABLE sessions 
    ADD CONSTRAINT fk_voting_policies 
    FOREIGN KEY(voting_policy_id) 
    REFERENCES voting_policies(id);

ALTER TABLE sessions
    ADD CONSTRAINT fk_voters_policy
    FOREIGN KEY(voters_policy_id)
    REFERENCES voter_policies(id);

ALTER TABLE sessions 
    ADD CONSTRAINT fk_candidates_policy_id
    FOREIGN KEY(candidate_policy_id)
    REFERENCES candidate_policies(id);

ALTER TABLE sessions
    ADD CONSTRAINT fk_create_user
    FOREIGN KEY(created_by)
    REFERENCES users(id);

ALTER TABLE candidates
    ADD CONSTRAINT fk_session
    FOREIGN KEY(session_id)
    REFERENCES sessions(id);

ALTER TABLE candidates
    ADD CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES users(id);

ALTER TABLE voters
    ADD CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES users(id);

ALTER TABLE voters
    ADD CONSTRAINT fk_session
    FOREIGN KEY(session_id)
    REFERENCES sessions(id);

ALTER TABLE votes
    ADD CONSTRAINT fk_voter
    FOREIGN KEY(voter_id)
    REFERENCES voters(id);

ALTER TABLE votes
    ADD CONSTRAINT fk_candidate
    FOREIGN KEY(candidate_id)
    REFERENCES candidates(id);

ALTER TABLE votes
    ADD CONSTRAINT fk_session
    FOREIGN KEY(session_id)
    REFERENCES sessions(id);

ALTER TABLE flows
    ADD CONSTRAINT fk_session
    FOREIGN KEY(session_id)
    REFERENCES sessions(id);

ALTER TABLE flows
    ADD CONSTRAINT fk_state
    FOREIGN KEY(state_id)
    REFERENCES session_states(id);