CREATE TABLE users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    display_name VARCHAR NOT NULL
);

CREATE TABLE users_spi (
    id BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    verified_at TIMESTAMP,

    ssn CHAR(16) NOT NULL UNIQUE,
    legal_name VARCHAR NOT NULL,
    dob DATE NOT NULL,

    tax_id VARCHAR
);

CREATE TABLE wallets (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    balance_idr BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE users_wallets (
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    id_user BIGINT NOT NULL REFERENCES users (id),
    id_wallet BIGINT NOT NULL REFERENCES wallets (id),
    PRIMARY KEY (id_user, id_wallet)
);

CREATE TYPE transaction_type AS ENUM ('topup', 'withdraw', 'transfer', 'payment');

CREATE TYPE transaction_status AS ENUM ('pending', 'success', 'failed');

CREATE TABLE transactions (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    type transaction_type NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',

    ref_internal VARCHAR NOT NULL,
    ref_external VARCHAR,
    provider VARCHAR,
    note VARCHAR
);

CREATE TABLE entries (
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    id_wallet BIGINT NOT NULL REFERENCES wallets (id),
    id_transaction BIGINT NOT NULL REFERENCES transactions (id),
    PRIMARY KEY (id_wallet, id_transaction),

    amount BIGINT NOT NULL,
    balance_idr_after BIGINT NOT NULL
);

--

CREATE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
      NEW.updated_at = CURRENT_TIMESTAMP;
      RETURN NEW;
   ELSE
      RETURN OLD;
   END IF;
END;
$$ language plpgsql;

CREATE TRIGGER users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

CREATE TRIGGER users_spi_updated_at
BEFORE UPDATE ON users_spi
FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

CREATE TRIGGER wallets_updated_at
BEFORE UPDATE ON wallets
FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

CREATE TRIGGER transactions_updated_at
BEFORE UPDATE ON transactions
FOR EACH ROW EXECUTE PROCEDURE update_updated_at();
