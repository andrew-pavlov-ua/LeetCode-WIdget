-- name: InsertStatsInfo :exec
INSERT INTO lc_stats (lc_user_id, username, easy_submits, medium_submits, hard_submits, total_submits, rank, created_at, updated_at)
       VALUES (@lc_user_id, @username, @easy_submits, @medium_submits, @hard_submits, @total_submits, @rank, @created_at, @updated_at);

-- name: UserGetStatsBySlug :one
SELECT lc_user_id AS userSlug,
       username,
       easy_submits,
       medium_submits,
       hard_submits,
       total_submits,
       rank,
       updated_at
FROM lc_stats
WHERE lc_user_id = @lc_user_id;

-- name: UpdateLcStats :exec
UPDATE lc_stats
SET easy_submits = @easy_submits,
    medium_submits = @medium_submits,
    hard_submits = @hard_submits,
    total_submits = @total_submits,
    updated_at = @updated_at
WHERE lc_user_id = @lc_user_id;
