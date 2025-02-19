-- INITIALIZE SQL DATABASE

BEGIN;

-- CREATE TABLES
CREATE TABLE IF NOT EXISTS appointments (
  id SERIAL PRIMARY KEY,
  trainer_id INT NOT NULL,
  user_id INT NOT NULL,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT appointments_trainer_start UNIQUE (trainer_id, starts_at)
);

-- CREATE INDEXES
CREATE INDEX IF NOT EXISTS idx_appointments_id ON appointments(id);
CREATE INDEX IF NOT EXISTS idx_appointments_trainer_id ON appointments(trainer_id);
CREATE INDEX IF NOT EXISTS idx_appointments_user_id ON appointments(user_id);
CREATE INDEX IF NOT EXISTS idx_appointments_starts_at ON appointments(starts_at);

COMMIT;

BEGIN;

-- INSERT DATA
INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at) VALUES
(1, 1, '2019-01-24 09:00:00-08', '2019-01-24 09:30:00-08');

INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at) VALUES
(1, 2, '2019-01-24 10:00:00-08', '2019-01-24 10:30:00-08');

INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at) VALUES
(1, 3, '2019-01-25 10:00:00-08', '2019-01-25 10:30:00-08');

INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at) VALUES
(1, 4, '2019-01-25 10:30:00-08', '2019-01-25 11:00:00-08');

INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at) VALUES
(1, 5, '2019-01-26 10:00:00-08', '2019-01-26 10:30:00-08');

INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at) VALUES
(2, 6, '2019-01-24 09:00:00-08', '2019-01-24 09:30:00-08');

INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at) VALUES
(2, 7, '2019-01-26 10:00:00-08', '2019-01-26 10:30:00-08');

INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at) VALUES
(3, 8, '2019-01-26 12:00:00-08', '2019-01-26 12:30:00-08');

INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at) VALUES
(3, 9, '2019-01-26 13:00:00-08', '2019-01-26 13:30:00-08');

INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at) VALUES
(3, 10, '2019-01-26 14:00:00-08', '2019-01-26 14:30:00-08');

COMMIT;
