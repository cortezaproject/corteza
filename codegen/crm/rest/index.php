<?php

$templates = array(
	"http_interfaces.tpl" => function($name, $api) {
		return strtolower($name) . ".go";
	},
	"http_request.tpl" => function($name, $api) {
		return strtolower($name) . "_requests.go"; 
	},
	"http_handlers.tpl" => function($name, $api) {
		return strtolower($name) . "_handlers.go";
	}
);

foreach ($templates as $template => $fn)
foreach ($apis as $api) {
	if (is_array($api['struct'])) {
		$name = ucfirst($api['interface']);
		$filename = $dirname . "/" . $fn($name, $api);

		$tpl->load($template);
		$tpl->assign($common);
		$tpl->assign("package", basename(__DIR__));
		$tpl->assign("name", $name);
		$tpl->assign("api", $api);
		$tpl->assign("self", strtolower(substr($name, 0, 1)));
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

		file_put_contents($filename, $contents);
		echo $filename . "\n";
	}
}
