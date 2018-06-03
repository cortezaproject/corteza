#!/usr/bin/php
<?php

error_reporting(E_ALL^E_NOTICE);

include("vendor/autoload.php");

function array_change_key_case_recursive($arr) {
	return array_map(function ($item) {
		if (is_array($item)) {
			$item = array_change_key_case_recursive($item);
		}
		return $item;
	}, array_change_key_case($arr));
}

function cleanUp($contents) {
	$lines = array_map("trim", explode("\n", $contents));
	$empty = true;
	foreach ($lines as $k => $v) {
		if ($v === "") {
			if ($empty) {
				unset($lines[$k]);
			}
			$empty = true;
			continue;
		}
		$empty = false;
	}
	$contents = implode("\n", $lines);
	$contents = str_replace("`\n|", "` |", $contents);
	$contents = preg_replace("/^# /sm", "\n\n\n# ", $contents);
	return trim($contents);
}		

$spec = json_decode(file_get_contents("src/spec.json"), true);

$apis = array();
foreach ($spec as $api) {
	$entrypoint = $api['entrypoint'];
	$filename = "src/spec/" . $entrypoint . ".json";
	$filename_md = "src/" . $entrypoint . "/index.md";
	$api = json_decode(file_get_contents($filename), true);
	$api = array_change_key_case_recursive($api);
	if (file_exists($filename_md)) {
		$api['description'] = file_get_contents($filename_md);
	}
	foreach ($api['apis'] as $call) {
		$name = $call['name'];
		$filename_md = "src/" . $entrypoint . "/" . $name . ".md";
		if (file_exists($filename_md)) {
			$call['description'] = file_get_contents($filename_md);
		}
		$api['apis'][$name] = $call;
	}
	$apis[] = $api;
}

$tpl = new Monotek\MiniTPL\Template;
$tpl->set_compile_location("/tmp", true);
$tpl->set_paths("./");
$tpl->load("README.tpl");
$tpl->assign("apis", $apis);
file_put_contents("README.md", cleanUp($tpl->get()));