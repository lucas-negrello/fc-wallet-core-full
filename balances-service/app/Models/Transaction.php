<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Scope;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Model;

#[Table('transactions', key: 'id', keyType: 'string', incrementing: false)]
#[Fillable(['id', 'account_id_from', 'account_id_to', 'amount'])]
class Transaction extends Model
{
    protected function casts(): array
    {
        return [
            'amount' => 'float'
        ];
    }

    #[Scope]
    protected function forAccount(Builder $query, string $account_id): Builder
    {
        return $query->where('account_id_from', $account_id)
            ->orWhere('account_id_to', $account_id)
            ->orderByDesc('created_at')
            ->limit(50);
    }
}
