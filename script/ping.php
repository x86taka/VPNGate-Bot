<?php
ini_set("display_errors", 0);
// Configuration values --------
$vpn_host = $argv[1];
$vpn_port = (int)$argv[2];
// -----------------------------
for($i = 0; $i < 5; $i++){
    $starttime = microtime(true);
    $fp = fsockopen($vpn_host, $vpn_port, $errno, $errstr, 2);
    if (!$fp) {
        echo "Host Down";
        exit;
    }
    echo (microtime(true) - $starttime)*1000;
    echo " ms\n";
    fclose($fp);
}
?>