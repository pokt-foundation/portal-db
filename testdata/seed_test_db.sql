INSERT INTO pay_plans (plan_type, daily_limit)
VALUES ('FREETIER_V0', 250000),
    ('PAY_AS_YOU_GO_V0', 0),
    ('ENTERPRISE', 0),
    ('TEST_PLAN_V0', 100),
    ('TEST_PLAN_10K', 10000),
    ('TEST_PLAN_90K', 90000);
INSERT INTO applications (
        application_id,
        name,
        status,
        url,
        user_id,
        dummy
    )
VALUES (
        'test_app_47hfnths73j2se',
        'pokt_app_123',
        'IN_SERVICE',
        'https://test.app123.io',
        'test_user_1dbffbdfeeb225',
        true
    ),
    (
        'test_app_5hdf7sh23jd828',
        'pokt_app_456',
        'IN_SERVICE',
        'https://test.app456.io',
        'test_user_04228205bd261a',
        true
    );
INSERT INTO app_limits (
        application_id,
        pay_plan,
        custom_limit
    )
VALUES (
        'test_app_47hfnths73j2se',
        'FREETIER_V0',
        null
    ),
    (
        'test_app_5hdf7sh23jd828',
        'ENTERPRISE',
        2000000
    );
INSERT INTO gateway_aat (
        application_id,
        address,
        public_key,
        signature,
        client_public_key,
        private_key
    )
VALUES (
        'test_app_47hfnths73j2se',
        'test_34715cae753e67c75fbb340442e7de8e',
        'test_11b8d394ca331d7c7a71ca1896d630f6',
        'test_89a3af6a587aec02cfade6f5000424c2',
        'test_1dc39a2e5a84a35bf030969a0b3231f7',
        'test_d2ce53f115f4ecb2208e9188800a85cf'
    ),
    (
        'test_app_5hdf7sh23jd828',
        'test_558c0225c7019e14ccf2e7379ad3eb50',
        'test_96c981db344ab6920b7e87853838e285',
        'test_1272a8ab4cbbf636f09bf4fa5395b885',
        'test_d709871777b89ed3051190f229ea3f01',
        'test_53e50765d8bc1fb41b3b0065dd8094de'
    );
INSERT INTO gateway_settings (
        application_id,
        secret_key,
        secret_key_required
    )
VALUES (
        'test_app_47hfnths73j2se',
        'test_40f482d91a5ef2300ebb4e2308c',
        true
    ),
    (
        'test_app_5hdf7sh23jd828',
        'test_90210ac4bdd3423e24877d1ff92',
        false
    );
INSERT INTO notification_settings (
        application_id,
        signed_up,
        on_quarter,
        on_half,
        on_three_quarters,
        on_full
    )
VALUES (
        'test_app_47hfnths73j2se',
        true,
        false,
        false,
        true,
        true
    ),
    (
        'test_app_5hdf7sh23jd828',
        true,
        false,
        false,
        true,
        true
    );
-- INSERT INTO lb_apps (lb_id, app_id)
-- VALUES (
--         'test_lb_34987u329rfn23f',
--         'test_app_47hfnths73j2se'
--     ),
--     (
--         'test_lb_3890ru23jfi32fj',
--         'test_app_5hdf7sh23jd828'
--     );
