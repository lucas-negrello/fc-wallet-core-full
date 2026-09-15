<?php

namespace App\Kafka;

use App\Models\Balance;
use App\Models\Transaction;
use Illuminate\Support\Facades\Log;

class WalletEventHandler
{
    public function handle(string $topic, array $event): void
    {
        $name = $event['Name'] ?? null;
        $payload = $event['Payload'] ?? null;

        if (!is_array($payload)) {
            Log::warning('Event without usable payload', compact('topic', 'event'));
            return;
        }

        match ($name) {
            'BalanceUpdated' => $this->balanceUpdated($payload),
            'TransactionCreated' => $this->transactionCreated($payload),
            default => Log::warning('Unknown event', compact('topic', 'name')),
        };
    }

    private function balanceUpdated(array $payload): void
    {
        $this->upsertBalance(
            $payload['account_id_from'] ?? null,
            $payload['balance_account_id_from'] ?? null,
        );

        $this->upsertBalance(
            $payload['account_id_to'] ?? null,
            $payload['balance_account_id_to'] ?? null,
        );
    }

    private function upsertBalance(?string $accountId, int|float|null $balance): void
    {
        if ($accountId === null || $balance === null) {
            return;
        }

        Balance::query()->updateOrCreate(
            ['account_id' => $accountId],
            ['balance' => $balance]
        );
    }

    private function transactionCreated(array $payload): void
    {
        if (!isset($payload['id'])) {
            Log::warning('TransactionCreated event without id', compact('payload'));
            return;
        }

        Transaction::query()->updateOrCreate(
            ['id' => $payload['id']],
            [
                'account_id_from' => $payload['account_id_from'] ?? null,
                'account_id_to' => $payload['account_id_to'] ?? null,
                'amount' => $payload['amount'] ?? 0,
            ],
        );
    }
}
