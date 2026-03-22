<?php
header('Content-Type: application/json; charset=utf-8');
echo json_encode([
    'php_version' => PHP_VERSION,
    'loaded_ini' => php_ini_loaded_file(),
    'scanned_ini' => php_ini_scanned_files(),
    'gd_loaded' => extension_loaded('gd'),
    'freetype' => function_exists('gd_info') ? (gd_info()['FreeType Support'] ?? null) : null,
    'extensions' => get_loaded_extensions(),
], JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);
