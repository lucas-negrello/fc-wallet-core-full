<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('transactions', function (Blueprint $table) {
            $table->string('id', 36)->primary();
            $table->string('account_id_from', 36);
            $table->string('account_id_to', 36);
            $table->decimal('amount', 15, 2);
            $table->timestamps();

            $table->index('account_id_from');
            $table->index('account_id_to');
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('transactions');
    }
};
