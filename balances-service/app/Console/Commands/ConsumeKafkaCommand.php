<?php

namespace App\Console\Commands;

use App\Kafka\WalletEventHandler;
use Illuminate\Console\Attributes\Description;
use Illuminate\Console\Attributes\Signature;
use Illuminate\Console\Command;
use RdKafka\KafkaConsumer;
use RdKafka\Message;
use RdKafka\Conf;
use Throwable;

#[Signature('kafka:consume')]
#[Description('Consume Wallet Core events and keep the database in sync')]
class ConsumeKafkaCommand extends Command
{
    private bool $shouldStop = false;
    /**
     * Execute the console command.
     */
    public function handle(WalletEventHandler $handler): int
    {
        $this->listenForShutdownSignals();

        $topics = array_values(config('kafka.topics'));

        $consumer = new KafkaConsumer($this->kafkaConf());
        $consumer->subscribe($topics);

        $this->info(sprintf(
            'Consuming [%s] on %s (group.id: %s)',
            implode(', ', $topics),
            config('kafka.brokers'),
            config('kafka.group_id')
        ));

        while (!$this->shouldStop) {
            $message = $consumer->consume(config('kafka.consume_timeout'));

            match ($message->err) {
                RD_KAFKA_RESP_ERR_NO_ERROR => $this->process($consumer, $handler, $message),

                RD_KAFKA_RESP_ERR__PARTITION_EOF,
                RD_KAFKA_RESP_ERR__TIMED_OUT => null,

                default => $this->error("Kafka: {$message->errstr()} ({$message->err})"),
            };
        }

        $consumer->close();
        $this->info('Kafka consumer finished.');

        return self::SUCCESS;
    }

    private function process(KafkaConsumer $consumer, WalletEventHandler $handler, Message $message): void
    {
        $event = json_decode($message->payload, true);
        if (! is_array($event)) {
            $this->warn('JSON payload is not valid, skipping message');
            $consumer->commit($message);
            return;
        }

        try {
            $handler->handle($message->topic_name, $event);
            $consumer->commit($message);
            $this->line(sprintf(
                '[%s] offset %d - %s',
                $message->topic_name,
                $message->offset,
                $event['Name'] ?? 'Unknown event'
            ));
        } catch (Throwable $e) {
            $this->error("Failed to proccess offset {$message->offset}: {$e->getMessage()}");
        }
    }

    private function kafkaConf(): Conf
    {
        $conf = new Conf;
        $conf->set('metadata.broker.list', config('kafka.brokers'));
        $conf->set('group.id', config('kafka.group_id'));
        $conf->set('auto.offset.reset', config('kafka.auto_offset_reset'));
        $conf->set('enable.auto.commit', 'false');
        $conf->set('enable.partition.eof', 'true');
        return $conf;
    }

    private function listenForShutdownSignals(): void
    {
        if (! function_exists('pcntl_async_signals')) {
            return;
        }

        pcntl_async_signals(true);

        foreach ([SIGTERM, SIGINT] as $signal) {
            pcntl_signal($signal, function () {
                $this->shouldStop = true;
                $this->info('Shutdown signal received, stopping consumer...');
            });
        }
    }
}
