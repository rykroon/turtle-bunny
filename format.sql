CREATE TABLE IF NOT EXISTS accounts (
    id TEXT PRIMARY KEY CHECK (is_uint128(id)),
    debits_posted TEXT NOT NULL CHECK (is_uint128(debits_posted)),
    credits_posted TEXT NOT NULL CHECK (is_uint128(credits_posted)),
    user_data_128 TEXT NOT NULL CHECK (is_uint128(user_data_128)),
    user_data_64 TEXT NOT NULL CHECK (is_uint64(user_data_64)),
    user_data_32 INTEGER NOT NULL CHECK (user_data_32 BETWEEN 0 AND 4294967295),
    ledger INTEGER NOT NULL CHECK (ledger BETWEEN 0 AND 4294967295),
    code INTEGER NOT NULL CHECK (code BETWEEN 0 AND 65535),
    debits_must_not_exceed_credits INTEGER NOT NULL CHECK (debits_must_not_exceed_credits IN (0,1)),
    credits_must_not_exceed_debits INTEGER NOT NULL CHECK (credits_must_not_exceed_debits IN (0,1)),
    imported INTEGER NOT NULL CHECK (imported IN (0,1)),
    timestamp TEXT NOT NULL UNIQUE CHECK (is_uint64(timestamp))
) STRICT, WITHOUT ROWID;


CREATE TRIGGER IF NOT EXISTS before_create_account BEFORE INSERT ON accounts
BEGIN
    SELECT
    CASE
        WHEN NEW.id = '0'
            THEN RAISE(ABORT, "id_must_not_be_zero")
        WHEN NEW.id = get_uint128_max()
            THEN RAISE(ABORT, "id_must_not_be_int_max")
        WHEN NEW.debits_posted != '0'
            THEN RAISE(ABORT, "debits_posted_must_be_zero")
        WHEN NEW.credits_posted != '0'
            THEN RAISE(ABORT, "credits_posted_must_be_zero")
        WHEN NEW.ledger = 0
            THEN RAISE(ABORT, "ledger_must_not_be_zero")
        WHEN NEW.code = 0
            THEN RAISE(ABORT, "code_must_not_be_zero")
        WHEN NEW.debits_must_not_exceed_credits AND NEW.credits_must_not_exceed_debits
            THEN RAISE(ABORT, "flags_are_mutually_exclusive")
        WHEN uint_cmp(NEW.timestamp, unix_nano()) = 1
            THEN RAISE(ABORT, "timestamp_must_not_advance")
        WHEN uint_cmp(NEW.timestamp, COALESCE((SELECT ts FROM last_account_timestamp), 0)) = -1
            THEN RAISE(ABORT, "timestamp_must_not_regress")
    END;
END;

CREATE TRIGGER IF NOT EXISTS before_update_account BEFORE UPDATE ON accounts
BEGIN
    SELECT
        CASE
        WHEN (
            (OLD.id != NEW.id) OR
            (OLD.user_data_128 != NEW.user_data_128) OR
            (OLD.user_data_64 != NEW.user_data_64) OR
            (OLD.user_data_32 != NEW.user_data_32) OR
            (OLD.ledger != NEW.ledger) OR
            (OLD.code != NEW.code) OR
            (OLD.debits_must_not_exceed_credits != NEW.debits_must_not_exceed_credits) OR
            (OLD.credits_must_not_exceed_debits != NEW.credits_must_not_exceed_debits) OR
            (OLD.imported != NEW.imported) OR
            (OLD.timestamp != NEW.timestamp)
        )
            THEN RAISE(ABORT, "account_cannot_be_modified")
        WHEN NEW.debits_must_not_exceed_credits AND NEW.debits_posted > NEW.credits_posted
            THEN RAISE(ABORT, "exceeds_credits")
        WHEN NEW.credits_must_not_exceed_debits AND NEW.credits_posted > NEW.debits_posted
            THEN RAISE(ABORT, "exceeds_debits")
    END;
END;

CREATE TRIGGER IF NOT EXISTS prevent_delete_on_accounts BEFORE DELETE ON accounts
BEGIN
    SELECT CASE WHEN true THEN RAISE(ABORT, "account_cannot_be_deleted") END;
END;


CREATE TABLE IF NOT EXISTS transfers (
    id TEXT PRIMARY KEY CHECK (is_uint128(id)),
    debit_account_id TEXT NOT NULL CHECK (is_uint128(debit_account_id)),
    credit_account_id TEXT NOT NULL CHECK (is_uint128(credit_account_id)),
    amount TEXT NOT NULL CHECK (is_uint128(amount)),
    user_data_128 TEXT NOT NULL CHECK (is_uint128(user_data_128)),
    user_data_64 TEXT NOT NULL CHECK (is_uint64(user_data_64)),
    user_data_32 INTEGER NOT NULL CHECK (user_data_32 BETWEEN 0 AND 4294967295),
    ledger INTEGER NOT NULL CHECK (ledger BETWEEN 0 AND 4294967295),
    code INTEGER NOT NULL CHECK (code BETWEEN 0 AND 65535),
    imported INTEGER NOT NULL CHECK (imported IN (0, 1)),
    timestamp TEXT NOT NULL UNIQUE CHECK (is_uint64(timestamp)),
    FOREIGN KEY (debit_account_id) REFERENCES accounts(id),
    FOREIGN KEY (credit_account_id) REFERENCES accounts(id)
) STRICT, WITHOUT ROWID;


CREATE TRIGGER IF NOT EXISTS before_create_transfer BEFORE INSERT ON transfers
BEGIN
    SELECT
    CASE
        WHEN NEW.id = '0' THEN RAISE(ABORT, "id_must_not_be_zero")

        WHEN NEW.id = get_uint128_max()
            THEN RAISE(ABORT, "id_must_not_be_int_max")

        WHEN NEW.debit_account_id = '0'
            THEN RAISE(ABORT, "debit_account_id_must_not_be_zero")
        
        WHEN NEW.debit_account_id = get_uint128_max()
            THEN RAISE(ABORT, "debit_account_id_must_not_be_int_max")

        WHEN NEW.credit_account_id = '0'
            THEN RAISE(ABORT, "credit_account_id_must_not_be_zero")
        
        WHEN NEW.credit_account_id = get_uint128_max()
            THEN RAISE(ABORT, "credit_account_id_must_not_be_int_max")

        WHEN NEW.debit_account_id = NEW.credit_account_id
            THEN RAISE(ABORT, "accounts_must_be_different")

        WHEN NEW.ledger = 0
            THEN RAISE(ABORT, "ledger_must_not_be_zero")

        WHEN NEW.code = 0
            THEN RAISE(ABORT, "code_must_not_be_zero")

        WHEN (SELECT id FROM accounts WHERE id = NEW.debit_account_id) IS NULL
            THEN RAISE(ABORT, "debit_account_not_found")

        WHEN (SELECT id FROM accounts WHERE id = NEW.credit_account_id) IS NULL
            THEN RAISE(ABORT, "credit_account_not_found")

        WHEN (SELECT ledger FROM accounts WHERE id = NEW.debit_account_id) != NEW.ledger
            THEN RAISE(ABORT, 'transfer_must_have_the_same_ledger_as_accounts')

        WHEN (SELECT ledger FROM accounts WHERE id = NEW.credit_account_id) != NEW.ledger
            THEN RAISE(ABORT, 'accounts_must_have_the_same_ledger')

        WHEN uint_cmp(NEW.timestamp, unix_nano()) = 1
            THEN RAISE(ABORT, "timestamp_must_not_advance")

        WHEN uint_cmp(NEW.timestamp, COALESCE((SELECT ts FROM last_transfer_timestamp), 0)) = -1
            THEN RAISE(ABORT, "timestamp_must_not_regress")
    END;
END;


CREATE TRIGGER IF NOT EXISTS after_create_transfer AFTER INSERT ON transfers
BEGIN
    UPDATE accounts
        SET debits_posted = uint_add(debits_posted, NEW.amount)
        WHERE id = NEW.debit_account_id;

    UPDATE accounts
        SET credits_posted = uint_add(credits_posted, NEW.amount)
        WHERE id = NEW.credit_account_id;
END;

CREATE TRIGGER IF NOT EXISTS prevent_update_on_transfers BEFORE UPDATE ON transfers
BEGIN
    SELECT CASE WHEN true THEN RAISE(ABORT, "transfers_cannot_be_modified") END;
END;

CREATE TRIGGER IF NOT EXISTS prevent_delete_on_transfers BEFORE DELETE ON transfers
BEGIN
    SELECT CASE WHEN true THEN RAISE(ABORT, "transfers_cannot_be_deleted") END;
END;

CREATE VIEW IF NOT EXISTS last_account_timestamp AS
    SELECT COALESCE(timestamp, 0) AS ts from accounts ORDER BY timestamp DESC LIMIT 1;

CREATE VIEW IF NOT EXISTS last_transfer_timestamp AS
    SELECT COALESCE(timestamp, 0) AS ts from transfers ORDER BY timestamp DESC LIMIT 1;
