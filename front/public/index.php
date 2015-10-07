<?php

require '../vendor/autoload.php';
$app = new \Slim\Slim();

$app->get('/', function() use($app) {
    $app->response->setStatus(200);
    echo "Welcome to Slim";
});

$app->run();
