INSERT INTO voting_policies(name) VALUES 
('Free'), 
('Single Vote Only');

INSERT INTO voter_policies(name) VALUES 
('Self Assign - Requires Account Registration'), 
('Through Invitation - Can be anonymous');

INSERT INTO candidate_policies(name) VALUES 
('Self Assign'), 
('Insert Manually');

INSERT INTO session_states(name) VALUES 
('Draft'), 
('Waiting for candidates'), 
('Voting'), 
('Suspended'),
('Cancelled'),
('Counting'),
('Concluded');

INSERT INTO users(public_id, name, email, password) VALUES
('c9a2a1f3-2b1c-4c0b-8a7e-9c2a1f3a8b1c', 'First User', 'firstuser@kaahvote.com', 'firstkaah');