<?php

$templates = array(
	"http_handlers_inline.tpl" => function($name, $api) {
		return strtolower($name) . ".go";
	},
);

foreach ($templates as $template => $fn)
foreach ($apis as $api) {
		$name = ucfirst($api['interface']);
		$filename = $dirname . "/" . $fn($name, $api);

		$tpl->load($template);
		$tpl->assign($common);
		$tpl->assign("package", basename(__DIR__));
		$tpl->assign("name", $name);
		$tpl->assign("api", $api);
		$tpl->assign("apis", $apis);
		$tpl->assign("self", strtolower(substr($name, 0, 1)));
		$tpl->assign("structs", $api['struct']);
		$imports = imports($api);
		$tpl->assign("imports", $imports);
		$tpl->assign("calls", $api['apis']);
		$contents = str_replace("\n\n}", "\n}", $tpl->get());

		file_put_contents($filename, $contents);
		echo $filename . "\n";
}
