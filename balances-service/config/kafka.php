<?php

return [
    'brokers' => env('KAFKA_BOOTSTRAP_SERVERS', 'kafka:29092'),

    'group_id' => env('KAFKA_GROUP_ID', 'balances-service'),

    'auto_offset_reset' => env('KAFKA_AUTO_OFFSET_RESET', 'earliest'),

    'consume_timeout_ms' => env('KAFKA_CONSUME_TIMEOUT_MS', 10000),

    'topics' => [
        'balances' => env('KAFKA_TOPIC_BALANCES', 'balances'),
        'transactions' => env('KAFKA_TOPIC_TRANSACTIONS', 'transactions'),
    ],
];
