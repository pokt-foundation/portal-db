-- name: SelectBlockchains :many
SELECT b.blockchain_id,
    b.altruist,
    b.blockchain,
    b.blockchain_aliases,
    b.chain_id,
    b.chain_id_check,
    b.description,
    b.enforce_result,
    b.log_limit_blocks,
    b.network,
    b.path,
    b.request_timeout,
    b.ticker,
    b.active,
    s.synccheck as s_sync_check,
    s.allowance as s_allowance,
    s.body as s_body,
    s.path as s_path,
    s.result_key as s_result_key,
    COALESCE(redirects.r, '[]') AS redirects
FROM blockchains as b
    LEFT JOIN sync_check_options AS s ON b.blockchain_id = s.blockchain_id
    LEFT JOIN LATERAL (
        SELECT json_agg(
                json_build_object(
                    'alias',
                    r.alias,
                    'loadBalancerID',
                    r.loadbalancer,
                    'domain',
                    eg.domain
                )
            ) AS r
        FROM redirects AS r
        WHERE b.blockchain_id = r.blockchain_id
    ) redirects ON true
ORDER BY b.blockchain_id;
-- name: InsertBlockchain :exec
INSERT into blockchains (
        blockchain_id,
        active,
        altruist,
        blockchain,
        blockchain_aliases,
        chain_id,
        chain_id_check,
        description,
        enforce_result,
        log_limit_blocks,
        network,
        path,
        request_timeout,
        ticker
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13,
        $14
    );
-- name: InsertRedirect :exec
INSERT into redirects (
        blockchain_id,
        alias,
        loadbalancer,
        domain
    )
VALUES (
        $1,
        $2,
        $3,
        $4
    );
-- name: InsertSyncCheckOptions :exec
INSERT into sync_check_options (
        blockchain_id,
        synccheck,
        allowance,
        body,
        path,
        result_key
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    );
-- name: ActivateBlockchain :exec
UPDATE blockchains
SET active = $2
WHERE blockchain_id = $1;
