CREATE TABLE analytics.bill
(
    id UUID,
    group_ip UUID,
    user_id UUID,
    created DateTime,
    total Float32
)
ENGINE = MergeTree()
PRIMARY KEY(group_id, user_id, id)
ORDER BY (id, group_id, user_id, created)
PARTITION BY toYYYYMM(created);

CREATE MATERIALIZED VIEW analytics.bill_group
ENGINE = AggregatingMergeTree() ORDER BY (group_id)
AS SELECT
          group_id,
          sumState(total)    AS total
FROM analytics.bill
GROUP BY group_id;

CREATE MATERIALIZED VIEW analytics.bill_group_yyyymm
ENGINE = AggregatingMergeTree() ORDER BY (group_id, yyyy, mm)
AS SELECT
    group_id,
    toYear(created) AS yyyy,
    toMonth(created) AS mm,
    sumState(total)    AS total
FROM analytics.bill
GROUP BY group_id, yyyy, mm;

CREATE MATERIALIZED VIEW analytics.bill_user
ENGINE = AggregatingMergeTree() ORDER BY (user_id)
AS SELECT
      user_id,
      sumState(total)    AS total
FROM analytics.bill
GROUP BY user_id;

CREATE MATERIALIZED VIEW analytics.bill_user_yyyymm
ENGINE = AggregatingMergeTree() ORDER BY (user_id, yyyy, mm)
AS SELECT
  user_id,
  toYear(created) AS yyyy,
  toMonth(created) AS mm,
  sumState(total)    AS total
FROM analytics.bill
GROUP BY user_id, yyyy, mm;