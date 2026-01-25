-- Refund history tracking table
-- Tracks all refund transactions for escrow bounties
CREATE TABLE IF NOT EXISTS refund_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  bounty_id BIGINT NOT NULL,
  amount BIGINT NOT NULL CHECK (amount > 0),
  recipient_address TEXT NOT NULL,
  refund_mode TEXT NOT NULL CHECK (refund_mode IN ('Full', 'Partial', 'Custom')),
  transaction_hash TEXT,
  ledger_number INT,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_refund_history_bounty_id ON refund_history(bounty_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_refund_history_recipient ON refund_history(recipient_address);
CREATE INDEX IF NOT EXISTS idx_refund_history_transaction ON refund_history(transaction_hash);

-- Refund approvals table
-- Tracks admin-approved refund requests (for pre-deadline refunds)
CREATE TABLE IF NOT EXISTS refund_approvals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  bounty_id BIGINT NOT NULL UNIQUE,
  amount BIGINT NOT NULL CHECK (amount > 0),
  recipient_address TEXT NOT NULL,
  refund_mode TEXT NOT NULL CHECK (refund_mode IN ('Full', 'Partial', 'Custom')),
  approved_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  approved_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  used_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_refund_approvals_bounty_id ON refund_approvals(bounty_id);
CREATE INDEX IF NOT EXISTS idx_refund_approvals_approved_by ON refund_approvals(approved_by);
CREATE INDEX IF NOT EXISTS idx_refund_approvals_unused ON refund_approvals(bounty_id) WHERE used_at IS NULL;
