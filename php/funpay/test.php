<?php

use FpDbTest\Database;
use FpDbTest\DatabaseTest;

try
{
    spl_autoload_register(function ($class)
    {
        $a = array_slice(explode('\\', $class), 1);
        if (!$a) {
            throw new Exception("e001");
        }
        $filename = implode('/', [__DIR__, ...$a]) . '.php';
        require_once $filename;
    });

//    $mysqli = @new mysqli('localhost', 'root', 'password', 'database', 3306);
    $mysqli = @new mysqli('localhost', 'root', '', 'database', 3306);
    if ($mysqli->connect_errno) {
        throw new Exception("e002".$mysqli->connect_error);
    }


    $db = new Database($mysqli);
    $test = new DatabaseTest($db);
    $test->testBuildQuery();
}
catch (Exception $e)
{
    exit("ERROR: ".$e->getMessage()."\n");
}

exit("OK\n");
