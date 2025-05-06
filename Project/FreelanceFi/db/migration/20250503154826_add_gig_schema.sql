-- +goose Up
BEGIN;

CREATE TABLE categories (
  id         SERIAL      PRIMARY KEY,
  name       VARCHAR(50) UNIQUE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE jobs (
  id           SERIAL      PRIMARY KEY,
  client_id    INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  category_id  INT         NOT NULL REFERENCES categories(id),
  title        VARCHAR(150) NOT NULL,
  description  TEXT        NOT NULL,
  budget_min   NUMERIC(10,2) NOT NULL,
  budget_max   NUMERIC(10,2),
  deadline     DATE,
  status       VARCHAR(20) NOT NULL DEFAULT 'open'
                   CHECK (status IN ('open','in_progress','completed','cancelled')),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE bids (
  id             SERIAL      PRIMARY KEY,
  job_id         INT         NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
  freelancer_id  INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  amount         NUMERIC(10,2) NOT NULL,
  cover_letter   TEXT,
  status         VARCHAR(20) NOT NULL DEFAULT 'pending'
                   CHECK (status IN ('pending','accepted','rejected')),
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
  id           SERIAL      PRIMARY KEY,
  user_id      INT         NOT NULL REFERENCES users(id),
  job_id       INT         NOT NULL REFERENCES jobs(id),
  bid_id       INT         NOT NULL REFERENCES bids(id),
  type         VARCHAR(20) NOT NULL CHECK (type IN ('escrow_deposit','payout')),
  amount       NUMERIC(10,2) NOT NULL,
  recorded_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_jobs_client         ON jobs(client_id);
CREATE INDEX idx_jobs_status_created ON jobs(status, created_at DESC);
CREATE INDEX idx_jobs_category       ON jobs(category_id);
CREATE INDEX idx_bids_job            ON bids(job_id);
CREATE INDEX idx_bids_freelancer     ON bids(freelancer_id);
CREATE INDEX idx_tx_user             ON transactions(user_id);
CREATE INDEX idx_tx_job              ON transactions(job_id);
CREATE INDEX idx_tx_user_type        ON transactions(user_id, type);

COMMIT;

-- +goose Down
BEGIN;

DROP INDEX idx_tx_user_type;
DROP INDEX idx_tx_job;
DROP INDEX idx_tx_user;
DROP INDEX idx_bids_freelancer;
DROP INDEX idx_bids_job;
DROP INDEX idx_jobs_category;
DROP INDEX idx_jobs_status_created;
DROP INDEX idx_jobs_client;

DROP TABLE transactions;
DROP TABLE bids;
DROP TABLE jobs;
DROP TABLE categories;

COMMIT;
