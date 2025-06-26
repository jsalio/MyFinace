-- CreateEnum
CREATE TYPE "WalletType" AS ENUM ('Debit', 'Credit');

-- Creating the wallets table to store user wallet information
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type "WalletType" NOT NULL,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    user_id INTEGER NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT unique_user_wallet_name UNIQUE (user_id, name)
);

-- Adding comments for better documentation
COMMENT ON TABLE wallets IS 'Stores information about user wallets';
COMMENT ON COLUMN wallets.id IS 'Unique identifier for the wallet';
COMMENT ON COLUMN wallets.name IS 'User-defined identifier for the wallet, unique per user';
COMMENT ON COLUMN wallets.type IS 'Type of wallet (e.g., checking, savings, credit)';
COMMENT ON COLUMN wallets.balance IS 'Current monetary amount in the wallet';
COMMENT ON COLUMN wallets.user_id IS 'Foreign key referencing the user who owns this wallet';