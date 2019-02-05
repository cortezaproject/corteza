<?php

foreach ($apis as $api) {
	$name = ucfirst($api['interface']);
	$filename = $dirname . "/" . str_replace("..", ".", strtolower($name) . ".go");

	$tpl->load("http_.tpl");
	$tpl->assign($common);
	$tpl->assign("package", basename(__DIR__));
	$tpl->assign("name", $name);
	$tpl->assign("api", $api);
	$tpl->assign("self", strtolower(substr($name, 0, 1)));
	$tpl->assign("structs", $api['struct']);
	$imports = imports($api);
	$tpl->assign("imports", $imports);
	$tpl->assign("calls", $api['apis']);
	$contents = str_replace("\n\n}", "\n}", $tpl->get());

	if (!file_exists($filename)) {
		file_put_contents($filename, $contents);
		echo $filename . "\n";
	}
}
