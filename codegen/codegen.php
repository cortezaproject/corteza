#!/usr/bin/env php
<?php

if (count($argv) < 2) {
	echo "Usage: ./codegen [path]\n";
	exit(255);
}

if (!is_dir(__DIR__ . "/" . $argv[1])) {
	echo "No folder " . $argv[1] . "\n";
	die;
}

$project = $argv[1];

error_reporting(E_ALL^E_NOTICE);

include(__DIR__ . "/vendor/autoload.php");

function capitalize($name) {
	$names = explode("/", $name);
	return implode("", array_map("ucfirst", $names));
}

function expose($name) {
	if ($name == "id") {
		return "ID";
	}
	return ucfirst($name);
}

function array_change_key_case_recursive($arr) {
	return array_map(function ($item) {
		if (is_array($item)) {
			$item = array_change_key_case_recursive($item);
		}
		return $item;
	}, array_change_key_case($arr));
}

$tpl = new Monotek\MiniTPL\Template;
$tpl->set_compile_location("/tmp", true);
$tpl->add_default("newline", "\n");

$generators = array();
exec("find -L " . __DIR__ . "/" . $project . " -name index.php", $generators);

$api_files = glob($project . "/docs/src/spec/*.json");
$apis = array_map(function($filename) {
	$api = array_change_key_case_recursive(json_decode(file_get_contents($filename), true));
	if (empty($api['parameters'])) {
		$api['parameters'] = array();
	}
	foreach ($api['apis'] as $kk => $call) {
		if (empty($call['parameters'])) {
			$call['parameters'] = array();
		}
		$call['parameters'] = array_merge($api['parameters'], $call['parameters']);
		$api['apis'][$kk] = $call;
	}
	return $api;
}, $api_files);

usort($apis, function($a, $b) {
	return strcmp($a['interface'], $b['interface']);
});

$parsers = array(
	"uint64" => "parseUInt64",
	"bool" => "parseBool",
);

foreach ($generators as $generator) {
	$tpl->set_paths(array(dirname($generator) . "/", __DIR__ . "/templates/"));

	$dirname = strstr(dirname($generator), $project. "/");
	if (empty($dirname)) {
		$dirname = $project;
	}
	// echo "generator=". dirname($generator) . " project=$project, dirname=$dirname\n";
	if (!is_dir($dirname) && !empty($dirname)) {
		mkdir($dirname, 0777, true);
	}
	$common = compact("parsers", "project");
	include($generator);
}

/* foreach (array("structs", "handlers", "interfaces", "request", "") as $type) {
	foreach ($apis as $api) {
		if (is_array($api['struct'])) {
			$name = ucfirst($api['interface']);
			$filename = str_replace("..", ".", strtolower($name) . "." . $type . ".go");

			$tpl->load("http_$type.tpl");
			$tpl->assign("parsers", 
			$tpl->assign("package", $api['package']);
			$tpl->assign("name", $name);
			$tpl->assign("self", strtolower(substr($name, 0, 1)));
			$tpl->assign("api", $api);
			$tpl->assign("structs", $api['struct']);
			$imports = array();
			foreach ($api['struct'] as $struct) {
				if (isset($struct['imports']))
				foreach ($struct['imports'] as $import) {
					$imports[] = $import;
				}
			}
			$tpl->assign("imports", $imports);
			$tpl->assign("calls", $api['apis']);
			$contents = str_replace("\n\n}", "\n}", $tpl->get());

			$save = true;
			if ($type === "" && file_exists($filename)) {
				$save = false;
			}
			if ($save) {
				file_put_contents($filename, $contents);
			}
		}
	}
}

foreach (array("routes") as $type) {
	$name = ucfirst($api['interface']);
	$filename = str_replace("..", ".", $type . ".go");

	$tpl->load("http_$type.tpl");
	$tpl->assign("package", reset($apis)['package']);
	$tpl->assign("apis", $apis);
	$contents = $tpl->get();

	file_put_contents($filename, $contents);
}
*/

// camel case to snake case
function decamel($input) {
  preg_match_all('!([A-Z][A-Z0-9]*(?=$|[A-Z][a-z0-9])|[A-Za-z][a-z0-9]+)!', $input, $matches);
  $ret = $matches[0];
  foreach ($ret as &$match) {
    $match = $match == strtoupper($match) ? strtolower($match) : lcfirst($match);
  }
  return implode('_', $ret);
}
