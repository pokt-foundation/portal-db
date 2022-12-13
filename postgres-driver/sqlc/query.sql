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
-- name: SelectPayPlans :many
SELECT plan_type,
    daily_limit
FROM pay_plans
ORDER BY plan_type ASC;
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
-- name: SelectApplications :many
SELECT a.application_id,
    a.contact_email,
    a.description,
    a.dummy,
    a.name,
    a.owner,
    a.status,
    a.url,
    a.user_id,
    a.first_date_surpassed,
    ga.address AS ga_address,
    ga.client_public_key AS ga_client_public_key,
    ga.private_key AS ga_private_key,
    ga.public_key AS ga_public_key,
    ga.signature AS ga_signature,
    ga.version AS ga_version,
    gs.secret_key,
    gs.secret_key_required,
    gs.whitelist_blockchains,
    gs.whitelist_contracts,
    gs.whitelist_methods,
    gs.whitelist_origins,
    gs.whitelist_user_agents,
    ns.signed_up,
    ns.on_quarter,
    ns.on_half,
    ns.on_three_quarters,
    ns.on_full,
    al.custom_limit,
    al.pay_plan,
    pp.daily_limit as plan_limit,
    a.created_at,
    a.updated_at
FROM applications AS a
    LEFT JOIN gateway_aat AS ga ON a.application_id = ga.application_id
    LEFT JOIN gateway_settings AS gs ON a.application_id = gs.application_id
    LEFT JOIN notification_settings AS ns ON a.application_id = ns.application_id
    LEFT JOIN app_limits AS al ON a.application_id = al.application_id
    LEFT JOIN pay_plans AS pp ON al.pay_plan = pp.plan_type
ORDER BY a.application_id ASC;
-- name: SelectAppLimit :one
SELECT application_id,
    pay_plan,
    custom_limit
FROM app_limits
WHERE application_id = $1;
-- name: SelectGatewaySettings :one
SELECT application_id,
    secret_key,
    secret_key_required,
    whitelist_blockchains,
    whitelist_contracts,
    whitelist_methods,
    whitelist_origins,
    whitelist_user_agents
FROM gateway_settings
WHERE application_id = $1;
-- name: SelectNotificationSettings :one
SELECT application_id,
    signed_up,
    on_quarter,
    on_half,
    on_three_quarters,
    on_full
FROM notification_settings
WHERE application_id = $1;
-- name: InsertApplication :exec
INSERT into applications (
        application_id,
        user_id,
        name,
        contact_email,
        description,
        owner,
        url,
        status,
        dummy
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
        $9
    );
-- name: InsertAppLimit :exec
INSERT into app_limits (application_id, pay_plan, custom_limit)
VALUES ($1, $2, $3);
-- name: InsertGatewayAAT :exec
INSERT into gateway_aat (
        application_id,
        address,
        client_public_key,
        private_key,
        public_key,
        signature,
        version
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7
    );
-- name: InsertGatewaySettings :exec
INSERT into gateway_settings (
        application_id,
        secret_key,
        secret_key_required,
        whitelist_contracts,
        whitelist_methods,
        whitelist_origins,
        whitelist_user_agents,
        whitelist_blockchains
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    );
-- name: InsertNotificationSettings :exec
INSERT into notification_settings (
        application_id,
        signed_up,
        on_quarter,
        on_half,
        on_three_quarters,
        on_full
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    );
-- name: UpsertApplication :exec
INSERT INTO applications (
        application_id,
        name,
        status,
        first_date_surpassed
    )
VALUES ($1, $2, $3, $4) ON CONFLICT (application_id) DO
UPDATE
SET name = EXCLUDED.name,
    status = EXCLUDED.status,
    first_date_surpassed = EXCLUDED.first_date_surpassed;
-- name: UpsertAppLimit :exec
INSERT INTO app_limits (
        application_id,
        pay_plan,
        custom_limit
    )
VALUES ($1, $2, $3) ON CONFLICT (application_id) DO
UPDATE
SET pay_plan = EXCLUDED.pay_plan,
    custom_limit = EXCLUDED.custom_limit;
-- name: UpsertGatewaySettings :exec
INSERT INTO gateway_settings (
        application_id,
        secret_key,
        secret_key_required,
        whitelist_contracts,
        whitelist_methods,
        whitelist_origins,
        whitelist_user_agents,
        whitelist_blockchains
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (application_id) DO
UPDATE
SET secret_key = EXCLUDED.secret_key,
    secret_key_required = EXCLUDED.secret_key_required,
    whitelist_contracts = EXCLUDED.whitelist_contracts,
    whitelist_methods = EXCLUDED.whitelist_methods,
    whitelist_origins = EXCLUDED.whitelist_origins,
    whitelist_user_agents = EXCLUDED.whitelist_user_agents,
    whitelist_blockchains = EXCLUDED.whitelist_blockchains;
-- name: UpsertNotificationSettings :exec
INSERT INTO notification_settings (
        application_id,
        signed_up,
        on_quarter,
        on_half,
        on_three_quarters,
        on_full
    )
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (application_id) DO
UPDATE
SET signed_up = EXCLUDED.signed_up,
    on_quarter = EXCLUDED.on_quarter,
    on_half = EXCLUDED.on_half,
    on_three_quarters = EXCLUDED.on_three_quarters,
    on_full = EXCLUDED.on_full;
-- name: UpdateFirstDateSurpassed :exec
UPDATE applications
SET first_date_surpassed = @first_date_surpassed
WHERE application_id IN (@application_ids::VARCHAR []);
-- name: RemoveApplication :exec
UPDATE applications
SET status = COALESCE($2, status)
WHERE application_id = $1;
-- name: SelectLoadBalancers :many
SELECT lb.lb_id,
    lb.name,
    lb.created_at,
    lb.updated_at,
    lb.request_timeout,
    lb.gigastake,
    lb.gigastake_redirect,
    lb.user_id,
    so.duration,
    so.sticky_max,
    so.stickiness,
    so.origins,
    STRING_AGG(la.app_id, ',') AS app_ids
FROM loadbalancers AS lb
    LEFT JOIN stickiness_options AS so ON lb.lb_id = so.lb_id
    LEFT JOIN lb_apps AS la ON lb.lb_id = la.lb_id
GROUP BY lb.lb_id,
    lb.lb_id,
    lb.name,
    lb.created_at,
    lb.updated_at,
    lb.request_timeout,
    lb.gigastake,
    lb.gigastake_redirect,
    lb.user_id,
    so.duration,
    so.sticky_max,
    so.stickiness,
    so.origins;
-- name: InsertLoadBalancer :exec
INSERT into loadbalancers (
        lb_id,
        name,
        user_id,
        request_timeout,
        gigastake,
        gigastake_redirect
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    );
-- name: InsertStickinessOptions :exec
INSERT INTO stickiness_options (
        lb_id,
        duration,
        sticky_max,
        stickiness,
        origins
    )
VALUES ($1, $2, $3, $4, $5);
-- name: UpsertStickinessOptions :exec
INSERT INTO stickiness_options (
        lb_id,
        duration,
        sticky_max,
        stickiness,
        origins
    )
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (lb_id) DO
UPDATE
SET duration = EXCLUDED.duration,
    sticky_max = EXCLUDED.sticky_max,
    stickiness = EXCLUDED.stickiness,
    origins = EXCLUDED.origins;
-- name: InsertLbApps :exec
INSERT into lb_apps (lb_id, app_id)
SELECT @lb_id,
    unnest(@app_ids::VARCHAR []);
-- name: UpdateLB :exec
UPDATE loadbalancers
SET name = $2
WHERE lb_id = $1;
-- name: RemoveLB :exec
UPDATE loadbalancers
SET user_id = ''
WHERE lb_id = $1;
