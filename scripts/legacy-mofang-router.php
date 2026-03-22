<?php

$publicRoot = 'C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public';
$originalRequestUri = $_SERVER['REQUEST_URI'] ?? '/';
$requestUri = urldecode(parse_url($originalRequestUri, PHP_URL_PATH) ?? '/');
$target = realpath($publicRoot . $requestUri);
$rootRealpath = realpath($publicRoot);
$adminStaticRoot = realpath($publicRoot . '/admin01');

if (
    $requestUri !== '/' &&
    $target !== false &&
    $rootRealpath !== false &&
    strpos($target, $rootRealpath) === 0 &&
    is_file($target)
) {
    return false;
}

if (
    $requestUri !== '/' &&
    $target !== false &&
    $rootRealpath !== false &&
    strpos($target, $rootRealpath) === 0 &&
    is_dir($target)
) {
    $indexHtml = $target . DIRECTORY_SEPARATOR . 'index.html';
    $indexPhp = $target . DIRECTORY_SEPARATOR . 'index.php';
    if (is_file($indexHtml) || is_file($indexPhp)) {
        return false;
    }
}

if ($adminStaticRoot !== false) {
    $adminAliasTarget = false;
    if (
        strpos($requestUri, '/css/') === 0 ||
        strpos($requestUri, '/js/') === 0 ||
        strpos($requestUri, '/lang/') === 0 ||
        strpos($requestUri, '/img/') === 0 ||
        strpos($requestUri, '/fonts/') === 0 ||
        strpos($requestUri, '/media/') === 0
    ) {
        $adminAliasTarget = realpath($adminStaticRoot . $requestUri);
    } elseif ($requestUri === '/favicon.ico') {
        $adminAliasTarget = realpath($adminStaticRoot . '/favicon.ico');
    }

    if (
        $adminAliasTarget !== false &&
        strpos($adminAliasTarget, $adminStaticRoot) === 0 &&
        is_file($adminAliasTarget)
    ) {
        $extension = strtolower(pathinfo($adminAliasTarget, PATHINFO_EXTENSION));
        $mimeMap = [
            'css' => 'text/css; charset=UTF-8',
            'js' => 'application/javascript; charset=UTF-8',
            'json' => 'application/json; charset=UTF-8',
            'ico' => 'image/x-icon',
            'png' => 'image/png',
            'jpg' => 'image/jpeg',
            'jpeg' => 'image/jpeg',
            'svg' => 'image/svg+xml',
            'woff' => 'font/woff',
            'woff2' => 'font/woff2',
            'ttf' => 'font/ttf',
            'eot' => 'application/vnd.ms-fontobject',
            'map' => 'application/json; charset=UTF-8',
        ];
        if (isset($mimeMap[$extension])) {
            header('Content-Type: ' . $mimeMap[$extension]);
        }
        readfile($adminAliasTarget);
        exit;
    }
}

$referer = $_SERVER['HTTP_REFERER'] ?? '';
if (
    $adminStaticRoot !== false &&
    $requestUri !== '/' &&
    strpos($requestUri, '/admin01') !== 0 &&
    strpos($referer, '/admin01') !== false
) {
    $requestUri = '/admin01' . $requestUri;
    $queryString = $_SERVER['QUERY_STRING'] ?? '';
    $_SERVER['REQUEST_URI'] = $requestUri . ($queryString !== '' ? '?' . $queryString : '');
}

$_SERVER['SCRIPT_NAME'] = '/index.php';
$_SERVER['PHP_SELF'] = '/index.php';
$_SERVER['SCRIPT_FILENAME'] = $publicRoot . '/index.php';

require $publicRoot . '/index.php';
