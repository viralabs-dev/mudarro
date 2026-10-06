<?php
require __DIR__ . '/../src/value.php';
if (sample_value() !== 42) { fwrite(STDERR, 'unexpected value'); exit(1); }
echo "PHP_SAMPLE_OK\n";
