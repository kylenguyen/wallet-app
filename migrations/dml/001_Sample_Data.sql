-- This script inserts sample data into the users, wallets, and transactions tables.
-- It uses Common Table Expressions (CTEs) to capture the generated UUIDs
-- and use them in subsequent inserts, making the script runnable in one go.

-- Clear existing data to make the script idempotent (optional)
-- DELETE FROM transactions;
-- DELETE FROM wallets;
-- DELETE FROM users;

DO $$
DECLARE
    -- Declare variables to hold the UUIDs of our sample users and wallets
    user_id_alice UUID;
    user_id_bob UUID;
    user_id_charlie UUID;

    wallet_id_alice_primary UUID;
    wallet_id_alice_savings UUID;
    wallet_id_bob_main UUID;
    wallet_id_charlie_spending UUID;
BEGIN

    -- =================================================================
    --  Insert Users
    -- =================================================================
    -- Insert 3 sample users and capture their new IDs into variables.

    INSERT INTO users (name, email) VALUES ('Alice Smith', 'alice@example.com') RETURNING id INTO user_id_alice;
    INSERT INTO users (name, email) VALUES ('Bob Johnson', 'bob@example.com') RETURNING id INTO user_id_bob;
    INSERT INTO users (name, email) VALUES ('Charlie Brown', 'charlie@example.com') RETURNING id INTO user_id_charlie;

    -- =================================================================
    --  Insert Wallets
    -- =================================================================
    -- Give Alice two wallets, and Bob and Charlie one wallet each.

    -- Alice's wallets
    INSERT INTO wallets (user_id, name, balance) VALUES (user_id_alice, 'Primary', 0) RETURNING id INTO wallet_id_alice_primary;
    INSERT INTO wallets (user_id, name, balance) VALUES (user_id_alice, 'Savings', 0) RETURNING id INTO wallet_id_alice_savings;

    -- Bob's wallet
    INSERT INTO wallets (user_id, name, balance) VALUES (user_id_bob, 'Main Account', 0) RETURNING id INTO wallet_id_bob_main;

    -- Charlie's wallet
    INSERT INTO wallets (user_id, name, balance) VALUES (user_id_charlie, 'Spending', 0) RETURNING id INTO wallet_id_charlie_spending;

    -- =================================================================
    --  Create Transactions (10 samples)
    -- =================================================================
    -- Simulate a series of deposits, withdrawals, and transfers.
    -- For each transaction, we insert a record into the transactions table
    -- and update the corresponding wallet balance(s).

    -- 1. Alice deposits $1000 into her Primary wallet.
    UPDATE wallets SET balance = balance + 1000.00 WHERE id = wallet_id_alice_primary;
    INSERT INTO transactions (wallet_id, type, amount) VALUES (wallet_id_alice_primary, 'deposit', 1000.00);

    -- 2. Alice deposits $5000 into her Savings wallet.
    UPDATE wallets SET balance = balance + 5000.00 WHERE id = wallet_id_alice_savings;
    INSERT INTO transactions (wallet_id, type, amount) VALUES (wallet_id_alice_savings, 'deposit', 5000.00);

    -- 3. Bob deposits $200 into his Main Account.
    UPDATE wallets SET balance = balance + 200.00 WHERE id = wallet_id_bob_main;
    INSERT INTO transactions (wallet_id, type, amount) VALUES (wallet_id_bob_main, 'deposit', 200.00);

    -- 4. Charlie deposits $500 into his Spending wallet.
    UPDATE wallets SET balance = balance + 500.00 WHERE id = wallet_id_charlie_spending;
    INSERT INTO transactions (wallet_id, type, amount) VALUES (wallet_id_charlie_spending, 'deposit', 500.00);

    -- 5. Alice withdraws $50 from her Primary wallet.
    UPDATE wallets SET balance = balance - 50.00 WHERE id = wallet_id_alice_primary;
    INSERT INTO transactions (wallet_id, type, amount) VALUES (wallet_id_alice_primary, 'withdrawal', 50.00);

    -- 6. Alice transfers $200 from her Primary wallet to Bob's Main Account.
    UPDATE wallets SET balance = balance - 200.00 WHERE id = wallet_id_alice_primary;
    UPDATE wallets SET balance = balance + 200.00 WHERE id = wallet_id_bob_main;
    INSERT INTO transactions (wallet_id, type, amount, related_wallet_id) VALUES (wallet_id_alice_primary, 'transfer', -200.00, wallet_id_bob_main);
    INSERT INTO transactions (wallet_id, type, amount, related_wallet_id) VALUES (wallet_id_bob_main, 'transfer', 200.00, wallet_id_alice_primary);

    -- 7. Bob transfers $75 from his Main Account to Charlie's Spending wallet.
    UPDATE wallets SET balance = balance - 75.00 WHERE id = wallet_id_bob_main;
    UPDATE wallets SET balance = balance + 75.00 WHERE id = wallet_id_charlie_spending;
    INSERT INTO transactions (wallet_id, type, amount, related_wallet_id) VALUES (wallet_id_bob_main, 'transfer', -75.00, wallet_id_charlie_spending);
    INSERT INTO transactions (wallet_id, type, amount, related_wallet_id) VALUES (wallet_id_charlie_spending, 'transfer', 75.00, wallet_id_bob_main);

    -- 8. Charlie withdraws $20 from his Spending wallet.
    UPDATE wallets SET balance = balance - 20.00 WHERE id = wallet_id_charlie_spending;
    INSERT INTO transactions (wallet_id, type, amount) VALUES (wallet_id_charlie_spending, 'withdrawal', 20.00);

    -- 9. Alice transfers $1000 from Savings to Primary.
    UPDATE wallets SET balance = balance - 1000.00 WHERE id = wallet_id_alice_savings;
    UPDATE wallets SET balance = balance + 1000.00 WHERE id = wallet_id_alice_primary;
    INSERT INTO transactions (wallet_id, type, amount, related_wallet_id) VALUES (wallet_id_alice_savings, 'transfer', -1000.00, wallet_id_alice_primary);
    INSERT INTO transactions (wallet_id, type, amount, related_wallet_id) VALUES (wallet_id_alice_primary, 'transfer', 1000.00, wallet_id_alice_savings);

    -- 10. Bob deposits another $150 into his Main Account.
    UPDATE wallets SET balance = balance + 150.00 WHERE id = wallet_id_bob_main;
    INSERT INTO transactions (wallet_id, type, amount) VALUES (wallet_id_bob_main, 'deposit', 150.00);

END $$;
