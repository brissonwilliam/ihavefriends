CREATE TABLE ihavefriends.bill
(
    id UUID,
    group_id UUID,
    user_id UUID,
    created DateTime,
    total Float32
)
ENGINE = MergeTree()
PRIMARY KEY(group_id, user_id, created)
ORDER BY (group_id, user_id, created, id)
PARTITION BY toYYYYMM(created);

CREATE MATERIALIZED VIEW ihavefriends.bill_group
ENGINE = AggregatingMergeTree() ORDER BY (group_id)
AS SELECT
          group_id,
          sumState(total)    AS total
FROM ihavefriends.bill
GROUP BY group_id;

CREATE MATERIALIZED VIEW ihavefriends.bill_group_yyyymm
ENGINE = AggregatingMergeTree() ORDER BY (group_id, yyyy, mm)
AS SELECT
    group_id,
    toYear(created) AS yyyy,
    toMonth(created) AS mm,
    sumState(total)    AS total
FROM ihavefriends.bill
GROUP BY group_id, yyyy, mm;

CREATE MATERIALIZED VIEW ihavefriends.bill_user
ENGINE = AggregatingMergeTree() ORDER BY (user_id)
AS SELECT
      user_id,
      sumState(total)    AS total
FROM ihavefriends.bill
GROUP BY user_id;

CREATE MATERIALIZED VIEW ihavefriends.bill_user_yyyymm
ENGINE = AggregatingMergeTree() ORDER BY (user_id, yyyy, mm)
AS SELECT
  user_id,
  toYear(created) AS yyyy,
  toMonth(created) AS mm,
  sumState(total)    AS total
FROM ihavefriends.bill
GROUP BY user_id, yyyy, mm;