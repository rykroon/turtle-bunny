CREATE TABLE IF NOT EXISTS accounts (
    id TEXT PRIMARY KEY CHECK (is_uint128(id)),
    debits_posted TEXT NOT NULL DEFAULT 0 CHECK (is_uint128(debits_posted)),
    credits_posted TEXT NOT NULL DEFAULT 0 CHECK (is_uint128(credits_posted)),
    user_data_128 TEXT NOT NULL CHECK (is_uint128(user_data_128)),
    user_data_64 TEXT NOT NULL CHECK (is_uint64(user_data_64)),
    user_data_32 INTEGER NOT NULL CHECK (user_data_32 BETWEEN 0 AND 4294967295),
    ledger INTEGER NOT NULL CHECK (ledger BETWEEN 0 AND 4294967295),
    code INTEGER NOT NULL CHECK (code BETWEEN 0 AND 65535),
    debits_must_not_exceed_credits INTEGER NOT NULL CHECK (debits_must_not_exceed_credits IN (0,1)),
    credits_must_not_exceed_debits INTEGER NOT NULL CHECK (credits_must_not_exceed_debits IN (0,1)),
    timestamp TEXT NOT NULL DEFAULT (unix_nano()) CHECK (is_uint64(timestamp))
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

CREATE TRIGGER IF NOT EXISTS before_update_debits_posted
BEFORE UPDATE OF debits_posted ON accounts
BEGIN
    SELECT
    CASE
    WHEN NEW.debits_must_not_exceed_credits AND NEW.debits_posted > NEW.credits_posted
        THEN RAISE(ABORT, "exceeds_credits")
    END;
END;

CREATE TRIGGER IF NOT EXISTS before_update_credits_posted
BEFORE UPDATE OF credits_posted ON accounts
BEGIN
    SELECT
    CASE
    WHEN NEW.credits_must_not_exceed_debits AND NEW.credits_posted > NEW.debits_posted
        THEN RAISE(ABORT, "exceeds_debits")
    END;
END;

CREATE TRIGGER IF NOT EXISTS before_update_account BEFORE UPDATE ON accounts
BEGIN
    SELECT
    CASE
        WHEN OLD.id != NEW.id
            THEN RAISE(ABORT, "field 'id' is immutable")
        WHEN OLD.user_data_128 != NEW.user_data_128
            THEN RAISE(ABORT, "field 'user_data_128' is immutable")
        WHEN OLD.user_data_64 != NEW.user_data_64
            THEN RAISE(ABORT, "field 'user_data_64' is immutable")
        WHEN OLD.user_data_32 != NEW.user_data_32
            THEN RAISE(ABORT, "field 'user_data_32' is immutable")
        WHEN OLD.ledger != NEW.ledger
            THEN RAISE(ABORT, "field 'ledger' is immutable")
        WHEN OLD.code != NEW.code
            THEN RAISE(ABORT, "field 'code' is immutable")
        WHEN OLD.debits_must_not_exceed_credits != NEW.debits_must_not_exceed_credits
            THEN RAISE(ABORT, "field 'debits_must_not_exceed_credits is immutable")
        WHEN OLD.credits_must_not_exceed_debits != NEW.credits_must_not_exceed_debits
            THEN RAISE(ABORT, "field 'credits_must_not_exceed_debits' is immutable")
        WHEN OLD.timestamp != NEW.timestamp
            THEN RAISE(ABORT, "field 'timestamp' is immutable")
    END;
END;

CREATE TRIGGER IF NOT EXISTS before_delete_account BEFORE DELETE ON accounts
BEGIN
    SELECT
    CASE
        WHEN true
        THEN RAISE(ABORT, "cannot delete accounts")
    END;
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
    timestamp TEXT NOT NULL DEFAULT (unix_nano()) CHECK(is_uint64(timestamp)),
    FOREIGN KEY (debit_account_id) REFERENCES accounts(id),
    FOREIGN KEY (credit_account_id) REFERENCES accounts(id)
) STRICT, WITHOUT ROWID;


CREATE TRIGGER IF NOT EXISTS before_create_transfer BEFORE INSERT ON transfers
BEGIN
    SELECT
    CASE
        WHEN NEW.id = '0'
            THEN RAISE(ABORT, "id_must_not_be_zero")

        WHEN NEW.id = '340282366920938463463374607431768211455'
            THEN RAISE(ABORT, "id_must_not_be_int_max")

        WHEN NEW.debit_account_id = '0'
            THEN RAISE(ABORT, "debit_account_id_must_not_be_zero")

        WHEN NEW.credit_account_id = '0'
            THEN RAISE(ABORT, "credit_account_id_must_not_be_zero")

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

    END;
END;


CREATE TRIGGER IF NOT EXISTS after_create_transfer AFTER INSERT ON transfers
BEGIN
    UPDATE accounts SET debits_posted = decimal_add(debits_posted, NEW.amount) WHERE id = NEW.debit_account_id;
    UPDATE accounts SET credits_posted = decimal_add(credits_posted, NEW.amount) WHERE id = NEW.credit_account_id;
END;

CREATE TRIGGER IF NOT EXISTS before_update_transfer BEFORE UPDATE ON transfers
BEGIN
    SELECT
    CASE
    WHEN true
        THEN RAISE(ABORT, "cannot update transfers")
    END;
END;

CREATE TRIGGER IF NOT EXISTS before_delete_transfer BEFORE DELETE ON transfers
BEGIN
    SELECT
    CASE
    WHEN true
        THEN RAISE(ABORT, "cannot delete transfers")
    END;
END;
