CREATE TABLE IF NOT EXISTS accounts (
    id TEXT PRIMARY KEY CHECK (
        id REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(id, '0') IN (0, 1) AND
        decimal_cmp(id, '340282366920938463463374607431768211455') IN (-1, 0)
    ),
    debits_posted TEXT NOT NULL DEFAULT 0 CHECK (
        debits_posted REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(debits_posted, '0') IN (0, 1) AND
        decimal_cmp(debits_posted, '340282366920938463463374607431768211455') IN (-1, 0)
    ),
    credits_posted TEXT NOT NULL DEFAULT 0 CHECK (
        credits_posted REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(credits_posted, '0') IN (0, 1) AND
        decimal_cmp(credits_posted, '340282366920938463463374607431768211455') IN (-1, 0)
    ),
    user_data_128 TEXT NOT NULL CHECK (
        user_data_128 REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(user_data_128, '0') IN (0, 1) AND
        decimal_cmp(user_data_128, '340282366920938463463374607431768211455') IN (-1, 0)
    ),
    user_data_64 TEXT NOT NULL CHECK (
        user_data_128 REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(user_data_128, '0') IN (0, 1) AND
        decimal_cmp(user_data_128, '18446744073709551615') IN (-1, 0)
    ),
    user_data_32 INTEGER NOT NULL CHECK (user_data_32 BETWEEN 0 AND 4294967295),
    ledger INTEGER NOT NULL CHECK (ledger BETWEEN 0 AND 4294967295),
    code INTEGER NOT NULL CHECK (code BETWEEN 0 AND 65535),
    debits_must_not_exceed_credits INTEGER NOT NULL CHECK (debits_must_not_exceed_credits IN (0,1)),
    credits_must_not_exceed_debits INTEGER NOT NULL CHECK (credits_must_not_exceed_debits IN (0,1)),
    timestamp TEXT NOT NULL UNIQUE DEFAULT ((decimal(unixepoch('subsec') * 1000000))) CHECK (
        timestamp REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(timestamp, '0') IN (0, 1) AND
        decimal_cmp(timestamp, '18446744073709551615') IN (-1, 0) AND
        -- make sure timestamp is less than or equal to current time
    )
) STRICT, WITHOUT ROWID;


CREATE TRIGGER IF NOT EXISTS before_create_account BEFORE INSERT ON accounts
BEGIN
    SELECT
    CASE
        WHEN NEW.id = '0'
            THEN RAISE(ABORT, "id_must_not_be_zero")
        WHEN NEW.id = '340282366920938463463374607431768211455'
            THEN RAISE(ABORT, "id_must_not_be_int_max")
        WHEN NEW.debits_must_not_exceed_credits AND NEW.credits_must_not_exceed_debits
            THEN RAISE(ABORT, "flags_are_mutually_exclusive")
        WHEN NEW.debits_posted != '0'
            THEN RAISE(ABORT, "debits_posted_must_be_zero")
        WHEN NEW.credits_posted != '0'
            THEN RAISE(ABORT, "credits_posted_must_be_zero")
        WHEN NEW.ledger = 0
            THEN RAISE(ABORT, "ledger_must_not_be_zero")
        WHEN NEW.code = 0
            THEN RAISE(ABORT, "code_must_not_be_zero")
    END;
END;

CREATE TRIGGER IF NOT EXISTS before_update_account BEFORE UPDATE ON accounts
BEGIN
    SELECT
        CASE
        WHEN
            (OLD.id != NEW.id) OR
            (OLD.user_data_128 != NEW.user_data_128) OR
            (OLD.user_data_64 != NEW.user_data_64) OR
            (OLD.user_data_32 != NEW.user_data_32) OR
            (OLD.ledger != NEW.ledger) OR
            (OLD.code != NEW.code) OR
            (OLD.debits_must_not_exceed_credits != NEW.debits_must_not_exceed_credits) OR
            (OLD.credits_must_not_exceed_debits != NEW.credits_must_not_exceed_debits) OR
            (OLD.timestamp != NEW.timestamp)
        THEN RAISE(ABORT, "accounts cannot be changed")
        WHEN NEW.credits_must_not_exceed_debits AND NEW.credits_posted > NEW.debits_posted
        THEN RAISE(ABORT, "exceeds_debits")
        WHEN NEW.credits_must_not_exceed_debits AND NEW.credits_posted > NEW.debits_posted
        THEN RAISE(ABORT, "exceeds_debits")
    END;
END;

CREATE TRIGGER IF NOT EXISTS prevent_delete_on_accounts BEFORE DELETE ON accounts
BEGIN
    SELECT CASE WHEN true THEN RAISE(ABORT, "accounts cannot be deleted") END;
END;


CREATE TABLE IF NOT EXISTS transfers (
    id TEXT PRIMARY KEY CHECK (
        id REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(id, '0') IN (0, 1) AND
        decimal_cmp(id, '340282366920938463463374607431768211455') IN (-1, 0)
    ),
    debit_account_id TEXT NOT NULL CHECK (
        debit_account_id REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(debit_account_id, '0') IN (0, 1) AND
        decimal_cmp(debit_account_id, '340282366920938463463374607431768211455') IN (-1, 0)
    ),
    credit_account_id TEXT NOT NULL CHECK (
        credit_account_id REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(credit_account_id, '0') IN (0, 1) AND
        decimal_cmp(credit_account_id, '340282366920938463463374607431768211455') IN (-1, 0)
    ),
    amount TEXT NOT NULL CHECK (
        amount REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(amount, '0') IN (0, 1) AND
        decimal_cmp(amount, '340282366920938463463374607431768211455') IN (-1, 0)
    ),
    user_data_128 TEXT NOT NULL CHECK (
        user_data_128 REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(user_data_128, '0') IN (0, 1) AND
        decimal_cmp(user_data_128, '340282366920938463463374607431768211455') IN (-1, 0)
    ),
    user_data_64 TEXT NOT NULL CHECK (
        user_data_128 REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(user_data_128, '0') IN (0, 1) AND
        decimal_cmp(user_data_128, '18446744073709551615') IN (-1, 0)
    ),
    user_data_32 INTEGER NOT NULL CHECK (user_data_32 BETWEEN 0 AND 4294967295),
    ledger INTEGER NOT NULL CHECK (ledger BETWEEN 0 AND 4294967295),
    code INTEGER NOT NULL CHECK (code BETWEEN 0 AND 65535),
    timestamp TEXT NOT NULL UNIQUE DEFAULT (decimal(unixepoch('subsec') * 1000000)) CHECK (
        timestamp REGEXP '^(0|[1-9][0-9]*)$' AND
        decimal_cmp(timestamp, '0') IN (0, 1) AND
        decimal_cmp(timestamp, '18446744073709551615') IN (-1, 0)
        -- make sure timestamp is less than or equal to current time
    ),
    FOREIGN KEY (debit_account_id) REFERENCES accounts(id),
    FOREIGN KEY (credit_account_id) REFERENCES accounts(id)
) STRICT, WITHOUT ROWID;


CREATE TRIGGER IF NOT EXISTS before_create_transfer BEFORE INSERT ON transfers
BEGIN
    SELECT
    CASE
        WHEN NEW.id = '0' THEN RAISE(ABORT, "id_must_not_be_zero")

        WHEN NEW.id = '340282366920938463463374607431768211455'
            THEN RAISE(ABORT, "id_must_not_be_int_max")

        WHEN NEW.debit_account_id = '0'
            THEN RAISE(ABORT, "debit_account_id_must_not_be_zero")
        
        WHEN NEW.debit_account_id = '340282366920938463463374607431768211455'
            THEN RAISE(ABORT, "debit_account_id_must_not_be_int_max")

        WHEN NEW.credit_account_id = '0'
            THEN RAISE(ABORT, "credit_account_id_must_not_be_zero")
        
        WHEN NEW.credit_account_id = '340282366920938463463374607431768211455'
            THEN RAISE(ABORT, "credit_account_id_must_not_be_int_max")

        WHEN NEW.debit_account_id = NEW.credit_account_id
            THEN RAISE(ABORT, "accounts_must_be_different")

        WHEN NEW.ledger = 0 THEN RAISE(ABORT, "ledger_must_not_be_zero")
        WHEN NEW.code = 0 THEN RAISE(ABORT, "code_must_not_be_zero")

        WHEN (SELECT id FROM accounts WHERE id = NEW.debit_account_id) IS NULL
            THEN RAISE(ABORT, "debit_account_not_found")

        WHEN (SELECT id FROM accounts WHERE id = NEW.credit_account_id) IS NULL
            THEN RAISE(ABORT, "credit_account_not_found")

        WHEN (SELECT ledger FROM accounts WHERE id = NEW.debit_account_id) != NEW.ledger
            THEN RAISE(ABORT, 'transfer_must_have_the_same_ledger_as_accounts')

        WHEN (SELECT ledger FROM accounts WHERE id = NEW.credit_account_id) != NEW.ledger
            THEN RAISE(ABORT, 'accounts_must_have_the_same_ledger')
    END;
END;


CREATE TRIGGER IF NOT EXISTS after_create_transfer AFTER INSERT ON transfers
BEGIN
    UPDATE accounts SET debits_posted = decimal_add(debits_posted, NEW.amount) WHERE id = NEW.debit_account_id;
    UPDATE accounts SET credits_posted = decimal_add(credits_posted, NEW.amount) WHERE id = NEW.credit_account_id;
END;

CREATE TRIGGER IF NOT EXISTS prevent_update_on_transfers BEFORE UPDATE ON transfers
BEGIN
    SELECT CASE WHEN true THEN RAISE(ABORT, "transfers cannot be updated") END;
END;

CREATE TRIGGER IF NOT EXISTS prevent_delete_on_transfers BEFORE DELETE ON transfers
BEGIN
    SELECT CASE WHEN true THEN RAISE(ABORT, "transfers cannot be deleted") END;
END;
