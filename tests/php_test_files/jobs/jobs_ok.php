<?php

/**
 * @var Goridge\RelayInterface $relay
 */

ini_set('display_errors', 'stderr');
require dirname(__DIR__) . "/vendor/autoload.php";

$consumer = new Spiral\RoadRunner\Jobs\Consumer();

while ($task = $consumer->waitTask()) {
	try {
		$task->ack();
	} catch (\Throwable $e) {
		$task->error((string)$e);
	}
}
