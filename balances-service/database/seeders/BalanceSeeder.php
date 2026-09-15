<?php

namespace Database\Seeders;

use App\Models\Balance;
use Illuminate\Database\Seeder;

class BalanceSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        if (Balance::query()->exists()) {
            $this->command->info('Balances already seeded. Skipping BalanceSeeder.');
            return;
        }

        Balance::query()->insert([
            [
                'account_id' => 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
                'balance' => 1000.00,
                'created_at' => now(),
                'updated_at' => now(),
            ],
            [
                'account_id' => 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
                'balance' => 500.00,
                'created_at' => now(),
                'updated_at' => now(),
            ]
        ]);

        $this->command->info('Balances seeded. (2 records)');
    }
}
