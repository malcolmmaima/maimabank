CREATE TABLE "exchange_rates" (
  "id" bigserial PRIMARY KEY,
  "base_currency" varchar NOT NULL,
  "target_currency" varchar NOT NULL,
  "exchange_rate" decimal(10, 2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "exchange_rates" ("base_currency");

CREATE INDEX ON "exchange_rates" ("target_currency");

-- add unique constraint on base_currency and target_currency
ALTER TABLE "exchange_rates" ADD CONSTRAINT "base_currency_target_currency_key" UNIQUE ("base_currency", "target_currency");