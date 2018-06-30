#!/usr/bin/env php
<?php

error_reporting(E_ALL^E_NOTICE);

include("docs/vendor/autoload.php");

$tpl = new Monotek\MiniTPL\Template;
$tpl->set_compile_location("/tmp", true);
$tpl->add_default("newline", "\n");

$apis = json_decode(file_get_contents("docs/src/spec.json"), true);
foreach ($apis as $k => $v) {
	foreach ($v['apis'] as $kk => $vv) {
		$vv['path'] = "/" . $vv['name'];
		$v['apis'][$kk] = $vv;
	}
	$v['path'] = "/" . $v['entrypoint'];
	$apis[$k] = $v;
}

foreach (array("structs", "handlers", "interfaces", "request", "") as $type) {
	foreach ($apis as $api) {
		if (!empty($api['struct'])) {
			$name = ucfirst($api['entrypoint']);
			$filename = str_replace("..", ".", strtolower($name) . "." . $type . ".go");

			$tpl->load("http_$type.tpl");
			$tpl->assign("parsers", array(
				"uint64" => "parseUInt64"
			));
			$tpl->assign("package", $api['package']);
			$tpl->assign("name", $name);
			$tpl->assign("self", strtolower(substr($name, 0, 1)));
			$tpl->assign("api", $api);
			$tpl->assign("fields", $api['struct']);
			$tpl->assign("calls", $api['apis']);
			$contents = $tpl->get();

			file_put_contents($filename, $contents);
		}
	}
}

foreach (array("routes") as $type) {
	$name = ucfirst($api['entrypoint']);
	$filename = str_replace("..", ".", $type . ".go");

	$tpl->load("http_$type.tpl");
	$tpl->assign("package", reset($apis)['package']);
	$tpl->assign("apis", $apis);
	$contents = $tpl->get();

	file_put_contents($filename, $contents);
}

passthru("go fmt");