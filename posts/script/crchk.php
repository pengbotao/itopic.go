<?php

$folder = "C:\www\asm\asm_admin\application";

scan($folder, 'chk_cr');

function scan($folder, $callback)
{
    if(! is_dir($folder)) {
        return false;
    }
    $files = scandir($folder);
    foreach($files as $val) {
        if($val == '.' || $val == '..') {
            continue;
        }
        $tmp_folder = trim($folder, "/").'/'.$val;
        if(is_dir($tmp_folder)) {
            scan($tmp_folder, $callback);
        } else {
            $filename = $tmp_folder;
            if(substr($filename, -4) == ".php") {
                call_user_func($callback, $filename);
            }
        }
    }
}

function chk_cr($filename)
{
    $fp = fopen($filename, "r");
    while(false !== ($char = fgetc($fp))) {
        if($char == "\t") {
            echo $filename . "\n";
            break;
        }
    }
}