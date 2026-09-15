<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Model;

#[Table('balances', key: 'account_id', keyType: 'string', incrementing: false)]
#[Fillable(['account_id', 'balance'])]
class Balance extends Model
{
    protected function casts(): array
    {
        return [
            'balance' => 'float',
        ];
    }
}
