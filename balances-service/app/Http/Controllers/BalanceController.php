<?php

namespace App\Http\Controllers;

use App\Models\Balance;
use App\Models\Transaction;
use Illuminate\Http\JsonResponse;

class BalanceController extends Controller
{
    public function show(Balance $balance): JsonResponse
    {
        return response()->json([
            'account_id' => $balance->account_id,
            'balance' => $balance->balance,
            'updated_at' => $balance->updated_at,
        ]);
    }

    public function transactions(Balance $balance): JsonResponse
    {
        $accountId = $balance->account_id;

        $transactions = Transaction::forAccount($accountId)
            ->get();

        return response()->json([
            'account_id' => $accountId,
            'transactions' => $transactions
        ]);
    }
}
