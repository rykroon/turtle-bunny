CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY CHECK (id >= 0),
    debits_posted INTEGER NOT NULL DEFAULT 0 CHECK (debits_posted >= 0),
    credits_posted INTEGER NOT NULL DEFAULT 0 CHECK (credits_posted >= 0),
    user_data_1 INTEGER NOT NULL DEFAULT 0 CHECK (user_data_1 >= 0),
    user_data_2 INTEGER NOT NULL DEFAULT 0 CHECK (user_data_2 >= 0),
    user_data_3 INTEGER NOT NULL DEFAULT 0 CHECK (user_data_3 >= 0),
    ledger INTEGER NOT NULL CHECK (ledger >= 0),
    code INTEGER NOT NULL CHECK (code >= 0),
    debits_must_not_exceed_credits INTEGER NOT NULL DEFAULT 0 CHECK (debits_must_not_exceed_credits IN (0,1)),
    credits_must_not_exceed_debits INTEGER NOT NULL DEFAULT 0 CHECK (credits_must_not_exceed_debits IN (0,1)),
    timestamp INTEGER NOT NULL DEFAULT (unixepoch()) CHECK (timestamp >= 0),
    CONSTRAINT id_must_not_be_zero CHECK (id != 0),
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
);

CREATE TABLE IF NOT EXISTS transfers (
    id INTEGER PRIMARY KEY CHECK (id >= 0),
    debit_account_id INTEGER NOT NULL,
    credit_account_id INTEGER NOT NULL,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    user_data_1 INTEGER NOT NULL DEFAULT 0 CHECK (user_data_1 >= 0),
    user_data_2 INTEGER NOT NULL DEFAULT 0 CHECK (user_data_2 >= 0),
    user_data_3 INTEGER NOT NULL DEFAULT 0 CHECK (user_data_3 >= 0),
    ledger INTEGER NOT NULL CHECK (ledger >= 0),
    code INTEGER NOT NULL CHECK (code >= 0),
    timestamp INTEGER NOT NULL DEFAULT (unixepoch()) CHECK (timestamp >= 0),
    FOREIGN KEY (debit_account_id) REFERENCES accounts(id),
    FOREIGN KEY (credit_account_id) REFERENCES accounts(id),
    CONSTRAINT id_must_not_be_zero CHECK (id != 0),
    CONSTRAINT debit_account_id_must_not_be_zero CHECK (debit_account_id != 0),
    CONSTRAINT credit_account_id_must_not_be_zero CHECK (credit_account_id != 0),
    CONSTRAINT accounts_must_be_different CHECK (debit_account_id != credit_account_id),
    CONSTRAINT ledger_must_not_be_zero CHECK (ledger != 0),
    CONSTRAINT code_must_not_be_zero CHECK (code != 0)
);


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
    UPDATE accounts SET debits_posted = debits_posted + NEW.amount WHERE id = NEW.debit_account_id;
    UPDATE accounts SET credits_posted = credits_posted + NEW.amount WHERE id = NEW.credit_account_id;
END;
