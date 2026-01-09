package turtlebunny

import (
	"testing"

	"lukechampine.com/uint128"
)

func TestIdMustNotBeZero(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	account := Account{
		Ledger: 1,
		Code:   1,
	}

	err = client.CreateAccount(account)

	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != "id_must_not_be_zero" {
		t.Errorf("expected %s, got %s", "id_must_not_be_zero", err.Error())
	}
}

func TestIdMustNotBeIntMax(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	account := Account{
		Id:     uint128.Max,
		Ledger: 1,
		Code:   1,
	}
	err = client.CreateAccount(account)

	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != "id_must_not_be_int_max" {
		t.Errorf("expected %s, got %s", "id_must_not_be_int_max", err.Error())
	}
}

func TestDebitsPostedMustBeZero(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	account := Account{
		Id:           uint128.From64(1),
		DebitsPosted: uint128.From64(1),
		Ledger:       1,
		Code:         1,
	}
	err = client.CreateAccount(account)
	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != "debits_posted_must_be_zero" {
		t.Errorf("expected %s, got %s", "debits_posted_must_be_zero", err.Error())
	}
}

func TestCreditsPostedMustBeZero(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	account := Account{
		Id:            uint128.From64(1),
		CreditsPosted: uint128.From64(1),
		Ledger:        1,
		Code:          1,
	}
	err = client.CreateAccount(account)
	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != "credits_posted_must_be_zero" {
		t.Errorf("expected %s, got %s", "credits_posted_must_be_zero", err.Error())
	}
}

func TestLedgerMustNotBeZero(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	account := Account{
		Id:     uint128.From64(1),
		Ledger: 0,
		Code:   1,
	}
	err = client.CreateAccount(account)

	if err == nil {
		t.Error("expected error")
	}

	if err.Error() != "ledger_must_not_be_zero" {
		t.Errorf("expected %s, got %s", "ledger_must_not_be_zero", err.Error())
	}
}

func TestCodeMustNotBeZero(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Fatal(err)
	}

	account := Account{
		Id:     uint128.From64(1),
		Ledger: 1,
		Code:   0,
	}
	err = client.CreateAccount(account)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "code_must_not_be_zero" {
		t.Fatalf("expected %s, got %s", "code_must_not_be_zero", err.Error())
	}
}

func TestFlagsAreMutuallyExclusive(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Fatal(err)
	}

	account := Account{
		Id:                         uint128.From64(1),
		Ledger:                     1,
		Code:                       1,
		DebitsMustNotExceedCredits: true,
		CreditsMustNotExceedDebits: true,
	}
	err = client.CreateAccount(account)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "flags_are_mutually_exclusive" {
		t.Errorf("expected %s, got %s", "flags_are_mutually_exclusive", err.Error())
	}
}

func TestTimestampMustNotAdvance(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Fatal(err)
	}

	account := Account{
		Id:        uint128.From64(1),
		Ledger:    1,
		Code:      1,
		Timestamp: uint64(unixNano()) + 1e9,
	}
	err = client.CreateAccount(account)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "timestamp_must_not_advance" {
		t.Errorf("expected %s, got %s", "timestamp_must_not_advance", err.Error())
	}
}

func TestTimestampMustNotRegress(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	account1 := Account{
		Id:     uint128.From64(1),
		Ledger: 1,
		Code:   1,
	}
	err = client.CreateAccount(account1)

	if err != nil {
		t.Error(err)
	}

	account2 := Account{
		Id:        uint128.From64(2),
		Ledger:    1,
		Code:      1,
		Timestamp: uint64(unixNano()) - 1e9,
	}
	err = client.CreateAccount(account2)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "timestamp_must_not_regress" {
		t.Errorf("expected %s, got %s", "timestamp_must_not_regress", err.Error())
	}

}
