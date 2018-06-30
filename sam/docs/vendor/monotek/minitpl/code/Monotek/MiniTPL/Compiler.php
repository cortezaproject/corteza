<?php

namespace Monotek\MiniTPL;

/*

Tit Petrič, Monotek d.o.o., (cc) tit.petric@monotek.net
http://creativecommons.org/licenses/by-sa/3.0/

*/

/** Template compiler class */
class Compiler
{
	function __construct()
	{
		$this->_tag_php_open = "<"."?php";
		$this->_tag_php_close = "?".">\n";
		$this->_global_variables = array();
		$this->_literals = array();
	}

	protected function load_contents($filename)
	{
		if (($contents = file_get_contents($filename)) !== false) {
			if (substr($contents, 0, 3) == "\xEF\xBB\xBF") {
				return substr($contents, 3);
			}
		}
		return $contents;
	}

	/** Compile template file into php code */
	function compile($filename, $output_filename, $find_path, $nocache)
	{
		$contents = $this->load_contents($filename);
		$r = 0;
		if ($contents!==false && $contents!=="") {
			while (preg_match_all("/\{include\ (.*?)\}/s", $contents, $matches)) {
				$matches = array_unique($matches[1]);
				foreach ($matches as $file) {
					$cn = "<!-- ".$file." -->";
					if (($fn = call_user_func($find_path, $file))!==false) {
						$cn = $this->load_contents($fn.$file);
					}
					$contents = str_replace("{include ".$file."}", $cn, $contents);
				}
			}
			while (preg_match_all("/\{load\ (.*?)\}/s", $contents, $matches)) {
				$matches = array_unique($matches[1]);
				foreach ($matches as $file) {
					$file_var = (substr($file,0,1) == '$') ? $this->_split_exp($file) : '"'.$file.'"';
					$cn = $this->_code('$this->push();$this->load('.$file_var.');$this->assign($_v);$this->render();$this->pop();');
					$contents = str_replace("{load ".$file."}", $cn, $contents);
				}
			}
			$nocache = $nocache ? $this->_code("@unlink(__FILE__);") : "";
			$contents = str_replace("{*nocache*}",$nocache,$contents);
			$contents = $this->_strip_comments($contents);
			$contents = $this->_parse_constants($contents);
			$contents = $this->_parse_functions($contents, $filename);
			$contents = $this->_parse_expressions($contents);
			$contents = $this->_parse_variables($contents);
			if (!empty($this->_global_variables)) {
				$globals = array_unique($this->_global_variables);
				$contents = $this->_code('global '.implode(", ",$globals).';').$contents;
			}
			$contents = $this->_template_cleanup($contents);
			$this->_r_mkdir(dirname($output_filename));
			if ($f = @fopen($output_filename,"w")) {
				fwrite($f, $contents);
				fclose($f);
				$r = 1;
			}
		}
		return $r;
	}

	function _r_mkdir($dir)
	{
		if (file_exists($dir)) return;
		if (!file_exists(dirname($dir))) $this->_r_mkdir(dirname($dir));
		@mkdir($dir);
	}

	/** Insert system configuration, clean up code */
	function _template_cleanup($contents)
	{
		// set up variables
		$contents = $this->_code('$_v=&$this->vars;') . $contents;
		// strip unnecessary php tags
		$contents = str_replace($this->_tag_php_close.$this->_tag_php_open, "", $contents);
		// strip new line whitespace between php code
		$contents = str_replace($this->_tag_php_close."\n".$this->_tag_php_open.' ', "", $contents);
		$contents = str_replace("echo ;", "", $contents);
		foreach ($this->_literals as $key=>$value) {
			$contents = str_replace("[[".$key."]]", $value, $contents);
		}
		return $contents;
	}

	/** Strip template style comments */
	function _strip_comments($contents)
	{
		$contents = preg_replace("/\{\*.+\*\}/sU", "", $contents);
		return $contents;
	}

	/** Replace constant definitions */
	function _parse_constants($contents)
	{
		$matches = array();
		if (preg_match_all("/\{(\_[a-zA-Z0-9\_]+)\}/", $contents, $matches)) {
			$matches = array_unique($matches[1]);
			foreach ($matches as $m) {
				$contents = str_replace("{".$m."}", $this->_code("echo ".$m.";"), $contents);
			}
		}
		return $contents;
	}

	/** Search and replace for function blocks and inline definitions */
	function _parse_functions($contents, $filename)
	{
		$inlines = $blocks = array();

		if (preg_match_all("/\{(block|inline)\ ([a-zA-Z0-9\_\-]+)\}(.*?)\{\/\\1\}/s", $contents, $matches)) {
			foreach ($matches[0] as $k=>$ma) {
				$m = array("content"=>trim($matches[3][$k]), "src"=>$ma);
				if ($matches[1][$k]=="block") {
					$blocks[$matches[2][$k]] = $m;
				} else {
					$inlines[$matches[2][$k]] = $m;
				}
			}
		}

		if (preg_match_all("/\<script\ ([^\>]+)\>(.*?)\<\/script\>/s", $contents, $matches)) {
			foreach ($matches[0] as $k=>$parameters) {
				if (strpos($parameters,"text/template")!==false || strpos($parameters,"text/x-jquery")!==false) {
					$key = count($this->_literals)."_literal";
					$this->_literals[$key] = $matches[2][$k];
					$contents = str_replace($matches[2][$k], "[[".$key."]]", $contents);
				}
			}
		}

		foreach ($blocks as $name=>$code) {
			$lambda = sprintf("%u", crc32($code['content']))."_".sprintf("%u", crc32($filename));
			$block_code = "if (!function_exists('".$name."_".$lambda."')) { function ".$name."_".$lambda."(\$_v) {".$this->_tag_php_close.$code['content'].$this->_tag_php_open." } }";
			$contents = str_replace($code['src'], $this->_code($block_code), $contents);
		}

		foreach ($inlines as $name=>$code) {
			$contents = str_replace($code['src'], '', $contents);
		}

		$matches = array();
		while (preg_match_all("/\{inline\:([a-zA-Z0-9\_\-]+)\}/s", $contents, $matches)) {
			foreach ($matches[0] as $k=>$ma) {
				$contents = str_replace($ma, $inlines[$matches[1][$k]]['content'], $contents);
			}
		}

		foreach ($blocks as $name=>$code) {
			$contents = str_replace("{block:".$name."}", $this->_code($name."_".$lambda."(&\$_v);"), $contents);
		}
		return $contents;
	}

	/** Parse expression syntax: if, elseif, foreach, else, for, eval, eval_literal */
	function _parse_expressions($contents)
	{
		// foreach parsing
		if (preg_match_all("/\{foreach (.+)\}/sU", $contents, $matches)) {
			foreach ($matches[1] as $k=>$exp) {
				$exp = trim(trim($exp,"()"));
				list($e_left, $e_right) = explode(" as ", $exp);
				$e_right = explode("=>", $e_right);

				$left_exp = $this->_split_exp($e_left);
				$code = "";
				if (substr($left_exp, 0, 5) !== "array") {
					$code = "if(!empty(".$left_exp."))";
				}
				$code .= "foreach(".$left_exp." as ".$this->_split_exp($e_right[0]);
				if (count($e_right)==2) {
					$code .= '=>'.$this->_split_exp($e_right[1]);
				}
				$code .= '){';
				$contents = str_replace($matches[0][$k], trim($this->_code($code)), $contents);
			}
		}
		// if & for & elseif parsing
		if (preg_match_all("/\{(if|elseif|for|while) (.+)\}/sU", $contents, $matches)) {
			foreach ($matches[1] as $k=>$v) {
				if ($v=="for") {
					$matches[2][$k] = trim($matches[2][$k],"()");
				}
				$code = $v."(".$this->_split_exp($matches[2][$k])."){";
				if ($v=="elseif") {
					$code = "}".$code;
				}
				$contents = str_replace($matches[0][$k], $this->_code($code), $contents);
			}
		}
		// eval & eval_literal parsing
		if (preg_match_all("/\{(eval|eval_literal) (.+)\}/sU", $contents, $matches)) {
			foreach ($matches[1] as $k=>$type) {
				$code = rtrim(trim($matches[2][$k]),';');
				if ($type=="eval") {
					$code = $this->_split_exp($code);
				}
				$code .= ";";
				$contents = str_replace($matches[0][$k], $this->_code($code), $contents);
			}
		}
		$contents = str_replace("{else}", $this->_code("}else{"), $contents);
		$contents = str_replace(array("{/foreach}","{/while}","{/for}","{/if}"), trim($this->_code("}")), $contents);
		return $contents;
	}

	/** Parse variables */
	function _parse_variables($contents)
	{
		$mycontent = preg_replace("/\<\?php.+\?\>/sU","",$contents);
		// [a-zA-Z\_\$\"\'\[\]\ ]
		if (preg_match_all("/\{([^\{]+)\}/sU", $mycontent, $matches)) {
			foreach ($matches[1] as $k=>$v) {
				if (strstr($v,"\n")===false && $v{0}!=" ") {
					if ($v{0}!='$' && !in_array($v{0}, array("'",'"'))) {
						// shorthand variables {v}
						$v = '$'.$v;
					}
					$code = "";
					if (strstr($v,"|")!==false) {
						list($left,$right) = explode("|",$v);
						$left = $this->_split_exp($left);
						switch ($right) {
							case "toupper": $right = "strtoupper"; break;
							case "tolower": $right = "strtolower"; break;
							case "escape": $code = "echo htmlspecialchars(".$left.", ENT_QUOTES);"; break;
						}
						if ($code=='') {
							$code = "echo ".$right."(".$left.");";
						}
					} else {
						$code = "echo ".$this->_split_exp($v).";";
					}
					$contents = str_replace($matches[0][$k], $this->_code($code), $contents);
				}
			}
		}
		return $contents;
	}

	/** Split up variables from a php expression and replace them with actual variable locations */
	function _split_exp($exp)
	{
		$code = str_replace(".","__1","<"."?php if (".$exp.") { ?".">");
		$tokens = token_get_all($code);
		$objects = array();
		$variable = false;
		$variables = array();
		$variable_continues = false;
		foreach ($tokens as $k=>$v) {
			if (is_array($v)) {
				if ($v[0] == T_OBJECT_OPERATOR) {
					$variable_continues = false;
					$objects[] = $variable;
				}
				if ($v[0] == T_VARIABLE) {
					if (!$variable_continues && isset($variable) && !in_array($variable,$objects)) {
						$variables[] = $variable;
					}
					$variable = $variable_continues ? $variable.$v[1] : $v[1];
					if (strstr($variable,"__1")!==false) {
						$variable = str_replace("__1",".",$variable);
						$variable_continues = false;
						if (substr($variable,-1)==".") {
							$variable_continues = true;
						}
					} else if ($variable_continues) {
						$variable_continues = false;
					}
				}
				$v[0] = token_name($v[0]);
				$tokens[$k] = $v;
			}
		}
		if (isset($variable) && !in_array($variable,$variables) && !in_array($variable,$objects)) {
			$variables[] = $variable;
		}
		// globalize objects
		foreach ($objects as $object) {
			if ($object!='$this' && is_object($GLOBALS[substr($object,1)])) {
				$this->_global_variables[] = $object;
			} else {
				$variables[] = $object;
			}
		}

		// closure to sort vars by length and alphabetically
		usort($variables, function($a, $b) {
					if (strlen($a)==strlen($b)) {
						if ($a==$b) {
							return 0;
						}
						return ($a<$b) ? 1 : -1;
					}
					return (strlen($a)<strlen($b)) ? 1 : -1;
				} );

		foreach ($variables as $var) {
			if ($var != '$this') {
				$exp = str_replace($var, $this->_get_var($var), $exp);
			}
		}
		return $exp;
	}

	/** Helper function for replacing tags into actual variable locations */
	function _get_var($var) {
		$left_modifier = substr($var,1); // remove $
		$retval = $var;
		if ($var{0}!='"' && $var{0}!="'") {
			$retval = '$_v';
			if (strstr($left_modifier,'.')!==false) {
				// we have ourselves a table index
				$table_indices = explode('.',$left_modifier);
				foreach ($table_indices as $v) {
					$retval .= (($v{0}=='$') ? "[".$this->_get_var($v)."]" : "['".$v."']");
				}
			} else {
				$retval .= (($left_modifier{0}=='$') ? "[".$this->_get_var($left_modifier)."]" : "['".$left_modifier."']");
			}
		}
		return $retval;
	}

	/** Helper function for php code shorthand syntax, optimizing compiler size */
	function _code($s) {
		return $this->_tag_php_open." ".$s.$this->_tag_php_close;
	}
}
