<?php

namespace FpDbTest;

use Exception;
use mysqli;

class Database implements DatabaseInterface
{
    private mysqli $mysqli;

    public function __construct(mysqli $mysqli)
    {
        $this->mysqli = $mysqli;
    }

    public function buildQuery(string $query, array $args = []): string
    {
        $queryMade = "";

// method for replace specifier
        $specifierReplace = function($specifier, $args, $argIndex): string {

// method for draw value
            $drawValue = function($rune, $flagEscape = false) use (&$drawValue) : string {

// check null
                if (is_null($rune) === true) {
                    return "NULL";
                }

// check bool
                if (is_bool($rune) === true) {
                    if ($rune === false) {
                        return "0";
                    }
                    return "1";
                }

// check string
                if (is_string($rune) === true) {
                    if ($flagEscape === true) {
                        return "'".$rune."'";
                    }
                    return "`".$rune."`";
                }

// check int
                if (is_int($rune) === true) {
                    settype($rune, "string");
                    return $rune;
                }

// check float
                if (is_float($rune) === true) {
                    settype($rune, "string");
                    return $rune;
                }

// check array
                if (is_array($rune) === true) {
                    $tmp = "";
                    $index = 0;
                    foreach ($rune as $key => $value) {
                        if ($index !== 0) {
                            $tmp .= ", ";
                        }

                        if (is_integer($key) === false) {
                            $tmp .= $drawValue($key);
                            $tmp .= " = ";
                            $tmp .= $drawValue($value, true);
                        } else {
                            $tmp .= $drawValue($value, false);
                        }
                        $index++;
                    }

                    return $tmp;
                }

// check object (special type for skip method)
                if (is_object($rune) === true) {
                    return "";
                }

                throw new Exception("arg type is not implemented \"".gettype($rune)."\", \"".print_r($rune, true)."\"");
            };

            if ($argIndex >= count($args)) {
                throw new Exception("Invalid args");
            }
            $rune = $args[$argIndex];

            if ($specifier == "?") {
                return $drawValue($rune, true);
            }

            if ($specifier == "?d") {
                return $drawValue($rune);
            }

            if ($specifier == "?f") {
                return $drawValue($rune);
            }

            if ($specifier == "?a") {
                if (is_array($rune) === false) {
                    throw new Exception("arg is not array \"".print_r($rune, true)."\"");
                }
                return $drawValue($rune);
            }

            if ($specifier == "?#") {
                return $drawValue($rune);
            }

            throw new Exception("Invalid specifier \"".$specifier."\"");
        };

// main logic
        $block = "";
        $flagBlock = false;
        $flagBlockOk = true;
        $argIndex = 0;
        $len = strlen($query);

        for ($i=0; $i < $len; $i++) {
            $rune = $query[$i];

            if ($rune === '?') {
                $specifier='?';
                if (($i + 1) < $len) {
                    $runeNext = $query[$i + 1];
                    for (;;) {
                        if (($runeNext === " ") || ($runeNext === "\n") || ($runeNext === "\r")) {
                            break;
                        }
                        $specifier .= $runeNext;
                        $i++;
                        break;
                    }
                }
                $rune = $specifierReplace($specifier, $args, $argIndex);
                $argIndex++;

                if ($flagBlock === true) {
                    if ($rune === "") {
                        $flagBlockOk = false;
                    }
                }
            }


            if ($rune === '{') {
                $flagBlock = true;
                continue;
            }


            if ($rune === '}') {
                if ($flagBlockOk === true) {
                    $queryMade .= $block;
                }

                $block = "";
                $flagBlock = false;
                $flagBlockOk = true;
                continue;
            }


            if ($flagBlock === false) {
                $queryMade .= $rune;
            } else {
                $block .= $rune;
            }
        }

        return $queryMade;
    }

    public function skip()
    {
        return (object)[];
    }
}
