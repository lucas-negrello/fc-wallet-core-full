<?php

use App\Http\Controllers\BalanceController;
use Illuminate\Support\Facades\Route;

Route::group([
    'prefix' => 'balances',
    'as' => 'balances.',
    'missing' => fn () => response()->json(['error' => 'Account not found'], 404)
], function () {

    Route::group([
        'prefix' => '{balance}',
        'as' => 'balance.',
    ], function () {

        Route::get('/', [BalanceController::class, 'show'])
            ->name('show');
        Route::get('/transactions', [BalanceController::class, 'transactions'])
            ->name('transactions');

    });
});
