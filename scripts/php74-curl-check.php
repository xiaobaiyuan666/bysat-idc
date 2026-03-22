<?php

$urls = [
    'https://www.idcsmart.com',
    'http://www.idcsmart.com',
    'https://my.idcsmart.com',
    'http://my.idcsmart.com',
    'https://w1.test.idcsmart.com',
    'http://w1.test.idcsmart.com',
];

foreach ($urls as $url) {
    $ch = curl_init($url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_TIMEOUT, 15);
    curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 10);
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
    curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, false);
    curl_exec($ch);

    echo $url, PHP_EOL;
    echo 'ERRNO=', curl_errno($ch), PHP_EOL;
    echo 'ERROR=', curl_error($ch), PHP_EOL;
    echo 'HTTP=', curl_getinfo($ch, CURLINFO_HTTP_CODE), PHP_EOL;
    echo '-----', PHP_EOL;

    curl_close($ch);
}
