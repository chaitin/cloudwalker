<?php
/**
 * socket server for php ast
 * @author Mingyang.Liu
 */

/* Top level error reporting for debug */
error_reporting(E_ALL);

/* Allow the script to hang around waiting for connections. */
set_time_limit(0);

/* Turn on implicit output flushing so we see what we're getting
 * as it comes in. */
ob_implicit_flush();

/**
 * Parse php code, generate ast, and serialize to json format
 * (Require php extension 'php-ast', available in pecl)
 * @return string returns a json string, which consists of a json object of key 'status' and others
 * 
 * while 'status' is 'failed', it contains key 'reason' for syntactical error of the code
 * while 'status' is 'successed', it contains key 'ast'
 * 'status' is guaranteed to be only 'failed' or 'succeeded'
 * 
 * key 'ast' contains a zend_ast :
 *     zend_ast: {kind: int, flags: int, lineno: int, children: array of zend_ast}
 * 
 *   kind property specifies the type of the node. 
 *     It is an integral value, which corresponds to one of the ast\AST_* constants, 
 *     for example ast\AST_STMT_LIST. [See the AST node kinds section for an overview of 
 *     the available node kinds](https://github.com/nikic/php-ast#ast-node-kinds).
 * 
 *   flags property contains node specific flags. It is always defined, 
 *     but for most nodes it is always zero. See the flags section for a list of flags 
 *     supported by the different node kinds.
 * 
 *   lineno property specifies the starting line number of the node.
 * 
 *   children property contains an array of child-nodes. 
 *     These children can be either other ast\Node objects or plain values. 
 *     There are two general categories of nodes: Normal AST nodes, which have a fixed set 
 *     of named child nodes, as well as list nodes, which have a variable number of children. 
 *     The AST node kinds section contains a list of the child names for the different node kinds.
 */
function parseToJson(string $code) : string {

    $ast = null;

    try {
        $ast = ast\parse_code($code, $version=50);
    }
    catch(Throwable $e) {
        $err_msg = $e->getMessage();

        return json_encode(
            array(
                "status" => "failed", 
                "reason" => $err_msg
            ),
            JSON_PARTIAL_OUTPUT_ON_ERROR
        );
    }

    $json = json_encode(
        array(
            "status" => "successed",
            "ast" => $ast
        ),
        JSON_PARTIAL_OUTPUT_ON_ERROR,
        4096
    );

    if($json === false) {
        return json_encode(
            array(
                "status" => "failed",
                "reason" => "json encode error: ".json_last_error_msg()
            )
        );
    }
    
    return $json;

}

/**
 * php ast server : socket communication
 */
class PhpAstServer {

    /**
     * get header data
     * @return length of data
     */
    function getHeader() : int {
        fscanf(STDIN, "%d%*c", $length);
        if (!$length) {
             throw new Exception('message header must be number');
        }
        return $length;
    }

    /**
     * get body data
     * @param body_length:int could given by getHeader()
     */
    function getBody(int $body_length) : string {
        $receive_length = 0;
        $receive_buffer = "";
        while($receive_length < $body_length) {
            $receive_buffer .= fread(STDIN, $body_length - $receive_length);
            $receive_length = strlen($receive_buffer);
        }
        return $receive_buffer;
    }

    /**
     * write data to client
     */
    public function write(string $msg) {
        $msg = $msg."\n";
        $len_str = strval(strlen($msg))."\n";
        fwrite(STDOUT, $len_str);
        fwrite(STDOUT, $msg);
    }

    /**
     * loop to deal with request
     */
    public function loop() {
        while (!feof(STDIN))  {
            try {
                $total_length = $this->getHeader();
                $src = $this->getBody($total_length);
                $talk_back = parseToJson($src);
                $this->write($talk_back);
            }
            catch(Exception $exception) {
                $this->write(json_encode(
                    array(
                        "status" => "failed",
                        "reason" => $exception->getMessage()
                    )
                ));
            }
        };
    }
}

function main() {
    $server = new PhpAstServer();
    $server->loop();
}

main();

?>
