<?php

class TemplateTest extends PHPUnit_Framework_TestCase
{
	/**
	 * @dataProvider compileProvider
	 */
	public function testCompile($template)
	{
		global $tpl;

		$destination = "test/compile/".$template;
		$compiled = "test/compiled/".$template;

		$tpl = new Monotek\MiniTPL\Template;

		$tpl->set_paths("test/templates/");
		$tpl->set_compile_location("test/compile/", false);
		$tpl->add_default("key", "val");

		$source = "test/templates/".$template;
		$return = $tpl->compile($source, $destination);

		$this->assertTrue((bool)$return);
		$this->assertFileEquals($destination, $compiled);
	}

	public function compileProvider()
	{
		$tests = array();
		$templates = glob("test/templates/*.tpl");
		sort($templates);
		foreach ($templates as $template) {
			$tests[] = array(basename($template));
		}
		return $tests;
	}

	public function testRendering()
	{
		$tpl = new Monotek\MiniTPL\Template;

		$tpl->set_paths("test/templates/");
		$tpl->set_compile_location("test/compile/", false);
		$tpl->add_default("key", "val");

		$retval = array();
		exec("rm -rf test/templates/test -rf", $retval);

		$this->assertTrue($tpl->load("08_utf8_bom.tpl"));

		touch("test/templates/test/compile/08_utf8_bom.tpl", filemtime("test/templates/08_utf8_bom.tpl") - 86400);

		$this->assertTrue($tpl->load("08_utf8_bom.tpl"));

		$items = array();
		$items[] = array("id" => $i++);
		$items[] = array("id" => $i++);
		$items[] = array("id" => $i++);
		$tpl->assign("items", $items);
		$tpl->assign(array("foo"=>"bar", "d" => array("burger")), "foo");
		$tpl->assign(".foo_foo", "baz");
		$tpl->assign(".foo_d", array("steak", "beef", "pork", "chicken"));

		$contents1 = $tpl->get();

		ob_start();
		$tpl->render();
		$contents2 = ob_get_contents();
		ob_end_clean();

		$this->assertEquals($contents1, $contents2);

		$this->assertFalse($tpl->_find_path("404.tpl"));

		$tpl->set_compile_location("test/compile/", false);
		$this->assertEquals("test/templates/test/compile/", $tpl->_compile_path("test/templates/"));
		$tpl->set_compile_location("/test/compile/", true);
		$this->assertEquals("/test/compile/test/templates/", $tpl->_compile_path("test/templates/"));

		$retval = array();
		exec("rm -rf test/templates/test -rf", $retval);
	}

	public function testException()
	{
		$this->setExpectedException("Exception");
		$tpl = new Monotek\MiniTPL\Template;
		$tpl->set_compile_location("test/compile/", false);
		$tpl->set_paths("test/templates2/");
		$tpl->load("fail_to_compile.tpl");
	}

	/**
	 * @dataProvider varsProvider
	 */
	public function testVars($expression, $expected, $description)
	{
		global $tpl;
		$tpl = new Monotek\MiniTPL\Compiler;

		$result = $tpl->_split_exp($expression);
		$this->assertEquals($expected, $result);
	}

	public function testFailure()
	{
		$this->setExpectedException("Exception");
		$tpl = new Monotek\MiniTPL\Template;
		$this->assertFalse($tpl->load("missing.tpl"));
		$tpl->render();
	}

	public function varsProvider()
	{
		$vars = array();
		$vars[] = array('news_section_news_list.tpl', 'news_section_news_list.tpl', "normal string");
		$vars[] = array('$var', "\$_v['var']", "variable");
		$vars[] = array('$var.netko', "\$_v['var']['netko']", "array index");
		$vars[] = array('$var . "netko"', "\$_v['var'] . \"netko\"", "string concat");
		$vars[] = array('$var1 . $var2', "\$_v['var1'] . \$_v['var2']", "variable concat");
		$vars[] = array('$var1.$var2', "\$_v['var1'][\$_v['var2']]", "array var index");
		$vars[] = array('$items.0', "\$_v['items']['0']", "array int index");
		$vars[] = array('$tpl->get()', "\$tpl->get()", "global function");
		$vars[] = array('$tplx->get()', "\$_v['tplx']->get()", "object function");
		return $vars;
	}
}
