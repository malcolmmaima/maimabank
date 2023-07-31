CREATE TABLE "exchange_rates" (
  "id" bigserial PRIMARY KEY,
  "base_currency" varchar NOT NULL,
  "target_currency" varchar NOT NULL,
  "exchange_rate" decimal NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "exchange_rates" ("base_currency");

CREATE INDEX ON "exchange_rates" ("target_currency");
