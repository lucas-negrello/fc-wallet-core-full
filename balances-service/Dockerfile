FROM dunglas/frankenphp:php8.5

WORKDIR /app

COPY --from=mlocati/php-extension-installer /usr/bin/install-php-extensions /usr/local/bin/

RUN install-php-extensions pdo_mysql pcntl rdkafka zip opcache

RUN apt-get update && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

COPY composer.json composer.lock ./
RUN composer install --no-dev --no-interaction --no-scripts --no-autoloader --prefer-dist

COPY . .

RUN composer dump-autoload --optimize --no-dev \
    && cp .env.example .env \
    && php artisan package:discover --ansi \
    && chmod -R 775 storage bootstrap/cache

ENV SERVER_NAME=:3003

EXPOSE 3003
