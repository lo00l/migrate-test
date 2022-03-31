CREATE TABLE test_v1
(
    `account_id` UInt64,
    `ts` DateTime DEFAULT now()
) ENGINE = MergeTree()
    ORDER BY account_id