-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE "payments" (
  "id" BIGSERIAL PRIMARY KEY,
  "version" int,
  "amount_cents" bigint,
  "currency_id" bigint,
  "beneficiary_party_id" bigint,
  "debtor_party_id" bigint,
  "payment_date" timestamptz,
  "reference" text
);

CREATE TABLE "currencies" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" text,
  "symbol" text
);

CREATE TABLE "parties" (
  "id" BIGSERIAL PRIMARY KEY,
  "account_name" text,
  "account_number" text,
  "account_number_code" text,
  "address" text,
  "bank_id" bigint,
  "name" text
);

CREATE TABLE "banks" (
  "id" BIGSERIAL PRIMARY KEY,
  "id_code" text,
  "name" text
);

ALTER TABLE "payments" ADD FOREIGN KEY ("currency_id") REFERENCES "currencies" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("beneficiary_party_id") REFERENCES "parties" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("debtor_party_id") REFERENCES "parties" ("id");

ALTER TABLE "parties" ADD FOREIGN KEY ("bank_id") REFERENCES "banks" ("id");


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE payments;
DROP TABLE currencies;
DROP TABLE parties;
DROP TABLE banks;