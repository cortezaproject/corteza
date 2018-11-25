<?php

$cleanUp = function($contents) {
	$lines = array_map("rtrim", explode("\n", $contents));
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
};

foreach ($apis as $k => $api) {
	$entrypoint = $api['entrypoint'];
	$filename_md = $api['dirname'] . $entrypoint . "/index.md";
	if (file_exists($filename_md)) {
		$api['description'] = file_get_contents($filename_md);
	}

	foreach ($api['apis'] as $key => $call) {
		$name = $call['name'];
		$filename_md = $api['dirname'] . $entrypoint . "/" . $name . ".md";
		if (file_exists($filename_md)) {
			$call['description'] = file_get_contents($filename_md);
		}
		$api['apis'][$key] = $call;
	}

	$apis[$k] = $api;
}

$tpl->load("README.tpl");
$tpl->assign("apis", $apis);
file_put_contents($dirname . "README.md", $cleanUp($tpl->get()));
echo $dirname . "README.md\n";