CREATE TABLE IF NOT EXISTS accounts (
    id TEXT PRIMARY KEY CHECK (is_uint128(id)),
    debits_posted TEXT NOT NULL DEFAULT 0 CHECK (is_uint128(debits_posted)),
    credits_posted TEXT NOT NULL DEFAULT 0 CHECK (is_uint128(credits_posted)),
    user_data_128 TEXT NOT NULL DEFAULT 0 CHECK (is_uint128(user_data_128)),
    user_data_64 TEXT NOT NULL DEFAULT 0 CHECK (is_uint64(user_data_64)),
    user_data_32 INTEGER NOT NULL DEFAULT 0 CHECK (user_data_32 BETWEEN 0 AND 4294967295),
    ledger INTEGER NOT NULL CHECK (ledger BETWEEN 0 AND 4294967295),
    code INTEGER NOT NULL CHECK (code BETWEEN 0 AND 65535),
    debits_must_not_exceed_credits INTEGER NOT NULL DEFAULT 0 CHECK (debits_must_not_exceed_credits IN (0,1)),
    credits_must_not_exceed_debits INTEGER NOT NULL DEFAULT 0 CHECK (credits_must_not_exceed_debits IN (0,1)),
    timestamp INTEGER NOT NULL DEFAULT (unix_nano()) CHECK (timestamp >= 0),
    CONSTRAINT id_must_not_be_zero CHECK (id != '0'),
    CONSTRAINT id_must_not_be_int_max CHECK (id != '340282366920938463463374607431768211455'),
    CONSTRAINT flags_are_mutually_exclusive CHECK (
        NOT (debits_must_not_exceed_credits AND credits_must_not_exceed_debits)
    )
    CONSTRAINT ledger_must_not_be_zero CHECK (ledger != 0),
    CONSTRAINT code_must_not_be_zero CHECK (code != 0),
    CONSTRAINT exceeds_credits CHECK (
        NOT debits_must_not_exceed_credits OR debits_posted <= credits_posted
    ),
    CONSTRAINT exceeds_debits CHECK (
        NOT credits_must_not_exceed_debits OR credits_posted <= debits_posted
    )
) STRICT, WITHOUT ROWID;


CREATE TRIGGER IF NOT EXISTS before_create_account BEFORE INSERT ON accounts
BEGIN
    SELECT
    CASE
        WHEN NEW.debits_posted != '0'
            THEN RAISE(ABORT, "debits_posted_must_be_zero")
        WHEN NEW.credits_posted != '0'
            THEN RAISE(ABORT, "credits_posted_must_be_zero")
    END;
END;


CREATE TABLE IF NOT EXISTS transfers (
    id TEXT PRIMARY KEY CHECK (is_uint128(id)),
    debit_account_id TEXT NOT NULL CHECK (is_uint128(debit_account_id)),
    credit_account_id TEXT NOT NULL CHECK (is_uint128(credit_account_id)),
    amount TEXT NOT NULL CHECK (is_uint128(amount)),
    user_data_128 TEXT NOT NULL DEFAULT 0 CHECK (is_uint128(user_data_128)),
    user_data_64 TEXT NOT NULL DEFAULT 0 CHECK (is_uint64(user_data_64)),
    user_data_32 INTEGER NOT NULL DEFAULT 0 CHECK (user_data_32 BETWEEN 0 AND 4294967295),
    ledger INTEGER NOT NULL CHECK (ledger BETWEEN 0 AND 4294967295),
    code INTEGER NOT NULL CHECK (code BETWEEN 0 AND 65535),
    timestamp INTEGER NOT NULL DEFAULT (unix_nano()) CHECK (timestamp >= 0),
    FOREIGN KEY (debit_account_id) REFERENCES accounts(id),
    FOREIGN KEY (credit_account_id) REFERENCES accounts(id),
    CONSTRAINT id_must_not_be_zero CHECK (id != '0'),
    CONSTRAINT id_must_not_be_int_max CHECK (id != '340282366920938463463374607431768211455'),
    CONSTRAINT debit_account_id_must_not_be_zero CHECK (debit_account_id != 0),
    CONSTRAINT credit_account_id_must_not_be_zero CHECK (credit_account_id != 0),
    CONSTRAINT accounts_must_be_different CHECK (debit_account_id != credit_account_id),
    CONSTRAINT ledger_must_not_be_zero CHECK (ledger != 0),
    CONSTRAINT code_must_not_be_zero CHECK (code != 0)
) STRICT, WITHOUT ROWID;


CREATE TRIGGER IF NOT EXISTS before_create_transfer BEFORE INSERT ON transfers
BEGIN
    SELECT
    CASE
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
