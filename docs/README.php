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
	return $contents;
}		

$spec = json_decode(file_get_contents("src/_spec.json"), true);

$apis = array();
foreach ($spec as $api) {
	$entrypoint = $api['entrypoint'];
	$filename = "src/" . $entrypoint . ".json";
	$filenames_md = array("src/" . $entrypoint . ".md", "src/" . $entrypoint . "/index.md");
	$api = json_decode(file_get_contents($filename), true);
	$api = array_change_key_case_recursive($api);
	foreach ($filenames_md as $filename_md) {
		if (file_exists($filename_md)) {
			$api['description'] = file_get_contents($filename_md);
			break;
		}
	}
	foreach ($api['apis'] as $name => $call) {
		$filename_md = "src/" . $entrypoint . "/" . $name . ".md";
		if (file_exists($filename_md)) {
			$call['description'] = file_get_contents($filename_md);
		}
		$api['apis'][$name] = $call;
	}
	$apis[] = $api;
}

$tpl = new Monotek\MiniTPL\Template;
$tpl->set_paths("./");
$tpl->load("README.tpl");
$tpl->assign("apis", $apis);
file_put_contents("README.md", cleanUp($tpl->get()));